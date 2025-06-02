//go:build !windows

package cdq

import (
	"context"
	"os/exec"
)

func newShellCmd(ctx context.Context, command string) *exec.Cmd {
	return exec.CommandContext(ctx, "sh", "-c", command)
}
