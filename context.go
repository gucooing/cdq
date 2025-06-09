package cdq

import "math"

const abortIndex int8 = math.MaxInt8 >> 1

type Context struct {
	Flags    FlagMap    // 解析后的参数
	Command  *Command   // 指令信息
	App      CommandRun // 触发接口
	handlers Handlers   // 执行函数
	index    int8
	writ     interface {
		Return(code int, msg string)
	}
}

func newContext(a CommandRun, cmd *Command, flags FlagMap) *Context {
	return &Context{
		index:    -1,
		handlers: cmd.Handlers,
		App:      a,
		Command:  cmd,
		Flags:    flags,
	}
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func (c *Context) Return(code int, msg string) {
	c.writ.Return(code, msg)
}

func (c *Context) GetFlags() FlagMap {
	return c.Flags
}
