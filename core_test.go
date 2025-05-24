package gop

import "testing"

func TestSecurePasswd(t *testing.T) {
	e := SecurePassword("Theblackcatsaltedthecamel00#", true, true)
	if e != nil {
		t.Error(e)
	}
}
