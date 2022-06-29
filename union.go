package main

type Union1 uint64

func (u *Union1) SetOffset(offset uint) {
	*u = Union1(offset)
}

func (u Union1) Offset() uint64 {
	return uint64(u)
}

func (u *Union1) SetAddr2(addr2 uint) {
	*u = Union1(addr2)
}

func (u Union1) Addr2() uint64 {
	return uint64(u)
}

type Union2 uint64

func (u *Union2) SetAddr(addr uint) {
	*u = Union2(addr)
}

func (u Union2) Addr() uint64 {
	return uint64(u)
}

func (u *Union2) SetSpliceOffIn(spliceOffIn uint64) {
	*u = Union2(spliceOffIn)
}

func (u Union2) SpliceOffIn() uint64 {
	return uint64(u)
}

type Union3 uint32

func (u *Union3) SetRWFlag(flag int32) {
	*u = Union3(flag)
}

func (u Union3) RWFlag() int32 {
	return int32(u)
}

func (u Union3) FsyncFlags() uint32 {
	return uint32(u)
}

func (u *Union3) SetFsyncFlags(flag uint32) {
	*u = Union3(flag)
}

func (u Union3) PollEvents() uint16 {
	return uint16(u)
}

func (u *Union3) SetPollEvents(event uint16) {
	*u = Union3(event)
}

func (u Union3) Poll32Events() uint32 {
	return uint32(u)
}

func (u *Union3) SetPoll32Events(event uint32) {
	*u = Union3(event)
}

func (u Union3) SyncRangeFlags() uint32 {
	return uint32(u)
}

func (u *Union3) SetSyncRangeFlags(flag uint32) {
	*u = Union3(flag)
}

func (u Union3) MsgFlags() uint32 {
	return uint32(u)
}

func (u *Union3) SetMsgFlags(flag uint32) {
	*u = Union3(flag)
}

func (u Union3) TimeoutFlags() uint32 {
	return uint32(u)
}

func (u *Union3) SetTimeoutFlags(flag uint32) {
	*u = Union3(flag)
}

func (u Union3) CancelFlags() uint32 {
	return uint32(u)
}

func (u *Union3) SetCancelFlags(flag uint32) {
	*u = Union3(flag)
}

func (u Union3) OpenFlags() uint32 {
	return uint32(u)
}

func (u *Union3) SetOpenFlags(flag uint32) {
	*u = Union3(flag)
}

func (u Union3) StatxFlags() uint32 {
	return uint32(u)
}

func (u *Union3) SetStatxFlags(flag uint32) {
	*u = Union3(flag)
}

func (u Union3) FadviseAdvice() uint32 {
	return uint32(u)
}

func (u *Union3) SetFadviseAdviceFlags(advice uint32) {
	*u = Union3(advice)
}

func (u Union3) SpliceFlags() uint32 {
	return uint32(u)
}

func (u *Union3) SetSpliceFlags(flag uint32) {
	*u = Union3(flag)
}

type UnionStruct struct {
	Union4
	Personality uint16
	SpliceFdIn  int32
	pad         [2]uint64
}

type Union4 uint16

func (u Union4) BufIndex() uint16 {
	return uint16(u)
}

func (u *Union4) SetBufIndex(bufIndex uint16) {
	*u = Union4(bufIndex)
}

func (u Union4) BufGroup() uint16 {
	return uint16(u)
}

func (u *Union4) SetBufGroup(bufGroup uint16) {
	*u = Union4(bufGroup)
}
