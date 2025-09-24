package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"pfaivre/ssem-simulator-go/ssem"
	"strings"
)

func main() {
	fmt.Println("SSEM Simulator, 2025 Pierre Faivre")
	fmt.Println()

	filePath := flag.String("file", "", "Path to an asm or snp file to load")
	pretty := flag.Bool("pretty", false, "Improve readability when printing the store")
	printFlag := flag.Bool("print", false, "Print the store at every cycle with a small pause")
	maxCycles := flag.Uint("max-cycles", 1000, "Maximum number of cycles to execute")
	flag.Parse()

	machine := ssem.NewSsem()

	if *pretty {
		ssem.BinaryDigitReplacer = strings.NewReplacer("0", ".", "1", "#")
	}

	if len(*filePath) > 0 {
		switch filepath.Ext(*filePath) {
		case ".asm":
			err := machine.ReadAsm(*filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		case ".snp":
			err := machine.ReadSnp(*filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		default:
			fmt.Fprint(os.Stderr, "unknown file format\n")
			os.Exit(1)
		}
	}

	fmt.Println(machine)

	var cycles uint
	var err error

	if *printFlag {
		cycles, err = machine.RunAndPrint(*maxCycles)
	} else {
		fmt.Printf("Computing up to %d cycles...\n\n", *maxCycles)
		cycles, err = machine.Run(*maxCycles)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Println(machine)
	fmt.Printf("Stopped after %d cycles\n", cycles)
}
