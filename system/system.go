package system

import (
	"net/http"
	"reflect"

	"github.com/go-jar/gohttp/controller"
	"github.com/go-jar/gohttp/router"
)

type RoutePathFunc func(req *http.Request) string

type System struct {
	router        router.Router
	routePathFunc RoutePathFunc
}

func NewSystem(rt router.Router) *System {
	sys := &System{
		router: rt,
	}

	sys.routePathFunc = sys.routePath
	return sys
}

func (s *System) routePath(req *http.Request) string {
	return req.URL.Path
}

func (s *System) SetRouteFunc(routePathFunc RoutePathFunc) {
	s.routePathFunc = routePathFunc
}

func (s *System) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := s.routePathFunc(req)
	route := s.router.FindRoute(path)
	if route == nil {
		http.NotFound(w, req)
		return
	}

	actionContext := route.Controller.NewActionContext(w, req)

	defer func() {
		if err := recover(); err != nil {
			jmpItem, ok := err.(*jumpItem)
			if !ok {
				panic(err)
			}
			jmpItem.jumpFunc(actionContext, jmpItem.args...)
		}
		_, _ = w.Write(actionContext.ResponseBody())
		actionContext.Destruct()
	}()

	actionContext.BeforeAction()
	route.ActionValue.Call(s.makeArgs(actionContext, route.Args))
	actionContext.AfterAction()
}

func (s *System) makeArgs(ctxt controller.ActionContext, args []string) []reflect.Value {
	argsValues := make([]reflect.Value, len(args)+1)
	argsValues[0] = reflect.ValueOf(ctxt)

	for i, arg := range args {
		argsValues[i+1] = reflect.ValueOf(arg)
	}

	return argsValues
}

type JumpFunction func(ctxt controller.ActionContext, args ...interface{})

type jumpItem struct {
	jumpFunc JumpFunction
	args     []interface{}
}

func JumpOutAction(jf JumpFunction, args ...interface{}) {
	jmpItem := &jumpItem{
		jumpFunc: jf,
		args:     args,
	}
	panic(jmpItem)
}

func Redirect302(url string) {
	JumpOutAction(redirect302, url)
}

func redirect302(ctxt controller.ActionContext, args ...interface{}) {
	http.Redirect(ctxt.ResponseWriter(), ctxt.Request(), args[0].(string), 302)
}
