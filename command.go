package cdq

type Command struct {
	Name        string                                         // 指令
	AliasList   []string                                       // 别名
	Description string                                         // 描述
	Permissions Permissions                                    // 需要的权限
	Options     []*CommandOption                               // 附加参数
	CommandFunc func(options map[string]*CommandOption) string // 执行函数
}

type CommandOption struct {
	Name        string // 指令
	Description string // 描述
	Option      string // 选项
	Required    bool   // 是否必要参数
}

// CommandRun 指令执行接口
type CommandRun interface {
	Run()
	Exit()
	GenCommandOption(input any, command *Command) (map[string]*CommandOption, error) // 生成附加参数
}
