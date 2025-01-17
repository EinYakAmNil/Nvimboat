package nvimboat

import (
	"fmt"
	"io"
	"log"
	"os"
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

func (nb *Nvimboat) Log(val ...any) {
	var msg string
	for _, v := range val {
		msg += fmt.Sprintf("%+v", v)
	}
	log.Println(msg)
	if nb.Nvim != nil {
		nb.Nvim.Command(`echo "` + msg + `"`)
	}
}
