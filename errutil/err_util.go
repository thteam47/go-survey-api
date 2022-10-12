package errutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func ErrorStackTrace(err error) bool {
	return FerrorStackTrace(os.Stderr, err)
}

func SerrorStackTrace(err error) string {
	buf := new(bytes.Buffer)
	ok := FerrorStackTrace(buf, err)
	if !ok {
		return ""
	}
	return buf.String()
}

func FerrorStackTrace(w io.Writer, err error) bool {
	if err == nil {
		return false
	}
	if customError, ok := err.(*customError); ok {
		fmt.Fprintln(w, customError.originStackErr)
		fmt.Fprintln(w)
		st := customError.originStackErr.(stackTracer).StackTrace()
		if len(st) > 2 {
			st = st[2 : len(st)-2]
		}
		fmt.Fprintf(w, "%+v", st)
		fmt.Println()
	} else {
		fmt.Fprintln(w, err)
		fmt.Fprintln(w)
	}
	return err != nil
}

func ExitIfError(err error) {
	if ErrorStackTrace(err) {
		os.Exit(1)
		return
	}
}

func Message(err error) string {
	if err == nil {
		return ""
	}
	if customError, ok := err.(*customError); ok {
		return fmt.Sprint(customError.currentErr)
	}
	return fmt.Sprint(err)
}

func PanicStackTrace(p interface{}) {
	// fmt.Println(p)
	// log.Println(string(debug.Stack()))
	fmt.Fprintln(os.Stderr, p)
	debug.PrintStack()

}
