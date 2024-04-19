package main

import (
	"context"

	sshtunnel "github.com/superwhys/venkit/ssh-tunnel"
)

func main() {
	tunnel := sshtunnel.NewTunnel(&sshtunnel.SshConfig{
		User:         "yong",
		HostName:     "10.15.15.123",
		IdentityFile: "/Users/yong/.ssh/id_rsa",
	})

	if err := tunnel.Forward(context.TODO(), "localhost:29950", "10.15.15.15.231:80"); err != nil {
		panic(err)
	}

	tunnel.Wait()
}
