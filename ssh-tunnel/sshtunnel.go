package sshtunnel

import (
	"context"
	"io"
	"io/ioutil"
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
	buff, _ := ioutil.ReadFile(filePath)
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

func (st *SshTunnel) Forward(ctx context.Context, localAddr, remoteAddr string) error {
	st.wg.Add(1)

	// start listen on local addr
	local, err := net.Listen("tcp", localAddr)
	if err != nil {
		return errors.Wrapf(err, "listen on local addr %s", localAddr)
	}

	go func() {
		defer func() {
			lg.Infoc(ctx, "disconnected forwarding %s to %s", localAddr, remoteAddr)
		}()
		defer st.wg.Done()
		defer local.Close()
		for {
			if err := ctx.Err(); err != nil {
				return
			}
			// accept connection from local listener
			client, err := local.Accept()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() && netErr.Temporary() {
					// continue if timeout
					continue
				}
				lg.Errorc(ctx, "local accept error: %v, Redialing...", err)
				if local != nil {
					local.Close()
				}
				newLocal, err := net.Listen("tcp", localAddr)
				if err != nil {
					lg.Errorc(ctx, "local listen error: %v", err)
					return
				}
				local = newLocal
				continue
			}
			lg.Debugc(ctx, "local %s accept connection from %s", client.LocalAddr().String(), client.RemoteAddr().String())

			// dial remote addr and handle local client connections data to remote server
			go func(client net.Conn) {
				defer client.Close()
				if st.sshClient == nil {
					lg.Errorc(ctx, "lost ssh connection")
					return
				}

				remote, err := st.sshClient.Dial("tcp", remoteAddr)
				if err != nil {
					lg.Errorc(ctx, "dial remote addr %s error: %v", remoteAddr, err)
					return
				}

				lg.Debugc(ctx, "start handle local %s connection to remote %s", client.LocalAddr().String(), remoteAddr)
				st.handleClient(ctx, client, remote)
				lg.Debugc(ctx, "end handle local %s connection to remote %s", client.LocalAddr().String(), remoteAddr)
			}(client)
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
			lg.Errorc(ctx, "remote -> local error: %v", err)
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
