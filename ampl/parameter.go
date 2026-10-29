package ampl

//#include "ampl/ampl_c.h"
import "C"
import (
  "runtime"
  "unsafe"
  //"errors"
  "reflect"
  "fmt"
)

type Parameter struct {
  Entity
}


// Set sets a numeric or string parameter
func (param *Parameter) Set(value interface{}) error {
	var errInfo *C.AMPL_ERRORINFO

	switch v := value.(type) {
	case float64, float32, int, int64, int32:
		// Convert the value to C.double
		numValue := C.double(reflect.ValueOf(v).Convert(reflect.TypeOf(float64(0))).Interface().(float64))
		errInfo = C.AMPL_ParameterSetNumeric(param.ampl.campl, C.CString(param.name), numValue)
	case string:
		errInfo = C.AMPL_ParameterSetString(param.ampl.campl, C.CString(param.name), C.CString(v))
	default:
		return fmt.Errorf("unsupported parameter type: %T", value)
	}

	return amplError(errInfo)
}

// SetSomeDoubleValues assigns values to the parameter instances identified
// by the given (single-dimensional) string indices.
func (param *Parameter) SetSomeDoubleValues(indices []string, values []float64) error {
  if len(indices) != len(values) {
    return fmt.Errorf("indices and values must have the same length")
  }

  tuples := make([]*Tuple, len(indices))
  cIndices := make([]*C.AMPL_TUPLE, len(indices))
  for i, index := range indices {
    tuple, err := NewStrings([]string{index})
    if err != nil {
      return err
    }
    tuples[i] = tuple
    cIndices[i] = tuple.ctuple
  }

  var cValues *C.double
  if len(values) > 0 {
    cValues = (*C.double)(unsafe.Pointer(&values[0]))
  }

  var cIndicesPtr **C.AMPL_TUPLE
  if len(cIndices) > 0 {
    cIndicesPtr = (**C.AMPL_TUPLE)(unsafe.Pointer(&cIndices[0]))
  }

  errInfo := C.AMPL_ParameterSetSomeDoubleValues(param.ampl.campl, C.CString(param.name), C.size_t(len(values)), cIndicesPtr, cValues)
  runtime.KeepAlive(tuples)
  return amplError(errInfo)
}

// SetAllDoubleValues assigns the given values to every instance of the
// parameter, in the same order as the indices of the entity.
func (param *Parameter) SetAllDoubleValues(values []float64) error {
  var cValues *C.double
  if len(values) > 0 {
    cValues = (*C.double)(unsafe.Pointer(&values[0]))
  }

  errInfo := C.AMPL_ParameterSetArgsDoubleValues(param.ampl.campl, C.CString(param.name), C.size_t(len(values)), cValues)
  return amplError(errInfo)
}
