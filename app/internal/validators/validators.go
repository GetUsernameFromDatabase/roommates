// Package validators holds validation functions
package validators

import (
	"fmt"
)

// ValidateCookieHashBlockPair see https://github.com/gorilla/securecookie#examples
//   - hashKey must be a key with 32 or 64 bytes
//   - blockKey valid lengths are 16 (will select AES-128), 24 (AES-192), or 32 (AES-256) bytes
func ValidateCookieHashBlockPair(hashKey, blockKey []byte) error {
	if len(hashKey) != 32 && len(hashKey) != 64 {
		return fmt.Errorf("hash key must be 32 or 64 bytes, got %d", len(hashKey))
	}
	switch len(blockKey) {
	case 16, 24, 32:
	default:
		return fmt.Errorf("block key must be 16, 24 or 32 bytes, got %d", len(blockKey))
	}
	return nil
}
