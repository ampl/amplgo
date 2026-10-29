package ampl

//#cgo LDFLAGS: -L./lib/amd64 -lampl
//#cgo CFLAGS: -I./include
//#include "ampl/ampl_c.h"
import "C"
import (
  //"runtime"
  //"unsafe"
  //"errors"
  //"reflect"
  //"fmt"
)

type ParameterInstance struct {
  Instance
}