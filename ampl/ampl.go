package ampl

/*
#cgo LDFLAGS: -lampl
#include "ampl/ampl_c.h"
#include <stdlib.h>
*/
import "C"
import (
  "fmt"
  "strconv"
	"unsafe"
  "runtime"
)

type Ampl struct {
  campl *C.AMPL
}

/*
NewModel
Description:

	Creates a new model from a given environment.
*/
func New() (*Ampl, error) {
	var campl *C.AMPL
  errInfo := C.AMPL_Create(&campl)
	if errInfo != nil {
		return nil, amplError(errInfo)
	}
  wrapper := &Ampl{campl: campl}
  runtime.SetFinalizer(wrapper, func(ampl *Ampl) {
    if ampl.campl != nil {
      C.AMPL_Free(&ampl.campl)
      ampl.campl = nil
    }
  })

	return wrapper, nil
}

func (ampl *Ampl) Eval(name string) error {
  errInfo := C.AMPL_Eval(ampl.campl, C.CString(name))
  return amplError(errInfo)
}

func (ampl *Ampl) Reset() error {
  errInfo := C.AMPL_Reset(ampl.campl)
  return amplError(errInfo)
}

func (ampl *Ampl) Close() error {
  errInfo := C.AMPL_Close(ampl.campl)
  return amplError(errInfo)
}

func (ampl *Ampl) IsRunning() (bool, error) {
  var running C.bool
  errInfo := C.AMPL_IsRunning(ampl.campl, &running)
  if errInfo != nil {
    return false, amplError(errInfo)
  }
  return running == C.bool(true), nil
}

func (ampl *Ampl) IsBusy() (bool, error) {
  var busy C.bool
  errInfo := C.AMPL_IsBusy(ampl.campl, &busy)
  if errInfo != nil {
    return false, amplError(errInfo)
  }
  return busy == C.bool(true), nil
}

func (ampl *Ampl) Read(filename string) error {
  errInfo := C.AMPL_Read(ampl.campl, C.CString(filename))
  return amplError(errInfo)
}

func (ampl *Ampl) ReadData(filename string) error {
  errInfo := C.AMPL_ReadData(ampl.campl, C.CString(filename))
  return amplError(errInfo)
}

func (ampl *Ampl) Solve(problem string, solver string) error {
  errInfo := C.AMPL_Solve(ampl.campl, C.CString(problem), C.CString(solver))
  return amplError(errInfo)
}

// Write invokes the solver on the current problem and writes the solver
// files (.nl, auxiliary files) without actually solving it.
func (ampl *Ampl) Write(filename string, auxfiles string) error {
  errInfo := C.AMPL_Write(ampl.campl, C.CString(filename), C.CString(auxfiles))
  return amplError(errInfo)
}

// SetData assigns the values in a DataFrame to the modelling entities in
// the current environment. If setName is not empty, the indices of the
// DataFrame are used to assign values to the set with that name.
func (ampl *Ampl) SetData(dataframe *DataFrame, setName string) error {
  errInfo := C.AMPL_SetData(ampl.campl, dataframe.cdataframe, C.CString(setName))
  return amplError(errInfo)
}

func setOption[T any](campl *C.AMPL, name string, value T) error {
  var errInfo *C.AMPL_ERRORINFO

  switch v := any(value).(type) {
  case string:
    errInfo = C.AMPL_SetOption(campl, C.CString(name), C.CString(v))
  case int:
    errInfo = C.AMPL_SetIntOption(campl, C.CString(name), C.int(v))
  case float64:
    errInfo = C.AMPL_SetDblOption(campl, C.CString(name), C.double(v))
  case bool:
    errInfo = C.AMPL_SetBoolOption(campl, C.CString(name), C.bool(v))
  default:
    return fmt.Errorf("unsupported type %T", value)
  }

  return amplError(errInfo)
}

// SetOption sets an option with a value.
func (ampl *Ampl) SetOption(name string, value any) error {
  return setOption(ampl.campl, name, value)
}

