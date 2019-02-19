package parser

import (
	"crypto/sha1"
)

func toSHA1(data []byte) string {
	hash := sha1.New()
	hash.Write(data)
	return string(hash.Sum(nil))
}