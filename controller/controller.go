package controller

import (
	"net/http"

	"github.com/goinbox/gomisc"
)

type ActionContext interface {
	Request() *http.Request
	ResponseWriter() http.ResponseWriter

	ResponseBody() []byte
	SetResonseBody(body []byte)
	AppendResopnseBody(body []byte)

	BeforeAction()
	AfterAction()
	Destruct()
}

type BaseContext struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	responseBody   []byte
}

func NewBaseContext(responseWriter http.ResponseWriter, request *http.Request) *BaseContext {
	return &BaseContext{
		request:        request,
		responseWriter: responseWriter,
	}
}

func (context *BaseContext) Request() *http.Request {
	return context.request
}

func (context *BaseContext) ResponseWriter() http.ResponseWriter {
	return context.responseWriter
}

func (context *BaseContext) ResponseBody() []byte {
	return context.responseBody
}

func (context *BaseContext) SetResonseBody(body []byte) {
	context.responseBody = body
}

func (context *BaseContext) AppendResopnseBody(body []byte) {
	context.responseBody = gomisc.AppendBytes(context.responseBody, body)
}

func (context *BaseContext) BeforeAction() {
}

func (context *BaseContext) AfterAction() {
}

func (context *BaseContext) Destruct() {
}

type Controller interface {
	NewActionContext(responseWriter http.ResponseWriter, request *http.Request) ActionContext
}
