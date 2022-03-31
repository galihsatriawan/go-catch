package catch

import (
	"github.com/rs/zerolog/log"
)

type CatchHandlerInterface interface {
	error(handler OnErrorHandler, err interface{})
	failure(handler OnFailureHandler, err interface{})
	success(handler OnSuccessHandler)
	finally(handler FinallyHandler)
}

type CatchHandler struct {
	onErrorHandler   *OnErrorHandler
	onFailureHandler *OnFailureHandler
	onSuccessHandler *OnSuccessHandler
	finallyHandler   *FinallyHandler
}

func OnError(callback func(err interface{})) func(*CatchHandler) {
	return func(ch *CatchHandler) {
		ch.onErrorHandler = &OnErrorHandler{
			callback: callback,
		}
	}
}
func OnFailure(dst interface{}, callback func(err interface{}) interface{}) func(*CatchHandler) {
	return func(ch *CatchHandler) {
		ch.onFailureHandler = &OnFailureHandler{
			handler: handler{
				dst: dst,
			},
			callback: callback,
		}
	}
}
func OnSuccess(dst interface{}, callback func() interface{}) func(*CatchHandler) {
	return func(ch *CatchHandler) {
		ch.onSuccessHandler = &OnSuccessHandler{
			handler: handler{
				dst: dst,
			},
			callback: callback,
		}
	}
}

func Finally(dst interface{}, callback func() interface{}) func(*CatchHandler) {
	return func(ch *CatchHandler) {
		ch.finallyHandler = &FinallyHandler{
			handler: handler{
				dst: dst,
			},
			callback: callback,
		}
	}
}

// error handle panic error
func (t *CatchHandler) error(handler OnErrorHandler, err interface{}) {
	t.onErrorHandler = &handler
	t.onErrorHandler.callback(err)
}

func (t *CatchHandler) failure(handler OnFailureHandler, err interface{}) {
	t.onFailureHandler = &handler
	t.onFailureHandler.callback(err)
}

func (t *CatchHandler) success(handler OnSuccessHandler) {
	t.onSuccessHandler = &handler
	t.onSuccessHandler.callback()
}

func (t *CatchHandler) finally(handler FinallyHandler) {
	t.finallyHandler = &handler
	t.finallyHandler.callback()
}

func catch(tCatchHandler CatchHandler, err *error) {
	if r := recover(); r != nil {
		rError := r.(error)
		*err = rError
		log.Error().
			Err(rError).
			Msg("[catch] panic error")

		tCatchHandler.onErrorHandler.callback(r)
	}
}

func DefaultCatchHandler() CatchHandler {
	return CatchHandler{
		onErrorHandler:   &defaultErrorFunctionHandling,
		onFailureHandler: &defaultFailureFunctionHandling,
		onSuccessHandler: &defaultSuccessFunctionHandling,
		finallyHandler:   &defaultFinally,
	}
}
func assignFunctionHandling(handlers ...func(*CatchHandler)) CatchHandler {
	defaultHandler := DefaultCatchHandler()

	for _, handler := range handlers {
		handler(&defaultHandler)
	}
	return defaultHandler
}
func Catch(fn func() error, handlers ...func(*CatchHandler)) (err error) {
	var errorHandler error
	handler := assignFunctionHandling(handlers...)
	defer func(catchHandler CatchHandler) {
		returnFinallyCallback := catchHandler.finallyHandler.callback()
		errorHandler = catchHandler.finallyHandler.Assign(returnFinallyCallback)
		if errorHandler != nil {
			log.Error().
				Err(errorHandler).
				Msg("error in finallyHandler when try to assign")
		}
	}(handler)
	err = func(catchHandler CatchHandler, err error) error {
		// only catch panic error for function wrapper
		defer catch(handler, &err)
		return fn()
	}(handler, err)
	if err != nil {
		returnOnFailureCallback := handler.onFailureHandler.callback(err)
		errorHandler = handler.onFailureHandler.Assign(returnOnFailureCallback)
		if errorHandler != nil {
			log.Error().
				Err(errorHandler).
				Msg("error in onFailureHandler when try to assign")
		}

	} else {
		returnOnSuccessCallback := handler.onSuccessHandler.callback()
		errorHandler = handler.onSuccessHandler.Assign(returnOnSuccessCallback)
		if errorHandler != nil {
			log.Error().
				Err(errorHandler).
				Msg("error in onSuccessHandler when try to assign")
		}
	}
	return
}
