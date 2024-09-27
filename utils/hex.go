package utils

import "encoding/hex"

func BytesCopy(dst []byte, data []byte) {
	dstLen := len(dst)
	dataLen := len(data)
	if dataLen > dstLen {
		data = data[dataLen-dstLen:]
	}

	copy(dst[dstLen-len(data):], data)
}

func HexToBytes(val []byte, receiver *[]byte, allowOdd bool) error {
	valLen := len(val)
	if valLen < 2 {
		return hex.ErrLength
	}

	if valLen == 2 {
		if string(val) == "0x" || string(val) == "0X" {
			*receiver = nil
			return nil
		}

		return hex.ErrLength
	}

	isOdd := len(val)%2 != 0

	if isOdd {
		if allowOdd {
			val = val[1:]
			val[0] = '0'
		} else {
			return hex.ErrLength
		}
	} else {
		val = val[2:]
	}

	data, err := hex.DecodeString(string(val))
	if err != nil {
		return err
	}

	*receiver = data
	return nil
}
