package hasher

import (
	"fmt"
	"testing"
)

func Test_Hash(t *testing.T) {
	s := New(10, "somesalt")
	fmt.Println(s.Hash("password"))
}
