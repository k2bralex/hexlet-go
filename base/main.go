package main

import (
	. "base/lessons/model"
	"fmt"
)

func main() {
	fmt.Println(CalcArea(&Rectangle{}))
	fmt.Println(CalcArea(&Circle{R: 3}))
}
