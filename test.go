package main

import (
	"net/http"

	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/exec"
	_ "github.com/henrylee2cn/pholcus_lib" // 此为公开维护的spider规则库
	// _ "pholcus_lib_pte" // 同样你也可以自由添加自己的规则库
)

func init() {
	Spider{
		Name:        "静态规则示例",
		Description: "静态规则示例 [Auto Page] [http://xxx.xxx.xxx]",
		// Pausetime: 300,
		// Limit:   LIMIT,
		// Keyin:   KEYIN,
		EnableCookie:    true,
		NotDefaultField: false,
		Namespace:       nil,
		SubNamespace:    nil,
		RuleTree: &RuleTree{
			Root: func(ctx *Context) {
				ctx.AddQueue(&request.Request{Url: "http://xxx.xxx.xxx", Rule: "登录页"})
			},
			Trunk: map[string]*Rule{
				"登录页": {
					ParseFunc: func(ctx *Context) {
						ctx.AddQueue(&request.Request{
							Url:      "http://xxx.xxx.xxx",
							Rule:     "登录后",
							Method:   "POST",
							PostData: "username=123456@qq.com&password=123456&login_btn=login_btn&submit=login_btn",
						})
					},
				},
				"登录后": {
					ParseFunc: func(ctx *Context) {
						ctx.Output(map[string]interface{}{
							"全部": ctx.GetText(),
						})
						ctx.AddQueue(&request.Request{
							Url:    "http://accounts.xxx.xxx/member",
							Rule:   "个人中心",
							Header: http.Header{"Referer": []string{ctx.GetUrl()}},
						})
					},
				},
				"个人中心": {
					ParseFunc: func(ctx *Context) {
						ctx.Output(map[string]interface{}{
							"全部": ctx.GetText(),
						})
					},
				},
			},
		},
	}.Register()
}
func main() {
	// 设置运行时默认操作界面，并开始运行
	// 运行软件前，可设置 -a_ui 参数为"web"、"gui"或"cmd"，指定本次运行的操作界面
	// 其中"gui"仅支持Windows系统
	exec.DefaultRun("web")
}
