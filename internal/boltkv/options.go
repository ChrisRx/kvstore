package boltkv

// BoltKVOption represents the function signature for any BoltKV functional
// options.
type BoltKVOption func(*BoltKV)

// DefaultBoltBucketName is the default bucket name used for the underlying
// boltdb for BoltKV.
const DefaultBoltBucketName = "kv"

// WithBucketName is a functional option for specifying a boltdb bucket name
// other than the default.
func WithBucketName(bucket string) BoltKVOption {
	return func(b *BoltKV) {
		b.bucket = []byte(bucket)
	}
}
