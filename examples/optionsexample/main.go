// examples/optionsexample/main.go
package main

import (
	"fmt"

	"github.com/ampl/amplgo/ampl"
)

func main() {
	model, err := ampl.New()
	if err != nil {
		panic(err)
	}

	// Get the value of the option presolve and print
	presolve, err := model.GetIntOption("presolve")
	if err != nil {
		panic(err)
	}
	fmt.Println("AMPL presolve is", presolve)

	// Set the value to false (maps to 0)
	if err := model.SetOption("presolve", false); err != nil {
		panic(err)
	}

	// Get the value of the option presolve and print
	presolve, err = model.GetIntOption("presolve")
	if err != nil {
		panic(err)
	}
	fmt.Println("AMPL presolve is now", presolve)

	// Check whether an option with a specified name exists
	value, err := model.GetOption("solver")
	if err == nil {
		fmt.Println("Option solver exists and has value:", value)
	}

	// Check again, this time failing
	_, err = model.GetOption("s_o_l_v_e_r")
	if err != nil {
		fmt.Println("Option s_o_l_v_e_r does not exist")
	}
}
