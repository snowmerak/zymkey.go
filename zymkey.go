package zymkey

// #cgo LDFLAGS: -lzk_app_utils
// #cgo CFLAGS: -I/usr/include/zymkey/
// #include <stdlib.h>
// #include "zk_app_utils.h"
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

type Zymkey struct {
	ctx C.zkCTX
}

// Read implements io.Reader.
func (z *Zymkey) Read(p []byte) (n int, err error) {
	randBytes, err := z.GenerateRandomBytes(len(p))
	if err != nil {
		return 0, err
	}

	copy(p, randBytes)
	return len(randBytes), nil
}

func NewZymkey() (*Zymkey, error) {
	var (
		ctx C.zkCTX
	)

	ret := C.zkOpen(&ctx)
	if ret < 0 {
		return nil, fmt.Errorf("zkOpen() failed with code %d", ret)
	}

	return &Zymkey{ctx: ctx}, nil
}

func (z *Zymkey) Close() error {
	ret := C.zkClose(z.ctx)
	if ret < 0 {
		return fmt.Errorf("zkClose() failed with code %d", ret)
	}

	return nil
}

func (z *Zymkey) GenerateRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, fmt.Errorf("length must be positive")
	}

	buffer := []byte(nil)
	var rdata *C.uint8_t
	ret := C.zkGetRandBytes(z.ctx, &rdata, C.int(length))
	if ret < 0 {
		return nil, fmt.Errorf("zkGetRandBytes() failed with code %d", ret)
	} else {
		buffer = C.GoBytes(unsafe.Pointer(rdata), C.int(length))
	}
	if rdata != nil {
		C.free(unsafe.Pointer(rdata))
	}

	return buffer, nil
}

func (z *Zymkey) LockBuffer(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data must not be empty")
	}

	src := (*C.uint8_t)(unsafe.Pointer(&data[0]))
	dst := (*C.uint8_t)(nil)
	var dstSize C.int

	ret := C.zkLockDataB2B(z.ctx, src, C.int(len(data)), &dst, &dstSize, C.bool(false))
	if ret < 0 {
		return nil, fmt.Errorf("zkLockDataB2B() failed with code %d", ret)
	}

	lockedData := C.GoBytes(unsafe.Pointer(dst), dstSize)
	if dst != nil {
		C.free(unsafe.Pointer(dst))
	}

	return lockedData, nil
}

func (Z *Zymkey) UnlockBuffer(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data must not be empty")
	}

	src := (*C.uint8_t)(unsafe.Pointer(&data[0]))
	srcSize := C.int(len(data))
	dst := (*C.uint8_t)(nil)
	var dstSize C.int

	ret := C.zkUnlockDataB2B(Z.ctx, src, srcSize, &dst, &dstSize, C.bool(false))
	if ret < 0 {
		return nil, fmt.Errorf("zkUnlockDataB2B() failed with code %d", ret)
	}

	unlockedData := C.GoBytes(unsafe.Pointer(dst), dstSize)
	if dst != nil {
		C.free(unsafe.Pointer(dst))
	}

	return unlockedData, nil
}

func (z *Zymkey) Now() (time.Time, error) {
	var epochTime C.uint32_t
	ret := C.zkGetTime(z.ctx, &epochTime, C.bool(true))
	if ret < 0 {
		return time.Time{}, fmt.Errorf("zkGetTime() failed with code %d", ret)
	}

	return time.Unix(int64(epochTime), 0), nil
}