// GetOption retrieves the current value of the specified option as type T.
// It takes an *Ampl instance as an argument and handles different types manually.
func (ampl *Ampl) GetOption(name string) (interface{}, error) {
	// Zero value to return in case of errors
	var zeroValue interface{}

	// Prepare C variables
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName)) // Free the allocated memory for cName

	var exists C.bool
	var value *C.char

	// Call the C API to get the option
	errInfo := C.AMPL_GetOption(ampl.campl, cName, &exists, &value)
	if errInfo != nil {
		return zeroValue, fmt.Errorf("error retrieving option '%s': %w", name, amplError(errInfo))
	}

	// Free the memory allocated for the C string `value`
	defer C.free(unsafe.Pointer(value))

	// Check if the option exists
	if !bool(exists) {
		return zeroValue, fmt.Errorf("option '%s' does not exist", name)
	}

	// Decode the C string to a Go string
	goValue := C.GoString(value)

	// Attempt to convert the value to different types
	if intValue, err := strconv.Atoi(goValue); err == nil {
		return intValue, nil
	} else if floatValue, err := strconv.ParseFloat(goValue, 64); err == nil {
		return floatValue, nil
	} else {
		return goValue, nil
	}
}

// GetIntOption retrieves the current value of the specified option as an int.
func (ampl *Ampl) GetIntOption(name string) (int, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var exists C.bool
	var value C.int

	errInfo := C.AMPL_GetIntOption(ampl.campl, cName, &exists, &value)
	if errInfo != nil {
		return 0, fmt.Errorf("error retrieving option '%s': %w", name, amplError(errInfo))
	}
	if !bool(exists) {
		return 0, fmt.Errorf("option '%s' does not exist", name)
	}

	return int(value), nil
}

func (ampl *Ampl) ExportModel(fileName string) (string, error) {
  var cOutput *C.char

  errInfo := C.AMPL_ExportModel(ampl.campl, C.CString(fileName), &cOutput)
  if errInfo != nil {
    return "", amplError(errInfo)
  }

  output := C.GoString(cOutput)
  C.AMPL_StringFree(&cOutput)
  return output, nil
}

func (ampl *Ampl) ExportData(fileName string) (string, error) {
  var cOutput *C.char

  errInfo := C.AMPL_ExportData(ampl.campl, C.CString(fileName), &cOutput)
  if errInfo != nil {
    return "", amplError(errInfo)
  }

  output := C.GoString(cOutput)
  C.AMPL_StringFree(&cOutput)
  return output, nil
}

func (ampl *Ampl) Cd() (string, error) {
  var cOutput *C.char

  errInfo := C.AMPL_Cd(ampl.campl, &cOutput)
  if errInfo != nil {
    return "", amplError(errInfo)
  }

  output := C.GoString(cOutput)
  C.AMPL_StringFree(&cOutput)
  return output, nil
}

func (ampl *Ampl) Cd2(path string) (string, error) {
  var cOutput *C.char

  errInfo := C.AMPL_Cd2(ampl.campl, C.CString(path), &cOutput)
  if errInfo != nil {
    return "", amplError(errInfo)
  }

  output := C.GoString(cOutput)
  C.AMPL_StringFree(&cOutput)
  return output, nil
}

func (ampl *Ampl) GetCurrentObjective() (string, error) {
  var cOutput *C.char

  errInfo := C.AMPL_GetCurrentObjective(ampl.campl, &cOutput)
  if errInfo != nil {
    return "", amplError(errInfo)
  }

  output := C.GoString(cOutput)
  C.AMPL_StringFree(&cOutput)
  return output, nil
}

func (ampl *Ampl) ReadTable(tablename string) error {
  errInfo := C.AMPL_ReadTable(ampl.campl, C.CString(tablename))
  return amplError(errInfo)
}

func (ampl *Ampl) WriteTable(tablename string) error {
  errInfo := C.AMPL_WriteTable(ampl.campl, C.CString(tablename))
  return amplError(errInfo)
}

