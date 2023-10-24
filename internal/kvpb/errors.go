package kvpb

import "errors"

// ErrKeyNotFound is an error returned when the provided key is not found in
// the underlying storage for the KV service.
var ErrKeyNotFound = errors.New("key not found")
