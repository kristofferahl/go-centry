// Code generated by go-bindata. DO NOT EDIT.
// sources:
// schemas/manifest.json (3.142kB)

package config

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _schemasManifestJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x56\x41\x8f\x9b\x3c\x10\xbd\xe7\x57\x58\xfe\xbe\xe3\x6e\x50\xaf\xb9\xf6\x54\xa9\x52\x7b\xaf\x56\x95\x03\x03\x78\x8b\x67\xe8\xd8\x49\x37\xaa\xf8\xef\x15\x86\x36\x60\xe3\xa4\x5b\xca\xa1\xa7\xdd\xbc\x99\x79\x66\x3c\xf3\x1e\x7c\xdf\x09\x21\xff\xb7\x79\x0d\x46\xc9\x83\x90\xb5\x73\xed\x21\xcb\x9e\x2d\xe1\xe3\x80\xee\x89\xab\x6c\xf8\xf7\x3f\xf9\xe0\xd3\x75\xf1\x33\xd5\x1e\xb2\xac\xd2\xae\x3e\x1d\xf7\x39\x99\xec\x0b\x6b\xeb\xa8\x2c\x81\x55\xdd\x64\x15\x3d\xe6\x80\x8e\x2f\x63\xb9\xcd\x8c\x42\x5d\x82\x75\xfb\x9e\x7f\x20\x73\x97\x16\x7a\x36\x3a\x3e\x43\xee\x06\xac\x65\x6a\x81\x9d\x06\x2b\x0f\xa2\x7f\x42\x21\xa4\xcd\x59\xb7\xee\x0a\x4c\x4a\x15\xb3\xba\xf8\x4a\x0f\x6b\x07\x66\x9a\x37\xc9\xb4\x8e\x35\x56\x72\x0c\x74\xfe\x6f\x37\x14\xca\x9c\x8c\x51\x58\xac\x3c\x61\xd2\xc6\x18\x59\x68\x66\x8c\xa0\x32\x10\x60\xf1\xb3\x3e\xcc\xa3\x46\xe3\x7b\xc0\xca\xd5\xf2\x20\xde\x4c\x42\xdd\x34\x4f\xb6\xca\x27\x6c\xc0\x5c\x43\xd3\x6e\xc3\x5c\xc0\x30\x62\x4d\xb8\xcd\x01\x0a\x91\x9c\xea\xf9\xc3\x49\xc4\xf3\x4b\xf7\xaf\x8b\x02\x6e\x3c\xe0\x91\xa8\x01\x85\x33\x82\xdd\x02\x95\x64\xf8\x7a\xd2\x0c\xbd\x92\x3e\x45\x4b\x11\x0f\xf3\x17\xf0\xb4\xb8\xbc\xd4\x86\x7d\x6d\xbc\xbb\x63\xcd\xab\xc6\x04\x78\x32\x41\xbb\x1e\x5f\xcc\x16\xc3\x5d\xc6\xa8\x46\x07\x15\x70\x1c\xb0\xd0\x04\xa3\xbb\xde\x96\x88\x06\xb9\x9d\xf8\x6c\x4d\xec\x56\x51\x87\x31\xf5\xf2\x1b\xc7\x02\x9e\x3f\x6f\xd7\xd4\xe6\xea\x2c\xa0\x54\xa7\x66\xdd\xbd\x25\xc9\x27\x62\x7b\x85\x6e\xff\x41\xf7\xf0\x44\xf1\x9e\xdf\x71\x8f\x9c\xb0\xd4\xd5\x92\x79\x04\x96\x90\x32\x84\x25\x35\xdd\x1a\x5c\x62\x6c\xd3\x0e\xd3\x0b\xb7\x92\xf8\x0c\x6c\xff\x3a\x69\x43\x55\x8a\x30\x72\xd5\x5b\xbe\xda\x53\xc1\x19\x9a\x08\xbe\x27\x84\xb4\xbd\xfa\xbb\x3c\x9e\xe2\x02\xef\xa5\x25\x2d\xe1\xdf\x14\xe3\x12\x0e\xcc\x14\x3b\xaf\x7f\x47\xa1\xce\x65\x80\x3f\xcd\x7e\x77\x81\x72\x5b\x86\x52\xbf\xfc\x49\xa3\x29\xcd\x4f\xc5\x92\x12\x0e\xe0\x59\x33\xa1\x01\x74\x1f\x39\x3e\x7f\xe5\x1a\xd4\xba\x80\x77\xe8\x80\x51\x35\x6f\xe3\xef\x49\x71\x53\xeb\x29\xa2\x0f\xd1\xbb\xfd\x0e\xcf\x2e\xe0\x5b\x76\x8b\x99\x33\x0c\xa3\xea\x2b\x7d\x55\x5c\x71\xfd\x3c\x9e\x7b\xc6\xae\xaf\xed\x7e\x04\x00\x00\xff\xff\x23\x6e\x7a\x73\x46\x0c\x00\x00")

func schemasManifestJsonBytes() ([]byte, error) {
	return bindataRead(
		_schemasManifestJson,
		"schemas/manifest.json",
	)
}

func schemasManifestJson() (*asset, error) {
	bytes, err := schemasManifestJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schemas/manifest.json", size: 3142, mode: os.FileMode(0644), modTime: time.Unix(1632472096, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x72, 0xb5, 0x4d, 0xc7, 0xa9, 0x2c, 0x8c, 0xdf, 0x1a, 0x42, 0x3c, 0x5f, 0xb2, 0x4f, 0xe9, 0x49, 0xb2, 0x2e, 0x2d, 0x3d, 0x77, 0x3d, 0xe2, 0x17, 0x38, 0x2e, 0x64, 0x32, 0x82, 0x1b, 0x6, 0x36}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"schemas/manifest.json": schemasManifestJson,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"schemas": {nil, map[string]*bintree{
		"manifest.json": {schemasManifestJson, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