func (ampl *Ampl) GetValue(scalarExpression string) (*Variant, error) {
  var cvariant *C.AMPL_VARIANT

  errInfo := C.AMPL_GetValue(ampl.campl, C.CString(scalarExpression), &cvariant)
  if errInfo != nil {
    return nil, amplError(errInfo)
  }

  return &Variant{cvariant: cvariant}, nil
}

func (ampl *Ampl) GetValueNumeric(scalarExpression string) (float64, error) {
  var value C.double

  errInfo := C.AMPL_GetValueNumeric(ampl.campl, C.CString(scalarExpression), &value)
  if errInfo != nil {
    return 0, amplError(errInfo)
  }

  return float64(value), nil
}

func (ampl *Ampl) GetValueString(scalarExpression string) (string, error) {
  var value *C.char

  errInfo := C.AMPL_GetValueString(ampl.campl, C.CString(scalarExpression), &value)
  if errInfo != nil {
    return "", amplError(errInfo)
  }

  return C.GoString(value), nil
}

func (ampl *Ampl) GetData(displayStatements []string) (*DataFrame, error) {
  var cdataframe *C.AMPL_DATAFRAME

  cStatements := makeCStringArray(displayStatements)
  defer freeCStringArray(cStatements, len(displayStatements))

  errInfo := C.AMPL_GetData(ampl.campl, cStatements, C.size_t(len(displayStatements)), &cdataframe)
  if errInfo != nil {
    return nil, amplError(errInfo)
  }

  return wrapDataFrame(cdataframe), nil
}

func (ampl *Ampl) GetVariables() ([]Variable, error) {
  var size C.size_t
  var cNames **C.char

  errInfo := C.AMPL_GetVariables(ampl.campl, &size, &cNames)
  if errInfo != nil {
    return nil, amplError(errInfo)
  }

  // Convert C strings (**char) to Go strings and create Variable structs
  variables := make([]Variable, size)
  cArray := (*[1 << 30]*C.char)(unsafe.Pointer(cNames))[:size:size]
  for i := 0; i < int(size); i++ {
    // Convert each C string to a Go string and assign it to a Variable struct
    name := C.GoString(cArray[i])
    variables[i] = Variable{
      Entity: Entity{
        ampl: ampl,
        name: name,
      },
    }
    // Free each string
    C.AMPL_StringFree(&cArray[i])
  }

  // Free the array of strings
  C.free(unsafe.Pointer(cNames))

  return variables, nil
}

func (ampl *Ampl) GetConstraints() ([]Constraint, error) {
  var size C.size_t
  var cNames **C.char

  errInfo := C.AMPL_GetConstraints(ampl.campl, &size, &cNames)
  if errInfo != nil {
    return nil, amplError(errInfo)
  }

  // Convert C strings (**char) to Go strings and create Constraint structs
  constraints := make([]Constraint, size)
  cArray := (*[1 << 30]*C.char)(unsafe.Pointer(cNames))[:size:size]
  for i := 0; i < int(size); i++ {
    // Convert each C string to a Go string and assign it to a Constraint struct
    name := C.GoString(cArray[i])
    constraints[i] = Constraint{
      Entity: Entity{
        ampl: ampl,
        name: name,
      },
    }
    // Free each string
    C.AMPL_StringFree(&cArray[i])
  }

  // Free the array of strings
  C.free(unsafe.Pointer(cNames))

  return constraints, nil
}

func (ampl *Ampl) GetParameters() ([]Parameter, error) {
  var size C.size_t
  var cNames **C.char

  errInfo := C.AMPL_GetParameters(ampl.campl, &size, &cNames)
  if errInfo != nil {
    return nil, amplError(errInfo)
  }

  // Convert C strings (**char) to Go strings and create Parameter structs
  parameters := make([]Parameter, size)
  cArray := (*[1 << 30]*C.char)(unsafe.Pointer(cNames))[:size:size]
  for i := 0; i < int(size); i++ {
    // Convert each C string to a Go string and assign it to a Parameter struct
    name := C.GoString(cArray[i])
    parameters[i] = Parameter{
      Entity: Entity{
        ampl: ampl,
        name: name,
      },
    }
    // Free each string
    C.AMPL_StringFree(&cArray[i])
  }

  // Free the array of strings
  C.free(unsafe.Pointer(cNames))

  return parameters, nil
}

