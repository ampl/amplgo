// examples/multidimensionalexample/main.go
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

	if err := model.Eval("set CITIES; set LINKS within (CITIES cross CITIES);"); err != nil {
		panic(err)
	}
	if err := model.Eval("param cost {LINKS} >= 0; param capacity {LINKS} >= 0;"); err != nil {
		panic(err)
	}
	if err := model.Eval("data; set CITIES := PITT NE SE BOS EWR BWI ATL MCO;"); err != nil {
		panic(err)
	}

	cost := []float64{2.5, 3.5, 1.7, 0.7, 1.3, 1.3, 0.8, 0.2, 2.1}
	capacity := []float64{250.0, 250.0, 100.0, 100.0, 100.0, 100.0, 100.0, 100.0, 100.0}
	linksFrom := []string{"PITT", "PITT", "NE", "NE", "NE", "SE", "SE", "SE", "SE"}
	linksTo := []string{"NE", "SE", "BOS", "EWR", "BWI", "EWR", "BWI", "ATL", "MCO"}

	df, err := ampl.NewDataFrame(2, 2, []string{"LINKSFrom", "LINKSTo", "cost", "capacity"})
	if err != nil {
		panic(err)
	}
	if err := df.SetColumnString("LINKSFrom", linksFrom); err != nil {
		panic(err)
	}
	if err := df.SetColumnString("LINKSTo", linksTo); err != nil {
		panic(err)
	}
	if err := df.SetColumnDouble("cost", cost); err != nil {
		panic(err)
	}
	if err := df.SetColumnDouble("capacity", capacity); err != nil {
		panic(err)
	}

	dfString, err := df.ToString()
	if err != nil {
		panic(err)
	}
	fmt.Println(dfString)

	if err := model.SetData(df, "LINKS"); err != nil {
		panic(err)
	}
}
