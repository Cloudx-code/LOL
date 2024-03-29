package admin

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"strings"
	"syscall"
)

func IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func MustRunWithAdmin() {
	if !IsAdmin() {
		RunMeElevated()
	} else {
		fmt.Println("已经处于管理员模式")
	}
}
func RunMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")
	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)
	var showCmd int32 = 1 // SW_NORMAL
	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(-2)
}
