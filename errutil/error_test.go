package errutil

import (
	"fmt"
	"net/http"
	"testing"
)

func TestWrap(t *testing.T) {
	err := Wrap(http.ErrHeaderTooLong, "api.GetData")
	err = Wrap(err, "app.Initialize")
	if !IsError(err, http.ErrHeaderTooLong) {
		t.Error("false negative")
	}
	if IsError(err, http.ErrHandlerTimeout) {
		t.Error("true negative")
	}
	// ErrorStackTrace(err)
	fmt.Println(err)
}
