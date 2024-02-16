package nvimboat

import (
	"fmt"
	"log"
	"os"
)

func SetupLogging(path string) error {
	logFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(logFile)
	log.SetFlags(0)
	return nil
}

func (nb *Nvimboat) Log(val ...any) {
	var msg string
	for _, v := range val {
		msg += fmt.Sprintf("%+v\n", v)
	}
	log.Println(msg)
	nb.Nvim.Exec("echo "+msg, false)
}
