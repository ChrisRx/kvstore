package boltkv_test

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"

	"github.com/ChrisRx/kvstore/internal/boltkv"
	"github.com/ChrisRx/kvstore/internal/kvpb"
)

func TestSetValue(t *testing.T) {
	addr := ":9090"

	b, err := boltkv.NewBoltKV("testdata/data.db")
	if err != nil {
		t.Fatal(err)
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer()
	kvpb.RegisterKVServer(s, b)

	go func() {
		s.Serve(l)
	}()

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	client := kvpb.NewKVClient(conn)

	if _, err := client.Set(context.Background(), &kvpb.SetRequest{
		Key:   "testing1",
		Value: "somedata",
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := client.Get(context.Background(), &kvpb.GetRequest{
		Key: "testing1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Value != "somedata" {
		t.Fatalf("expected key testing1 value to be 'somedata', received %q", resp.Value)
	}
}
