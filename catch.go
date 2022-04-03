package catch

import (
	"github.com/rs/zerolog/log"
)

type CatchHandler struct {
	onPanicHandler   *OnPanicHandler
	onFailureHandler *OnFailureHandler
	onSuccessHandler *OnSuccessHandler
	finallyHandler   *FinallyHandler
}

// OnPanic is handler that will be executed when there is panic error in function wrapper
func OnPanic(callback func(err interface{})) func(*CatchHandler) {
	return func(ch *CatchHandler) {
		ch.onPanicHandler = &OnPanicHandler{
			callback: callback,
		}
	}
}

// OnFailure is handler that will be executed when code in function wrapper return error or do panic error
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

// OnSuccess is handler that will be executed when code in `function wrapper` doesn't return error or doesn't do panic error
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

// Finally is handler that will be executed in the last time (although there is a panic error in a other handler)
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

// catch will use to catch panic error in `function wrapper`
func catch(tCatchHandler CatchHandler, err *error) {
	if r := recover(); r != nil {
		rError := r.(error)
		tCatchHandler.onPanicHandler.SetError(rError)
		*err = rError
		log.Error().
			Err(rError).
			Msg("[catch] panic error")
	}
}

// DefaultCatchHandler assign default handler
func DefaultCatchHandler() CatchHandler {
	return CatchHandler{
		onFailureHandler: &defaultFailureFunctionHandling,
		onSuccessHandler: &defaultSuccessFunctionHandling,
		finallyHandler:   &defaultFinally,
	}
}

// assignFunctionHandling is used for assign custom handler
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

	//function wrapper
	errFunction := func(catchHandler CatchHandler, err *error) error {
		// only catch panic error for function wrapper and if needed, sometime for routine worker
		if catchHandler.onPanicHandler != nil {
			defer catch(catchHandler, err)
		}

		return fn()
	}(handler, &err)

	// if there a panic error will try to assign with errFunction
	// panic error will be first priority to return
	if err == nil {
		err = errFunction
	}
	if err != nil {
		returnOnFailureCallback := handler.onFailureHandler.callback(err)
		errorHandler = handler.onFailureHandler.Assign(returnOnFailureCallback)
		if errorHandler != nil {
			handler.onFailureHandler.SetError(errorHandler)
			log.Error().
				Err(errorHandler).
				Msg("error in onFailureHandler when try to assign")
		}
	} else {
		returnOnSuccessCallback := handler.onSuccessHandler.callback()
		errorHandler = handler.onSuccessHandler.Assign(returnOnSuccessCallback)
		if errorHandler != nil {
			handler.onSuccessHandler.SetError(errorHandler)
			log.Error().
				Err(errorHandler).
				Msg("error in onSuccessHandler when try to assign")
		}
	}

	if handler.onPanicHandler != nil && handler.onPanicHandler.err != nil {
		return
	}

	// will execute if there is no panic error in function wrapper and will skip if onPanicHandler defined
	returnFinallyCallback := handler.finallyHandler.callback()
	errorHandler = handler.finallyHandler.Assign(returnFinallyCallback)
	if errorHandler != nil {
		log.Error().
			Err(errorHandler).
			Msg("error in finallyHandler when try to assign")
	}
	return
}
