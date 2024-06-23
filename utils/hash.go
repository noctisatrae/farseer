package utils

import (
	"encoding/hex"
)

func BytesToHex(bytes []byte) string {
	hexString := hex.EncodeToString(bytes)
	if hexString == "" {
		return ""
	}
	return "0x" + hexString
}