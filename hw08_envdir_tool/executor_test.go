package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("TEST-E1. Вызов невыполнимой команды", func(t *testing.T) {
		env := make(Environment)
		retVal := RunCmd([]string{"run_unknown_command", "param"}, env)
		require.Equal(t, -1, retVal, "Некорректный код возврата")
	})
	t.Run("TEST-E2. Ненулевой код возврата", func(t *testing.T) {
		env := make(Environment)
		retVal := RunCmd([]string{"cat", "unknown_file"}, env)
		require.NotEqual(t, 0, retVal, "Некорректный код возврата")
	})
}
