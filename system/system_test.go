package system

import (
	"fmt"
	"gohttp/controller"
	"gohttp/gracehttp"
	"gohttp/router"
	"net/http"
	"testing"
)

func TestSystem(t *testing.T) {
	demoController := new(DemoController)
	router := router.NewSimpleRouter()
	router.RegistRoutes(demoController)
	system := NewSystem(router)
	_ = gracehttp.ListenAndServe(":8010", system)
}

type DemoController struct {
}

func (demo *DemoController) NewActionContext(responseWriter http.ResponseWriter, request *http.Request) controller.ActionContext {
	return &DemoContext{
		controller.NewBaseContext(responseWriter, request),
	}
}

type DemoContext struct {
	*controller.BaseContext
}

func (context *DemoContext) BeforeAction() {
	context.AppendResopnseBody([]byte("before demo action\n"))
}

func (context *DemoContext) AfterAction() {
	context.AppendResopnseBody([]byte("after demo action\n"))
}

func (demoContext *DemoContext) Destruct() {
	fmt.Println("desctuct demo context")
}

func (demoController *DemoController) DescribeDemoAction(context *DemoContext) {
	context.AppendResopnseBody([]byte("DescribeDemo\n"))
}

func (demoController *DemoController) RedirectAction(context *DemoContext) {
	Redirect302("https://baidu.com")
}
