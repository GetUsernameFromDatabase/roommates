package utils

import (
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
	"runtime"
	"strconv"
)

func GetFileAndLine() string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return "unknown:0"
	}
	filename := filepath.Base(file)
	return filename + ":" + strconv.Itoa(line)
}

// random id for html element
func RandomHtmlID(prefix string) string {
	b := make([]byte, 3) // 3 bytes = 6 hex chars
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return prefix + "-" + hex.EncodeToString(b)
}
