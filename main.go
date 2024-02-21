package main

import (
	"fmt"
	"gokula/vm"
	"os"
)

type startupInfo struct {
	method string
	path   string
}

func main() {
	startupInfo, err := readArgs()
	if err != nil {
		info()
	} else {
		var err error
		vm.CompiledFileInstance, err = vm.Load(startupInfo.path)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
		if startupInfo.method == "show" {
			fmt.Println(vm.CompiledFileInstance)
		} else if startupInfo.method == "run" {
			err = vm.CompiledFileInstance.Run()
			if err != nil {
				fmt.Println("error: ", err)
			}
		}
	}
}

func readArgs() (startupInfo, error) {
	si := startupInfo{}
	err := fmt.Errorf("illegal args")
	if len(os.Args) < 2 {
		return si, err
	}
	for index, str := range os.Args {
		if index == 0 {
			continue
		} else if index == 1 {
			if str == "-r" || str == "--run" {
				si.method = "run"
			} else if str == "-s" || str == "--show" {
				si.method = "show"
			} else {
				return si, err
			}
		} else if index == 2 {
			si.path = str
		}
	}
	return si, nil
}

func info() {
	str := `Usage:	gokula <command> <*.kulac> [<args>]

	-r, --run	Run a kula-compiled-file in release mode
	-s, --show	Output a kula-compiled-file in bytecode format`
	fmt.Println(str)
}
