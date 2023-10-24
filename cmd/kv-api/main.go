package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/ChrisRx/kvstore/internal/kvpb"
	"github.com/ChrisRx/kvstore/internal/restapi"
)

var opts struct {
	Addr     string
	KVAddr   string
	Insecure bool

	// gRPC TLS client auth
	KeyFile    string
	CertFile   string
	CACertFile string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kv-api",
		Short:         "Run KV REST API server",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var copts []grpc.DialOption
			if opts.Insecure {
				copts = append(copts, grpc.WithInsecure())
			}
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
				}
				cert, err := tls.LoadX509KeyPair(opts.CertFile, opts.KeyFile)
				if err != nil {
					return err
				}
				tlsConfig.Certificates = []tls.Certificate{cert}
				copts = append(copts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
			}
			conn, err := grpc.Dial(opts.KVAddr, copts...)
			if err != nil {
				return err
			}
			e, err := restapi.New(kvpb.NewKVClient(conn))
			if err != nil {
				return err
			}
			e.Logger.Fatal(e.Start(opts.Addr))
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.Addr, "addr", ":8080", "Address for API server")
	cmd.Flags().StringVar(&opts.KVAddr, "kv-addr", ":9090", "Client address for KV gRPC server")
	cmd.Flags().StringVar(&opts.CertFile, "cert-file", "", "TLS auth cert file")
	cmd.Flags().StringVar(&opts.KeyFile, "key-file", "", "TLS auth key file")
	cmd.Flags().StringVar(&opts.CACertFile, "ca-file", "", "TLS auth CA cert file")
	cmd.Flags().BoolVar(&opts.Insecure, "insecure", false, "Allow insecure gRPC client connection")
	return cmd
}

func main() {
	if err := NewCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
