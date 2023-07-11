package http

import (
	"testing"
)

func TestNewHTTPResponse(t *testing.T) {
	t.Log(NewResponseBytes(200, []byte("{}")))
}
