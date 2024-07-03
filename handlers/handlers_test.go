package handlers_test

import (
	"testing"

	"github.com/noctisatrae/farseer/handlers"
	"github.com/stretchr/testify/assert"
)

func TestInitParams(t *testing.T) {
	// new handler
	dummy := handlers.Handler{
		InitHandler: func(params map[string]interface{}) error {
			params["hello"] = "world"
			params["counter"] = 1
			return nil
		},
	}

	params := map[string]interface{}{
		"foo": "bar",
	}

	err := dummy.InitHandler(params)
	assert.NoError(t, err)

	assert.Equal(t, map[string]interface{}{
		"hello": "world",
	}["hello"], params["hello"])
}
