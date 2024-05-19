package main

import (
	"fmt"
	"time"
)

const FARCASTER_EPOCH int64 = 1609459200000

type HubError struct {
	Type    BadRequestType
	Message string
}

type BadRequestType int

const (
	InvalidParam BadRequestType = iota
)

func (e *HubError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type.String(), e.Message)
}

func (b BadRequestType) String() string {
	return [...]string{"InvalidParam"}[b]
}

func getFarcasterTime() (uint32, error) {
	currentMillis := time.Now().UnixMilli()
	return toFarcasterTime(currentMillis)
}

func toFarcasterTime(time int64) (uint32, error) {
	if time < FARCASTER_EPOCH {
		return 0, &HubError{
			Type:    InvalidParam,
			Message: "time must be after Farcaster epoch (01/01/2021)",
		}
	}
	secondsSinceEpoch := (time - FARCASTER_EPOCH) / 1000
	if secondsSinceEpoch > (1<<32)-1 {
		return 0, &HubError{
			Type:    InvalidParam,
			Message: "time too far in future",
		}
	}
	return uint32(secondsSinceEpoch), nil
}

// func fromFarcasterTime(time uint32) (int64, error) {
// 	return int64(time)*1000 + FARCASTER_EPOCH, nil
// }
