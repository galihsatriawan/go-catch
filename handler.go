package catch

import (
	"fmt"
	"reflect"
)

var (
	ERR_DESTINATION_UNADDRESSABLE_VALUE = "destination using unaddressable value, %s must be a pointer or nil"
	ERR_NOT_ASSIGNABLE_TO_SPECIFIC_TYPE = "value of type %v is not assignable to type %v"
	ERR_INVALID_MEMORY_ADDRESS          = "invalid memory address or nil pointer"
)

type handler struct {
	dst    interface{}
	result interface{}
}
type HandlerInterface interface {
	Assign(result interface{}) (err error)
}

type OnFailureHandler struct {
	handler
	callback func(err interface{}) (dst interface{})
}

// OnSuccessHandler success action
type OnSuccessHandler struct {
	handler
	callback func() (dst interface{})
}

type FinallyHandler struct {
	handler
	callback func() (dst interface{})
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

// Currently, because the several handler has same function action for `Assign`, so we define `Assign` to parent
func (h *handler) Assign(result interface{}) (err error) {
	h.result = result

	if h.dst == nil {
		return
	}
	if reflect.ValueOf(h.dst).Kind() != reflect.Ptr {
		err = fmt.Errorf(ERR_DESTINATION_UNADDRESSABLE_VALUE, "`dst`")
		return
	}
	destinationType := reflect.TypeOf(h.dst)
	resultType := reflect.TypeOf(h.result)
	if resultType == nil {
		dstElem := reflect.ValueOf(h.dst).Elem()
		dstElem.Set(reflect.Zero(dstElem.Type()))
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
	setValue := reflect.ValueOf(h.result)
	dstValue := reflect.ValueOf(h.dst)
	if setValue.Kind() == reflect.Ptr {
		setValue = setValue.Elem()
	}
	dstValue.Elem().Set(setValue)
	return
}
