package utils

import (
	"encoding/hex"
)

func BytesToHex(bytes []byte) string {
	hexString := hex.EncodeToString(bytes)
	return "0x" + hexString
}