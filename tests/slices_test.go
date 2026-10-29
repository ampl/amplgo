package tests

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/ampl/amplgo/ampl"
)

func TestConstraintSlices(t *testing.T) {
  model, err := ampl.New()
  assert.Nil(t, err)
  model.Eval("var x;")
  model.Eval("s.t. c_x: x = 0;")
  constraints, _ := model.GetConstraints()
  assert.Equal(t, len(constraints), 1)

  model.Eval("reset;")
  constraints2, _ := model.GetConstraints()
  assert.Equal(t, len(constraints2), 0)

  model.Eval("var x;")
  model.Eval("s.t. c_x: x = 0;")
  constraints3, _ := model.GetConstraints()
  assert.Equal(t, len(constraints3), 1)

  model.Eval("var y{1..2};")
  model.Eval("s.t. c_y{i in 1..2}: y[i] = i;")
  constraints4, _ := model.GetConstraints()
  assert.Equal(t, len(constraints4), 2)

  model.Eval(`var z{1..2, {"a", "b"}};`)
  model.Eval(`s.t. c_z{(i, j) in {1..2, {"a", "b"}}}: z[i,j] = i;`)
  constraints5, _ := model.GetConstraints()
  assert.Equal(t, len(constraints5), 3)

  model.Eval("var xx{1..2,4..5};")
  model.Eval("s.t. c_xx{(i, j) in {1..2,4..5}}: xx[i,j] = i+j;")
  constraints6, _ := model.GetConstraints()
  assert.Equal(t, len(constraints6), 4)

  c_x, _ := model.GetConstraint("c_x")
  c_x_numinstances, _ := c_x.GetNumInstances()
  assert.Equal(t, c_x_numinstances, 1)

  c_y, _ := model.GetConstraint("c_y")
  c_y_numinstances, _ := c_y.GetNumInstances()
  assert.Equal(t, c_y_numinstances, 2)

  c_z, _ := model.GetConstraint("c_z")
  c_z_numinstances, _ := c_z.GetNumInstances()
  assert.Equal(t, c_z_numinstances, 4)

  c_xx, _ := model.GetConstraint("c_xx")
  c_xx_numinstances, _ := c_xx.GetNumInstances()
  assert.Equal(t, c_xx_numinstances, 4)
}

func TestVariableSlices(t *testing.T) {
  model, err := ampl.New()
  assert.Nil(t, err)
  model.Eval("var x;")
  variables, _ := model.GetVariables()
  assert.Equal(t, len(variables), 1)

  model.Eval("reset;")
  variables2, _ := model.GetVariables()
  assert.Equal(t, len(variables2), 0)

  model.Eval("var x;")
  variables3, _ := model.GetVariables()
  assert.Equal(t, len(variables3), 1)

  model.Eval("var y{1..2};")
  variables4, _ := model.GetVariables()
  assert.Equal(t, len(variables4), 2)

  model.Eval(`var z{1..2, {"a", "b"}};`)
  variables5, _ := model.GetVariables()
  assert.Equal(t, len(variables5), 3)

  model.Eval("var xx{1..2,4..5};")
  variables6, _ := model.GetVariables()
  assert.Equal(t, len(variables6), 4)

  model.Eval("var yy{1..2,4..5,7..8};")
  variables7, _ := model.GetVariables()
  assert.Equal(t, len(variables7), 5)

  model.Eval(`var zz{1..2,{"a", "b"},5..6};`)
  variables8, _ := model.GetVariables()
  assert.Equal(t, len(variables8), 6)

  x, _ := model.GetVariable("x")
  x_numinstances, _ := x.GetNumInstances()
  assert.Equal(t, x_numinstances, 1)

  y, _ := model.GetVariable("y")
  y_numinstances, _ := y.GetNumInstances()
  assert.Equal(t, y_numinstances, 2)

  z, _ := model.GetVariable("z")
  z_numinstances, _ := z.GetNumInstances()
  assert.Equal(t, z_numinstances, 4)

  xx, _ := model.GetVariable("xx")
  xx_numinstances, _ := xx.GetNumInstances()
  assert.Equal(t, xx_numinstances, 4)

  yy, _ := model.GetVariable("yy")
  yy_numinstances, _ := yy.GetNumInstances()
  assert.Equal(t, yy_numinstances, 8)

  zz, _ := model.GetVariable("zz")
  zz_numinstances, _ := zz.GetNumInstances()
  assert.Equal(t, zz_numinstances, 8)
}

func TestParameterSlices(t *testing.T) {
  model, err := ampl.New()
  assert.Nil(t, err)
  model.Eval("param x;")
  params, _ := model.GetParameters()
  assert.Equal(t, len(params), 1)

  model.Eval("reset;")
  params2, _ := model.GetParameters()
  assert.Equal(t, len(params2), 0)

  model.Eval("param x;")
  params3, _ := model.GetParameters()
  assert.Equal(t, len(params3), 1)

  model.Eval("param y{1..2};")
  params4, _ := model.GetParameters()
  assert.Equal(t, len(params4), 2)

  model.Eval(`param z{1..2, {"a", "b"}};`)
  params5, _ := model.GetParameters()
  assert.Equal(t, len(params5), 3)

  model.Eval("param xx{1..2,4..5};")
  params6, _ := model.GetParameters()
  assert.Equal(t, len(params6), 4)

  model.Eval("param yy{1..2,4..5,7..8};")
  params7, _ := model.GetParameters()
  assert.Equal(t, len(params7), 5)

  model.Eval(`param zz{1..2,{"a", "b"},5..6};`)
  params8, _ := model.GetParameters()
  assert.Equal(t, len(params8), 6)

  x, _ := model.GetParameter("x")
  x_numinstances, _ := x.GetNumInstances()
  assert.Equal(t, x_numinstances, 1)

  y, _ := model.GetParameter("y")
  y_numinstances, _ := y.GetNumInstances()
  assert.Equal(t, y_numinstances, 2)

  z, _ := model.GetParameter("z")
  z_numinstances, _ := z.GetNumInstances()
  assert.Equal(t, z_numinstances, 4)

  xx, _ := model.GetParameter("xx")
  xx_numinstances, _ := xx.GetNumInstances()
  assert.Equal(t, xx_numinstances, 4)

  yy, _ := model.GetParameter("yy")
  yy_numinstances, _ := yy.GetNumInstances()
  assert.Equal(t, yy_numinstances, 8)

  zz, _ := model.GetParameter("zz")
  zz_numinstances, _ := zz.GetNumInstances()
  assert.Equal(t, zz_numinstances, 8)
}

