package client

import (
	"testing"
)

func TestAPISetDeleteValue(t *testing.T) {
	expected := randomString(10)
	url := "http://localhost:8080/testing1"

	// set a new key/value
	if err := post(url, "testing1", expected); err != nil {
		t.Fatal(err)
	}

	// ensure it set successfully
	v, err := get(url)
	if err != nil {
		t.Fatal(err)
	}
	if v != expected {
		t.Fatalf("expected key testing1 value to be %q, received %q", expected, v)
	}

	// delete key
	if err := delete(url, "testing1"); err != nil {
		t.Fatal(err)
	}

	// check if key delete was successfully
	v, err = get(url)
	if err != nil {
		t.Fatal(err)
	}
	if v != "" {
		t.Fatalf("key was not deleted")
	}
}
