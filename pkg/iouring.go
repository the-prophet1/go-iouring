package pkg

import (
	"errors"
	"reflect"
	"sync/atomic"
	"syscall"
	"unsafe"
)

const (
	IORING_OP_NOP = iota
	IORING_OP_READV
	IORING_OP_WRITEV
	IORING_OP_FSYNC
	IORING_OP_READ_FIXED
	IORING_OP_WRITE_FIXED
	IORING_OP_POLL_ADD
	IORING_OP_POLL_REMOVE
	IORING_OP_SYNC_FILE_RANGE
	IORING_OP_SENDMSG
	IORING_OP_RECVMSG
	IORING_OP_TIMEOUT
	IORING_OP_TIMEOUT_REMOVE
	IORING_OP_ACCEPT
	IORING_OP_ASYNC_CANCEL
	IORING_OP_LINK_TIMEOUT
	IORING_OP_CONNECT
	IORING_OP_FALLOCATE
	IORING_OP_OPENAT
	IORING_OP_CLOSE
	IORING_OP_FILES_UPDATE
	IORING_OP_STATX
	IORING_OP_READ
	IORING_OP_WRITE
	IORING_OP_FADVISE
	IORING_OP_MADVISE
	IORING_OP_SEND
	IORING_OP_RECV
	IORING_OP_OPENAT2
	IORING_OP_EPOLL_CTL
	IORING_OP_SPLICE
	IORING_OP_PROVIDE_BUFFERS
	IORING_OP_REMOVE_BUFFERS
	IORING_OP_TEE

	/* this goes last, obviously */
	IORING_OP_LAST
)

const (
	_NR_IO_URING_SETUP    uintptr = 425
	_NR_IO_URING_ENTER    uintptr = 426
	_NR_IO_URING_REGISTER uintptr = 427
)

const (
	IORING_SETUP_IOPOLL    = 1 << iota /* io_context is polled */
	IORING_SETUP_SQPOLL    = 1 << iota /* SQ poll thread */
	IORING_SETUP_SQ_AFF    = 1 << iota /* sq_thread_cpu is valid */
	IORING_SETUP_CQSIZE    = 1 << iota /* app defines CQ size */
	IORING_SETUP_CLAMP     = 1 << iota /* clamp SQ/CQ ring sizes */
	IORING_SETUP_ATTACH_WQ = 1 << iota /* attach to existing wq */

)

const (
	//IORING_FEATURE_SINGLE_MMAP use single mapping memory
	//cq and sq shared a mapping memory
	IORING_FEATURE_SINGLE_MMAP     = 1 << iota
	IORING_FEATURE_NODROP          = 1 << iota
	IORING_FEATURE_SUBMIT_STABLE   = 1 << iota
	IORING_FEATURE_RW_CUR_POS      = 1 << iota
	IORING_FEATURE_CUR_PERSONALITY = 1 << iota
	IORING_FEATURE_FAST_POLL       = 1 << iota
	IORING_FEATURE_POLL_32BITS     = 1 << iota
)

const (
	NSIG uintptr = 33
)

const (
	IORING_OFFSET_SQ_RING = 0
	IORING_OFFSET_CQ_RING = 0x8000000
	IORING_OFFSET_SQES    = 0x10000000
)

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

type IoCqringOffsets struct {
	Head        uint32
	Tail        uint32
	RingMask    uint32
	RingEntries uint32
	Overflow    uint32
	Cqes        uint32
	Flags       uint32
	Resv1       uint32
	Resv2       uint64
}

//IoUringSqe is iouring submission queue entry
type IoUringSqe struct {
	Opcode            uint8
	Flags             uint8
	IoPrIo            uint16
	Fd                int32
	OffOrAddr2        Union1
	AddrOrSpliceOffIn Union2
	Len               uint32
	FlagsOrEvents     Union3
	UserData          uint64
	UnionStruct
}

func (sqe *IoUringSqe) PrepReadv(fd int, iovecs []syscall.Iovec, offset uint64) {
	sqe.prepRW(IORING_OP_READV, int32(fd), reflect.ValueOf(iovecs).Pointer(), uint32(len(iovecs)), offset)
}

func (sqe *IoUringSqe) SetUserData(data interface{}) {
	dataPtr := unsafe.Pointer(&data)
	sqe.UserData = uint64(uintptr(dataPtr))
}

