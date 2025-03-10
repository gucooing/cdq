package cdq

import (
	"fmt"
)

// applicationCommandHelp 默认指令
func (c *CDQ) applicationCommandHelp() {
	c.ApplicationCommand(&Command{
		Name:        "help",
		AliasList:   []string{"help", "h"},
		Description: "有关某个命令的详细信息，请键入 help 命令名",
		Permissions: Guest,
		CommandFunc: c.help,
		Options: []*CommandOption{
			{
				Name:        "c",
				Description: "显示该命令的帮助信息。",
				Required:    false,
			},
		},
	})
}

func (c *CDQ) help(options map[string]*CommandOption) (string, error) {
	command := options["c"]
	var returnstr string
	if command == nil {
		returnstr += "有关某个命令的详细信息，请键入 help c:命令名\n"
		for _, comm := range c.commandList {
			returnstr += fmt.Sprintf(
				"%s---别名:%s------%s\n",
				comm.Name,
				comm.AliasList,
				comm.Description,
			)
		}
	} else {
		comm, ok := c.commandMap[command.Option]
		if !ok {
			returnstr += "不支持此命令\n"
		} else {
			returnstr += comm.Description + "\n"
			returnstr += fmt.Sprintf("别名:%s\n", comm.AliasList)
			example := comm.Name
			var opt string
			for _, option := range comm.Options {
				if !option.Required {
					example += fmt.Sprintf(" [%s:msg]", option.Name)
				} else {
					example += fmt.Sprintf(" %s", option.Name)
				}
				opt += fmt.Sprintf("      %s - %s\n", option.Name, option.Description)
			}
			returnstr += fmt.Sprintf("\n%s", example)
			returnstr += fmt.Sprintf("\n%s", opt)
		}
	}

	return returnstr, nil
}
