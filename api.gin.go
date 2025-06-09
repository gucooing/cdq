package cdq

import (
	ctx "context"
	"errors"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"net/http"
	"runtime"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

type GinApi struct {
	c      *CDQ
	Router *gin.Engine
	ApiKey map[string]bool
	server *http.Server
}

func NewGinApi(c *CDQ) *GinApi {
	a := &GinApi{
		c:      c,
		ApiKey: make(map[string]bool),
	}

	c.Log.Debug("启用GinApi指令")
	return a
}

func (a *GinApi) NewRouter(addr string, debug bool) {
	gin.SetMode(gin.ReleaseMode)
	if debug {
		a.Router = gin.Default()
	} else {
		a.Router = gin.New()
	}
	a.Router.Use(gin.Recovery())
	a.SetRouter(a.Router)
	a.server = &http.Server{Addr: addr, Handler: a.Router}
}

func (a *GinApi) SetRouter(router *gin.Engine) {
	a.Router = router
	a.Router.GET("/cdq/api", a.AutoGucooingApi, a.GetApi)
	a.Router.GET("/cdq/api/help", a.AutoGucooingApi, a.Help)
	a.Router.GET("/cdq/api/shell", a.AutoGucooingApi, a.shell)
}

func (a *GinApi) SetApiKey(key ...string) {
	for _, k := range key {
		a.ApiKey[k] = true
	}
}

func (a *GinApi) Run() {
	if a.server == nil {
		return
	}
	err := a.server.ListenAndServe()
	if err != nil {
		a.c.Log.Error("gin api 服务器运行失败:%s", err.Error())
		return
	}
}

func (a *GinApi) Exit() {
	if a.server == nil {
		return
	}
	a.server.Close()
}

type GinApiResponse struct {
	Code GinApiCode `json:"code"`
	Msg  string     `json:"msg"`
}

type GinApiCode = int

const (
	ApiCodeOk            GinApiCode = iota
	ApiCodeCmdUnknown               // 不存在该指令
	ApiCodeOptionUnknown            // 参数错误
	GinApiCodeErr                   // 其他错误,详情看msg
)

func (a *GinApi) GetApi(ginc *gin.Context) {
	c := &GinApiContext{
		c: ginc,
	}
	cmd := ginc.Query("cmd")
	command := a.c.commandMap[cmd]
	if command == nil {
		c.Return(ApiCodeCmdUnknown, fmt.Sprintf("不存在命令:%s", cmd))
		return
	}
	// 附加参数解析
	ctxs, err := a.GenCommandOption(ginc, command)
	if err != nil {
		c.Return(ApiCodeOptionUnknown, err.Error())
		return
	}
	c.Context = ctxs
	ctxs.writ = c
	ctxs.Next()
}

type GinApiContext struct {
	*Context
	c *gin.Context
}

func (a *GinApi) GenCommandOption(input any, command *Command) (*Context, error) {
	ginc := input.(*gin.Context)
	flags := make(FlagMap)
	for _, op := range command.Options {
		v := orString(ginc.Query(op.Name), ginc.Query(op.Alias))
		if v == "" && op.Required {
			return nil, errors.New(fmt.Sprintf("缺少必要参数:%s", op.Name))
		}
		fi, err := op.genFlagMapItem(op.Alias, v)
		if err != nil {
			return nil, err
		}
		flags[op.Name] = fi
	}

	return newContext(a, command, flags), nil
}

func (c *GinApiContext) Return(code int, msg string) {
	c.c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
	})
}

func (a *GinApi) AutoGucooingApi(c *gin.Context) {
	if len(a.ApiKey) == 0 ||
		a.ApiKey[c.GetHeader("Authorization")] {
		return
	} else {
		c.String(401, "Unauthorized")
		c.Abort()
	}
}

func (a *GinApi) Help(c *gin.Context) {
	var returnstr string
	for _, comm := range a.c.commandList {
		returnstr += "----------------------------------\n"
		returnstr += fmt.Sprintf("命令:%s 描述:%s 别名:%s", comm.Name, comm.Description, comm.AliasList)
		example := fmt.Sprintf("用法:/cdq/api?cmd=%s", comm.Name)
		var opt string
		for _, op := range comm.Options {
			example += fmt.Sprintf("&%s=msg", op.Name)
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
	c.String(200, returnstr)
}

func (a *GinApi) shell(c *gin.Context) {
	command := c.Query("shell")
	if command == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'cmd' query parameter"})
		return
	}

	ctxs, cancel := ctx.WithTimeout(ctx.Background(), 10*time.Second)
	defer cancel()

	cmd := newShellCmd(ctxs, command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) > 0 {
			utf8Output := convertToUTF8(output)
			c.String(http.StatusOK, utf8Output)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Command execution failed: " + err.Error(),
		})
		return
	}
	utf8Output := convertToUTF8(output)
	c.String(http.StatusOK, utf8Output)
}

func convertToUTF8(data []byte) string {
	if utf8.Valid(data) {
		return string(data)
	}

	var decoder *encoding.Decoder
	switch runtime.GOOS {
	case "windows":
		decoder = simplifiedchinese.GBK.NewDecoder()
	default:
		return string(data)
	}

	result, _, err := transform.Bytes(decoder, data)
	if err != nil {
		return string(data)
	}
	return string(result)
}
