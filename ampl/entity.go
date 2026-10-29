package ampl

//#include "ampl/ampl_c.h"
import "C"
import (
  "unsafe"
)

type Entity struct {
  ampl *Ampl
  name string
}

func (entity Entity) GetName() string {
  return entity.name
}

func (entity Entity) GetIndexarity() (int, error) {
  var arity C.size_t
  errInfo := C.AMPL_EntityGetIndexarity(entity.ampl.campl, C.CString(entity.name), &arity)
  if errInfo != nil {
    return 0, amplError(errInfo)
  }
  return int(arity), nil
}

func (entity Entity) GetNumInstances() (int, error) {
  var size C.size_t
  errInfo := C.AMPL_EntityGetNumInstances(entity.ampl.campl, C.CString(entity.name), &size)
  if errInfo != nil {
    return 0, amplError(errInfo)
  }
  return int(size), nil
}

func (entity Entity) GetInstances() ([]Instance, error) {
  var size C.size_t
  var cIndicies **C.AMPL_TUPLE

  errInfo := C.AMPL_EntityGetTuples(entity.ampl.campl, C.CString(entity.name), &cIndicies, &size)
  if errInfo != nil {
    return nil, amplError(errInfo)
  }

  instances := make([]Instance, size)
  for i := 0; i < int(size); i++ {
    // Perform pointer arithmetic to get the i-th tuple
		cTuplePtr := (**C.AMPL_TUPLE)(unsafe.Pointer(uintptr(unsafe.Pointer(cIndicies)) + uintptr(i)*unsafe.Sizeof(*cIndicies)))

    instances[i] = Instance{
      entity: &Entity{
        ampl: entity.ampl,
        name: entity.name,
      },
      ctuple: *cTuplePtr,
    }
  }
  return instances, nil
}

func (entity Entity) GetValues(suffixes []string) (*DataFrame, error) {
  var cdataframe *C.AMPL_DATAFRAME

  var cSuffixes **C.char
  if len(suffixes) > 0 {
    cSuffixes = makeCStringArray(suffixes)
    defer freeCStringArray(cSuffixes, len(suffixes))
  }

  errInfo := C.AMPL_EntityGetValues(entity.ampl.campl, C.CString(entity.name), cSuffixes, C.size_t(len(suffixes)), &cdataframe)
  if errInfo != nil {
    return nil, amplError(errInfo)
  }

  return wrapDataFrame(cdataframe), nil
}

func (entity Entity) getDoubleSuffix(suffix C.AMPL_NUMERICSUFFIX) (float64, error) {
  var value C.double
  errInfo := C.AMPL_InstanceGetDoubleSuffix(entity.ampl.campl, C.CString(entity.name), nil, suffix, &value)
  if errInfo != nil {
    return 0, amplError(errInfo)
  }
  return float64(value), nil
}

func (entity Entity) getIntSuffix(suffix C.AMPL_NUMERICSUFFIX) (int, error) {
  var value C.int
  errInfo := C.AMPL_InstanceGetIntSuffix(entity.ampl.campl, C.CString(entity.name), nil, suffix, &value)
  if errInfo != nil {
    return 0, amplError(errInfo)
  }
  return int(value), nil
}

func (entity Entity) getStringSuffix(suffix C.AMPL_STRINGSUFFIX) (string, error) {
  var value *C.char
  errInfo := C.AMPL_InstanceGetStringSuffix(entity.ampl.campl, C.CString(entity.name), nil, suffix, &value)
  if errInfo != nil {
    return "", amplError(errInfo)
  }
  result := C.GoString(value)
  C.AMPL_StringFree(&value)
  return result, nil
}

func (entity Entity) GetInstance(index Tuple) (Instance, error) {
  instance := Instance{
    entity: &Entity{
      ampl: entity.ampl,
      name: entity.name,
    },
    ctuple: index.ctuple,
  }

  return instance, nil
}
