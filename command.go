package cdq

import (
	"errors"
	"fmt"
)

type Handler func(c *Context)
type Handlers []Handler

type Command struct {
	Name        string           // 指令
	AliasList   []string         // 别名
	Description string           // 描述
	Permissions Permissions      // 需要的权限
	Options     []*CommandOption // 附加参数
	Handlers    Handlers         // 执行函数
}

type CommandOption struct {
	Name        string          // 指令
	Alias       string          // 别名
	Description string          // 描述
	Required    bool            // 是否必要参数
	Default     string          // 默认值
	ExpectedS   []string        // 预期值
	expected    map[string]bool // 预期值
}

// CommandRun 指令执行接口
type CommandRun interface {
	Run()
	Exit()
	GenCommandOption(input any, command *Command) (*Context, error) // 生成附加参数
}

func AddHandlers(ojbs ...Handler) Handlers {
	return ojbs
}

func (co *CommandOption) genFlagMapItem(v, k string) (*FlagMapItem, error) {
	if co.Required && k == "" {
		return nil, errors.New(fmt.Sprintf("缺少必要参数:%s", v))
	}
	fi := &FlagMapItem{}
	if k == "" {
		k = co.Default
	}
	if len(co.expected) > 0 && !co.expected[k] {
		return nil, errors.New(fmt.Sprintf("参数%s 传入:%s 不存在该类型", v, k))
	}
	fi.Value = k
	return fi, nil
}

func orString(strs ...string) string {
	for _, str := range strs {
		if str != "" {
			return str
		}
	}
	return ""
}
