package router

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-jar/gohttp/controller"
	"github.com/go-jar/golog"
)

type SimpleRouter struct {
	controllerRegex *regexp.Regexp
	actionRegex     *regexp.Regexp

	controllerTableDefined []*controllerItemDefined
	controllerTable        map[string]*controllerItem

	logger golog.ILogger
}

type controllerItem struct {
	ctrl      controller.Controller
	ctrlValue *reflect.Value
	ctrlType  reflect.Type

	controllerName string
	actionValueMap map[string]*actionItem
}

type controllerItemDefined struct {
	pathRegex *regexp.Regexp

	controllerName string
	actionName     string
}

type actionItem struct {
	argsNum     int
	actionValue *reflect.Value
}

func NewSimpleRouter(logger golog.ILogger) *SimpleRouter {
	if logger == nil {
		logger = new(golog.NoopLogger)
	}

	return &SimpleRouter{
		controllerRegex: regexp.MustCompile("([A-Z][A-Za-z0-9]*)Controller$"),
		actionRegex:     regexp.MustCompile("^([A-Z][A-Za-z0-9]*)Action$"),
		controllerTable: make(map[string]*controllerItem),
		logger:          logger,
	}
}

func (sr *SimpleRouter) RegisterRoutes(cls ...controller.Controller) {
	for _, ctrl := range cls {
		sr.registerRoute(ctrl)
	}

	sr.logRoutes()
}

func (sr *SimpleRouter) DefineRoute(pattern string, ctrl controller.Controller, actionName string) {
	methodName := strings.Title(actionName) + "Action"
	actionName = strings.ToLower(methodName)
	if actionName == "" {
		return
	}

	ctrlItem := sr.getOrInitControllerItem(ctrl)
	if ctrlItem == nil {
		return
	}

	action, ok := ctrlItem.ctrlType.MethodByName(methodName)
	if !ok {
		return
	}

	actionArgsNum := sr.getActionArgsNum(action, ctrlItem.ctrlType)
	if actionArgsNum == -1 {
		return
	}

	actionValue := ctrlItem.ctrlValue.MethodByName(methodName)
	ctrlItem.actionValueMap[actionName] = &actionItem{
		argsNum:     actionArgsNum,
		actionValue: &actionValue,
	}

	sr.controllerTableDefined = append(sr.controllerTableDefined, &controllerItemDefined{
		pathRegex: regexp.MustCompile(pattern),

		controllerName: strings.ToLower(ctrlItem.controllerName),
		actionName:     strings.ToLower(actionName),
	})
}

func (sr *SimpleRouter) registerRoute(ctrl controller.Controller) {
	ctrlItem := sr.getOrInitControllerItem(ctrl)
	if ctrlItem == nil {
		return
	}

	controllerType := ctrlItem.ctrlType
	controllerValue := ctrlItem.ctrlValue

	for i := 0; i < controllerValue.NumMethod(); i++ {
		actionT := controllerType.Method(i)
		actionName := sr.getActionName(actionT.Name)
		if actionName == "" {
			continue
		}

		actionArgsNum := sr.getActionArgsNum(actionT, controllerType)
		if actionArgsNum == -1 {
			continue
		}

		actionV := controllerValue.Method(i)
		ctrlItem.actionValueMap[actionName] = &actionItem{
			argsNum:     actionArgsNum,
			actionValue: &actionV,
		}
	}
}

func (sr *SimpleRouter) getOrInitControllerItem(ctrl controller.Controller) *controllerItem {
	controllerValue := reflect.ValueOf(ctrl)
	controllerType := controllerValue.Type()

	controllerName := sr.getControllerName(controllerType.String())
	if controllerName == "" {
		return nil
	}

	ctrlItem, ok := sr.controllerTable[controllerName]
	if !ok {
		ctrlItem = &controllerItem{
			ctrl:      ctrl,
			ctrlType:  controllerType,
			ctrlValue: &controllerValue,

			controllerName: controllerName,
			actionValueMap: make(map[string]*actionItem),
		}
		sr.controllerTable[controllerName] = ctrlItem
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

func (sr *SimpleRouter) getActionArgsNum(actionMethod reflect.Method, controllerType reflect.Type) int {
	n := actionMethod.Type.NumIn()
	if n < 2 {
		return -1
	}

	if actionMethod.Type.In(0).String() != controllerType.String() {
		return -1
	}

	if n > 2 {
		valid := true
		for i := 2; i < n; i++ {
			if actionMethod.Type.In(i).String() != "string" {
				valid = false
				break
			}
		}
		if !valid {
			return -1
		}
	}

	return n - 2 // delete sr and context
}

func (sr *SimpleRouter) logRoutes() {
	sr.logger.Debug([]byte("routers registered:"))

	for _, route := range sr.controllerTableDefined {
		sr.logger.Debug([]byte(route.controllerName + "." + route.actionName))
	}

	for ctrlName, actions := range sr.controllerTable {
		for actionName, _ := range actions.actionValueMap {
			sr.logger.Debug([]byte(ctrlName + "." + actionName))
		}
	}
}

func (sr *SimpleRouter) FindRoute(path string) *Route {
	path = strings.ToLower(path)

	route := sr.findRouteDefined(path)
	if route == nil {
		route = sr.findRouteGeneral(path)
	}

	return route
}

func (sr *SimpleRouter) findRouteDefined(path string) *Route {
	for _, ctrlDefined := range sr.controllerTableDefined {
		matches := ctrlDefined.pathRegex.FindStringSubmatch(path)
		if matches == nil {
			continue
		}

		route, argsNum := sr.getRoute(ctrlDefined.controllerName, ctrlDefined.actionName)
		route.Args = sr.makeActionArgs(matches[1:], argsNum)
		return route
	}

	return nil
}

func (sr *SimpleRouter) findRouteGeneral(path string) *Route {
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

	route, _ := sr.getRoute(controllerName, actionName)
	return route
}

func (sr *SimpleRouter) getRoute(controllerName, actionName string) (*Route, int) {
	controllerName = strings.ToLower(controllerName)
	actionName = strings.ToLower(actionName)

	ctrlItem, ok := sr.controllerTable[controllerName]
	if !ok {
		sr.logger.Error([]byte(controllerName + "not found"))
		return nil, 0
	}

	actItem, ok := ctrlItem.actionValueMap[actionName]
	if !ok {
		sr.logger.Error([]byte(actionName + "not found"))
		return nil, 0
	}

	return &Route{
		Controller:  ctrlItem.ctrl,
		ActionValue: actItem.actionValue,
	}, actItem.argsNum
}

func (sr *SimpleRouter) makeActionArgs(args []string, validArgsNum int) []string {
	rgArgsNum := len(args)
	missArgsNum := validArgsNum - rgArgsNum

	switch {
	case missArgsNum == 0:
	case missArgsNum > 0:
		for i := 0; i < missArgsNum; i++ {
			args = append(args, "")
		}
	case missArgsNum < 0:
		args = args[:validArgsNum]
	}

	return args
}
