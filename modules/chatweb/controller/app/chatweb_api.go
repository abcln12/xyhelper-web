package app

import (
	"context"

	"github.com/cool-team-official/cool-admin-go/cool"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type ChatwebApiController struct {
	*cool.ControllerSimple
}

func init() {
	var chatweb_api_controller = &ChatwebApiController{
		&cool.ControllerSimple{
			Perfix: "/app/chatweb/api",
		},
	}
	// 注册路由
	cool.RegisterControllerSimple(chatweb_api_controller)
}

// 增加 Welcome 演示 方法
type ChatwebApiWelcomeReq struct {
	g.Meta `path:"/welcome" method:"GET"`
}
type ChatwebApiWelcomeRes struct {
	*cool.BaseRes
	Data interface{} `json:"data"`
}

func (c *ChatwebApiController) Welcome(ctx context.Context, req *ChatwebApiWelcomeReq) (res *ChatwebApiWelcomeRes, err error) {
	res = &ChatwebApiWelcomeRes{
		BaseRes: cool.Ok("Welcome to Cool Admin Go"),
		Data:    gjson.New(`{"name": "Cool Admin Go", "age":0}`),
	}
	return
}

// SessionReq 请求参数
type ChatwebApiSessionReq struct {
	g.Meta `path:"/session" method:"POST"`
}

// SessionRes 返回参数
type ChatwebApiSessionRes struct {
	Status  string `json:"status"`  // 状态
	Message string `json:"message"` // 消息
	Data    *struct {
		Auth  bool   `json:"auth"`  // 认证
		Model string `json:"model"` // 模型
	} `json:"data"` // 数据
}

// Session 会话
func (c *ChatwebApiController) Session(ctx context.Context, req *ChatwebApiSessionReq) (res *ChatwebApiSessionRes, err error) {
	res = &ChatwebApiSessionRes{
		Status:  "Success",
		Message: "",
		Data: &struct {
			Auth  bool   `json:"auth"`
			Model string `json:"model"`
		}{
			Auth:  false,
			Model: "ChatGPTUnofficialProxyAPI",
		},
	}
	return
}
