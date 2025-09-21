package main

import (
	"fmt"
	"os"
	"path/filepath"
	"pfaivre/ssem-simulator-go/ssem"
)

func main() {
	fmt.Println("SSEM Simulator, 2025 Pierre Faivre")
	fmt.Println("(Work in progress, not functional yet)")
	fmt.Println()

	machine := ssem.Ssem{}

	if len(os.Args) > 1 {
		path := os.Args[1]

		switch filepath.Ext(path) {
		case ".asm":
			err := machine.ReadAsm(path)
			if err != nil {
				panic(err)
			}
		case ".snp":
			err := machine.ReadSnp(path)
			if err != nil {
				panic(err)
			}
		default:
			panic(fmt.Errorf("unknown file format"))
		}
	}

	machine.Run(100000)

	fmt.Println(machine)
}
