package main

import (
	"context"

	sshtunnel "github.com/superwhys/venkit/ssh-tunnel"
	"github.com/superwhys/venkit/vflags"
)

func main() {
	vflags.Parse()
	tunnel := sshtunnel.NewTunnel(&sshtunnel.SshConfig{
		User:         "hoven",
		HostName:     "10.11.43.115",
		IdentityFile: "/Users/yong/.ssh/id_rsa_cnns",
	})

	defer tunnel.Close()

	if err := tunnel.Reverse(context.TODO(), ":28081", "localhost:8080"); err != nil {
		panic(err)
	}

	tunnel.Wait()
}
