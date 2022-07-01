package pkg

import (
	"syscall"
	"unsafe"
)

type Sigset struct {
	Val [16]int64
}

type IoUringParams struct {
	SqEntries    uint32
	CqEntries    uint32
	Flags        uint32
	SqThreadCpu  uint32
	SqThreadIdle uint32
	Features     uint32
	WqFd         uint32
	Resv         [3]uint32
	SqOff        IoSqringOffsets
	CqOff        IoCqringOffsets
}

func IoUringSetup(entries int, params *IoUringParams) (fd int, err error) {
	res, _, e := syscall.Syscall(_NR_IO_URING_SETUP,
		uintptr(entries), uintptr(unsafe.Pointer(params)),
		0)
	if e != 0 {
		err = e
	}
	return int(res), err
}

func IoUringEnter(fd int, toSubmit int, minComplete uint, flags uint, sig *Sigset) (err error) {

	_, _, e := syscall.Syscall6(_NR_IO_URING_ENTER,
		uintptr(fd), uintptr(toSubmit), uintptr(minComplete), uintptr(flags), uintptr(unsafe.Pointer(sig)),
		NSIG/8)
	if e != 0 {
		err = e
	}
	return err
}

func IoUringRegister(fd int, opcode uint, arg uintptr, nrArgs uint) (err error) {
	_, _, e := syscall.Syscall6(_NR_IO_URING_REGISTER, uintptr(fd), uintptr(opcode), arg, uintptr(nrArgs), 0, 0)
	if e != 0 {
		err = e
	}
	return err
}
