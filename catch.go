package catch

import (
	"fmt"
	"log"
)

type CatchHandlerInterface interface {
	error(fn func(r interface{}), err interface{})
	success(fn func())
	finally(fn func())
}
type CatchHandler struct {
	OnError   func(err interface{})
	OnSuccess func()
	Finally   func()
}

func (t *CatchHandler) error(fn func(r interface{}), err interface{}) {
	t.OnError = fn
	t.OnError(err)
}
func (t *CatchHandler) success(fn func()) {
	t.OnSuccess = fn
	t.OnSuccess()
}

func (t *CatchHandler) finally(fn func()) {
	t.Finally = fn
	t.Finally()
}

var defaultErrorFunctionHandling = func(err interface{}) {
	fmt.Println(err)
}
var defaultSuccessFunctionHandling = func() {
}

var defaultFinally = func() {
}

func catch(tCatchHandler CatchHandler) {
	if r := recover(); r != nil {
		tCatchHandler.OnError(r)
	}
}

func DefaultCatchHandler() CatchHandler {
	return CatchHandler{
		OnError:   defaultErrorFunctionHandling,
		OnSuccess: defaultSuccessFunctionHandling,
		Finally:   defaultFinally,
	}
}
func assignFunctionHandling(catchHandlerInterface CatchHandlerInterface) CatchHandler {
	defaultHandler := DefaultCatchHandler()
	if catchHandlerInterface == nil {
		return defaultHandler
	}
	handler := catchHandlerInterface.(*CatchHandler)
	if handler.OnError == nil {
		handler.OnError = defaultErrorFunctionHandling
	}
	if handler.OnSuccess == nil {
		handler.OnSuccess = defaultSuccessFunctionHandling
	}
	if handler.Finally == nil {
		handler.Finally = defaultFinally
	}
	return *handler
}
func Catch(catchHandler CatchHandlerInterface, err error, msg string) {
	handler := assignFunctionHandling(catchHandler)
	func() {
		defer catch(handler)
		if err != nil {
			log.Panicf("%s: %s", msg, err)
		} else {
			handler.OnSuccess()
		}
	}()
	handler.Finally()
}
