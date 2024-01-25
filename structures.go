package main

import "golang.org/x/sys/windows"

const banner = `

 ██████   ██████  ██       ██████  ██      ███████ ██████   ██████   ██████  ███████ 
██       ██    ██ ██      ██    ██ ██      ██      ██   ██ ██    ██ ██    ██ ██      
██   ███ ██    ██ ██      ██    ██ ██      ███████ ██████  ██    ██ ██    ██ █████   
██    ██ ██    ██ ██      ██    ██ ██           ██ ██      ██    ██ ██    ██ ██      
 ██████   ██████  ███████  ██████  ███████ ███████ ██       ██████   ██████  ██

An interactive shell to spoof some LOLBins
try !help
`
const helpcmd = `
!exit    -> Exit
!cls     -> Clear the screen
!help    -> This help message
`

var (
	ntdll               = windows.NewLazySystemDLL("ntdll.dll")
	ntQueryInfoProc     = ntdll.NewProc("NtQueryInformationProcess")
	kernel32            = windows.NewLazySystemDLL("kernel32.dll")
	readProcessMemory   = kernel32.NewProc("ReadProcessMemory")
	writeProcessMemory  = kernel32.NewProc("WriteProcessMemory")
	ntResumeThread      = ntdll.NewProc("NtResumeThread")
	waitForSingleObject = kernel32.NewProc("WaitForSingleObject")
)

const (
	PROCESS_BASIC_INFORMATION_CLASS = 0
	infinite                        = 0xFFFFFFFF
)

type PROCESS_BASIC_INFORMATION struct {
	Reserved1       uintptr
	PebBaseAddress  uintptr
	Reserved2       [2]uintptr
	UniqueProcessId uintptr
	Reserved3       uintptr
}
