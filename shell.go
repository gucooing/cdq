package cdq

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Shell struct {
	c    *CDQ
	done chan bool
}

func NewShell(c *CDQ) *Shell {
	s := &Shell{
		c:    c,
		done: make(chan bool),
	}
	c.Log.Debug("启用Shell指令")
	c.ApplicationCommand(
		&Command{
			Name:        "exit",
			AliasList:   []string{"quit"},
			Description: "shell 内置指令,退出shell程序",
			Options:     nil,
			Handlers: AddHandlers(func(c *Context) {
				s.Exit()
			}),
		})
	return s
}

func (s *Shell) Exit() {
	s.done <- true
}

func (s *Shell) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-s.done:
			return
		default:
		}
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		input = strings.TrimSpace(input)
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}
		c := &ShellContext{
			s: s,
		}
		command := s.c.commandMap[parts[0]]
		if command == nil {
			c.Return(ApiCodeCmdUnknown, fmt.Sprintf("不存在命令:%s", parts[0]), nil)
			continue
		}
		// 附加参数解析
		ctxs, err := s.GenCommandOption(parts[1:], command)
		if err != nil {
			c.Return(ApiCodeOptionUnknown, err.Error(), nil)
			continue
		}
		c.Context = ctxs
		ctxs.writ = c
		ctxs.Next()
	}
}

type ShellContext struct {
	*Context
	s *Shell
}

func (s *Shell) GenCommandOption(args any, command *Command) (*Context, error) {
	options := make(map[string]string, 0)
	parts := args.([]string)
	flags := make(FlagMap)
	for i := 0; i < len(parts)/2; i++ {
		index := i * 2
		switch string(parts[index][0]) {
		case "-":
			options[parts[index][1:]] = parts[index+1]
		default:
			options[parts[index]] = parts[index+1]
		}
	}
	for _, op := range command.Options {
		v := orString(options[op.Alias], options[op.Name])
		if v == "" && op.Required {
			return nil, errors.New(fmt.Sprintf("缺少必要参数:%s", op.Name))
		}
		fi, err := op.genFlagMapItem(op.Alias, v)
		if err != nil {
			return nil, err
		}
		flags[op.Name] = fi
	}

	return newContext(s, command, flags), nil
}

func (s *ShellContext) Return(code int, message string, data interface{}) {
	if code != ApiCodeOk {
		s.s.c.Log.Debug(message)
	} else {
		s.s.c.Log.Info(message)
	}
}
