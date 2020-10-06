package router

import (
	"gohttp/controller"
	"reflect"
)

type Route struct {
	Controller  controller.Controller
	ActionValue *reflect.Value
}

type Router interface {
	RegisterRoutes(cls ...controller.Controller)
	FindRoute(path string) *Route
}
