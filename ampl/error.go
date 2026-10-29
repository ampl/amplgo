package ampl

/*
#cgo LDFLAGS: -L./lib/amd64 -lampl
#cgo CFLAGS: -I./include
#include "ampl/ampl_c.h"
*/
import "C"
import "errors"

// amplError converts an AMPL_ERRORINFO into a Go error. Returns nil if
// errInfo is nil.
func amplError(errInfo *C.AMPL_ERRORINFO) error {
  if errInfo == nil {
    return nil
  }
  return errors.New(C.GoString(C.AMPL_ErrorInfoGetMessage(errInfo)))
}
