package cdq

import (
	"gucooing/cdq/logger"
)

type CDQ struct {
	Shell          bool                // 是否启用内置的shell执行器
	Log            logger.Logger       // log
	commandList    []*Command          // 注册的指令
	commandRunList []CommandRun        // 注册的指令器
	commandMap     map[string]*Command // 注册指令集合 - all name
}

// New c: cdq实例,commandList:注册的外置执行器
func New(c *CDQ, commandList ...CommandRun) *CDQ {
	if c == nil {
		c = new(CDQ)
	}
	if c.Log == nil {
		c.Log = logger.NewLog(logger.LevelDebug, nil)
	}
	if c.commandRunList == nil {
		c.commandRunList = make([]CommandRun, 0)
	}

	if c.Shell { // 是否启用内置的Shell
		s := c.newShell()
		c.commandRunList = append(c.commandRunList, s)
	}

	for _, cmd := range commandList { // 启动外置的执行器
		cmd.New(c)
		c.commandRunList = append(c.commandRunList, cmd)
		go cmd.Run()
	}
	// 注册默认命令
	c.ApplicationCommandHelp()

	return c
}

func (c *CDQ) ApplicationCommand(command *Command) {
	if c.commandMap == nil {
		c.commandMap = make(map[string]*Command)
	}
	if c.commandList == nil {
		c.commandList = make([]*Command, 0)
	}
	c.commandList = append(c.commandList, command)
	for _, name := range command.AliasList {
		if _, ok := c.commandMap[name]; ok {
			c.Log.Error("指令:%s ,别名:%s 重复注册", command.Name, name)
			continue
		}
		c.commandMap[name] = command
	}

	c.Log.Debug("注册指令名称:%s,别名:%s,描述:%s", command.Name, command.AliasList, command.Description)
}
