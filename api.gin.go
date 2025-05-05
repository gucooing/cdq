package cdq

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type GinApi struct {
	c       *CDQ
	Router  *gin.Engine
	ApiKey  string
	server  *http.Server
	ginSync *sync.Mutex
}

func NewGinApi(c *CDQ, ginSync *sync.Mutex) *GinApi {
	a := &GinApi{
		c:       c,
		ginSync: ginSync,
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
	a.Router.GET("/cdq/api", a.AutoGucooingApi, a.GetApi)
	a.Router.Use(gin.Recovery())
	a.server = &http.Server{Addr: addr, Handler: a.Router}
}

func (a *GinApi) SetRouter(router *gin.Engine) {
	a.Router = router
	a.Router.GET("/cdq/api", a.AutoGucooingApi, a.GetApi)
}

func (a *GinApi) SetApiKey(key string) {
	a.ApiKey = key
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

func (a *GinApi) GenCommandOption(input any, command *Command) (map[string]*CommandOption, error) {
	c := input.(*gin.Context)
	options := make(map[string]*CommandOption, 0)
	for _, op := range command.Options {
		part := c.Query(op.Name)
		if op.Required {
			if part == "" {
				return nil, errors.New(fmt.Sprintf("缺少必要参数:%s", op.Name))
			}
		} else {
			if part == "" {
				continue
			}
		}
		options[op.Name] = &CommandOption{
			Name:   op.Name,
			Option: part,
		}
	}

	return options, nil
}

func (a *GinApi) AutoGucooingApi(c *gin.Context) {
	if a.ApiKey == "" ||
		c.GetHeader("Authorization") == a.ApiKey {
		return
	} else {
		c.String(401, "Unauthorized")
		c.Abort()
	}
}
