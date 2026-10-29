package ampl

//#include "ampl/ampl_c.h"
import "C"
import (
	"runtime"
	"errors"
	"unsafe"
	//"fmt"
  )

type Tuple struct {
  ctuple *C.AMPL_TUPLE
}

func NewNumerics(values []float64) (*Tuple, error) {
  cValues := (*C.double)(unsafe.Pointer(&values[0]))

  var cTuple *C.AMPL_TUPLE

  result := C.AMPL_TupleCreateNumeric(&cTuple, C.size_t(len(values)), cValues)
  if result != 0 {
    return nil, errors.New("AMPL_TupleCreateNumeric failed")
  }

  wrapper := &Tuple{ctuple: cTuple}

  runtime.SetFinalizer(wrapper, func(tuple *Tuple) {
    if tuple.ctuple != nil {
      C.AMPL_TupleFree(&tuple.ctuple)
      tuple.ctuple = nil
    }
  })

  return wrapper, nil
}

func NewStrings(values []string) (*Tuple, error) {
  cValues := makeCStringArray(values)
  defer freeCStringArray(cValues, len(values))

  var cTuple *C.AMPL_TUPLE

  result := C.AMPL_TupleCreateString(&cTuple, C.size_t(len(values)), cValues)
  if result != 0 {
    return nil, errors.New("AMPL_TupleCreateString failed")
  }

  wrapper := &Tuple{ctuple: cTuple}

  runtime.SetFinalizer(wrapper, func(tuple *Tuple) {
    if tuple.ctuple != nil {
      C.AMPL_TupleFree(&tuple.ctuple)
      tuple.ctuple = nil
    }
  })

  return wrapper, nil
}
