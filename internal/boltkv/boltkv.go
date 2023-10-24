// Package boltkv provides an implementation of the gRPC service KV using
// boltdb as a backend.
package boltkv

import (
	bolt "go.etcd.io/bbolt"
)

// BoltKV represents a boltdb connection and implements the KVServer interface,
// allowing it to be used by the KV gRPC service.
type BoltKV struct {
	*bolt.DB

	bucket []byte
}

// NewBoltKV takes the path and any functional options, and constructs a new
// BoltKV.
func NewBoltKV(path string, options ...BoltKVOption) (_ *BoltKV, err error) {
	b := &BoltKV{
		bucket: []byte(DefaultBoltBucketName),
	}
	for _, o := range options {
		o(b)
	}
	b.DB, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	if err := b.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(b.bucket)
		return err
	}); err != nil {
		return nil, err
	}
	return b, nil
}
