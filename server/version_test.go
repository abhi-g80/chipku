//version_test.gp
package server

import (
	"testing"
)

var version string = "0.2.0"

func TestVersion(t *testing.T) {
	if version != Version {
		t.Errorf("expected version %s, got %s", version, Version)
	}
}
