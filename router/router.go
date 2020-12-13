package router

import (
	"gohttp/controller"
	"reflect"
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
