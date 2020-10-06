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
	ctrl           controller.Controller
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

func (sr *SimpleRouter) RegisterRoutes(cls ...controller.Controller) {
	for _, ctrl := range cls {
		sr.registerRoute(ctrl)
	}
}

func (sr *SimpleRouter) registerRoute(ctrl controller.Controller) {
	if sr.getControllerItem(ctrl) != nil {
		return
	}

	controllerType := reflect.TypeOf(ctrl)
	controllerValue := reflect.ValueOf(ctrl)
	actItem := make(map[string]*actionItem)

	for i := 0; i < controllerValue.NumMethod(); i++ {
		actionName := controllerType.Method(i).Name
		actionName = sr.getActionName(actionName)
		if actionName == "" {
			continue
		}

		action := controllerValue.Method(i)
		actItem[actionName] = &actionItem{
			actionValue: &action,
		}
	}

	controllerName := controllerType.String()
	controllerName = sr.getControllerName(controllerName)
	sr.controllerTable[controllerName] = &controllerItem{
		ctrl:           ctrl,
		actionValueMap: actItem,
	}
}

func (sr *SimpleRouter) getControllerItem(ctrl controller.Controller) *controllerItem {
	controllerType := reflect.TypeOf(ctrl)
	controllerName := sr.getControllerName(controllerType.String())
	if controllerName == "" {
		return nil
	}

	ctrlItem, ok := sr.controllerTable[controllerName]
	if !ok {
		return nil
	}
	return ctrlItem
}

func (sr *SimpleRouter) getControllerName(controllerName string) string {
	matches := sr.controllerRegex.FindStringSubmatch(controllerName)
	if matches == nil {
		return ""
	}

	return strings.ToLower(matches[1])
}

func (sr *SimpleRouter) getActionName(actionName string) string {
	matches := sr.actionRegex.FindStringSubmatch(actionName)
	if matches == nil {
		return ""
	}

	actionName = strings.ToLower(actionName)
	if actionName == "before" || actionName == "after" {
		return ""
	}

	return strings.ToLower(matches[1])
}

func (sr *SimpleRouter) FindRoute(path string) *Route {
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

	route := sr.getRoute(controllerName, actionName)
	return route
}

func (sr *SimpleRouter) getRoute(controllerName, actionName string) *Route {
	controllerName = strings.ToLower(controllerName)
	actionName = strings.ToLower(actionName)
	ctrlItem, ok := sr.controllerTable[controllerName]
	if !ok {
		return nil
	}

	actItem, ok := ctrlItem.actionValueMap[actionName]
	if !ok {
		return nil
	}
	return &Route{
		Controller:  ctrlItem.ctrl,
		ActionValue: actItem.actionValue,
	}
}
