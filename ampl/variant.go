package ampl

//#include "ampl/ampl_c.h"
import "C"
import (
	"errors"
	"runtime"
  )

type Variant struct {
  cvariant *C.AMPL_VARIANT
}

func wrapVariant(cvariant *C.AMPL_VARIANT) (*Variant) {
  wrapper := &Variant{cvariant: cvariant}

  runtime.SetFinalizer(wrapper, func(variant *Variant) {
    if variant.cvariant != nil {
      C.AMPL_VariantFree(&variant.cvariant)
      variant.cvariant = nil
    }
  })

  return wrapper
}

func NewEmpty() (*Variant, error) {
  var cvariant *C.AMPL_VARIANT
  if result := C.AMPL_VariantCreateEmpty(&cvariant); result != 0 {
    return nil, errors.New("AMPL_VariantCreateEmpty failed")
  }

  return wrapVariant(cvariant), nil
}

func NewNumeric(value float64) (*Variant, error) {
  var cvariant *C.AMPL_VARIANT
  if result := C.AMPL_VariantCreateNumeric(&cvariant, C.double(value)); result != 0 {
    return nil, errors.New("AMPL_VariantCreateNumeric failed")
  }

  return wrapVariant(cvariant), nil
}

func NewString(value string) (*Variant, error) {
  var cvariant *C.AMPL_VARIANT
  if result := C.AMPL_VariantCreateString(&cvariant, C.CString(value)); result != 0 {
    return nil, errors.New("AMPL_VariantCreateString failed")
  }

  return wrapVariant(cvariant), nil
}

func (variant *Variant) GetNumericValue() (float64, error) {
  var value C.double
  if result := C.AMPL_VariantGetNumericValue(variant.cvariant, &value); result != 0 {
    return 0, errors.New("AMPL_VariantGetNumericValue failed")
  }

  return float64(value), nil
}

func (variant *Variant) GetStringValue() (string, error) {
  var value *C.char
  if result := C.AMPL_VariantGetStringValue(variant.cvariant, &value); result != 0 {
    return "", errors.New("AMPL_VariantGetStringValue failed")
  }

  return C.GoString(value), nil
}

func (variant *Variant) GetFormat() (string, error) {
  var value *C.char
  if result := C.AMPL_VariantFormat(variant.cvariant, &value); result != 0 {
    return "", errors.New("AMPL_VariantFormat failed")
  }

  return C.GoString(value), nil
}
