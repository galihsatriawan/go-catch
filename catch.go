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
	ErrorHandling   func(err interface{})
	SuccessHandling func()
	FinallyHandling func()
}

func (t *CatchHandler) error(fn func(r interface{}), err interface{}) {
	t.ErrorHandling = fn
	t.ErrorHandling(err)
}
func (t *CatchHandler) success(fn func()) {
	t.SuccessHandling = fn
	t.SuccessHandling()
}

func (t *CatchHandler) finally(fn func()) {
	t.FinallyHandling = fn
	fmt.Println("Finally")
	t.FinallyHandling()
}

var defaultErrorFunctionHandling = func(err interface{}) {
	fmt.Println(err)
}
var defaultSuccessFunctionHandling = func() {
	fmt.Println("")
}

var defaultFinallyHandling = func() {
	fmt.Println()
}

func catch(tCatchHandler CatchHandler) {
	if r := recover(); r != nil {
		tCatchHandler.ErrorHandling(r)
	}
}

func DefaultCatchHandler() CatchHandler {
	return CatchHandler{
		ErrorHandling:   defaultErrorFunctionHandling,
		SuccessHandling: defaultSuccessFunctionHandling,
		FinallyHandling: defaultFinallyHandling,
	}
}
func assignFunctionHandling(catchHandlerInterface CatchHandlerInterface) CatchHandler {
	defaultHandler := DefaultCatchHandler()
	if catchHandlerInterface == nil {
		return defaultHandler
	}
	handler := catchHandlerInterface.(*CatchHandler)
	if handler.ErrorHandling == nil {
		handler.ErrorHandling = defaultErrorFunctionHandling
	}
	if handler.SuccessHandling == nil {
		handler.SuccessHandling = defaultSuccessFunctionHandling
	}
	if handler.FinallyHandling == nil {
		handler.FinallyHandling = defaultFinallyHandling
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
			handler.SuccessHandling()
		}
	}()
	handler.FinallyHandling()
}
