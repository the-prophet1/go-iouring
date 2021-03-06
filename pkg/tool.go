package pkg

import (
	"reflect"
	"unsafe"
)

func SliceByteAddrEqual(slice1, slice2 []byte) bool {
	return reflect.ValueOf(slice1).Pointer() == reflect.ValueOf(slice2).Pointer()
}

func SliceByteAddr(slice []byte) uintptr {
	return reflect.ValueOf(slice).Pointer()
}

func SliceBytePoint(slice []byte) unsafe.Pointer {
	return unsafe.Pointer(SliceByteAddr(slice))
}

func ByteSliceToSqes(slice []byte) []IoUringSqe {
	sliceHdr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	var res []IoUringSqe
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&res))
	hdr.Cap = sliceHdr.Cap / int(unsafe.Sizeof(IoUringSqe{}))
	hdr.Len = sliceHdr.Len / int(unsafe.Sizeof(IoUringSqe{}))
	hdr.Data = sliceHdr.Data
	return res
}

func PtrToUint32Slice(ptr uintptr, len, cap int) []uint32 {
	var res []uint32
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&res))
	hdr.Data = ptr
	hdr.Len = len
	hdr.Cap = cap
	return res
}

func PtrToCqes(ptr uintptr, len, cap int) []IoUringCqe {
	var res []IoUringCqe
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&res))
	hdr.Data = ptr
	hdr.Len = len
	hdr.Cap = cap
	return res
}
