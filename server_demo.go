package main

import (
	"fmt"
	"gohttp/controller"
	"gohttp/gracehttp"
	"gohttp/router"
	"gohttp/system"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	demoController := new(DemoController)
	testController := new(TestController)

	simpleRouter := router.NewSimpleRouter()

	simpleRouter.DefineRoute("/test/args/([0-9]+)$", testController, "Args")
	simpleRouter.DefineRoute("/demo/args/([0-9]+)$", testController, "Args")
	simpleRouter.RegisterRoutes(demoController)

	sys := system.NewSystem(simpleRouter)
	gracehttp.ListenAndServe(":8010", sys)
}

type DemoController struct {
}

func (dc *DemoController) NewActionContext(w http.ResponseWriter, req *http.Request) controller.ActionContext {
	return &DemoContext{
		controller.NewBaseContext(w, req),
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

func (tc *TestController) NewActionContext(w http.ResponseWriter, req *http.Request) controller.ActionContext {
	return &DemoContext{
		controller.NewBaseContext(w, req),
	}
}

func (tc *TestController) ArgsAction(dc *DemoContext, id string) {
	msg := "Hi, Client! Your data is: " + id + "\n"
	dc.ResponseWriter().Write([]byte(msg))
	fmt.Print(msg)
}
