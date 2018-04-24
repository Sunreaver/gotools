package googleauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"time"
)

// MakeGoogleAuthenticator 获取key&t对应的验证码
// key 秘钥
// t 1970年的秒
func MakeGoogleAuthenticator(key string, t int64) (string, error) {
	hs, e := hmacSha1(key, t/30)
	if e != nil {
		return "", e
	}
	snum := lastBit4byte(hs)
	d := snum % 1000000
	return fmt.Sprintf("%06d", d), nil
}

// MakeGoogleAuthenticatorForNow 获取key对应的验证码
func MakeGoogleAuthenticatorForNow(key string) (string, error) {
	return MakeGoogleAuthenticator(key, time.Now().Unix())
}

func lastBit4byte(hmacSha1 []byte) int32 {
	if len(hmacSha1) != sha1.Size {
		return 0
	}
	offsetBits := int8(hmacSha1[len(hmacSha1)-1]) & 0x0f
	p := (int32(hmacSha1[offsetBits]) << 24) | (int32(hmacSha1[offsetBits+1]) << 16) | (int32(hmacSha1[offsetBits+2]) << 8) | (int32(hmacSha1[offsetBits+3]) << 0)
	return (p & 0x7fffffff)
}

func hmacSha1(key string, t int64) ([]byte, error) {
	decodeKey, err := base32.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	cData := make([]byte, 8)
	binary.BigEndian.PutUint64(cData, uint64(t))

	h1 := hmac.New(sha1.New, decodeKey)
	_, e := h1.Write(cData)
	if e != nil {
		return nil, e
	}
	return h1.Sum(nil), nil
}
