package nvimboat

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
)

func SetupLogging(logPath string) (err error) {
	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	logOutputs := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(logOutputs)
	log.SetFlags(0)

	return
}

func Log(val ...any) {
	var (
		msg string
		w   any
	)
	for _, v := range val {
		if reflect.ValueOf(v).Kind() == reflect.Pointer {
			w = reflect.ValueOf(v).Elem().Interface()
		} else {
			w = v
		}
		if reflect.ValueOf(w).Kind() == reflect.Struct {
			msg += fmt.Sprintf("%+v\n", prettyStruct(w))
		} else {
			msg += fmt.Sprintf("%+v\n", prettyStruct(w))
		}
	}
	log.Println(msg)
	if Nvim != nil {
		Nvim.Command(`echo "` + msg + `"`)
	}
}
