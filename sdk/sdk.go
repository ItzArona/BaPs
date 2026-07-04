package sdk

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/config"
)

type SDK struct {
	router *gin.Engine
}

func New(router *gin.Engine) *SDK {
	s := &SDK{
		router: router,
	}

	s.initRouter()
	return s
}

func (s *SDK) initRouter() {
	s.router.LoadHTMLGlob(fmt.Sprintf("%s/templates/*", config.GetConfig().DataPath))

	// 根路径: 连接性检测 / HTML 首页
	s.router.Any("/", s.rootHandler)
	s.router.Any("/index", handleIndex)

	// 连接性检测: prod-clientpatch.bluearchiveyostar.com/test.txt
	s.router.GET("/test.txt", func(c *gin.Context) {
		c.String(200, "ok")
	})

	// 服务器列表 ServerInfo (/r:url) + 旧版 prod/index.json
	s.router.GET("/r:url", s.connectionGroups)
	s.router.GET("/prod/index.json", index)

	account := s.router.Group("/account")
	{
		account.POST("/yostar_auth_request", s.YostarAuthRequest)
		account.POST("/yostar_auth_submit", s.YostarAuthSubmit)
	}
	user := s.router.Group("/user")
	{
		user.POST("/yostar_createlogin", s.YostarCreatelogin)
		user.POST("/login", s.YostarLogin)
		// 新版客户端 SDK 快速登录
		user.POST("/quick-login", s.QuickLogin)
		user.Any("/agreement", agreement)
	}
	// 新版客户端 SDK 路由: /common/*
	common := s.router.Group("/common")
	{
		common.Any("/config", commonConfig)
		common.Any("/version", commonVersion)
	}
	app := s.router.Group("/app")
	{
		app.Any("/getSettings", getSettings)
		app.Any("/getCode", getCode)
	}
}

// rootHandler 根路径处理器,根据请求来源区分:
// - X-Original-Host 含 yostarplat (连接性检测) -> 返回 {"Code":200,"Data":"Hello world","Msg":"OK"}
// - 其他 -> 返回 HTML 首页
// X-Original-Host 由 mitmproxy 重定向脚本注入,保留重定向前的原始域名
func (s *SDK) rootHandler(c *gin.Context) {
	origHost := c.GetHeader("X-Original-Host")
	if strings.Contains(origHost, "yostarplat.com") ||
		strings.Contains(origHost, "yostar-serverinfo") {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(200, `{"Code":200,"Data":"Hello world","Msg":"OK"}`)
		return
	}
	handleIndex(c)
}

// commonVersion 新版客户端版本检查
// 真实响应: {"Code":200,"Data":{"Agreement":[...],"ErrorCode":"5.5"},"Msg":"OK"}
func commonVersion(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(200, `{"Code":200,"Data":{"Agreement":[{"Version":"0.1","Type":"user_agreement","Title":"用户协议","Content":"","Lang":"ja"},{"Version":"0.1","Type":"privacy_agreement","Title":"隐私政策","Content":"","Lang":"ja"}],"ErrorCode":"5.5"},"Msg":"OK"}`)
}

func handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":  "Ba Ps!",
		"github": "https://github.com/gucooing/BaPs",
	})
}
