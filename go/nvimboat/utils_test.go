package nvimboat

import (
	"fmt"
	"testing"
)

func TestSliceDelete(t *testing.T) {
	slice1 := []int{0, 1, 2, 3, 4, 5, 6, 7}
	slice2 := []string{"0", "1", "2", "3", "4", "hello", "world"}

	fmt.Println(slice1, 2, 4, 5)
	fmt.Println(sliceDelete(slice1, 2, 1, 6, 4, 5))
	fmt.Println(slice2, 0, 4)
	fmt.Println(sliceDelete(slice2, 0, 4))
}
