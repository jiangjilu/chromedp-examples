package main

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// 创建 Chrome 上下文
func createChromeContext() (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), // 无头模式
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx)
	return ctx, cancel
}

// 访问网页
func navigateHandler(c context.Context, ctx *app.RequestContext) {
	url := string(ctx.QueryArgs().Peek("url"))
	if url == "" {
		ctx.JSON(consts.StatusBadRequest, map[string]string{"error": "Missing 'url' parameter"})
		return
	}

	// 创建 Chrome 上下文
	chromeCtx, cancel := createChromeContext()
	defer cancel()

	err := chromedp.Run(chromeCtx, chromedp.Navigate(url))
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"message": "Page loaded successfully",
		"url":     url,
	})
}

// 获取网页标题
func titleHandler(c context.Context, ctx *app.RequestContext) {
	url := string(ctx.QueryArgs().Peek("url"))
	if url == "" {
		ctx.JSON(consts.StatusBadRequest, map[string]string{"error": "Missing 'url' parameter"})
		return
	}

	// 创建 Chrome 上下文
	chromeCtx, cancel := createChromeContext()
	defer cancel()

	var title string
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(url),
		chromedp.Title(&title),
	)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"url":   url,
		"title": title,
	})
}

// 获取网页截图
func screenshotHandler(c context.Context, ctx *app.RequestContext) {
	url := string(ctx.QueryArgs().Peek("url"))
	if url == "" {
		ctx.JSON(consts.StatusBadRequest, map[string]string{"error": "Missing 'url' parameter"})
		return
	}

	// 创建 Chrome 上下文
	chromeCtx, cancel := createChromeContext()
	defer cancel()

	var screenshot []byte
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second), // 等待页面加载
		chromedp.CaptureScreenshot(&screenshot),
	)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	encoded := base64.StdEncoding.EncodeToString(screenshot)
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"url":        url,
		"screenshot": encoded,
	})
}

func main() {
	h := server.New(server.WithHostPorts(":8080"))

	h.GET("/navigate", navigateHandler)
	h.GET("/title", titleHandler)
	h.GET("/screenshot", screenshotHandler)

	hlog.Infof("Server running at http://localhost:8080")
	h.Spin()
}
