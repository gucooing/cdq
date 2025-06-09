//go:build !windows

package cdq

import (
	ctx "context"
	"os/exec"
)

func newShellCmd(ctx ctx.Context, command string) *exec.Cmd {
	return exec.CommandContext(ctx, "sh", "-c", command)
}
