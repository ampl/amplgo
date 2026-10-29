package ampl

//#cgo LDFLAGS: -L./lib/amd64 -lampl
//#cgo CFLAGS: -I./include
//#include "ampl/ampl_c.h"
import "C"
//import (
  //"runtime"
  //"unsafe"
  //"errors"
  //"fmt"
//)

type Set struct {
  Entity
}

func (set Set) GetArity() (int, error) {
  var arity C.size_t
  errInfo := C.AMPL_SetGetArity(set.ampl.campl, C.CString(set.name), &arity)
  if errInfo != nil {
    return 0, amplError(errInfo)
  }
  return int(arity), nil
}

