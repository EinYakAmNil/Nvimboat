package nvimboat

import (
	"fmt"
	"log"
	"os"
)

func (nb *Nvimboat) setupLogging() {
	var err error

	nb.LogFile, err = os.OpenFile(nb.Config["log"].(string), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	log.SetOutput(nb.LogFile)
	log.SetFlags(0)
}

func (nb *Nvimboat) Log(val ...any) {
	fmt.Println(val...)
	log.Println(val...)
	msg := fmt.Sprintf(`echo "%v"`, val)
	nb.Nvim.Plugin.Nvim.Exec(msg, false)
}
