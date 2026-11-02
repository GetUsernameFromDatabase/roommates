// Package assets [embed.FS]
package assets

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"io/fs"
	"sync"
)

//go:embed all:public favicon/favicon.ico favicon/favicon.svg favicon/*.png pwa
var PublicAssetsFS embed.FS

// hashes maps every embedded path to a short sha256 of its content, the files never change while the server runs
var hashes = sync.OnceValue(func() map[string]string {
	m := map[string]string{}
	err := fs.WalkDir(PublicAssetsFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := PublicAssetsFS.ReadFile(p)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(b)
		m[p] = hex.EncodeToString(sum[:8])
		return nil
	})
	if err != nil {
		panic(err)
	}
	return m
})

// Hash is the content hash of an embedded path such as "public/output.css", "" when it does not exist
func Hash(path string) string {
	return hashes()[path]
}

// URL is the address of a file in public with its content hash, so the browser can keep it until it changes
func URL(name string) string {
	if h := Hash("public/" + name); h != "" {
		return "/assets/" + name + "?v=" + h
	}
	return "/assets/" + name
}
