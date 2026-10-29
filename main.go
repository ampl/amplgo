// main.go in root directory
package main

import (
    "fmt"
    "github.com/ampl/ampl.go/ampl"
)

func main() {
  model, err := ampl.New()
  if err != nil {
    panic(err)
  }

  if err := model.Read("./models/diet/diet.mod"); err != nil {
    panic(err)
  }
  if err := model.ReadData("./models/diet/diet.dat"); err != nil {
    panic(err)
  }

  if err := model.Solve("", "gurobi"); err != nil {
    panic(err)
  }

  variant, err := model.GetValue("Total_Cost")
  if err != nil {
    panic(err)
  }

  totalCost, err := variant.GetNumericValue()
  if err != nil {
    panic(err)
  }
  fmt.Println("Total Cost:", totalCost)

  vars, err := model.GetVariables()
  if err != nil {
    panic(err)
  }
  for _, variable := range vars {
    fmt.Printf("Variables: %s\n", variable.GetName())
  }

  params, err := model.GetParameters()
  if err != nil {
    panic(err)
  }
  for _, param := range params {
    indexarity, _ := param.GetIndexarity()
    fmt.Printf("Parameters: %s, %d\n", param.GetName(), indexarity)
  }

  //model.Free()

  fmt.Println("AMPL operations completed.")
}
