package encrypt

import "crypto/sha256"

func Crypt(s string) []byte {
	a := sha256.New()
	a.Write([]byte(s))
	b := a.Sum(nil)
	return b
}
