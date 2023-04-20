package api

import (
	"io"
	"net/http"
	"os"
	"time"
	"xyhelper-web/config"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/xyhelper/chatgpt-go"
)

// Session
func Session(c *gin.Context) {
	auth := false
	if os.Getenv("AUTH_SECRET_KEY") != "" {
		auth = true
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "",
		"data": gin.H{
			"auth":  auth,
			"model": "ChatGPTUnofficialProxyAPI",
		},
	})
}

// VerifyRequest
type VerifyRequest struct {
	Token string `json:"token" binding:"required"`
}

// Verify
func Verify(c *gin.Context) {
	req := &VerifyRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})

		return
	}
	if req.Token == os.Getenv("AUTH_SECRET_KEY") {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Success",
			"message": "",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Error",
			"message": "Token 错误",
			"data":    nil,
		})
	}
}

// ChatProcessRequest
type ChatProcessRequest struct {
	Prompt string `json:"prompt" binding:"required"`
	Optins *struct {
		ConversationId  string `json:"conversationId"`  // 会话ID
		ParentMessageId string `json:"parentMessageId"` // 父消息ID
	} `json:"options"` // 选项
	BaseURI     string `json:"baseURI"`     // 基础URI
	AccessToken string `json:"accessToken"` // 访问令牌
	IsGPT4      bool   `json:"isGPT4"`      // 是否为GPT4
}

// ChatProcessResponse
type ChatProcessResponse struct {
	Role            string `json:"role"`            // 角色
	Id              string `json:"id"`              // 消息ID
	ParentMessageId string `json:"parentMessageId"` // 父消息ID
	ConversationId  string `json:"conversationId"`  // 会话ID
	Text            string `json:"text"`            // 消息内容
}

// ChatProcess 响应
func ChatProcess(c *gin.Context) {
	if os.Getenv("AUTH_SECRET_KEY") != "" {
		Authorization := c.GetHeader("Authorization")
		if Authorization != "Bearer "+os.Getenv("AUTH_SECRET_KEY") {
			c.JSON(http.StatusOK, gin.H{
				"status":  "Unauthorized",
				"message": "Token 错误",
				"data":    nil,
			})
			return
		}
	}

	var req ChatProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})

		return
	}
	// g.DumpWithType(req)

	cli := chatgpt.NewClient(
		chatgpt.WithAccessToken(req.AccessToken),
		chatgpt.WithTimeout(time.Duration(config.TimeOutMs*1000*1000)),
		chatgpt.WithBaseURI(req.BaseURI),
	)
	if req.IsGPT4 {
		cli.SetModel("gpt-4")
	}
	stream, err := cli.GetChatStream(req.Prompt, req.Optins.ConversationId, req.Optins.ParentMessageId)
	// 如果返回404，说明会话不存在，重新获取会话
	if err != nil {
		if err.Error() == "send message failed: 404 Not Found" {
			stream, err = cli.GetChatStream(req.Prompt)
		}
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	res := &ChatProcessResponse{}

	// 使用 Stream 方法向客户端发送 SSE 数据
	c.Stream(func(w io.Writer) bool {
		for text := range stream.Stream {
			// g.DumpWithType(text)
			res.Id = text.MessageID
			res.Text = text.Content
			res.Role = "assistant"
			res.ConversationId = text.ConversationID
			res.ParentMessageId = req.Optins.ParentMessageId
			data := gjson.New(res).MustToJson()
			writeSSEData(w, data)

		}
		return false
	})

}

func writeSSEData(w io.Writer, data []byte) error {
	_, err := w.Write(append(data, byte('\n')))
	if err != nil {
		return err
	}
	w.(http.Flusher).Flush()

	return nil
}

// Message
func Config(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "",
		"data": gin.H{
			"apiModel":     "ChatGPTUnofficialProxyAPI",
			"reverseProxy": "https://freechat.xyhelper.cn/backend-api/conversation",
			"timeoutMs":    gconv.String(config.TimeOutMs/1000) + "秒",
			"socksProxy":   "-",
			"httpsProxy":   "-",
			"balance":      "-",
			"version":      config.Version,
		},
	})
}
