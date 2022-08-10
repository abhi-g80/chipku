//version_test.gp
package main

import (
	"testing"
)

var version string = "0.1.1"

func TestVersion(t *testing.T) {
	if version != Version {
		t.Errorf("expected version %s, got %s", version, Version)
	}
}
