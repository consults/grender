package main

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"strconv"
	"strings"
	"time"
)

func WaitLoadElement(page *rod.Page, xpath string, timeout int) bool {
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
			fmt.Printf("等待加载次数：【%d】 url：【%s】 \n", timeout, "url")
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
func ClickNextPage(page *rod.Page, xpath string, timeout int) {
	LoadDone := WaitLoadElement(page, xpath, timeout)
	if LoadDone == false {
		fmt.Println("点击按钮加载失败")

	}
	button, err := page.ElementX(xpath)
	if err != nil {
		panic("获取点击按钮失败")
	}
	ClickErr := button.Click(proto.InputMouseButtonLeft, 1)
	if ClickErr != nil {
		panic("点击下一页失败")
	}
}
func Add(elements rod.Elements, count *int) {
	for _, element := range elements {
		n, err := strconv.Atoi(strings.ReplaceAll(element.MustText(), "\n", ""))
		if err != nil {
			panic(err)
		}
		*count = *count + n
		fmt.Println(n)
	}
}
func main() {
	startUrl := "https://www.python-spider.com/challenge/23"
	TargetXpath := "//*[@class='info']"
	ButtonXpath := "//*[@class='xl-nextPage']"
	WaitElementTimeOUt := time.Millisecond * 1000
	l := launcher.New().NoSandbox(true).Headless(false)
	browser := rod.New().ControlURL(l.MustLaunch())
	browser.MustConnect()

	expr := proto.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour).Unix())
	cooks := make([]*proto.NetworkCookieParam, 0)
	cook := &proto.NetworkCookieParam{
		Name:     "sessionid",
		Value:    "2ppulij5db9ghllsnrt4bvo3m68uxluh",
		Domain:   "www.python-spider.com",
		HTTPOnly: true,
		Expires:  expr,
	}
	cooks = append(cooks, cook)
	cookErr := browser.SetCookies(cooks)
	if cookErr != nil {
		fmt.Println("设置cook失败")
		panic(cookErr.Error())
	}
	count := 0
	var page *rod.Page
	for i := 1; i < 101; i++ {
		var result bool
		if i == 1 {
			page = browser.MustPage(startUrl)

			time.Sleep(time.Second * 3)
			//page.HandleDialog()
			result = WaitLoadElement(page, TargetXpath, 10)
			if result == false {
				panic("渲染页面失败")
			}
			eList := page.MustElementsX(TargetXpath)
			Add(eList, &count)
			fmt.Printf("当前总和：【%d】 \n", count)
			ClickNextPage(page, ButtonXpath, 10)
			time.Sleep(WaitElementTimeOUt)
			fmt.Println("当前页：【1】")
			continue
		}
		result = WaitLoadElement(page, TargetXpath, 10)
		if result == false {
			panic(fmt.Sprintf("渲染页面:【%d】 失败 \n", i))
		}
		eList := page.MustElementsX(TargetXpath)
		Add(eList, &count)
		fmt.Printf("当前总和：【%d】 \n", count)
		ClickNextPage(page, ButtonXpath, 10)
		fmt.Printf("当前页：【%d】 \n", i)
		time.Sleep(WaitElementTimeOUt)

	}
	fmt.Println("end")
	fmt.Println(count)
	fmt.Println("end")
	utils.Pause()
}
