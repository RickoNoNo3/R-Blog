package bloghelper

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func GetAdminHash() string {
	hash := sha256.Sum256([]byte(time.Now().Format(time.RFC3339)))
	hashCode := hex.EncodeToString(hash[:])
	return hashCode
}