func TestObjectiveSlices(t *testing.T) {
  model, err := ampl.New()
  assert.Nil(t, err)
  model.Eval("var x;")
  model.Eval("maximize c_x: x;")
  objectives, _ := model.GetObjectives()
  assert.Equal(t, len(objectives), 1)

  model.Eval("reset;")
  objectives2, _ := model.GetObjectives()
  assert.Equal(t, len(objectives2), 0)

  model.Eval("var x;")
  model.Eval("maximize c_x: x;")
  objectives3, _ := model.GetObjectives()
  assert.Equal(t, len(objectives3), 1)

  model.Eval("var y{1..2};")
  model.Eval("maximize c_y{i in 1..2}: y[i];")
  objectives4, _ := model.GetObjectives()
  assert.Equal(t, len(objectives4), 2)

  model.Eval(`var z{1..2, {"a", "b"}};`)
  model.Eval(`maximize c_z{(i, j) in {1..2, {"a", "b"}}}: z[i,j];`)
  objectives5, _ := model.GetObjectives()
  assert.Equal(t, len(objectives5), 3)

  model.Eval("var xx{1..2,4..5};")
  model.Eval("maximize c_xx{(i, j) in {1..2,4..5}}: xx[i,j];")
  objectives6, _ := model.GetObjectives()
  assert.Equal(t, len(objectives6), 4)

  c_x, _ := model.GetObjective("c_x")
  c_x_numinstances, _ := c_x.GetNumInstances()
  assert.Equal(t, c_x_numinstances, 1)

  c_y, _ := model.GetObjective("c_y")
  c_y_numinstances, _ := c_y.GetNumInstances()
  assert.Equal(t, c_y_numinstances, 2)

  c_z, _ := model.GetObjective("c_z")
  c_z_numinstances, _ := c_z.GetNumInstances()
  assert.Equal(t, c_z_numinstances, 4)

  c_xx, _ := model.GetObjective("c_xx")
  c_xx_numinstances, _ := c_xx.GetNumInstances()
  assert.Equal(t, c_xx_numinstances, 4)
}

func TestSetSlices(t *testing.T) {
  model, err := ampl.New()
  assert.Nil(t, err)
  model.Eval("set x;")
  sets, _ := model.GetSets()
  assert.Equal(t, len(sets), 1)

  model.Eval("reset;")
  sets2, _ := model.GetSets()
  assert.Equal(t, len(sets2), 0)

  model.Eval("set x;")
  sets3, _ := model.GetSets()
  assert.Equal(t, len(sets3), 1)

  model.Eval("set y{1..2};")
  sets4, _ := model.GetSets()
  assert.Equal(t, len(sets4), 2)

  model.Eval(`set z{1..2, {"a", "b"}};`)
  sets5, _ := model.GetSets()
  assert.Equal(t, len(sets5), 3)

  model.Eval("set xx{1..2,4..5};")
  sets6, _ := model.GetSets()
  assert.Equal(t, len(sets6), 4)

  model.Eval("set yy{1..2,4..5,7..8};")
  sets7, _ := model.GetSets()
  assert.Equal(t, len(sets7), 5)

  model.Eval(`set zz{1..2,{"a", "b"},5..6};`)
  sets8, _ := model.GetSets()
  assert.Equal(t, len(sets8), 6)
}

//        self.assertEqual(list(ampl.get_set("x").members()), list(ampl.get_set("x")))
//        self.assertEqual(ampl.get_set("x").num_instances(), 1)
//        self.assertEqual(ampl.get_set("y").num_instances(), 2)
 //       self.assertEqual(list(ampl.set["y"][1].members()), list(ampl.set["y"][1]))
//        self.assertEqual(ampl.get_set("z").num_instances(), 4)
//        self.assertEqual(
//            list(ampl.set["z"][1, "a"].members()), list(ampl.set["z"][1, "a"])
//        )
//        self.assertEqual(ampl.get_set("xx").num_instances(), 4)
//        self.assertEqual(
//            list(ampl.set["xx"][1, 4].members()), list(ampl.set["xx"][1, 4])
//        )
//        self.assertEqual(ampl.get_set("yy").num_instances(), 8)
 //       self.assertEqual(
 //           list(ampl.set["yy"][1, 4, 7].members()), list(ampl.set["yy"][1, 4, 7])
//        )
//        self.assertEqual(ampl.get_set("zz").num_instances(), 8)
//        self.assertEqual(
//            list(ampl.set["zz"][1, "a", 5].members()), list(ampl.set["zz"][1, "a", 5])
//        )