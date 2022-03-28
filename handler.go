package catch

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/rs/zerolog/log"
)

var (
	ERR_UNADDRESSABLE_VALUE             = "destination using unaddressable value, %s must be a pointer"
	ERR_NOT_ASSIGNABLE_TO_SPECIFIC_TYPE = "value of type %v is not assignable to type %v"
	ERR_INVALID_MEMORY_ADDRESS          = "invalid memory address or nil pointer"
)

type HandlerInterface interface {
	Assign(result interface{}) (err error)
}
type OnErrorHandler struct {
	callback func(err interface{})
}

type OnFailureHandler struct {
	dst      interface{}
	callback func(err interface{}) (dst interface{})
	result   interface{}
}

// OnSuccessHandler success action
type OnSuccessHandler struct {
	dst      interface{}
	callback func() (dst interface{})
	result   interface{}
}

type FinallyHandler struct {
	dst      interface{}
	callback func() (dst interface{})
	result   interface{}
}

var defaultErrorFunctionHandling = OnErrorHandler{
	callback: func(err interface{}) {
	},
}
var defaultSuccessFunctionHandling = OnSuccessHandler{
	callback: func() (dst interface{}) {
		return
	},
}

var defaultFailureFunctionHandling = OnFailureHandler{
	callback: func(err interface{}) (dst interface{}) {
		return
	},
}

var defaultFinally = FinallyHandler{
	callback: func() (dst interface{}) {
		return
	},
}

func (h *OnErrorHandler) Assign(result interface{}) (err error) {
	return
}

func (h *OnFailureHandler) Assign(result interface{}) (err error) {
	h.result = result

	if h.dst == nil {
		return
	}
	if reflect.ValueOf(h.dst).Kind() != reflect.Ptr {
		return fmt.Errorf(ERR_UNADDRESSABLE_VALUE, "`dst`")
	}
	destinationType := reflect.TypeOf(h.dst)
	resultType := reflect.TypeOf(h.result)
	if resultType == nil {
		log.Warn().Err(errors.New(ERR_INVALID_MEMORY_ADDRESS)).
			Msg("the `result` is nil")
		return
	}
	if !resultType.AssignableTo(destinationType) {
		err = fmt.Errorf(ERR_NOT_ASSIGNABLE_TO_SPECIFIC_TYPE, resultType, destinationType)
		// give compatibility although `result` is pointer or not
		destinationElem := destinationType.Elem()
		if destinationElem != resultType {
			return
		}
		err = nil
	}
	// give compatibility although `result` is pointer or not
	onFailureSetValue := reflect.ValueOf(h.result)
	dstValue := reflect.ValueOf(h.dst)
	if onFailureSetValue.Kind() == reflect.Ptr {
		onFailureSetValue = onFailureSetValue.Elem()
	}
	dstValue.Elem().Set(onFailureSetValue)
	return
}

func (h *OnSuccessHandler) Assign(result interface{}) (err error) {
	h.result = result

	if h.dst == nil {
		return
	}
	if reflect.ValueOf(h.dst).Kind() != reflect.Ptr {
		return fmt.Errorf(ERR_UNADDRESSABLE_VALUE, "`dst`")
	}
	destinationType := reflect.TypeOf(h.dst)
	resultType := reflect.TypeOf(h.result)
	if resultType == nil {
		log.Warn().Err(errors.New(ERR_INVALID_MEMORY_ADDRESS)).
			Msg("the `result` is nil")
		return
	}
	if !resultType.AssignableTo(destinationType) {
		err = fmt.Errorf(ERR_NOT_ASSIGNABLE_TO_SPECIFIC_TYPE, resultType, destinationType)
		// give compatibility although `result` is pointer or not
		destinationElem := destinationType.Elem()
		if destinationElem != resultType {
			return
		}
		err = nil
	}
	// give compatibility although `result` is pointer or not
	onSuccessSetValue := reflect.ValueOf(h.result)
	dstValue := reflect.ValueOf(h.dst)
	if onSuccessSetValue.Kind() == reflect.Ptr {
		onSuccessSetValue = onSuccessSetValue.Elem()
	}

	dstValue.Elem().Set(onSuccessSetValue)
	return
}

func (h *FinallyHandler) Assign(result interface{}) (err error) {
	h.result = result

	if h.dst == nil {
		return
	}
	if reflect.ValueOf(h.dst).Kind() != reflect.Ptr {
		return fmt.Errorf(ERR_UNADDRESSABLE_VALUE, "`dst`")
	}
	destinationType := reflect.TypeOf(h.dst)
	resultType := reflect.TypeOf(h.result)
	if resultType == nil {
		log.Warn().Err(errors.New(ERR_INVALID_MEMORY_ADDRESS)).
			Msg("the `result` is nil")
		return
	}
	if !resultType.AssignableTo(destinationType) {
		err = fmt.Errorf(ERR_NOT_ASSIGNABLE_TO_SPECIFIC_TYPE, resultType, destinationType)
		// give compatibility although `result` is pointer or not
		destinationElem := destinationType.Elem()
		if destinationElem != resultType {
			return
		}
		err = nil
	}
	// give compatibility although `result` is pointer or not
	finallySetValue := reflect.ValueOf(h.result)
	dstValue := reflect.ValueOf(h.dst)
	if finallySetValue.Kind() == reflect.Ptr {
		finallySetValue = finallySetValue.Elem()
	}
	dstValue.Elem().Set(finallySetValue)
	return
}
