package cmd_test

import (
	"testing"

	"github.com/abhi-g80/chipku/cmd"
)

func TestRootCmd(t *testing.T) {
	err := cmd.Execute()

	if err != nil {
		t.Errorf("got error: %v", err)
	}
}
