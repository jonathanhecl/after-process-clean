package process

/*

Tested:

* Windows 11 x64 - OK
* Windows 10 x64 - OK
* Windows 7 x64	- OK
* Windows 7 x32 - OK


*/

import (
	"os"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

type ProcessStruct struct {
	PID      int
	Filename string
	Path     string
}

type drivesStruct struct {
	drive  string
	device string
}

var drives []drivesStruct
var process []ProcessStruct

const (
	MAX_PATH = 260

	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
	PROCESS_QUERY_INFORMATION         = 0x0400
	PROCESS_VM_READ                   = 0x0010
)

var (
	modKernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procCloseHandle              = modKernel32.NewProc("CloseHandle")
	procCreateToolhelp32Snapshot = modKernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modKernel32.NewProc("Process32FirstW")
	procProcess32Next            = modKernel32.NewProc("Process32NextW")
	procOpenProcess              = modKernel32.NewProc("OpenProcess")
	procQueryDosDeviceW          = modKernel32.NewProc("QueryDosDeviceW")
	procSuspendThread            = modKernel32.NewProc("SuspendThread")
	procGetModuleFileNameExW     = modKernel32.NewProc("GetModuleFileNameExW")

	modPsapi                     = syscall.NewLazyDLL("psapi.dll")
	procGetProcessImageFileNameW = modPsapi.NewProc("GetProcessImageFileNameW")
)

type PROCESSENTRY32 struct {
	Size              uint32
	CntUsage          uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	CntThreads        uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [MAX_PATH]uint16
}

func listDevices() {
	drives = nil

	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		_, err := os.Open(string(drive) + ":\\")
		if err == nil {
			realDevice := make([]uint16, MAX_PATH)
			r1, _, _ := procQueryDosDeviceW.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(string(drive)+":"))), uintptr(unsafe.Pointer(&realDevice[0])), uintptr(MAX_PATH))
			var nD uint32 = uint32(r1)
			if nD != 0 {
				realDeviceStr := syscall.UTF16ToString(realDevice[0:nD])
				drives = append(drives, drivesStruct{drive: string(drive) + ":", device: realDeviceStr})
			}
		}
	}

	return
}

func List() (list []ProcessStruct) {
	process = nil

	handle, _, _ := procCreateToolhelp32Snapshot.Call(0x00000002, 0)
	if handle < 0 {
		return
	}
	defer procCloseHandle.Call(handle)

	var entry PROCESSENTRY32
	entry.Size = uint32(unsafe.Sizeof(entry))
	ret, _, _ := procProcess32First.Call(handle, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		return
	}
	for {
		process = append(process, getInfoProcess(&entry))
		ret, _, _ := procProcess32Next.Call(handle, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
	}

	return process
}

func getInfoProcess(e *PROCESSENTRY32) ProcessStruct {
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	filename := syscall.UTF16ToString(e.ExeFile[:end])

	var sincpath chan string = make(chan string)
	go getProcessPath(e.ProcessID, sincpath, filename)

	return ProcessStruct{
		PID:      int(e.ProcessID),
		Filename: filename,
		Path:     <-sincpath,
	}
}

func getProcessPath(pid uint32, sincpath chan string, filename string) {
	if pid == 0 || filename == "System" || filename == "Secure System" || filename == "Registry" {
		sincpath <- ""
	}

	handle, _, _ := procOpenProcess.Call(PROCESS_QUERY_INFORMATION+PROCESS_VM_READ+PROCESS_QUERY_LIMITED_INFORMATION, uintptr(0), uintptr(pid))
	defer procCloseHandle.Call(handle)
	if handle == 0 {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE,
			`SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\`+filename,
			registry.READ)
		if err != nil {
			k, err = registry.OpenKey(registry.CURRENT_USER,
				`SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\`+filename,
				registry.READ)
			if err != nil {
				sincpath <- ""
			}
		}
		defer k.Close()

		val, _, err1 := k.GetStringValue("")
		if err1 != nil {
			sincpath <- "nada"
		}
		sincpath <- val
	}

	path := make([]uint16, MAX_PATH)
	r0, _, _ := procGetProcessImageFileNameW.Call(uintptr(handle), uintptr(unsafe.Pointer(&path[0])), uintptr(MAX_PATH)) // Return like \Device\HarddiskVolume2\...
	var n uint32 = uint32(r0)
	if n == 0 {
		sincpath <- ""
	}

	pathStr := syscall.UTF16ToString(path[0:n])

	if drives == nil {
		listDevices()
	}

	if strings.Index(pathStr, "\\Device\\") == 0 {
		for _, drive := range drives {
			if pathStr[0:len(drive.device)] == drive.device {
				pathStr = drive.drive + pathStr[len(drive.device):len(pathStr)]
				break
			}
		}
	}

	sincpath <- pathStr
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
