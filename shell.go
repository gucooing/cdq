package cdq

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type shell struct {
	c *CDQ
}

func (c *CDQ) newShell() *shell {
	s := &shell{c: c}
	go s.Run()
	return s
}

func (s *shell) New(c *CDQ) {
	s.c = c
}

func (s *shell) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
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
		options := s.GenCommandOption(input, command)
		s.c.Log.Info("执行指令:%s,\n%s", command.Name, command.CommandFunc(options))
	}
}

func (s *shell) GenCommandOption(input string, command *Command) map[string]*CommandOption {
	options := make(map[string]*CommandOption, 0)
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return options
	}
	// 建立索引
	partMap := make(map[string]string)
	for index := 0; index < len(parts); index++ {
		part := parts[index]
		ids := strings.Split(part, ":")
		if len(ids) != 2 {
			continue
		}
		if _, ok := partMap[ids[0]]; ok {
			s.c.Log.Debug("重复的附加参数:%s", part)
			continue
		}
		partMap[ids[0]] = ids[1]
	}
	for _, op := range command.Options {
		part, ok := partMap[op.Name]
		if !ok {
			if op.Required {
				s.c.Log.Error("缺失必要的附加参数:%s", op.Name)
			}
			continue
		}
		options[op.Name] = &CommandOption{
			Name:   op.Name,
			Option: part,
		}
	}

	return options
}
