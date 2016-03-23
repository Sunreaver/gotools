package goTools

import (
	"crypto/hmac"
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/sunreaver/goTools/base32"
)

// MakeAuth 获取key&t对应的验证码
// key 秘钥
// t 1970年的秒
func MakeAuth(key string, t int64) (string, error) {
	hs := hmacSha1(key, t/30)
	if hs == nil {
		return "", errors.New("输入有误")
	}
	snum := dt(hs)
	fmt.Println("snum = ", snum)
	d := snum % 1000000
	return fmt.Sprintf("%06d", d), nil
}

// MakeAuthNow 获取key对应的验证码
func MakeAuthNow(key string) (string, error) {
	return MakeAuth(key, time.Now().Unix())
}

func dt(hmacSha1 []byte) int32 {
	if len(hmacSha1) != sha1.Size {
		return 0
	}
	offsetBits := int8(hmacSha1[len(hmacSha1)-1]) & 0x0f
	p := (int32(hmacSha1[offsetBits]) << 24) | (int32(hmacSha1[offsetBits+1]) << 16) | (int32(hmacSha1[offsetBits+2]) << 8) | (int32(hmacSha1[offsetBits+3]) << 0)
	return (p & 0x7fffffff)
}

func hmacSha1(key string, t int64) []byte {
	decodeKey := base32.Decode(key)

	cData := make([]byte, 8)
	cData[4] = byte(t >> 24 & 0xff)
	cData[5] = byte(t >> 16 & 0xff)
	cData[6] = byte(t >> 8 & 0xff)
	cData[7] = byte(t >> 0 & 0xff)

	h1 := hmac.New(sha1.New, decodeKey)
	_, e := h1.Write(cData)
	if e != nil {
		return nil
	}
	return h1.Sum(nil)
}
