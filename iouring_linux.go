package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

const (
	NR_IO_URING_SETUP    uintptr = 425
	NR_IO_URING_ENTER    uintptr = 426
	NR_IO_URING_REGISTER uintptr = 427
)

const (
	NSIG uintptr = 33
)

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
	CqOff        IoSqringOffsets
}

type IoSqringOffsets struct {
	Head        uint32
	Tail        uint32
	RingMask    uint32
	RingEntries uint32
	Flags       uint32
	Dropped     uint32
	Array       uint32
	Resv1       uint32
	Resv2       uint64
}

type Sigset struct {
	Val [16]int64
}

func IoUringSetup(entries int, params *IoUringParams) (err error) {
	_, _, e := syscall.RawSyscall(NR_IO_URING_SETUP,
		uintptr(entries), uintptr(unsafe.Pointer(params)),
		0)
	if e != 0 {
		err = e
	}
	return err
}

func IoUringEnter(fd int, toSubmit int, minComplete uint, flags uint, sig *Sigset) (err error) {
	_, _, e := syscall.Syscall6(NR_IO_URING_ENTER,
		uintptr(fd), uintptr(toSubmit), uintptr(minComplete), uintptr(flags), uintptr(unsafe.Pointer(sig)),
		NSIG/8)
	if e != 0 {
		err = e
	}
	return err
}

func IoUringRegister(fd int, opcode uint, arg uintptr, nrArgs uint) (err error) {
	_, _, e := syscall.RawSyscall6(NR_IO_URING_REGISTER, uintptr(fd), uintptr(opcode), arg, uintptr(nrArgs), 0, 0)
	if e != 0 {
		err = e
	}
	return err
}

func main() {
	var ring IoUringParams
	log.Println(IoUringSetup(1024, &ring))
	fmt.Println(ring)
}
