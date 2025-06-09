package cdq

import (
	"fmt"
)

// applicationCommandHelp 默认指令
func (c *CDQ) applicationCommandHelp() {
	c.ApplicationCommand(&Command{
		Name:        "help",
		AliasList:   []string{"h"},
		Description: "有关某个命令的详细信息，请键入 help 命令名",
		Permissions: Guest,
		Handlers:    AddHandlers(c.Help),
		Options: []*CommandOption{
			{
				Name:        "c",
				Description: "显示该命令的帮助信息。",
				Required:    false,
			},
		},
	})
}

func (c *CDQ) Help(ctx *Context) {
	var returnstr string
	if commSrt := ctx.GetFlags().String("c"); commSrt == "" {
		returnstr += "有关某个命令的详细信息，请键入 help -c 命令名\n"
		for _, comm := range c.commandList {
			returnstr += fmt.Sprintf(
				"%s---别名:%s------%s\n",
				comm.Name,
				comm.AliasList,
				comm.Description,
			)
		}
	} else {
		comm, ok := c.commandMap[commSrt]
		if !ok {
			returnstr += "不支持此命令\n"
		} else {
			returnstr += fmt.Sprintf("命令:%s 描述:%s 别名:%s", comm.Name, comm.Description, comm.AliasList)
			example := comm.Name
			var opt string
			for _, op := range comm.Options {
				example += fmt.Sprintf(" -%s msg", op.Name)
				opt += fmt.Sprintf("      %s - 描述:%s -别名:%s", op.Name, op.Description, op.Alias)
				if op.Required {
					opt += " -必要参数"
				} else {
					opt += " -非必要参数"
				}
				if op.Default != "" {
					opt += fmt.Sprintf(" -默认值:%s", op.Default)
				}
				if len(op.ExpectedS) > 0 {
					opt += fmt.Sprintf(" -可选参数:%s", op.ExpectedS)
				}
				opt += "\n"
			}
			returnstr += fmt.Sprintf("\n%s", example)
			returnstr += fmt.Sprintf("\n%s\n", opt)
		}
	}
	ctx.Return(0, returnstr)
}
