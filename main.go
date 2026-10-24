package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"pfaivre/ssem-simulator-go/ssem"
	"syscall"
)

func main() {
	fmt.Println("SSEM Simulator, 2026 Pierre Faivre")
	fmt.Println()

	filePath := flag.String("file", "", "Path to an asm or snp file to load")
	pretty := flag.Bool("pretty", true, "Improve readability when printing the store")
	cyclesPerSec := flag.Uint("speed", 700, "Approximate target of number of instructions to execute per second. Set 0 for no limit.")
	maxCycles := flag.Uint("max-cycles", 10000, "Maximum number of cycles to execute")
	flag.Parse()

	machine := ssem.NewSsem()

	if *pretty {
		ssem.BinaryDigits = [2]byte{'.', '#'}
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

	cyclesChan := make(chan uint)
	stopChan := make(chan bool)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go machine.Run(cyclesChan, *maxCycles, *cyclesPerSec)
	go machine.Printer(stopChan)

	var cycles uint
	select {
	case cycles = <-cyclesChan:
		// Normal completion: Run finished
		machine.Stop()
		<-stopChan
		fmt.Printf("Stopped after %d cycles\n", cycles)
	case s := <-signalChan:
		// Signal received: stop the machine
		fmt.Printf("Received %s signal, terminating\n", s)
		machine.Stop()
		cycles = <-cyclesChan
		<-stopChan
		fmt.Printf("Stopped after %d cycles\n", cycles)
		os.Exit(0)
	}
}
