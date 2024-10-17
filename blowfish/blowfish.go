//go:build !solution

package blowfish

// #cgo pkg-config: libcrypto
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <openssl/blowfish.h>
import "C"

import (
	"crypto/cipher"
	"unsafe"
)

var _ cipher.Block = (*Blowfish)(nil)

type Blowfish struct {
	key C.BF_KEY
}

func New(key []byte) *Blowfish {
	if len(key) < 1 || len(key) > 56 { // максимальный размер ключа для Blowfish
		panic("invalid key length")
	}
	bf := &Blowfish{}
	C.BF_set_key(&bf.key, C.int(len(key)), (*C.uchar)(unsafe.Pointer(&key[0])))
	return bf
}

func (bf *Blowfish) BlockSize() int {
	return 8
}

func (bf *Blowfish) Encrypt(dst, src []byte) {
	if len(src) != 8 || len(dst) != 8 {
		panic("data must be exactly 8 bytes")
	}

	C.BF_ecb_encrypt(
		(*C.uchar)(unsafe.Pointer(&src[0])),
		(*C.uchar)(unsafe.Pointer(&dst[0])),
		&bf.key,
		C.int(1),
	)
}

func (bf *Blowfish) Decrypt(dst, src []byte) {
	if len(src) != 8 || len(dst) != 8 {
		panic("data must be exactly 8 bytes")
	}

	C.BF_ecb_encrypt(
		(*C.uchar)(unsafe.Pointer(&src[0])),
		(*C.uchar)(unsafe.Pointer(&dst[0])),
		&bf.key,
		C.int(0),
	)
}
