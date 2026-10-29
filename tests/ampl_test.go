package tests

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/ampl/ampl.go/ampl"
)

func TestSetOption(t *testing.T) {
  model, err := ampl.New()
  assert.Nil(t, err)
  model.SetOption("a", "s")
  model.SetOption("b", 123)
  model.SetOption("c", 1.23)

  a, a_error := model.GetOption("a")
  b, b_error := model.GetOption("b")
  c, c_error := model.GetOption("c")

  assert.Equal(t, a, "s")
  assert.Equal(t, b, 123)
  assert.Equal(t, c, 1.23)
  assert.Nil(t, a_error)
  assert.Nil(t, b_error)
  assert.Nil(t, c_error)
}

func TestSolveArguments(t *testing.T) {
  model, err := ampl.New()
  assert.Nil(t, err)
  model.Eval(`
    param n integer > 0; # N-queens
    var Row {1..n} integer >= 1 <= n;
    s.t. row_attacks: alldiff ({j in 1..n} Row[j]);
    s.t. diag_attacks: alldiff ({j in 1..n} Row[j]+j);
    s.t. rdiag_attacks: alldiff ({j in 1..n} Row[j]-j);
  `)
  paramN, errorN := model.GetParameter("n")
  assert.Nil(t, errorN)
  paramN.Set(5)
  model.SetOption("highs_options", "")

  //output := 
  model.Solve("highs", "")
  //assert.False(t, strings.Contains(output, "outlev"))

  model.SetOption("highs_options", "outlev=1")
  //output =
  model.Solve("highs", "")
  //assert.True(t, strings.Contains(output, "outlev"))
}

func TestInstances(t *testing.T) {
  model, err := ampl.New()
  assert.Nil(t, err)
  model.Eval("param y{1..5};")

  y, _ := model.GetParameter("y")
  y_numinstances, _ := y.GetNumInstances()
  y_instance, _ := y.GetInstances()

  assert.Equal(t, y_numinstances, len(y_instance))
  assert.Equal(t, len(y_instance), 5)

  name0, _ := y_instance[0].GetName()
  name1, _ := y_instance[1].GetName()
  name2, _ := y_instance[2].GetName()
  name3, _ := y_instance[3].GetName()
  name4, _ := y_instance[4].GetName()
  assert.Equal(t, name0, "y[1]")
  assert.Equal(t, name1, "y[2]")
  assert.Equal(t, name2, "y[3]")
  assert.Equal(t, name3, "y[4]")
  assert.Equal(t, name4, "y[5]")
}