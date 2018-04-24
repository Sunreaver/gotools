// Package base32 google两步验证中用到的basedecode
//
package base32

import "strings"

var (
	encodeingTable = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	paddingTable   = []uint{0, 6, 4, 3, 1}
	decodingTable  = []byte{}
	padding        = byte('=')
	nullEncode     = byte(0xff)
)

func init() {
	decodingTable = make([]byte, 256)

	for index := 0; index < len(decodingTable); index++ {
		decodingTable[index] = nullEncode
	}

	lowEncodeTable := strings.ToLower(encodeingTable)
	for index := 0; index < len(encodeingTable); index++ {
		decodingTable[encodeingTable[index]] = byte(index)
		decodingTable[lowEncodeTable[index]] = byte(index)
	}

	decodingTable[padding] = 0
}

// Decode 将key编码为base32并返回
// key 待编码字符串
// return 编码后的值
func Decode(key string) []byte {
	paddingAdjustment := []int{0, 1, 1, 1, 2, 3, 3, 4}
	key = strings.Replace(key, string(padding), "", -1)

	encodedBytes := []byte(key)
	encodedLength := len(key)
	encodedBlocks := encodedLength * 5 / 40
	if encodedLength%8 != 0 {
		encodedBlocks++
	}

	expectedDataLength := encodedBlocks * 5

	decodedBytes := make([]byte, expectedDataLength)

	var encodedByte1, encodedByte2, encodedByte3, encodedByte4 byte
	var encodedByte5, encodedByte6, encodedByte7, encodedByte8 byte

	encodedBytesToProcess := encodedLength
	encodedBaseIndex := 0
	decodedBaseIndex := 0
	encodedBlock := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	encodedBlockIndex := 0
	var c byte
	for encodedBytesToProcess >= 1 {
		encodedBytesToProcess--

		c = encodedBytes[encodedBaseIndex]
		encodedBaseIndex++
		if c == padding {
			break
		}
		c = decodingTable[c]
		if c == nullEncode {
			continue
		}
		encodedBlock[encodedBlockIndex] = c
		encodedBlockIndex++
		if encodedBlockIndex == 8 {
			encodedByte1 = encodedBlock[0]
			encodedByte2 = encodedBlock[1]
			encodedByte3 = encodedBlock[2]
			encodedByte4 = encodedBlock[3]
			encodedByte5 = encodedBlock[4]
			encodedByte6 = encodedBlock[5]
			encodedByte7 = encodedBlock[6]
			encodedByte8 = encodedBlock[7]

			decodedBytes[decodedBaseIndex+0] = (encodedByte1 << 3 & 0xf8) | (encodedByte2 >> 2 & 0x07)
			decodedBytes[decodedBaseIndex+1] = (encodedByte2 << 6 & 0xc0) | (encodedByte3 << 1 & 0x3e) | (encodedByte4 >> 4 & 0x01)
			decodedBytes[decodedBaseIndex+2] = (encodedByte4 << 4 & 0xf0) | (encodedByte5 >> 1 & 0x0f)
			decodedBytes[decodedBaseIndex+3] = (encodedByte5 << 7 & 0x80) | (encodedByte6 << 2 & 0x7c) | (encodedByte7 >> 3 & 0x03)
			decodedBytes[decodedBaseIndex+4] = (encodedByte7 << 5 & 0xe0) | (encodedByte8 & 0x1f)

			decodedBaseIndex += 5
			encodedBlockIndex = 0
		}
	}
	encodedByte2 = 0
	encodedByte3 = 0
	encodedByte4 = 0
	encodedByte5 = 0
	encodedByte6 = 0
	encodedByte7 = 0

	switch encodedBlockIndex {
	case 7:
		encodedByte7 = encodedBlock[6]
		fallthrough
	case 6:
		encodedByte6 = encodedBlock[5]
		fallthrough
	case 5:
		encodedByte5 = encodedBlock[4]
		fallthrough
	case 4:
		encodedByte4 = encodedBlock[3]
		fallthrough
	case 3:
		encodedByte3 = encodedBlock[2]
		fallthrough
	case 2:
		encodedByte2 = encodedBlock[1]
		fallthrough
	case 1:
		encodedByte1 = encodedBlock[0]
		decodedBytes[decodedBaseIndex] = byte((int(encodedByte1) << 3 & 0xf8) | (int(encodedByte2) >> 2 & 0x07))
		decodedBytes[decodedBaseIndex+1] = byte((int(encodedByte2) << 6 & 0xc0) | (int(encodedByte3) << 1 & 0x3e) | (int(encodedByte4) >> 4 & 0x01))
		decodedBytes[decodedBaseIndex+2] = byte((int(encodedByte4) << 4 & 0xf0) | (int(encodedByte5) >> 1 & 0x0f))
		decodedBytes[decodedBaseIndex+3] = byte((int(encodedByte5) << 7 & 0x80) | (int(encodedByte6) << 2 & 0x7c) | (int(encodedByte7) >> 3 & 0x03))
		decodedBytes[decodedBaseIndex+4] = byte(int(encodedByte7) << 5 & 0xe0)
	default:
		break
	}
	decodedBaseIndex += paddingAdjustment[encodedBlockIndex]
	return decodedBytes[0:decodedBaseIndex]
}
