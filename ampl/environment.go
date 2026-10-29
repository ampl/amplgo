package ampl

//#include "ampl/ampl_c.h"
import "C"
import "errors"

type Environment struct {
	cenv *C.AMPL_ENVIRONMENT
}

func NewEnv(binarydir string, binaryname string) (*Environment, error) {
  var cenv *C.AMPL_ENVIRONMENT
  if result := C.AMPL_EnvironmentCreate(&cenv, C.CString(binarydir), C.CString(binaryname)); result != 0 {
    return nil, errors.New("AMPL_EnvironmentCreate failed")
  }

  return &Environment{cenv: cenv}, nil
}

func (env *Environment) FreeEnv() error {
  if env == nil {
    return nil
  }
  if result := C.AMPL_EnvironmentFree(&env.cenv); result != 0 {
    return errors.New("AMPL_EnvironmentFree failed")
  }
  return nil
}

func (env *Environment) AddEnvironmentVariable(name string, value string) error {
  if result := C.AMPL_EnvironmentAddEnvironmentVariable(env.cenv, C.CString(name), C.CString(value)); result != 0 {
    return errors.New("AMPL_EnvironmentAddEnvironmentVariable failed")
  }
  return nil
}

func (env *Environment) GetBinaryDirectory() (string, error) {
  var cOutput *C.char
  if result := C.AMPL_EnvironmentGetBinaryDirectory(env.cenv, &cOutput); result != 0 {
    return "", errors.New("AMPL_EnvironmentGetBinaryDirectory failed")
  }

  return C.GoString(cOutput), nil
}

func (env *Environment) GetAMPLCommand() (string, error) {
  var cOutput *C.char
  if result := C.AMPL_EnvironmentGetAMPLCommand(env.cenv, &cOutput); result != 0 {
    return "", errors.New("AMPL_EnvironmentGetAMPLCommand failed")
  }

  return C.GoString(cOutput), nil
}

func (env *Environment) SetBinaryDirectory(binarydir string) error {
  if result := C.AMPL_EnvironmentSetBinaryDirectory(env.cenv, C.CString(binarydir)); result != 0 {
    return errors.New("AMPL_EnvironmentSetBinaryDirectory failed")
  }
  return nil
}

func (env *Environment) GetBinaryName() (string, error) {
  var cOutput *C.char
  if result := C.AMPL_EnvironmentGetBinaryName(env.cenv, &cOutput); result != 0 {
    return "", errors.New("AMPL_EnvironmentGetBinaryName failed")
  }

  return C.GoString(cOutput), nil
}

func (env *Environment) SetBinaryName(binaryname string) error {
  if result := C.AMPL_EnvironmentSetBinaryName(env.cenv, C.CString(binaryname)); result != 0 {
    return errors.New("AMPL_EnvironmentSetBinaryName failed")
  }
  return nil
}

func (env *Environment) ToString() (string, error) {
  var cOutput *C.char
  if result := C.AMPL_EnvironmentToString(env.cenv, &cOutput); result != 0 {
    return "", errors.New("AMPL_EnvironmentToString failed")
  }

  return C.GoString(cOutput), nil
}
