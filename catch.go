package catch

import (
	"reflect"
)

type CatchHandlerInterface interface {
	error(handler OnErrorHandler, err interface{})
	failure(handler OnFailureHandler, err interface{})
	success(handler OnSuccessHandler)
	finally(handler FinallyHandler)
}

type CatchHandler struct {
	OnError   *OnErrorHandler
	OnFailure *OnFailureHandler
	OnSuccess *OnSuccessHandler
	Finally   *FinallyHandler
}

// error handle panic error
func (t *CatchHandler) error(handler OnErrorHandler, err interface{}) {
	t.OnError = &handler
	t.OnError.Callback(err)
}

func (t *CatchHandler) failure(handler OnFailureHandler, err interface{}) {
	t.OnFailure = &handler
	t.OnFailure.Callback(err)
}

func (t *CatchHandler) success(handler OnSuccessHandler) {
	t.OnSuccess = &handler
	t.OnSuccess.Callback()
}

func (t *CatchHandler) finally(handler FinallyHandler) {
	t.Finally = &handler
	t.Finally.Callback()
}

func catch(tCatchHandler CatchHandler) {
	if r := recover(); r != nil {
		tCatchHandler.OnError.Callback(r)
	}
}

func DefaultCatchHandler() CatchHandler {
	return CatchHandler{
		OnError:   &defaultErrorFunctionHandling,
		OnFailure: &defaultFailureFunctionHandling,
		OnSuccess: &defaultSuccessFunctionHandling,
		Finally:   &defaultFinally,
	}
}
func assignFunctionHandling(catchHandlerInterface CatchHandlerInterface) CatchHandler {
	defaultHandler := DefaultCatchHandler()
	if catchHandlerInterface == nil {
		return defaultHandler
	}
	handler := catchHandlerInterface.(*CatchHandler)
	if handler.OnError == nil {
		handler.OnError = &defaultErrorFunctionHandling
	}

	if handler.OnFailure == nil {
		handler.OnFailure = &defaultFailureFunctionHandling
	}

	if handler.OnSuccess == nil {
		handler.OnSuccess = &defaultSuccessFunctionHandling
	}
	if handler.Finally == nil {
		handler.Finally = &defaultFinally
	}
	return *handler
}
func Catch(Callback func() error, catchHandler CatchHandlerInterface) (err error) {
	handler := assignFunctionHandling(catchHandler)
	func() {
		defer catch(handler)
		err = Callback()
		if err != nil {
			returnOnFailureCallback := handler.OnFailure.Callback(err)
			handler.OnFailure.SetResult(returnOnFailureCallback)

			err = handler.OnFailure.Validate()
			// panic error
			if err != nil {
				handler.OnError.Callback(err)
				return
			}

			// give compatibility although `result` is pointer or not
			onFailureSetValue := reflect.ValueOf(returnOnFailureCallback)
			if onFailureSetValue.Kind() == reflect.Ptr {
				onFailureSetValue = onFailureSetValue.Elem()
			}
			reflect.ValueOf(handler.OnFailure.Dst).Elem().Set(onFailureSetValue)
		} else {
			returnOnSuccessCallback := handler.OnSuccess.Callback()
			handler.OnSuccess.SetResult(returnOnSuccessCallback)

			err = handler.OnSuccess.Validate()
			// panic error
			if err != nil {
				handler.OnError.Callback(err)
				return
			}
			if handler.Finally.Dst != nil {
				// give compatibility although `result` is pointer or not
				onSuccessSetValue := reflect.ValueOf(returnOnSuccessCallback)
				if onSuccessSetValue.Kind() == reflect.Ptr {
					onSuccessSetValue = onSuccessSetValue.Elem()
				}
				reflect.ValueOf(handler.OnSuccess.Dst).Elem().Set(onSuccessSetValue)
			}
		}
	}()

	returnFinallyCallback := handler.Finally.Callback()
	handler.Finally.SetResult(returnFinallyCallback)

	err = handler.Finally.Validate()

	// panic error
	if err != nil {
		handler.OnError.Callback(err)
		return
	}
	if handler.Finally.Dst != nil {
		// give compatibility although `result` is pointer or not
		finallySetValue := reflect.ValueOf(returnFinallyCallback)
		if finallySetValue.Kind() == reflect.Ptr {
			finallySetValue = finallySetValue.Elem()
		}
		reflect.ValueOf(handler.Finally.Dst).Elem().Set(finallySetValue)
	}
	return
}
