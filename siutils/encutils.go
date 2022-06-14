package siutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

func HmacSah256(message []byte, hmacKey []byte) hash.Hash {
	mac := hmac.New(sha256.New, hmacKey)
	mac.Write(message)
	return mac
}

func HmacSah256HexStr(message []byte, hmacKey []byte) string {
	mac := HmacSah256(message, hmacKey)
	return hex.EncodeToString(mac.Sum(nil))
}
