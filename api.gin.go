package cdq

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"net/http"
	"runtime"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

type GinApi struct {
	c       *CDQ
	Router  *gin.Engine
	ApiKey  map[string]bool
	server  *http.Server
	ginSync *sync.Mutex
}

func NewGinApi(c *CDQ, ginSync *sync.Mutex) *GinApi {
	a := &GinApi{
		c:       c,
		ginSync: ginSync,
		ApiKey:  make(map[string]bool),
	}
	if a.ginSync == nil {
		a.ginSync = &sync.Mutex{}
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
	GinApiCodeOk        GinApiCode = iota
	GinApiCodeCmdErr               // 不存在该指令
	GinApiCodeOptionErr            // 参数错误
	GinApiCodeErr                  // 其他错误,详情看msg
)

func (a *GinApi) GetApi(c *gin.Context) {
	a.ginSync.Lock()
	defer a.ginSync.Unlock()
	resp := &GinApiResponse{
		Code: GinApiCodeOk,
		Msg:  "",
	}
	defer c.JSON(200, resp)
	cmd := c.Query("cmd")
	command := a.c.commandMap[cmd]
	if command == nil {
		resp.Code = GinApiCodeCmdErr
		resp.Msg = fmt.Sprintf("不存在命令:%s", cmd)
		return
	}
	// 附加参数解析
	options, err := a.GenCommandOption(c, command)
	if err != nil {
		resp.Code = GinApiCodeOptionErr
		resp.Msg = err.Error()
		return
	}
	msg, err := command.CommandFunc(options)
	if err != nil {
		resp.Code = GinApiCodeErr
		resp.Msg = err.Error()
		return
	}
	resp.Msg = fmt.Sprintf("%s", msg)
}

func (a *GinApi) GenCommandOption(input any, command *Command) (map[string]string, error) {
	c := input.(*gin.Context)
	options := make(map[string]string, 0)
	for _, op := range command.Options {
		part := c.Query(op.Name)
		if op.Required {
			if part == "" {
				return nil, errors.New(fmt.Sprintf("缺少必要参数:%s", op.Name))
			}
		}
		options[op.Name] = part
	}

	return options, nil
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

func (a *GinApi) shell(c *gin.Context) {
	command := c.Query("shell")
	if command == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'cmd' query parameter"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := newShellCmd(ctx, command)

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
