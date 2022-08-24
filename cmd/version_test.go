package cmd_test

import (
	"testing"

	"github.com/abhi-g80/chipku/cmd"
)

func TestVersionCmd(t *testing.T) {
	var args []string = []string{"version"}
	cmd.RootCmd.SetArgs(args)
	err := cmd.RootCmd.Execute()

	if err != nil {
		t.Errorf("got error: %v", err)
	}
}
