package boltkv

import (
	"context"
	"fmt"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ChrisRx/kvstore/internal/kvpb"
)

// Get takes the provided GetRequest, containing a key to lookup, and returns
// the key and value from the underlying boltdb storage. If the key is not set,
// an error is returned.
func (b *BoltKV) Get(ctx context.Context, req *kvpb.GetRequest) (*kvpb.GetResponse, error) {
	resp := &kvpb.GetResponse{
		Key: req.Key,
	}
	if err := b.View(func(tx *bolt.Tx) error {
		value := tx.Bucket(b.bucket).Get([]byte(req.Key))
		if value == nil {
			return fmt.Errorf("%w: %s", kvpb.ErrKeyNotFound, req.Key)
		}
		resp.Value = string(value)
		return nil
	}); err != nil {
		return nil, err
	}
	return resp, nil
}

// Set takes the provided SetRequest, containing a key and the desired value,
// and sets this in the underlying boltdb storage. If a key is already set, the
// value will be overwritten.
func (b *BoltKV) Set(ctx context.Context, req *kvpb.SetRequest) (*emptypb.Empty, error) {
	err := b.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(b.bucket).Put([]byte(req.Key), []byte(req.Value))
	})
	return &emptypb.Empty{}, err
}

// Delete takes the provided DeleteRequest, containing a key to lookup, and
// deletes the key/value from the underlying boltdb storage. This will not
// check if a key is set before attempting to delete.
func (b *BoltKV) Delete(ctx context.Context, req *kvpb.DeleteRequest) (*emptypb.Empty, error) {
	err := b.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(b.bucket).Delete([]byte(req.Key))
	})
	return &emptypb.Empty{}, err
}
