package sshtunnel

import (
	"context"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/superwhys/venkit/lg"
	"golang.org/x/crypto/ssh"
)

type SshConfig struct {
	HostName     string
	User         string
	IdentityFile string
}

func (sc *SshConfig) SetDefaults() {
	if sc.IdentityFile == "" {
		sc.IdentityFile = os.Getenv("HOME") + "/.ssh/id_rsa"
	}
	if !strings.Contains(sc.HostName, ":") {
		sc.HostName += ":22"
	}
	if sc.User == "" {
		sc.User = os.Getenv("USER")
	}
}

func getIdentifyKey(filePath string) (ssh.Signer, error) {
	buff, _ := os.ReadFile(filePath)
	return ssh.ParsePrivateKey(buff)
}

func (sc *SshConfig) ParseClientConfig() (*ssh.ClientConfig, error) {
	key, err := getIdentifyKey(sc.IdentityFile)
	if err != nil {
		return nil, err
	}

	return &ssh.ClientConfig{
		User: sc.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

type SshTunnel struct {
	conf      *SshConfig
	sshClient *ssh.Client
	wg        sync.WaitGroup
}

func NewTunnel(cf *SshConfig) *SshTunnel {
	cf.SetDefaults()

	tunnel := &SshTunnel{
		conf: cf,
	}
	client, err := tunnel.dial()
	lg.PanicError(err)
	tunnel.sshClient = client
	go tunnel.keepAlive()

	return tunnel
}

func (st *SshTunnel) Wait() {
	st.wg.Wait()
}

func (st *SshTunnel) dial() (*ssh.Client, error) {
	clientConf, err := st.conf.ParseClientConfig()
	if err != nil {
		return nil, err
	}

	client, err := ssh.Dial("tcp", st.conf.HostName, clientConf)
	if err != nil {
		return nil, err
	}
	lg.Infof("SSH dial [%v] success", st.conf.HostName)
	return client, nil
}

func (st *SshTunnel) keepAlive() {
	tick := time.NewTicker(15 * time.Second)
	defer tick.Stop()

	for range tick.C {
		// send keep alive request
		_, _, err := st.sshClient.SendRequest("keepalive@golang.org", true, nil)
		if err != nil {
			// if error in send request
			// try to reconnect
			if st.sshClient != nil {
				// before retry connection
				// close old connection if exists
				st.sshClient.Close()
			}
			st.sshClient = nil

			for {
				lg.Error("ssh connection lost, try to reconnect", err)
				newClient, err := st.dial()
				if err != nil {
					lg.Error("dial ssh server", err)
					time.Sleep(time.Second * 3)
					continue
				}
				st.sshClient = newClient
				break
			}
		}
	}
}

func (st *SshTunnel) Reverse(ctx context.Context, remoteAddr, localAddr string) error {
	st.wg.Add(1)

	ctx = lg.With(ctx, "[SSHReverse]")

	if strings.HasPrefix(remoteAddr, ":") {
		remoteAddr = "0.0.0.0" + remoteAddr
	}

	remoteLst, err := st.sshClient.Listen("tcp", remoteAddr)
	if err != nil {
		return errors.Wrapf(err, "listen on remote addr %s", remoteAddr)
	}

	lg.Infof("listen remote %v success", remoteAddr)

	go func() {
		defer func() {
			lg.Infoc(ctx, "disconnected reversing %s to %s", remoteAddr, localAddr)
		}()
		defer st.wg.Done()
		defer remoteLst.Close()

		for {
			if err := ctx.Err(); err != nil {
				return
			}

			if remoteLst == nil {
				if st.sshClient == nil {
					lg.Warnc(ctx, "SSH connections lost")
					continue
				}

				newLst, err := st.sshClient.Listen("tcp", remoteAddr)
				if err != nil {
					lg.Errorc(ctx, "SSH listen redial failed, err: %v", err)
					continue
				}
				lg.Debugc(ctx, "SSH listen redial success -> %v", remoteAddr)
				remoteLst = newLst
			}

			remote, err := remoteLst.Accept()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					lg.Debugc(ctx, "remote timeout listening")
					continue
				}
				lg.Warnc(ctx, "Remote listening accept error, Redial... : %v", err)
				if remoteLst != nil {
					remoteLst.Close()
				}
				remoteLst = nil
				continue
			}
			lg.Debugc(ctx, "remote %s accept connection from %s", remote.LocalAddr().String(), remote.RemoteAddr())

			go func(remote net.Conn) {
				local, err := net.Dial("tcp", localAddr)
				if err != nil {
					lg.Errorc(ctx, "dial local addr %s error: %v", localAddr, err)
					return
				}

				lg.Debugc(ctx, "start handle remote %s via %s to local %s", remote.RemoteAddr(), remote.LocalAddr(), localAddr)
				st.handleClient(ctx, remote, local)
				lg.Debugc(ctx, "end handle remote %s via %v to local %s", remote.RemoteAddr(), remote.LocalAddr(), localAddr)
			}(remote)
		}
	}()

	return nil
}

