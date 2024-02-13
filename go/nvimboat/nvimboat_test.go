package nvimboat

import (
	"fmt"
	"testing"
)

func TestTypedef(t *testing.T) {
	nb := new(Nvimboat)
	fmt.Printf("%+v\n", nb)
	fmt.Printf("%+v\n", nb.ChanExecDB)
}
