package lib

import "crypto/sha256"

const (
	charMax = 12
)

var charTable = []rune("abcdefghijklmnop")

func NameHash(key, value string) string {
	keyBs := []byte(key)
	valueBs := []byte(value)

	bs := make([]byte, 0, len(keyBs)+len(valueBs)+1)
	bs = append(bs, keyBs...)
	bs = append(bs, 0)
	bs = append(bs, valueBs...)

	hash := sha256.Sum256(bs)

	res := make([]rune, 0)
	for i, h := range split4Bit(hash[:]) {
		if i == charMax {
			break
		}

		res = append(res, charTable[int(h)])
	}

	return string(res)
}

func split4Bit(bs []byte) []byte {
	res := make([]byte, 0, 2*len(bs))

	for _, b := range bs {
		res = append(res, b>>4, b&0x0f)
	}

	return res
}
