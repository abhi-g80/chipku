package cmd

import (
	"testing"
)

func TestRootCmd(t *testing.T) {
	t.Run("TestRootCmdCall", func(t *testing.T) {
		err := Execute()

		if err != nil {
			t.Errorf("got error: %v", err)
		}
	})
}
