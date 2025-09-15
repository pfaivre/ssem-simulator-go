package main

import (
	"fmt"
	"os"
	"pfaivre/ssem-simulator-go/ssem"
)

func main() {
	fmt.Println("SSEM Simulator, 2025 Pierre Faivre")
	fmt.Println("(Work in progress, not functional yet)")
	fmt.Println()

	machine := ssem.Ssem{}

	if len(os.Args) > 1 {
		err := machine.ReadAsm(os.Args[1])
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(machine)
}