func (sqe *IoUringSqe) GetUserData() interface{} {
	dataPtr := unsafe.Pointer(uintptr(sqe.UserData))
	return *(*interface{})(dataPtr)
}

func (sqe *IoUringSqe) prepRW(op uint8, fd int32, addr uintptr, length uint32, offset uint64) {
	sqe.Opcode = op
	sqe.Flags = 0
	sqe.IoPrIo = 0
	sqe.Fd = fd
	sqe.OffOrAddr2.SetOffset(offset)
	sqe.AddrOrSpliceOffIn.SetAddr(uint64(addr))
	sqe.Len = length
	sqe.FlagsOrEvents.SetRWFlag(0)
	sqe.UserData = 0
	sqe.pad[0] = 0
	sqe.pad[1] = 0
	sqe.pad[2] = 0
}

//IoUringSq is iouring submission queue
type IoUringSq struct {
	KHead        *uint32
	KTail        *uint32
	KRingMask    *uint32
	KRingEntries *uint32
	KFlags       *uint32
	KDropped     *uint32
	Array        []uint32
	Sqes         []IoUringSqe
	SqeHead      uint32
	SqeTail      uint32
	RingSize     int
	Ring         []byte
}

//IoUringCqe is iouring completion queue entry
type IoUringCqe struct {
	UserData uint64
	Res      int32
	Flags    uint32
}

//IoUringCq is iouring completion queue
type IoUringCq struct {
	KHead        *uint32
	KTail        *uint32
	KRingMask    *uint32
	KRingEntries *uint32
	KFlags       *uint32
	KOverflow    *uint32
	Cqes         []IoUringCqe
	RingSize     int
	Ring         []byte
}

type IoUring struct {
	Sq     IoUringSq
	Cq     IoUringCq
	Flags  uint32
	RingFd int
}

//set shared memory info to uring
func (ring *IoUring) ioUringMmap(fd int, params *IoUringParams) error {
	var (
		size int
		err  error
		sq   = &ring.Sq
		cq   = &ring.Cq
	)

	sq.RingSize = int(params.SqOff.Array) + int(params.SqEntries)*int(unsafe.Sizeof(uint32(0)))
	cq.RingSize = int(params.CqOff.Cqes) + int(params.CqEntries)*int(unsafe.Sizeof(IoUringCqe{}))

	if params.Features&IORING_FEATURE_SINGLE_MMAP == 1 {
		if cq.RingSize > sq.RingSize {
			sq.RingSize = cq.RingSize
		}
		cq.RingSize = sq.RingSize
	}

	//get the memory mapping area of sq
	sq.Ring, err = syscall.Mmap(
		fd,
		IORING_OFFSET_SQ_RING,
		sq.RingSize,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED|syscall.MAP_POPULATE)
	if err != nil {
		return err
	}

	if params.Features&IORING_FEATURE_SINGLE_MMAP == 1 {
		cq.Ring = sq.Ring
	} else {
		//get the memory mapping area of cq
		cq.Ring, err = syscall.Mmap(fd,
			IORING_OFFSET_CQ_RING,
			cq.RingSize,
			syscall.PROT_READ|syscall.PROT_WRITE,
			syscall.MAP_SHARED|syscall.MAP_POPULATE,
		)
		if err != nil {
			ioUringUnmapRings(sq, cq)
			return err
		}
	}
	sq.KHead = (*uint32)(unsafe.Add(SliceBytePoint(sq.Ring), params.SqOff.Head))
	sq.KTail = (*uint32)(unsafe.Add(SliceBytePoint(sq.Ring), params.SqOff.Tail))
	sq.KRingMask = (*uint32)(unsafe.Add(SliceBytePoint(sq.Ring), params.SqOff.RingMask))
	sq.KRingEntries = (*uint32)(unsafe.Add(SliceBytePoint(sq.Ring), params.SqOff.RingEntries))
	sq.KFlags = (*uint32)(unsafe.Add(SliceBytePoint(sq.Ring), params.SqOff.Flags))
	sq.KDropped = (*uint32)(unsafe.Add(SliceBytePoint(sq.Ring), params.SqOff.Dropped))
	sq.Array = PtrToUint32Slice(
		SliceByteAddr(sq.Ring)+uintptr(params.SqOff.Array),
		int(*sq.KRingEntries),
		int(*sq.KRingEntries))

	size = int(params.SqEntries) * int(unsafe.Sizeof(IoUringSqe{}))

	//get the memory mapping area of sqes
	byteSqes, err := syscall.Mmap(fd,
		IORING_OFFSET_SQES,
		size,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED|syscall.MAP_POPULATE,
	)
	if err != nil {
		ioUringUnmapRings(sq, cq)
		return err
	}
	sq.Sqes = ByteSliceToSqes(byteSqes)

	cq.KHead = (*uint32)(unsafe.Add(SliceBytePoint(cq.Ring), params.CqOff.Head))
	cq.KTail = (*uint32)(unsafe.Add(SliceBytePoint(cq.Ring), params.CqOff.Tail))
	cq.KRingMask = (*uint32)(unsafe.Add(SliceBytePoint(cq.Ring), params.CqOff.RingMask))
	cq.KRingEntries = (*uint32)(unsafe.Add(SliceBytePoint(cq.Ring), params.CqOff.RingEntries))
	cq.KOverflow = (*uint32)(unsafe.Add(SliceBytePoint(cq.Ring), params.CqOff.Overflow))
	cq.Cqes = PtrToCqes(SliceByteAddr(cq.Ring)+uintptr(params.CqOff.Cqes), int(*cq.KRingEntries), int(*cq.KRingEntries))
	if params.CqOff.Flags != 0 {
		cq.KFlags = (*uint32)(unsafe.Add(SliceBytePoint(cq.Ring), params.CqOff.Flags))
	}

	return nil
}

