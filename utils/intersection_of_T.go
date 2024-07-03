package utils

import protos "github.com/noctisatrae/farseer/protos"

func MsgTypeToInt(msgType protos.MessageType) int64 {
	return int64(msgType.Number())
}

func IntersectionOfArrays[T comparable](a []T, b []T) []T {
	set := make([]T, 0)

	for _, v := range a {
		if containsGeneric(b, v) {
			set = append(set, v)
		}
	}

	return set
}

func containsGeneric[T comparable](b []T, e T) bool {
	for _, v := range b {
		if v == e {
			return true
		}
	}
	return false
}
