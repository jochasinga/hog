package hog

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
)

type Hash uint

const (
	MD5    Hash = 1 + iota // import crypto/md5
	SHA1                   // import crypto/sha1
	SHA256                 // import crypto/sha256
)

var hashes = [...]string{
	"MD5",
	"SHA1",
	"SHA256",
}

func (h Hash) String() string {
	return hashes[h]
}

type Hog interface {
	New() []byte
}

type Salt struct {
	Func   Hash
	Size   uint
	Secret string
}

type Combination struct {
	Func   Hash
	Secret string
	Salt   []byte
}

func (s Salt) New() []byte {

	fn := s.Func
	var size uint

	switch fn {
	case MD5:
		size = md5.Size
	case SHA1:
		size = sha1.Size
	case SHA256:
		size = sha256.Size
	default:
		size = sha1.Size
	}

	buf := make([]byte, s.Size, s.Size+size)
	_, err := io.ReadFull(rand.Reader, buf)

	if err != nil {
		fmt.Printf("random read failed: %v", err)
		os.Exit(1)
	}

	var h hash.Hash

	switch fn {
	case MD5:
		h = md5.New()
	case SHA1:
		h = sha1.New()
	case SHA256:
		h = sha256.New()
	default:
		h = sha1.New()
	}

	h.Write(buf)
	h.Write([]byte(s.Secret))

	return h.Sum(buf)
}

func (c Combination) New() []byte {

	if len(c.Salt) == 0 {
		c.Salt = generateSalt(c.Secret, 16)
	}

	salt := c.Salt
	comb := c.Secret + string(salt)

	fn := c.Func

	var pwhash hash.Hash

	switch fn {
	case MD5:
		pwhash = md5.New()
	case SHA1:
		pwhash = sha1.New()
	case SHA256:
		pwhash = sha256.New()
	default:
		pwhash = sha1.New()
	}

	io.WriteString(pwhash, comb)

	return pwhash.Sum(nil)
}

/***** Helper Functions *****/
func CreateHash(h Hog) []byte {
	return h.New()
}

func CreateHashFromString(secret string) []byte {

	salt := generateSalt(secret, 16)

	c := Combination{
		Func: SHA1,
		Secret: secret,
		Salt: salt,
	}

	return c.New()
}

func Match(h1, h2 Hog) bool {
	match := bytes.Equal(h1.New(), h2.New())
	return match
}

func generateSalt(secret string, saltSize uint) []byte {
	sc := []byte(secret)
	buf := make([]byte, saltSize, saltSize + sha1.Size)
	_, err := io.ReadFull(rand.Reader, buf)

	if err != nil {
		fmt.Printf("random read failed: %v", err)
		os.Exit(1)
	}

	hash := sha1.New()
	hash.Write(buf)
	hash.Write(sc)

	return hash.Sum(buf)
}
