package utils

import (
	"encoding/hex"
)

func BytesToHex(bytes []byte) string {
	if bytes == nil {
		return ""
	}
	hexString := hex.EncodeToString(bytes)
	if hexString == "" {
		return ""
	}
	return "0x" + hexString
}
