package internal

import (
	"encoding/json"
	"fmt"
	"github.com/Vadim992/avito/pkg/logger"
	"net/http"
	"runtime/debug"
)

type ErrorStruct struct {
	Err string
}

func NewErrorStruct(err error) *ErrorStruct {
	return &ErrorStruct{Err: err.Error()}
}

func ServerErr(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%e\n%s", err, debug.Stack())

	if err := logger.ErrLog.Output(2, trace); err != nil {
		logger.ErrLog.Println("failed to show error stack trace: %v", err)
	}

	funcErr := responseJSONWithErr(w, err, http.StatusInternalServerError)

	if funcErr != nil {
		logger.ErrLog.Println("err with marshaling JSON")
	}

}

func ClientErr(w http.ResponseWriter, code int, err error) {
	logger.ErrLog.Println(err)

	if code == http.StatusBadRequest {

		funcErr := responseJSONWithErr(w, err, http.StatusBadRequest)

		if funcErr != nil {
			logger.ErrLog.Println("err with marshaling JSON")

			return
		}

		return
	}

	w.WriteHeader(code)
}

func responseJSONWithErr(w http.ResponseWriter, err error, code int) error {
	w.WriteHeader(code)

	errStruct := NewErrorStruct(err)

	content, err := json.Marshal(errStruct)

	if err != nil {
		return err
	}

	if _, err := w.Write(content); err != nil {
		return err
	}

	return nil
}
