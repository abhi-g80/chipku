// version_test.gp
package server

import (
	"testing"
)

func TestVersion(t *testing.T) {
	var version string = "1.2.1"

	if version != Version {
		t.Errorf("expected version %s, got %s", version, Version)
	}
}
