package render

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
	log "github.com/sirupsen/logrus"
	"grender/core/configReader"
	"time"
)

var RodRender *Render

type Render struct {
	Launcher *launcher.Launcher
	Browser  *rod.Browser
	PagePool chan *rod.Page
}

func InitRender() {
	RodRender = &Render{}
	if configReader.Config.Render.Local == true {
		log.Warningln("使用本地浏览器")
		RodRender.Launcher = launcher.New().NoSandbox(true).Headless(false)
	} else {
		log.Warningf("使用远程浏览器：【%s】\n", configReader.Config.Render.RodAddress)
		RodRender.Launcher = launcher.MustNewManaged(configReader.Config.Render.RodAddress)
	}
	if configReader.Config.Proxy.Open {
		log.Warningf("代理地址：【%s】 \n", configReader.Config.Proxy.ProxyAddress)
		RodRender.Launcher.Proxy(configReader.Config.Proxy.ProxyAddress)
	}
	browser := rod.New().ControlURL(RodRender.Launcher.MustLaunch())
	RodRender.Browser = browser
	browser.Timeout(time.Second * 20)
	browser.MustConnect()
	//defer browser.MustClose()
	// 添加代理
	if configReader.Config.Proxy.Open {
		log.Warningf("代理认证：【%s】【%s】 \n", configReader.Config.Proxy.ProxyUser, configReader.Config.Proxy.ProxyPassword)
		go browser.HandleAuth(configReader.Config.Proxy.ProxyUser, configReader.Config.Proxy.ProxyPassword)()
	}
	RodRender.PagePool = make(chan *rod.Page, configReader.Config.Render.PoolSize)
	for i := 0; i < 10; i++ {
		page := stealth.MustPage(browser)
		RodRender.PagePool <- page
	}
}
func WaitLoadElement(page *rod.Page, url string, xpath string, timeout int) bool {
	fmt.Println(url)
	page.MustNavigate(url)
	page.MustWaitLoad()
	if timeout == 0 {
		timeout = 10
	}
	jsContent := `()=>{
	result = document.evaluate("%s", document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null);
	element = result.singleNodeValue;
	if (element===null){
	return false
	}else{
	return true
	}
}`
	result := false
	for {
		execJs := fmt.Sprintf(jsContent, xpath)
		e, findErr := page.Eval(execJs)
		if findErr != nil {
			fmt.Println(findErr.Error())
			fmt.Println("执行js失败")
			time.Sleep(time.Second * 1)
			break
		}
		if e.Value.String() == "false" {
			fmt.Printf("等待加载次数：【%d】 url：【%s】 \n", timeout, url)
			time.Sleep(time.Second * 1)
			timeout = timeout - 1
			if timeout <= 0 {
				break
			}
			continue
		}
		if e.Value.String() == "true" {
			result = true
			break
		}

	}
	return result
}

func GetHtml(page *rod.Page, url string, xpath string, timeout int) string {
	defer func() {
		RodRender.PagePool <- page
	}()
	html := ""
	renderResult := WaitLoadElement(page, url, xpath, timeout)
	if renderResult {
		fmt.Printf("渲染成功：%s \n", url)
		html = page.MustHTML()
	} else {
		fmt.Printf("渲染失败：%s \n", url)
	}
	return html
}
