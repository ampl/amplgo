package ampl

//#include "ampl/ampl_c.h"
import "C"

type Variable struct {
  Entity
}


// Get the current value of this variable.
func (variable Variable) Value() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_VALUE)
}

// Fix all instances of this variable to their current value otherwise.
func (variable Variable) Fix() error {
  errInfo := C.AMPL_VariableFix(variable.ampl.campl, C.CString(variable.name))
  return amplError(errInfo)
}

// Fix all instances of this variable to a value.
func (variable Variable) FixWithValue(value float64) error {
  errInfo := C.AMPL_VariableFixWithValue(variable.ampl.campl, C.CString(variable.name), C.double(value))
  return amplError(errInfo)
}

// Unfix all instances of this variable.
func (variable Variable) Unfix() error {
  errInfo := C.AMPL_VariableUnfix(variable.ampl.campl, C.CString(variable.name))
  return amplError(errInfo)
}

// Set the current value of this variable (does not fix it), equivalent to the AMPL command `let`.
func (variable Variable) SetValue(value float64) error {
  errInfo := C.AMPL_VariableSetValue(variable.ampl.campl, C.CString(variable.name), C.double(value))
  return amplError(errInfo)
}

// Get the AMPL status (fixed, presolved, or substituted out).
func (variable Variable) Astatus() (string, error) {
  return variable.getStringSuffix(C.AMPL_ASTATUS)
}

// Get the index in `_con` of "defining constraint" used to substitute variable out.
func (variable Variable) Defeqn() (int, error) {
  return variable.getIntSuffix(C.AMPL_DEFEQN)
}

// Get the dual value on defining constraint of variable substituted out.
func (variable Variable) Dual() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_DUAL)
}

// Get the current initial guess.
func (variable Variable) Init() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_INIT)
}

// Get the original initial guess (set by `:=` or`default` or by a data statement).
func (variable Variable) Init0() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_INIT0)
}

// Returns the current lower bound.
func (variable Variable) Lb() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_LB)
}

// Returns the current upper bound.
func (variable Variable) Ub() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_UB)
}

// Returns the initial lower bounds, from the var declaration.
func (variable Variable) Lb0() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_LB0)
}

// Returns the initial upper bound, from the var declaration.
func (variable Variable) Ub0() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_UB0)
}

// Returns the weaker lower bound from AMPL's presolve phase.
func (variable Variable) Lb1() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_LB1)
}

// Returns the weaker upper bound from AMPL's presolve phase.
func (variable Variable) Ub1() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_UB1)
}

// Returns the stronger lower bound from AMPL's presolve phase.
func (variable Variable) Lb2() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_LB2)
}

// Returns the stronger upper bound from AMPL's presolve phase.
func (variable Variable) Ub2() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_UB2)
}

// Returns the reduced cost at lower bound.
func (variable Variable) Lrc() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_LRC)
}

// Returns the reduced cost at upper bound.
func (variable Variable) Urc() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_URC)
}

// Return the slack at lower bound (``val - lb``).
func (variable Variable) Lslack() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_LSLACK)
}

// Return the slack at upper bound (``ub - val``).
func (variable Variable) Uslack() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_USLACK)
}

// Get the reduced cost (at the nearer bound).
func (variable Variable) Rc() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_RC)
}

// Returns the bound slack which is the lesser of `Lslack` and `Uslack`.
func (variable Variable) Slack() (float64, error) {
  return variable.getDoubleSuffix(C.AMPL_SLACK)
}

// Solver status (basis status of variable).
func (variable Variable) Sstatus() (string, error) {
  return variable.getStringSuffix(C.AMPL_SSTATUS)
}

// AMPL status if not `in`, otherwise solver status.
func (variable Variable) Status() (string, error) {
  return variable.getStringSuffix(C.AMPL_STATUS)
}
