package nvimboat

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
)

func SetupLogging(logPath string) (err error) {
	logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/SetupLogging"))
		return
	}
	logOutputs := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(logOutputs)
	log.SetFlags(0)

	return
}

func Log(val ...any) {
	var msg string
	for _, v := range val {
		if reflect.ValueOf(v).Kind() == reflect.Struct {
			msg += fmt.Sprintf("%+v\n", prettyStruct(v))
		} else {
			msg += fmt.Sprintf("%+v\n", v)
		}
	}
	log.Println(msg)
	if Nvim != nil {
		Nvim.Command(`echo "` + msg + `"`)
	}
}