//GetSqe get a sqe from the sqe ring queue in iouring
func (ring *IoUring) GetSqe() (*IoUringSqe, error) {
	sq := &ring.Sq
	head := atomic.LoadUint32(sq.KHead)

	next := sq.SqeTail + 1
	if next-head <= *sq.KRingEntries {
		sqe := &sq.Sqes[sq.SqeTail&*sq.KRingMask]
		sq.SqeTail = next
		return sqe, nil
	}
	return nil, errors.New("sqe ring is full")
}

func (ring *IoUring) Submit() error {
	return ring.submitAndWait(0)
}

func (ring *IoUring) submitAndWait(waitNumber int) error {
	return ring.submit(waitNumber)
}

func (ring *IoUring) submit(waitNumber int) error {

}

func (ring *IoUring) flushSQ() uint32 {
	var (
		sq       = &ring.Sq
		mask     = *sq.KRingMask
		ktail    uint32
		toSubmit uint32
	)

	if sq.SqeHead == sq.SqeTail {
		ktail = *sq.KTail
		return ktail - *sq.KHead
	}

	/*
		Fill in sqes that we have queued up, adding them to the kernel ring
			ktail = *sq->ktail;
			to_submit = sq->sqe_tail - sq->sqe_head;
			while (to_submit--) {
				sq->array[ktail & mask] = sq->sqe_head & mask;
				ktail++;
				sq->sqe_head++;
			}
	*/
	ktail = *sq.KTail
	toSubmit = sq.SqeTail - sq.SqeHead
	for ; toSubmit != 0; toSubmit-- {
		sq.Array[ktail&mask] = sq.SqeHead & mask
		ktail++
		sq.SqeHead++
	}
	atomic.StoreUint32(sq.KTail, ktail)
	return ktail - *sq.KHead
}

func ioUringUnmapRings(sq *IoUringSq, cq *IoUringCq) {
	_ = syscall.Munmap(sq.Ring)
	if cq.Ring != nil && SliceByteAddrEqual(cq.Ring, sq.Ring) {
		_ = syscall.Munmap(cq.Ring)
	}
}

// NewIoUringQueue create an iouring queue.
// entries set the number of elements in the ring,entries needs to be 2^n.
// If entries is not 2^n, it will automatically grow to 2^n.
// flags is used to set the flag bits of the kernel ring.
func NewIoUringQueue(entries int, flags uint32) (*IoUring, error) {
	var params IoUringParams
	params.Flags = flags

	return NewIoUringQueueParams(entries, &params)
}

func NewIoUringQueueParams(entries int, params *IoUringParams) (*IoUring, error) {
	fd, err := IoUringSetup(entries, params)
	if err != nil {
		return nil, err
	}

	return IoUringQueueMmap(fd, params)
}

func IoUringQueueMmap(fd int, params *IoUringParams) (*IoUring, error) {
	var ring IoUring
	ring.RingFd = fd
	ring.Flags = params.Flags

	if err := ring.ioUringMmap(fd, params); err != nil {
		return nil, err
	}

	return &ring, nil
}
