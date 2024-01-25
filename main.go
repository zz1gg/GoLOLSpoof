package main

import (
	"fmt"
	"os/exec"
	"strings"
)

//for Windows

func main() {

	fmt.Println(banner)
	// Start the interactive shell loop
	for {
		// Get and parse command
		fmt.Print("[GoLOLSpoof]> ")
		cmdline := readLineFromStdin()
		cmdline = strings.TrimSpace(cmdline)
		if cmdline == "" {
			continue
		}

		// Handle special command
		if strings.HasPrefix(cmdline, "!") {
			processInput(cmdline)
			continue
		}

		foundPath, err := exec.LookPath(strings.Fields(cmdline)[0])
		if err != nil {
			fmt.Println(err)
		}

		cmdlineSeq := strings.Fields(cmdline)
		cmdlineSeq[0] = foundPath
		cmdline = strings.Join(cmdlineSeq, " ")

		// Fire in the hole!
		if !executeSpoofedLolbin(cmdline) {
			fmt.Printf("[-] Could not spoof binary: %s\n", strings.Fields(cmdline)[0])
		}
	}
}
