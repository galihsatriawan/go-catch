package catch

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ERR_UNADDRESSABLE_VALUE             = "using unaddressable value"
	ERR_NOT_ASSIGNABLE_TO_SPECIFIC_TYPE = "value of type %v is not assignable to type %v"
)

type HandlerInterface interface {
	Validate() (err error)
	SetResult(result interface{})
}

type OnErrorHandler struct {
	Callback func(err interface{})
}

type OnFailureHandler struct {
	Dst      interface{}
	Callback func(err interface{}) (dst interface{})
	result   interface{}
}

// OnSuccessHandler success action
type OnSuccessHandler struct {
	Dst      interface{}
	Callback func() (dst interface{})
	result   interface{}
}

type FinallyHandler struct {
	Dst      interface{}
	Callback func() (dst interface{})
	result   interface{}
}

var defaultErrorFunctionHandling = OnErrorHandler{
	Callback: func(err interface{}) {
		fmt.Println(err)
	},
}
var defaultSuccessFunctionHandling = OnSuccessHandler{
	Callback: func() (dst interface{}) {
		return
	},
}

var defaultFailureFunctionHandling = OnFailureHandler{
	Callback: func(err interface{}) (dst interface{}) {
		return
	},
}

var defaultFinally = FinallyHandler{
	Callback: func() (dst interface{}) {
		return
	},
}

func (h *OnFailureHandler) Validate() (err error) {
	if h.Dst == nil {
		return
	}
	if reflect.ValueOf(h.Dst).Kind() != reflect.Ptr {
		return errors.New(ERR_UNADDRESSABLE_VALUE)
	}

	destinationType := reflect.ValueOf(h.Dst).Type()
	returnOnFailureCallbackType := reflect.ValueOf(h.result).Type()
	if destinationType != returnOnFailureCallbackType {
		err = fmt.Errorf(ERR_NOT_ASSIGNABLE_TO_SPECIFIC_TYPE, destinationType, returnOnFailureCallbackType)

		// give compatibility although `result` is pointer or not
		destinationElem := destinationType.Elem()

		if destinationElem == returnOnFailureCallbackType {
			err = nil
		}
	}
	return
}

func (h *OnFailureHandler) SetResult(result interface{}) {
	h.result = result
}

func (h *OnSuccessHandler) Validate() (err error) {
	if h.Dst == nil {
		return
	}
	if reflect.ValueOf(h.Dst).Kind() != reflect.Ptr {
		return errors.New(ERR_UNADDRESSABLE_VALUE)
	}

	destinationType := reflect.ValueOf(h.Dst).Type()
	returnOnSuccessCallbackType := reflect.ValueOf(h.result).Type()
	if destinationType != returnOnSuccessCallbackType {
		err = fmt.Errorf(ERR_NOT_ASSIGNABLE_TO_SPECIFIC_TYPE, destinationType, returnOnSuccessCallbackType)

		// give compatibility although `result` is pointer or not
		destinationElem := destinationType.Elem()

		if destinationElem == returnOnSuccessCallbackType {
			err = nil
		}
	}
	return
}

func (h *OnSuccessHandler) SetResult(result interface{}) {
	h.result = result
}

func (h *FinallyHandler) Validate() (err error) {
	if h.Dst == nil {
		return
	}
	if reflect.ValueOf(h.Dst).Kind() != reflect.Ptr {
		return errors.New(ERR_UNADDRESSABLE_VALUE)
	}

	destinationType := reflect.ValueOf(h.Dst).Type()
	returnFinallyCallbackType := reflect.ValueOf(h.result).Type()
	if destinationType != returnFinallyCallbackType {
		err = fmt.Errorf(ERR_NOT_ASSIGNABLE_TO_SPECIFIC_TYPE, destinationType, returnFinallyCallbackType)

		// give compatibility although `result` is pointer or not
		destinationElem := destinationType.Elem()

		if destinationElem == returnFinallyCallbackType {
			err = nil
		}
	}
	return
}

func (h *FinallyHandler) SetResult(result interface{}) {
	h.result = result
}
