package cdq

import (
	"github.com/gucooing/cdq/logger"
)

type CDQ struct {
	Log            logger.Logger       // log
	commandList    []*Command          // 注册的指令
	commandRunList []CommandRun        // 注册的指令器
	commandMap     map[string]*Command // 注册指令集合 - all name
}

// New c: cdq实例,commandList:注册的外置执行器
func New(c *CDQ) *CDQ {
	if c == nil {
		c = new(CDQ)
	}
	if c.Log == nil {
		c.Log = logger.NewLog(logger.LevelDebug, nil)
	}

	// 注册默认命令
	c.applicationCommandHelp()

	return c
}

func (c *CDQ) AddCommandRun(commandList ...CommandRun) {
	if c.commandRunList == nil {
		c.commandRunList = make([]CommandRun, 0)
	}
	for _, cmd := range commandList { // 启动外置的执行器
		c.commandRunList = append(c.commandRunList, cmd)
		go cmd.Run()
	}
}

func (c *CDQ) ApplicationCommand(command *Command) {
	if c.commandMap == nil {
		c.commandMap = make(map[string]*Command)
	}
	if c.commandList == nil {
		c.commandList = make([]*Command, 0)
	}
	c.commandList = append(c.commandList, command)
	addName := func(name string) {
		if _, ok := c.commandMap[name]; ok {
			c.Log.Error("指令:%s ,别名:%s 重复注册", command.Name, name)
			return
		}
		c.commandMap[name] = command
	}
	addName(command.Name)
	for _, name := range command.AliasList {
		addName(name)
	}
	for _, op := range command.Options {
		op.expected = make(map[string]bool)
		for _, x := range op.ExpectedS {
			op.expected[x] = true
		}
	}

	c.Log.Debug("注册指令名称:%s,别名:%s,描述:%s", command.Name, command.AliasList, command.Description)
}
