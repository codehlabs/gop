package gop

import (
	"fmt"
	"testing"
)

func TestUID(t *testing.T) {
	id, err := unique_id(16)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(id)
}