func (st *SshTunnel) Forward(ctx context.Context, localAddr, remoteAddr string) error {
	st.wg.Add(1)

	ctx = lg.With(ctx, "[SSHForward]")

	// start listen on local addr
	localLst, err := net.Listen("tcp", localAddr)
	if err != nil {
		return errors.Wrapf(err, "listen on local addr %s", localAddr)
	}

	go func() {
		defer func() {
			lg.Infoc(ctx, "disconnected forwarding %s to %s", localAddr, remoteAddr)
		}()
		defer st.wg.Done()
		defer localLst.Close()
		for {
			if err := ctx.Err(); err != nil {
				return
			}

			if localLst == nil {
				if st.sshClient == nil {
					lg.Warnc(ctx, "SSH connections lost")
					continue
				}

				newLocalLst, err := net.Listen("tcp", localAddr)
				if err != nil {
					lg.Errorc(ctx, "local listen redial failed, err: %v", err)
					return
				}
				localLst = newLocalLst
			}

			local, err := localLst.Accept()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				lg.Warnc(ctx, "Local listening accept error, Redial... : %v", err)
				if localLst != nil {
					localLst.Close()
				}
				localLst = nil
				continue
			}
			lg.Debugc(ctx, "local %s accept connection from %s", local.LocalAddr(), local.RemoteAddr())

			go func(local net.Conn) {
				defer local.Close()
				if st.sshClient == nil {
					lg.Errorc(ctx, "lost ssh connection")
					return
				}

				remote, err := st.sshClient.Dial("tcp", remoteAddr)
				if err != nil {
					lg.Errorc(ctx, "dial remote addr %s error: %v", remoteAddr, err)
					return
				}

				lg.Debugc(ctx, "start handle local %s connection to remote %s", local.LocalAddr(), remoteAddr)
				st.handleClient(ctx, local, remote)
				lg.Debugc(ctx, "end handle local %s connection to remote %s", local.LocalAddr(), remoteAddr)
			}(local)
		}
	}()
	return nil
}

func (st *SshTunnel) handleClient(ctx context.Context, local, remote net.Conn) {
	defer local.Close()
	defer remote.Close()

	ctx, cancel := context.WithCancel(ctx)

	// remote -> local transfer
	go func() {
		_, err := io.Copy(local, remote)
		if err != nil {
			lg.Warnc(ctx, "remote -> local error: %v", err)
		}
		cancel()
	}()

	// local -> remote transfer
	go func() {
		_, err := io.Copy(remote, local)
		if err != nil {
			lg.Errorc(ctx, "local -> remote error: %v", err)
		}
		cancel()
	}()
	<-ctx.Done()
}

func (st *SshTunnel) Close() error {
	lg.Debugf("SSH tunnel [%v] close", st.conf.HostName)
	st.wg.Done()

	return st.sshClient.Close()
}
