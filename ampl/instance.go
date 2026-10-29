package ampl

//#include "ampl/ampl_c.h"
import "C"

type Instance struct {
  entity *Entity
  ctuple *C.AMPL_TUPLE
}

type NumericSuffix int64

const (
  // VARIABLES
	Value    NumericSuffix = 0
	Defeqn   NumericSuffix = 1
	Dual     NumericSuffix = 2
	Init     NumericSuffix = 3
  Init0    NumericSuffix = 4
  Lb       NumericSuffix = 5
  Ub       NumericSuffix = 6
  Lb0      NumericSuffix = 7
  Ub0      NumericSuffix = 8
  Lb1      NumericSuffix = 9
  Ub1      NumericSuffix = 10
  Lb2      NumericSuffix = 11
  Ub2      NumericSuffix = 12
  Lrc      NumericSuffix = 13
  Urc      NumericSuffix = 14
  Lslack   NumericSuffix = 15
  Uslack   NumericSuffix = 16
  Rc       NumericSuffix = 17
  Slack    NumericSuffix = 18

  // CONSTRAINTS
  Body     NumericSuffix = 19
  Defvar   NumericSuffix = 20
  Dinit    NumericSuffix = 21
  Dinit0   NumericSuffix = 22
  Lbs      NumericSuffix = 23
  Ubs      NumericSuffix = 24
  Ldual    NumericSuffix = 25
  Udual    NumericSuffix = 26
  Val      NumericSuffix = 27  // for logical constraints

  // OBJECTIVES
  Exitcode NumericSuffix  = 28
)

func (instance *Instance) GetName() (string, error) {
  var cName *C.char

  errInfo := C.AMPL_InstanceGetName(instance.entity.ampl.campl, C.CString(instance.entity.name), instance.ctuple, &cName)
  if errInfo != nil {
    return "", amplError(errInfo)
  }

  name := C.GoString(cName)
  C.AMPL_StringFree(&cName)
  return name, nil
}

func (instance *Instance) ToString() (string, error) {
  var cBuffer *C.char

  errInfo := C.AMPL_InstanceToString(instance.entity.ampl.campl, C.CString(instance.entity.name), instance.ctuple, &cBuffer)
  if errInfo != nil {
    return "", amplError(errInfo)
  }

  buffer := C.GoString(cBuffer)
  C.AMPL_StringFree(&cBuffer)
  return buffer, nil
}
