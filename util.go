package main

import (
	"crypto/sha256"
	"encoding/base64"
)

func generateKey() string {
	return "this.is.sparta"
}

func hash(s string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}
