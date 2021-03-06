package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"errors"
)

const AesKeyPepper = "bridgx"

var (
	ErrEncryptFailed = errors.New("encrypt failed")
	ErrDecryptFailed = errors.New("decrypt failed")
)

func AESEncrypt(key, plaintext string) (text string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrEncryptFailed
			return
		}
	}()
	keyB := ensureKeyLength(key)
	block, err := aes.NewCipher(keyB)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	iv := make([]byte, blockSize)
	origData := padding([]byte(plaintext), blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cryptText := make([]byte, len(origData))
	blockMode.CryptBlocks(cryptText, origData)
	return base64.StdEncoding.EncodeToString(cryptText), nil
}

func AESDecrypt(key string, ct16 string) (text string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrDecryptFailed
			return
		}
	}()
	ciphertext, err := base64.StdEncoding.DecodeString(ct16)
	if err != nil {
		return "", err
	}
	keyB := ensureKeyLength(key)
	block, err := aes.NewCipher(keyB)
	if err != nil {
		return "", err
	}
	iv := make([]byte, block.BlockSize())
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(origData, ciphertext)
	origData = unPadding(origData)
	return string(origData), nil
}

func ensureKeyLength(key string) []byte {
	keyB := md5.Sum([]byte(key))
	return keyB[:]
}

func padding(src []byte, blockSize int) []byte {
	pad := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(src, padText...)
}

func unPadding(src []byte) []byte {
	length := len(src)
	unPad := int(src[length-1])
	return src[:(length - unPad)]
}
