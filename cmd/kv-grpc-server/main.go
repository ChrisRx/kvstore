package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/ChrisRx/kvstore/internal/boltkv"
	"github.com/ChrisRx/kvstore/internal/kvpb"
)

var opts struct {
	Addr   string
	DBFile string

	// gRPC TLS server auth
	KeyFile    string
	CertFile   string
	CACertFile string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kv-grpc-server",
		Short:         "Run KV service gRPC server",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := boltkv.NewBoltKV(opts.DBFile)
			if err != nil {
				return err
			}
			l, err := net.Listen("tcp", opts.Addr)
			if err != nil {
				return err
			}
			var sopts []grpc.ServerOption
			if opts.KeyFile != "" && opts.CertFile != "" {
				tlsConfig := &tls.Config{
					ClientCAs: x509.NewCertPool(),
				}

				if opts.CACertFile != "" {
					caCert, err := os.ReadFile(opts.CACertFile)
					if err != nil {
						return err
					}
					if !tlsConfig.ClientCAs.AppendCertsFromPEM(caCert) {
						return fmt.Errorf("failed to add ca cert file: %s", opts.CACertFile)
					}

					// In the case that a CA certificate is provided, this must
					// be used to verify the clients.
					tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
				}
				cert, err := tls.LoadX509KeyPair(opts.CertFile, opts.KeyFile)
				if err != nil {
					return err
				}
				tlsConfig.Certificates = []tls.Certificate{cert}
				sopts = append(sopts, grpc.Creds(credentials.NewTLS(tlsConfig)))
			}
			s := grpc.NewServer(sopts...)
			kvpb.RegisterKVServer(s, b)
			fmt.Printf("Running kv-grpc-server on %s ...\n", opts.Addr)
			s.Serve(l)
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.Addr, "addr", ":9090", "Address for gRPC server")
	cmd.Flags().StringVar(&opts.DBFile, "db-file", "data.db", "Name of boltdb file")
	cmd.Flags().StringVar(&opts.CertFile, "cert-file", "", "TLS auth cert file")
	cmd.Flags().StringVar(&opts.KeyFile, "key-file", "", "TLS auth key file")
	cmd.Flags().StringVar(&opts.CACertFile, "ca-file", "", "TLS auth CA cert file")
	return cmd
}

func main() {
	if err := NewCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
