package ampl

//#cgo LDFLAGS: -L./lib/amd64 -lampl
//#cgo CFLAGS: -I./include
//#include "ampl/ampl_c.h"
import "C"

type Objective struct {
  Entity
}

// Get the value of the objective.
func (objective Objective) Value() (float64, error) {
  return objective.getDoubleSuffix(C.AMPL_VALUE)
}

// Return the AMPL status.
func (objective Objective) Astatus() (string, error) {
  return objective.getStringSuffix(C.AMPL_ASTATUS)
}

// Return the solver status.
func (objective Objective) Sstatus() (string, error) {
  return objective.getStringSuffix(C.AMPL_SSTATUS)
}

// Exit code returned by solver after most recent solve with this objective.
func (objective Objective) Exitcode() (int, error) {
  return objective.getIntSuffix(C.AMPL_EXITCODE)
}

// Result message returned by solver after most recent solve with this objective.
func (objective Objective) Message() (string, error) {
  return objective.getStringSuffix(C.AMPL_MESSAGE)
}

// Result string returned by solver after most recent solve with this objective.
func (objective Objective) Result() (string, error) {
  return objective.getStringSuffix(C.AMPL_RESULT)
}

// Drop this objective.
func (objective Objective) Drop() error {
  errInfo := C.AMPL_EntityDrop(objective.ampl.campl, C.CString(objective.name))
  return amplError(errInfo)
}

// Restore this objective (if it had been dropped, no effect otherwise).
func (objective Objective) Restore() error {
  errInfo := C.AMPL_EntityRestore(objective.ampl.campl, C.CString(objective.name))
  return amplError(errInfo)
}

// Get the sense of this objective. Returns true if minimize, false if maximize.
func (objective Objective) Minimization() (bool, error) {
  value, err := objective.getStringSuffix(C.AMPL_SENSE)
  if err != nil {
    return false, err
  }
  return value == "minimize", nil
}
