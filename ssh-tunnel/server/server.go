package server

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	sshtunnel "github.com/superwhys/venkit/ssh-tunnel"
	"github.com/superwhys/venkit/ssh-tunnel/sshtunnelpb"
)

var _ sshtunnelpb.SshTunnelServer = (*Server)(nil)

type ConnectType string

const (
	Forward ConnectType = "Forward"
	Reverse ConnectType = "Reverse"
)

type connect struct {
	cancel      func()
	typ         ConnectType
	local       netip.AddrPort
	localAlias  string
	remote      netip.AddrPort
	remoteAlias string
}

func (c *connect) Close() {
	c.cancel()
}

type Server struct {
	sshtunnelpb.UnimplementedSshTunnelServer

	sshTunnel *sshtunnel.SshTunnel
	cache     map[string]*connect
}

func NewSSHTunnelServer(tunnel *sshtunnel.SshTunnel) *Server {
	return &Server{
		sshTunnel: tunnel,
		cache:     make(map[string]*connect),
	}
}

func (s *Server) parseAddress(addr string) netip.AddrPort {
	return netip.MustParseAddrPort(addr)
}

func (s *Server) Forward(ctx context.Context, in *sshtunnelpb.ConnectRequest) (*sshtunnelpb.ForwardReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	c := &connect{
		cancel:      cancel,
		typ:         Forward,
		local:       s.parseAddress(in.Local),
		localAlias:  in.LocalAlias,
		remote:      s.parseAddress(in.Remote),
		remoteAlias: in.RemoteAlias,
	}

	if err := s.sshTunnel.Forward(ctx, c.local.String(), c.remote.String()); err != nil {
		return nil, errors.Wrap(err, "Forward")
	}

	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Wrap(err, "NewUUID")
	}

	s.cache[uid.String()] = c

	return &sshtunnelpb.ForwardReply{
		Uuid: uid.String(),
	}, nil
}

func (s *Server) Reverse(ctx context.Context, in *sshtunnelpb.ConnectRequest) (*sshtunnelpb.ReverseReply, error) {
	ctx, cancel := context.WithCancel(ctx)

	c := &connect{
		cancel:      cancel,
		typ:         Reverse,
		local:       s.parseAddress(in.Local),
		localAlias:  in.LocalAlias,
		remote:      s.parseAddress(in.Remote),
		remoteAlias: in.RemoteAlias,
	}

	if err := s.sshTunnel.Reverse(ctx, c.remote.String(), c.local.String()); err != nil {
		return nil, errors.Wrap(err, "Forward")
	}

	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Wrap(err, "NewUUID")
	}

	s.cache[uid.String()] = c

	return &sshtunnelpb.ReverseReply{
		Uuid: uid.String(),
	}, nil
}

func (s *Server) Disconnect(ctx context.Context, in *sshtunnelpb.DisconnectRequest) (*sshtunnelpb.DisconnectReply, error) {
	c, exists := s.cache[in.Uuid]
	if !exists {
		return nil, fmt.Errorf("Uid: %v connect not exists", in.Uuid)
	}

	c.Close()

	return &sshtunnelpb.DisconnectReply{}, nil
}
