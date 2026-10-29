// examples/firstexample/main.go
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
	if err := model.ReadData(fmt.Sprintf("%s/diet/diet.dat", modelDir)); err != nil {
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
	fmt.Println("Objective is:", totalCost)

	// Reassign data - specific instances
	cost, err := model.GetParameter("cost")
	if err != nil {
		panic(err)
	}
	if err := cost.SetSomeDoubleValues([]string{"BEEF", "HAM"}, []float64{5.01, 4.55}); err != nil {
		panic(err)
	}
	fmt.Println("Increased costs of beef and ham.")

	if err := model.Solve("", ""); err != nil {
		panic(err)
	}
	totalCost, err = objective.Value()
	if err != nil {
		panic(err)
	}
	fmt.Println("New objective value:", totalCost)

	// Reassign data - all instances
	if err := cost.SetAllDoubleValues([]float64{3.0, 5.0, 5.0, 6.0, 1.0, 2.0, 5.01, 4.55}); err != nil {
		panic(err)
	}
	fmt.Println("Updated all costs.")

	if err := model.Solve("", ""); err != nil {
		panic(err)
	}
	totalCost, err = objective.Value()
	if err != nil {
		panic(err)
	}
	fmt.Println("New objective value:", totalCost)

	// Get the values of the variable Buy in a DataFrame
	buy, err := model.GetVariable("Buy")
	if err != nil {
		panic(err)
	}
	df, err := buy.GetValues(nil)
	if err != nil {
		panic(err)
	}
	dfString, err := df.ToString()
	if err != nil {
		panic(err)
	}
	fmt.Println(dfString)

	// Get the values of an expression into a DataFrame
	df2, err := model.GetData([]string{"{j in FOOD} 100*Buy[j]/Buy[j].ub"})
	if err != nil {
		panic(err)
	}
	df2String, err := df2.ToString()
	if err != nil {
		panic(err)
	}
	fmt.Println(df2String)
}
