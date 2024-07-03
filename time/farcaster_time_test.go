package time_test

import (
	"testing"

	"github.com/noctisatrae/farseer/time"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestFromFarcasterTime(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	unix, err := time.FromFarcasterTime(107778482)
	assert.NoError(t, err)

	log.Debug(unix)
}

func TestToFarcasterTime(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	fTime, err := time.ToFarcasterTime(1717237682000)
	assert.NoError(t, err)

	assert.Equal(t, int64(107778482), fTime)
}
