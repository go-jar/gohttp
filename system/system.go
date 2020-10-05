package system

import (
	"gohttp/controller"
	"gohttp/router"
	"net/http"
	"reflect"
)

type RoutePathFunc func(request *http.Request) string

type System struct {
	router        router.Router
	routePathFunc RoutePathFunc
}

func NewSystem(router router.Router) *System {
	system := &System{
		router: router,
	}

	system.routePathFunc = system.getRoutePathFunc
	return system
}

func (system *System) getRoutePathFunc(request *http.Request) string {
	return request.URL.Path
}

func (system *System) SetRouteFunc(routePathFunc RoutePathFunc) {
	system.routePathFunc = routePathFunc
}

func (system *System) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	path := system.routePathFunc(request)
	route := system.router.FindRoute(path)
	if route == nil {
		http.NotFound(responseWriter, request)
		return
	}

	context := route.Controller.NewActionContext(responseWriter, request)

	defer func() {
		if err := recover(); err != nil {
			jmpItem, ok := err.(*jumpItem)
			if !ok {
				panic(err)
			}
			jmpItem.jumpFunc(context, jmpItem.args...)
		}
		_, _ = responseWriter.Write(context.ResponseBody())
		context.Destruct()
	}()

	context.BeforeAction()
	route.ActionValue.Call(system.makeArgs(context))
	context.AfterAction()
}

func (system *System) makeArgs(context controller.ActionContext) []reflect.Value {
	argsValues := make([]reflect.Value, 1)
	argsValues[0] = reflect.ValueOf(context)
	return argsValues
}

type JumpFunction func(context controller.ActionContext, args ...interface{})

type jumpItem struct {
	jumpFunc JumpFunction
	args     []interface{}
}

func JumpOutAction(jumpFunc JumpFunction, args ...interface{}) {
	jmpItem := &jumpItem{
		jumpFunc: jumpFunc,
		args:     args,
	}
	panic(jmpItem)
}

func Redirect302(url string) {
	JumpOutAction(redirect302, url)
}

func redirect302(context controller.ActionContext, args ...interface{}) {
	http.Redirect(context.ResponseWriter(), context.Request(), args[0].(string), 302)
}
