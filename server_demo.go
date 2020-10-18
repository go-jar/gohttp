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
	simpleRouter := router.NewSimpleRouter()
	simpleRouter.RegisterRoutes(demoController)
	sys := system.NewSystem(simpleRouter)
	_ = gracehttp.ListenAndServe(":8010", sys)
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
