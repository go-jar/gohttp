package router

import (
	"gohttp/controller"
	"reflect"
	"regexp"
	"strings"
)

type SimpleRouter struct {
	controllerRegex *regexp.Regexp
	actionRegex     *regexp.Regexp
	controllerTable map[string]*controllerItem
}

type controllerItem struct {
	controller     controller.Controller
	actionValueMap map[string]*actionItem
}

type actionItem struct {
	actionValue *reflect.Value
}

func NewSimpleRouter() *SimpleRouter {
	return &SimpleRouter{
		controllerRegex: regexp.MustCompile("([A-Z][A-Za-z0-9]*)Controller$"),
		actionRegex:     regexp.MustCompile("^([A-Z][A-Za-z0-9]*)Action$"),
		controllerTable: make(map[string]*controllerItem),
	}
}

func (router *SimpleRouter) RegistRoutes(controllers ...controller.Controller) {
	for _, controller := range controllers {
		router.registRoute(controller)
	}
}

func (router *SimpleRouter) registRoute(controller controller.Controller) {
	if router.getControllerItem(controller) != nil {
		return
	}

	controllerType := reflect.TypeOf(controller)
	controllerValue := reflect.ValueOf(controller)
	actItem := make(map[string]*actionItem)

	for i := 0; i < controllerValue.NumMethod(); i++ {
		actionName := controllerType.Method(i).Name
		if actionName = router.getActionName(actionName); actionName == "" {
			continue
		}
		action := controllerValue.Method(i)
		actItem[actionName] = &actionItem{
			actionValue: &action,
		}
	}

	controllerName := controllerType.String()
	controllerName = router.getControllerName(controllerName)
	router.controllerTable[controllerName] = &controllerItem{
		controller:     controller,
		actionValueMap: actItem,
	}
}

func (router *SimpleRouter) getControllerItem(controller controller.Controller) *controllerItem {
	controllerType := reflect.TypeOf(controller)
	controllerName := router.getControllerName(controllerType.String())
	if controllerName == "" {
		return nil
	}

	ctrlItem, ok := router.controllerTable[controllerName]
	if !ok {
		return nil
	}
	return ctrlItem
}

func (router *SimpleRouter) getControllerName(controllerName string) string {
	matches := router.controllerRegex.FindStringSubmatch(controllerName)
	if matches == nil {
		return ""
	}

	return strings.ToLower(matches[1])
}

func (router *SimpleRouter) getActionName(actionName string) string {
	matches := router.actionRegex.FindStringSubmatch(actionName)
	if matches == nil {
		return ""
	}

	actionName = strings.ToLower(actionName)
	if actionName == "before" || actionName == "after" {
		return ""
	}

	return strings.ToLower(matches[1])
}

func (router *SimpleRouter) FindRoute(path string) *Route {
	path = strings.Trim(path, "/")
	pathSplit := strings.Split(path, "/")

	if len(pathSplit) != 2 {
		return nil
	}

	controllerName := strings.TrimSpace(pathSplit[0])
	actionName := strings.TrimSpace(pathSplit[1])
	if controllerName == "" || actionName == "" {
		return nil
	}

	route := router.getRoute(controllerName, actionName)
	return route
}

func (router *SimpleRouter) getRoute(controllerName, actionName string) *Route {
	controllerName = strings.ToLower(controllerName)
	actionName = strings.ToLower(actionName)
	ctrlItem, ok := router.controllerTable[controllerName]
	if !ok {
		return nil
	}

	actItem, ok := ctrlItem.actionValueMap[actionName]
	if !ok {
		return nil
	}
	return &Route{
		Controller:  ctrlItem.controller,
		ActionValue: actItem.actionValue,
	}
}