func (ampl *Ampl) GetObjectives() ([]Objective, error) {
  var size C.size_t
  var cNames **C.char

  errInfo := C.AMPL_GetObjectives(ampl.campl, &size, &cNames)
  if errInfo != nil {
    return nil, amplError(errInfo)
  }

  // Convert C strings (**char) to Go strings and create Objective structs
  objectives := make([]Objective, size)
  cArray := (*[1 << 30]*C.char)(unsafe.Pointer(cNames))[:size:size]
  for i := 0; i < int(size); i++ {
    // Convert each C string to a Go string and assign it to an Objective struct
    name := C.GoString(cArray[i])
    objectives[i] = Objective{
      Entity: Entity{
        ampl: ampl,
        name: name,
      },
    }
    // Free each string
    C.AMPL_StringFree(&cArray[i])
  }

  // Free the array of strings
  C.free(unsafe.Pointer(cNames))

  return objectives, nil
}

func (ampl *Ampl) GetSets() ([]Set, error) {
  var size C.size_t
  var cNames **C.char

  errInfo := C.AMPL_GetSets(ampl.campl, &size, &cNames)
  if errInfo != nil {
    return nil, amplError(errInfo)
  }

  // Convert C strings (**char) to Go strings and create Set structs
  sets := make([]Set, size)
  cArray := (*[1 << 30]*C.char)(unsafe.Pointer(cNames))[:size:size]
  for i := 0; i < int(size); i++ {
    // Convert each C string to a Go string and assign it to a Set struct
    name := C.GoString(cArray[i])
    sets[i] = Set{
      Entity: Entity{
        ampl: ampl,
        name: name,
      },
    }
    // Free each string
    C.AMPL_StringFree(&cArray[i])
  }

  // Free the array of strings
  C.free(unsafe.Pointer(cNames))

  return sets, nil
}

func (ampl *Ampl) getEntity(name string, wantType C.AMPL_ENTITYTYPE, typeName string) error {
  var cType C.AMPL_ENTITYTYPE

  errInfo := C.AMPL_EntityGetType(ampl.campl, C.CString(name), &cType)
  if errInfo != nil {
    return amplError(errInfo)
  }

  if cType != wantType {
    return fmt.Errorf("%s is not %s", name, typeName)
  }

  return nil
}

func (ampl *Ampl) GetVariable(name string) (*Variable, error) {
  if err := ampl.getEntity(name, C.AMPL_VARIABLE, "a variable"); err != nil {
    return nil, err
  }

  return &Variable{
    Entity: Entity{
      ampl: ampl,
      name: name,
    },
  }, nil
}

func (ampl *Ampl) GetConstraint(name string) (*Constraint, error) {
  if err := ampl.getEntity(name, C.AMPL_CONSTRAINT, "a constraint"); err != nil {
    return nil, err
  }

  return &Constraint{
    Entity: Entity{
      ampl: ampl,
      name: name,
    },
  }, nil
}

func (ampl *Ampl) GetParameter(name string) (*Parameter, error) {
  if err := ampl.getEntity(name, C.AMPL_PARAMETER, "a parameter"); err != nil {
    return nil, err
  }

  return &Parameter{
    Entity: Entity{
      ampl: ampl,
      name: name,
    },
  }, nil
}

func (ampl *Ampl) GetObjective(name string) (*Objective, error) {
  if err := ampl.getEntity(name, C.AMPL_OBJECTIVE, "an objective"); err != nil {
    return nil, err
  }

  return &Objective{
    Entity: Entity{
      ampl: ampl,
      name: name,
    },
  }, nil
}

func (ampl *Ampl) GetSet(name string) (*Set, error) {
  if err := ampl.getEntity(name, C.AMPL_SET, "a set"); err != nil {
    return nil, err
  }

  return &Set{
    Entity: Entity{
      ampl: ampl,
      name: name,
    },
  }, nil
}
