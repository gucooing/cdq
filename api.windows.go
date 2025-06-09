//go:build windows

package cdq

import (
	ctx "context"
	"os/exec"
	"syscall"
)

func newShellCmd(ctx ctx.Context, command string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "cmd", "/C", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	return cmd
}
