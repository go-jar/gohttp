package main

import (
	"fmt"
	"github.com/go-jar/golog"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-jar/gohttp/controller"
	"github.com/go-jar/gohttp/gracehttp"
	"github.com/go-jar/gohttp/router"
	"github.com/go-jar/gohttp/system"
)

func main() {
	demoController := new(DemoController)
	testController := new(TestController)

	logger, err := golog.NewConsoleLogger(golog.LevelDebug)
	if err != nil {
		fmt.Println(err)
	}
	simpleRouter := router.NewSimpleRouter(logger)

	simpleRouter.DefineRoute("/test/args/([0-9]+)$", testController, "Args")
	simpleRouter.DefineRoute("/demo/args/([0-9]+)$", testController, "Args")
	simpleRouter.RegisterRoutes(demoController)

	sys := system.NewSystem(simpleRouter)
	_ = gracehttp.ListenAndServe(":8010", sys)
}

type DemoController struct {
}

func (dc *DemoController) NewActionContext(req *http.Request, w http.ResponseWriter) controller.ActionContext {
	return &DemoContext{
		controller.NewBaseContext(req, w),
	}
}

type DemoContext struct {
	*controller.BaseContext
}

func (c *DemoContext) BeforeAction() {
	c.AppendResponseBody([]byte("before demo action\n"))
}

func (c *DemoContext) AfterAction() {
	c.AppendResponseBody([]byte("after demo action\n"))
}

func (c *DemoContext) Destruct() {
	fmt.Println("destruct demo context")
}

func (dc *DemoController) RedirectAction(c *DemoContext) {
	system.Redirect302("https://baidu.com")
}

func (dc *DemoController) DescribeDemoAction(c *DemoContext) {
	c.AppendResponseBody([]byte("DescribeDemo\n"))
}

func (dc *DemoController) ProcessPostAction(c *DemoContext) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Println("Read failed:", err)
	}

	defer c.Request().Body.Close()

	msg := "Hi, Client! Your data is: " + string(body) + "\n"
	c.ResponseWriter().Write([]byte(msg))

	fmt.Println(string(body))
}

func (dc *DemoController) ArgsAction(c *DemoContext, id string) {
	msg := "Hi, Client! Your data is: " + id + "\n"
	c.ResponseWriter().Write([]byte(msg))
	fmt.Print(msg)
}

type TestController struct {
}

func (tc *TestController) NewActionContext(req *http.Request, w http.ResponseWriter) controller.ActionContext {
	return &DemoContext{
		controller.NewBaseContext(req, w),
	}
}

func (tc *TestController) ArgsAction(dc *DemoContext, id string) {
	msg := "Hi, Client! Your data is: " + id + "\n"
	dc.ResponseWriter().Write([]byte(msg))
	fmt.Print(msg)
}
