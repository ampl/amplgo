package ampl

//#cgo LDFLAGS: -L./lib/amd64 -lampl
//#cgo CFLAGS: -I./include
//#include "ampl/ampl_c.h"
import "C"
import "fmt"

type Constraint struct {
  Entity
}


// Check if the constraint is a logical constraint. Returns true iff logical.
func (constraint Constraint) IsLogical() (bool, error) {
  var cValue C.bool

  errInfo := C.AMPL_ConstraintIsLogical(constraint.ampl.campl, C.CString(constraint.name), &cValue)
  if errInfo != nil {
    return false, amplError(errInfo)
  }

  return bool(cValue), nil
}

// Drop this constraint.
func (constraint Constraint) Drop() error {
  errInfo := C.AMPL_EntityDrop(constraint.ampl.campl, C.CString(constraint.name))
  return amplError(errInfo)
}

// Restore this constraint (if it had been dropped, no effect otherwise).
func (constraint Constraint) Restore() error {
  errInfo := C.AMPL_EntityRestore(constraint.ampl.campl, C.CString(constraint.name))
  return amplError(errInfo)
}

// Get the current value of the constraint's body.
func (constraint Constraint) Body() (float64, error) {
  return constraint.getDoubleSuffix(C.AMPL_BODY)
}

// Get the current AMPL status (dropped, presolved, or substituted out).
func (constraint Constraint) Astatus() (string, error) {
  return constraint.getStringSuffix(C.AMPL_ASTATUS)
}

// Get the index in `_var` of "defined variable" substituted out by the constraint.
func (constraint Constraint) Defvar() (int, error) {
  return constraint.getIntSuffix(C.AMPL_DEFVAR)
}

// Get the current initial guess for the constraint's dual variable.
func (constraint Constraint) Dinit() (float64, error) {
  return constraint.getDoubleSuffix(C.AMPL_DINIT)
}

// Get the original initial guess for the constraint's dual variable.
func (constraint Constraint) Dinit0() (float64, error) {
  return constraint.getDoubleSuffix(C.AMPL_DINIT0)
}

// Get the current value of the constraint's dual variable.
func (constraint Constraint) Dual() (float64, error) {
  return constraint.getDoubleSuffix(C.AMPL_DUAL)
}

// Get the current value of the constraint's lower bound.
func (constraint Constraint) Lb() (float64, error) {
  return constraint.getDoubleSuffix(C.AMPL_LB)
}

// Get the current value of the constraint's upper bound.
func (constraint Constraint) Ub() (float64, error) {
  return constraint.getDoubleSuffix(C.AMPL_UB)
}

// Get the constraint lower bound sent to the solver (reflecting adjustment for fixed variables).
func (constraint Constraint) Lbs() (float64, error) {
  return constraint.getDoubleSuffix(C.AMPL_LBS)
}

// Get the constraint upper bound sent to the solver (reflecting adjustment for fixed variables).
func (constraint Constraint) Ubs() (float64, error) {
  return constraint.getDoubleSuffix(C.AMPL_UBS)
}

// Get the current dual value associated with the lower bound.
func (constraint Constraint) Ldual() (float64, error) {
  return constraint.getDoubleSuffix(C.AMPL_LDUAL)
}

// Get the current dual value associated with the upper bounds.
func (constraint Constraint) Udual() (float64, error) {
  return constraint.getDoubleSuffix(C.AMPL_UDUAL)
}

// Get the slack at lower bound `body - lb`.
func (constraint Constraint) Lslack() (float64, error) {
  return constraint.getDoubleSuffix(C.AMPL_LSLACK)
}

// Get the slack at upper bound `ub - body`.
func (constraint Constraint) Uslack() (float64, error) {
  return constraint.getDoubleSuffix(C.AMPL_USLACK)
}

// Constraint slack (the lesser of lslack and uslack).
func (constraint Constraint) Slack() (float64, error) {
  return constraint.getDoubleSuffix(C.AMPL_SLACK)
}

// Get the solver status (basis status of constraint's slack or artificial variable).
func (constraint Constraint) Sstatus() (string, error) {
  return constraint.getStringSuffix(C.AMPL_SSTATUS)
}

// Get the AMPL status if not `in`, otherwise solver status.
func (constraint Constraint) Status() (string, error) {
  return constraint.getStringSuffix(C.AMPL_STATUS)
}

// Set the value of the dual variable associated to this constraint (valid only if the constraint is scalar).
func (constraint Constraint) SetDual(dual float64) error {
  errInfo := C.AMPL_ConstraintSetDual(constraint.ampl.campl, C.CString(constraint.name), C.double(dual))
  return amplError(errInfo)
}

// Get the AMPL val suffix. Valid only for logical constraints.
func (constraint Constraint) Val() (float64, error) {
	isLogical, err := constraint.IsLogical()
	if err != nil {
		return 0, err
	}

	if !isLogical {
		return 0, fmt.Errorf("%s is not a logical constraint", constraint.name)
	}

	return constraint.getDoubleSuffix(C.AMPL_VAL)
}
