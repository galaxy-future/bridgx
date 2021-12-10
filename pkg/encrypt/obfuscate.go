package encrypt

import (
	"encoding/base64"
	"errors"
)

var (
	ErrWrongPepperOrSalt = errors.New("wrong pepper or salt")
	ErrRestoreTextFailed = errors.New("restore text failed")
)

// ObfuscateText will obfuscate text with pepper and salt.
func ObfuscateText(pepper, text, salt string) string {
	messedTextTop, messedTextDown := messUpOrder(text)
	messedPepperTop, messedPepperDown := messUpOrder(pepper)
	messedSaltTop, messedSaltDown := messUpOrder(salt)
	return messedPepperTop + messedTextTop + messedSaltTop + messedPepperDown + messedSaltDown + messedTextDown
}

// RestoreText will restore obfuscated text with pepper and salt.
func RestoreText(pepper, obfuscated, salt string) (string, error) {
	messedPepperTop, messedPepperDown := messUpOrder(pepper)
	messedSaltTop, messedSaltDown := messUpOrder(salt)
	textLen := len(obfuscated) - len(messedPepperTop) - len(messedPepperDown) - len(messedSaltTop) - len(messedSaltDown)
	tlHalf := textLen / 2
	restText, err := splitHead(obfuscated, messedPepperTop)
	if err != nil {
		return "", err
	}
	restText, tTop, err := extractHead(restText, tlHalf)
	if err != nil {
		return "", err
	}
	restText, err = splitHead(restText, messedSaltTop)
	if err != nil {
		return "", err
	}
	restText, err = splitHead(restText, messedPepperDown)
	if err != nil {
		return "", err
	}
	tDown, err := splitHead(restText, messedSaltDown)
	if err != nil {
		return "", err
	}
	return restoreText(tTop + tDown)
}

func splitHead(t, head string) (string, error) {
	hl := len(head)
	if hl > len(t) {
		return "", ErrRestoreTextFailed
	}

	if h := t[:hl]; h != head {
		return "", ErrWrongPepperOrSalt
	}
	return t[hl:], nil
}

func extractHead(t string, hl int) (rest, head string, err error) {
	if hl > len(t) {
		return "", "", ErrRestoreTextFailed
	}
	return t[hl:], t[:hl], nil
}

func messUpOrder(t string) (string, string) {
	b64 := base64.StdEncoding.EncodeToString([]byte(t))
	messed := swapOddEven(b64)
	half := len(messed) / 2
	return messed[:half], messed[half:]
}

func restoreText(t string) (string, error) {
	r := swapOddEven(t)
	dr, err := base64.StdEncoding.DecodeString(r)
	if err != nil {
		return "", err
	}
	return string(dr), err
}

func swapOddEven(t string) string {
	r := []rune(t)
	for i := 0; i <= len(r); i += 2 {
		if i+1 > len(r)-1 {
			return string(r)
		}
		r[i], r[i+1] = r[i+1], r[i]
	}
	return string(r)
}
