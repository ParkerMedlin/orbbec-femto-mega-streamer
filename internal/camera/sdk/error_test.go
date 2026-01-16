//go:build windows && cgo

package sdk

import "testing"

func TestCheckErrorNil(t *testing.T) {
	if err := CheckError(nil); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
