package service

import "testing"

func TestEncryptDecryptAccount(t *testing.T) {
	pepper := "test_pepper"
	key := "this_is_key"
	text := "this_is_text_pwd"
	salt, err := generateSalt()
	if err != nil {
		t.Errorf("generateSalt failed.err :[%s]", err.Error())
		return
	}
	dec, err := EncryptAccount(pepper, salt, key, text)
	if err != nil {
		t.Errorf("EncryptAccount failed.err :[%s]", err.Error())
		return
	}
	gotText, err := DecryptAccount(pepper, salt, key, dec)
	if err != nil {
		t.Errorf("DecryptAccount failed.err :[%s]", err.Error())
		return
	}
	if gotText != text {
		t.Errorf("Encrypt Decrypt failed")
		return
	}
}
