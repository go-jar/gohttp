package controller

import (
	"net/http"

	"github.com/goinbox/gomisc"
)

type ActionContext interface {
	Request() *http.Request
	ResponseWriter() http.ResponseWriter

	ResponseBody() []byte
	SetResponseBody(body []byte)
	AppendResponseBody(body []byte)

	BeforeAction()
	AfterAction()
	Destruct()
}

type BaseContext struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	responseBody   []byte
}

func NewBaseContext(w http.ResponseWriter, req *http.Request) *BaseContext {
	return &BaseContext{
		request:        req,
		responseWriter: w,
	}
}

func (bc *BaseContext) Request() *http.Request {
	return bc.request
}

func (bc *BaseContext) ResponseWriter() http.ResponseWriter {
	return bc.responseWriter
}

func (bc *BaseContext) ResponseBody() []byte {
	return bc.responseBody
}

func (bc *BaseContext) SetResponseBody(body []byte) {
	bc.responseBody = body
}

func (bc *BaseContext) AppendResponseBody(body []byte) {
	bc.responseBody = gomisc.AppendBytes(bc.responseBody, body)
}

func (bc *BaseContext) BeforeAction() {
}

func (bc *BaseContext) AfterAction() {
}

func (bc *BaseContext) Destruct() {
}

type Controller interface {
	NewActionContext(w http.ResponseWriter, req *http.Request) ActionContext
}
