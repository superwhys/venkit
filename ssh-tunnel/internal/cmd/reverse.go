/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/service"
	"github.com/superwhys/venkit/ssh-tunnel/server"
	"github.com/superwhys/venkit/ssh-tunnel/sshtunnelpb"
	"github.com/superwhys/venkit/vflags"
	"google.golang.org/grpc"
)

// ReverseCmd represents the list command
var ReverseCmd = &cobra.Command{
	Use:   "reverse -r <address:port> [flags] -l <address:port> [flags]",
	Short: "Proxy a remote accessible address to a local address port",

	RunE: func(cmd *cobra.Command, args []string) error {
		vflags.Parse()

		tunnel := newTunnel()
		s := server.NewSSHTunnelServer(tunnel)

		ms := service.NewVkService(
			service.WithServiceName("ssh-proxy"),
			service.WithGrpcServer(func(srv *grpc.Server) {
				sshtunnelpb.RegisterSshTunnelServer(srv, s)
			}),
			service.WithGrpcUI(),
			service.WithWorker(func(ctx context.Context) error {
				in := &sshtunnelpb.ConnectRequest{
					Local:  localAddr,
					Remote: remoteAddr,
				}

				resp, err := s.Reverse(ctx, in)
				if err != nil {
					return err
				}

				table := make(map[string][]string)
				table[resp.Uuid] = append(
					table[resp.Uuid],
					[]string{string(server.Reverse), fmt.Sprintf("%v -> %v", remoteAddr, localAddr)}...,
				)
				lg.Info("Connected services\n" + prettyMaps(table))

				<-ctx.Done()
				tunnel.Close()
				return nil
			}),
		)
		ms.Run(port())

		return nil
	},
}

func init() {
	ReverseCmd.Flags().StringVarP(&localAddr, "local", "l", "", "Local accessible address port")
	ReverseCmd.Flags().StringVarP(&remoteAddr, "remote", "r", "", "Remote address port")
	ReverseCmd.MarkFlagsRequiredTogether("local", "remote")
}
