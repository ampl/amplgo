// examples/trackingmodel/main.go
package main

import (
	"fmt"
	"os"

	"github.com/ampl/amplgo/ampl"
)

func main() {
	args := os.Args

	model, err := ampl.New()
	if err != nil {
		panic(err)
	}

	if len(args) > 1 {
		if err := model.SetOption("solver", args[1]); err != nil {
			panic(err)
		}
	}

	modelDir := "models/tracking"
	if len(args) == 3 {
		modelDir = args[2]
	}

	if err := model.Read(fmt.Sprintf("%s/tracking.mod", modelDir)); err != nil {
		panic(err)
	}
	if err := model.ReadData(fmt.Sprintf("%s/tracking.dat", modelDir)); err != nil {
		panic(err)
	}
	if err := model.Read(fmt.Sprintf("%s/trackingbit.run", modelDir)); err != nil {
		panic(err)
	}

	// Set the scalar string parameter used by the table declarations
	if err := model.Eval(fmt.Sprintf("let data_dir := '%s';", modelDir)); err != nil {
		panic(err)
	}

	if err := model.ReadTable("assets"); err != nil {
		panic(err)
	}
	if err := model.ReadTable("indret"); err != nil {
		panic(err)
	}
	if err := model.ReadTable("returns"); err != nil {
		panic(err)
	}

	hold, err := model.GetVariable("hold")
	if err != nil {
		panic(err)
	}
	ifinuniverse, err := model.GetParameter("ifinuniverse")
	if err != nil {
		panic(err)
	}

	// Relax integrality and solve the QP relaxation
	if err := model.SetOption("relax_integrality", true); err != nil {
		panic(err)
	}
	if err := model.Solve("", ""); err != nil {
		panic(err)
	}
	objectives, err := model.GetObjectives()
	if err != nil {
		panic(err)
	}
	qpObjective, err := objectives[0].Value()
	if err != nil {
		panic(err)
	}
	fmt.Println("QP objective value", qpObjective)

	const lowCutoff = 0.04
	const highCutoff = 0.1

	// For each asset: force in (2), leave free (1), or force out (0)
	holdDf, err := hold.GetValues(nil)
	if err != nil {
		panic(err)
	}
	colIndex, err := holdDf.GetColumnIndex("hold.val")
	if err != nil {
		panic(err)
	}
	numRows, err := holdDf.GetNumRows()
	if err != nil {
		panic(err)
	}

	toHold := make([]float64, numRows)
	for i := 0; i < numRows; i++ {
		variant, err := holdDf.Element(i, colIndex)
		if err != nil {
			panic(err)
		}
		x, err := variant.GetNumericValue()
		if err != nil {
			panic(err)
		}
		switch {
		case x < lowCutoff:
			toHold[i] = 0.0
		case x > highCutoff:
			toHold[i] = 2.0
		default:
			toHold[i] = 1.0
		}
	}

	if err := ifinuniverse.SetAllDoubleValues(toHold); err != nil {
		panic(err)
	}

	// Solve the integer problem
	if err := model.SetOption("relax_integrality", false); err != nil {
		panic(err)
	}
	if err := model.Solve("", ""); err != nil {
		panic(err)
	}
	objectives, err = model.GetObjectives()
	if err != nil {
		panic(err)
	}
	qmipObjective, err := objectives[0].Value()
	if err != nil {
		panic(err)
	}
	fmt.Println("QMIP objective value", qmipObjective)
}
