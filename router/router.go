package router

import (
	"reflect"

	"github.com/go-jar/gohttp/controller"
)

type Route struct {
	Controller  controller.Controller
	ActionValue *reflect.Value
	Args        []string
}

type Router interface {
	RegisterRoutes(cls ...controller.Controller)
	DefineRoute(pattern string, ctrl controller.Controller, actionName string)

	FindRoute(path string) *Route
}
