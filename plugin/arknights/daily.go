package arknights

import (
	"github.com/playwright-community/playwright-go"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func daily(ctx *zero.Ctx) {
	url := "https://prts.wiki/"
	pw, err := playwright.Run()
	if err != nil {
		playwright.Install()
		playwright.Run()
	}
	defer pw.Stop()
	browser, err := pw.Chromium.Launch()
	if err != nil {
		playwright.Install()
	}
	page, err := browser.NewPage()
	if err != nil {
		return
	}
	_, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(10000),
	})
	if err != nil {
		return
	}
	page.SetViewportSize(1920, 1080)
	result, _ := page.QuerySelector("div[class=\"mp-container-left\"]")
	page.WaitForSelector("img")
	result.ScrollIntoViewIfNeeded()
	box, _ := result.BoundingBox()
	PageScreenshotOptions := playwright.PageScreenshotOptions{
		Clip: &playwright.PageScreenshotOptionsClip{
			X:      playwright.Float(float64(box.X - 25)),
			Y:      playwright.Float(float64(box.Y - 55)),
			Width:  playwright.Float(float64(box.Width - 130)),
			Height: playwright.Float(float64(box.Width + 70)),
		},
	}
	b, err := page.Screenshot(PageScreenshotOptions)
	if err != nil {
		ctx.SendChain(message.Text("出错了，稍后再试试吧~"), message.Text("\n - 数据来源于 https://prts.wiki/"))
	}
	ctx.SendChain(message.ImageBytes(b))
}
