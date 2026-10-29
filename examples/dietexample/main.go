// examples/dietexample/main.go
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

	modelDir := "models"
	if len(args) == 3 {
		modelDir = args[2]
	}

	if err := model.Read(fmt.Sprintf("%s/diet/diet.mod", modelDir)); err != nil {
		panic(err)
	}

	// Food set and its parameters — all columns declared up front
	foods := []string{"BEEF", "CHK", "FISH", "HAM", "MCH", "MTL", "SPG", "TUR"}
	costs := []float64{3.59, 2.59, 2.29, 2.89, 1.89, 1.99, 1.99, 2.49}
	fmin := []float64{2.0, 2.0, 2.0, 2.0, 2.0, 2.0, 2.0, 2.0}
	fmax := []float64{10.0, 10.0, 10.0, 10.0, 10.0, 10.0, 10.0, 10.0}

	dfFood, err := ampl.NewDataFrame(1, 3, []string{"FOOD", "cost", "f_min", "f_max"})
	if err != nil {
		panic(err)
	}
	if err := dfFood.SetColumnString("FOOD", foods); err != nil {
		panic(err)
	}
	if err := dfFood.SetColumnDouble("cost", costs); err != nil {
		panic(err)
	}
	if err := dfFood.SetColumnDouble("f_min", fmin); err != nil {
		panic(err)
	}
	if err := dfFood.SetColumnDouble("f_max", fmax); err != nil {
		panic(err)
	}
	if err := model.SetData(dfFood, "FOOD"); err != nil {
		panic(err)
	}

	// Nutrient set and its parameters — all columns declared up front
	nutrients := []string{"A", "C", "B1", "B2", "NA", "CAL"}
	nmin := []float64{700.0, 700.0, 700.0, 700.0, 0.0, 16000.0}
	nmax := []float64{20000.0, 20000.0, 20000.0, 20000.0, 50000.0, 24000.0}

	dfNutr, err := ampl.NewDataFrame(1, 2, []string{"NUTR", "n_min", "n_max"})
	if err != nil {
		panic(err)
	}
	if err := dfNutr.SetColumnString("NUTR", nutrients); err != nil {
		panic(err)
	}
	if err := dfNutr.SetColumnDouble("n_min", nmin); err != nil {
		panic(err)
	}
	if err := dfNutr.SetColumnDouble("n_max", nmax); err != nil {
		panic(err)
	}
	if err := model.SetData(dfNutr, "NUTR"); err != nil {
		panic(err)
	}

	// Amount matrix: rows = nutrients (6), cols = foods (8), row-major
	amounts := []float64{
		60.0, 8.0, 8.0, 40.0, 15.0, 70.0, 25.0, 60.0, // A
		20.0, 0.0, 10.0, 40.0, 35.0, 30.0, 50.0, 20.0, // C
		10.0, 20.0, 15.0, 35.0, 15.0, 15.0, 25.0, 15.0, // B1
		15.0, 20.0, 10.0, 10.0, 15.0, 15.0, 15.0, 10.0, // B2
		928.0, 2180.0, 945.0, 278.0, 1182.0, 896.0, 1329.0, 1397.0, // NA
		295.0, 770.0, 440.0, 430.0, 315.0, 400.0, 379.0, 450.0, // CAL
	}

	dfAmt, err := ampl.NewDataFrame(2, 1, []string{"NUTR", "FOOD", "amt"})
	if err != nil {
		panic(err)
	}
	if err := dfAmt.SetMatrixDoubles(nutrients, foods, amounts); err != nil {
		panic(err)
	}
	if err := model.SetData(dfAmt, ""); err != nil {
		panic(err)
	}

	if err := model.Solve("", ""); err != nil {
		panic(err)
	}

	objective, err := model.GetObjective("Total_Cost")
	if err != nil {
		panic(err)
	}
	totalCost, err := objective.Value()
	if err != nil {
		panic(err)
	}
	fmt.Println("Objective:", totalCost)
}
