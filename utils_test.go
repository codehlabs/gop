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

func TestSecurePsswd(t *testing.T) {
	err := SecurePassword(
		"thecatandthelionwentonawalkdowntheriverside",
		true,
		false,
	)
	if err != nil {
		t.Error(err)
	}
}
