// examples/writemodel/main.go
package main

import (
	"fmt"

	"github.com/ampl/amplgo/ampl"
)

func makeAmplModel(model *ampl.Ampl, numvars int, unfeasible bool, unbounded bool) {
	if unfeasible && unbounded {
		panic("unfeasible and unbounded cannot both be true")
	}

	if err := model.Eval(
		"set varIndexes within Integers; set constraintIndexes within Integers; " +
			"param varCosts{ varIndexes }; " +
			"param rightHandSides{ constraintIndexes }; " +
			"var x{ i in varIndexes } >= 0; " +
			"maximize objectiveFunction : sum{ i in varIndexes } varCosts[i] * x[i];"); err != nil {
		panic(err)
	}

	indexes := make([]float64, numvars)
	costs := make([]float64, numvars)
	for i := 0; i < numvars; i++ {
		indexes[i] = float64(i)
		costs[i] = float64(i)
	}
	d, err := ampl.NewDataFrame(1, 1, []string{"varIndexes", "varCosts"})
	if err != nil {
		panic(err)
	}
	if err := d.SetColumnDouble("varIndexes", indexes); err != nil {
		panic(err)
	}
	if err := d.SetColumnDouble("varCosts", costs); err != nil {
		panic(err)
	}
	if err := model.SetData(d, "varIndexes"); err != nil {
		panic(err)
	}

	if !unbounded {
		if err := model.Eval("mainConstraints{j in constraintIndexes}: x[j] + x[j+1] <= rightHandSides[j];"); err != nil {
			panic(err)
		}

		ci := make([]float64, numvars-1)
		rhs := make([]float64, numvars-1)
		for i := 0; i < numvars-1; i++ {
			ci[i] = float64(i)
			if unfeasible {
				rhs[i] = -1.0
			} else {
				rhs[i] = 1.0
			}
		}
		d2, err := ampl.NewDataFrame(1, 1, []string{"constraintIndexes", "rightHandSides"})
		if err != nil {
			panic(err)
		}
		if err := d2.SetColumnDouble("constraintIndexes", ci); err != nil {
			panic(err)
		}
		if err := d2.SetColumnDouble("rightHandSides", rhs); err != nil {
			panic(err)
		}
		if err := model.SetData(d2, "constraintIndexes"); err != nil {
			panic(err)
		}
	}
}

// tryWrite attempts to write the model; it prints and returns true if AMPL
// raises an error.
func tryWrite(model *ampl.Ampl, filename string, auxfiles string) bool {
	if err := model.Write(filename, auxfiles); err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func main() {
	model, err := ampl.New()
	if err != nil {
		panic(err)
	}

	makeAmplModel(model, 10, true, false)

	// Infeasible model — presolver detects infeasibility during write
	if !tryWrite(model, "gtest", "") {
		panic("expected an error for infeasible model")
	}

	// Wrong file type — AMPL rejects the format
	if !tryWrite(model, "htest", "") {
		panic("expected an error for wrong file type")
	}

	// Disable presolve so the infeasible model is written as-is — this works
	if err := model.SetOption("presolve", 0); err != nil {
		panic(err)
	}
	if err := model.Write("gtest", "cr"); err != nil {
		panic(err)
	}

	// Trivial model — presolver raises an error
	if err := model.Reset(); err != nil {
		panic(err)
	}
	if err := model.SetOption("presolve", 10); err != nil {
		panic(err)
	}
	if err := model.Eval("var x; maximize z : x; c: x = 5;"); err != nil {
		panic(err)
	}
	if !tryWrite(model, "gtest", "") {
		panic("expected an error for trivial model")
	}
}
