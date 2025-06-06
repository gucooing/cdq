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
			AliasList:   []string{"exit"},
			Description: "shell 内置指令,退出shell程序",
			Options:     nil,
			CommandFunc: func(options map[string]string) (string, error) {
				s.Exit()
				return "exit", nil
			},
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
		command := s.c.commandMap[parts[0]]
		if command == nil {
			s.c.Log.Error("不存在命令:%s", parts[0])
			continue
		}
		// 附加参数解析
		options, err := s.GenCommandOption(input, command)
		if err != nil {
			s.c.Log.Error(err.Error())
			continue
		}
		msg, err := command.CommandFunc(options)
		if err != nil {
			s.c.Log.Error("执行指令:%s 失败,错误%s", command.Name, err.Error())
			continue
		}
		s.c.Log.Info("执行指令:%s 成功,%s", command.Name, msg)
	}
}

func (s *Shell) GenCommandOption(input any, command *Command) (map[string]string, error) {
	options := make(map[string]string, 0)
	parts := strings.Fields(input.(string))
	for index, op := range command.Options {
		if op.Required && len(parts) < index+2 {
			return nil, errors.New(fmt.Sprintf("缺少必要参数:%s", op.Name))
		}
		if op.Required {
			if len(parts) < index+2 {
				return nil, errors.New(fmt.Sprintf("缺少必要参数:%s", op.Name))
			}
			options[op.Name] = parts[index+1]
		} else {
			if len(parts) < index+2 {
				continue
			}
			ids := strings.Split(parts[index+1], ":")
			if len(ids) != 2 {
				continue
			}
			options[ids[0]] = ids[1]
		}
	}

	return options, nil
}
