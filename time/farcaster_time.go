package time

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

const FARCASTER_EPOCH int64 = 1609459200000 // January 1, 2021 UTC

// Get the current Farcaster time.
func GetFarcasterTime() (int64, error) {
	return ToFarcasterTime(time.Now().UnixMilli())
}

// Converts from a Unix to Farcaster timestamp.
func ToFarcasterTime(timeMillis int64) (int64, error) {
	if timeMillis < FARCASTER_EPOCH {
		return 0, errors.New("bad_request.invalid_param: time must be after Farcaster epoch (01/01/2022)")
	}
	secondsSinceEpoch := (timeMillis - FARCASTER_EPOCH) / 1000
	if secondsSinceEpoch > math.MaxUint32 {
		return 0, errors.New("bad_request.invalid_param: time too far in future")
	}
	return secondsSinceEpoch, nil
}

// Converts from a Farcaster to Unix timestamp.
func FromFarcasterTime(timeSeconds int64) (int64, error) {
	return timeSeconds*1000 + FARCASTER_EPOCH, nil
}

// Extracts the timestamp from an event ID.
func ExtractEventTimestamp(eventId int64) (int64, error) {
	binaryEventId := fmt.Sprintf("%064b", eventId)
	const SEQUENCE_BITS = 12
	if len(binaryEventId) <= SEQUENCE_BITS {
		return 0, errors.New("invalid event ID")
	}
	binaryTimestamp := binaryEventId[:len(binaryEventId)-SEQUENCE_BITS]
	timestamp, err := strconv.ParseInt(binaryTimestamp, 2, 64)
	if err != nil {
		return 0, err
	}
	return timestamp + FARCASTER_EPOCH, nil
}
