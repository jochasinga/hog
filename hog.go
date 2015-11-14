/*
Package hog is a wrapper around Go `crypto` package for quick hashing of passwords
with salts appended.

Currently only supporting `md5`, `sha1`, and `sha256` packages for both the
salt and the combination.

Examples

Create an instance of `Combination` and supply the hash function type
constant (MD5, SHA1, SHA256) and a password string. Without providing
a salt, it will be created with SHA1 with `Salt.Size = 16` and append to
the password before being hashed.

	package main

		import (
			"fmt"
			"github.com/jochasinga/hog"
		)

		func main() {
			pw := hog.Combination{
				Func: hog.SHA1,
				Secret: "supersecretpassword",
			}

			hogged := pw.New()           // return a []byte of hashed combination
			fmt.Println(string(hogged))  // 0a12a5f5046bdea95f1ba8f8a570d7361922c0a9
		}

A `Salt` can be created separately before using in the `Combination`

	func main() {

		password := "superStrongPassword321"

		salt := hog.Salt{
			Func: hog.MD5,
			Secret: password,
			Size: 16,

		}

		bsalt := salt.New()

		pw := hog.Combination{
			Func: hog.SHA256,
			Secret: password,
			Salt: bsalt
		}

		hogged := pw.New()
	}

Or use helper functions to create hash

	func main() {

		// Same variables as previous

		bsalt := hog.CreateHash(salt)
		hogged := hog.CreateHash(pw)

		// or just create default combination from password string
		fin := hog.CreateHashFromString(password)

	}
*/
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

/* Helper Functions */

// CreateHash will create a hashed `[]byte` from a `Salt` or `Combination`
// instance.
func CreateHash(h Hog) []byte {
	return h.New()
}

// CreateHashFromString will create a hashed combination based on a provided
// password string.
func CreateHashFromString(secret string) []byte {

	salt := generateSalt(secret, 16)

	c := Combination{
		Func:   SHA1,
		Secret: secret,
		Salt:   salt,
	}

	return c.New()
}

// Match two `Hog` instance i.e. password matching
func Match(h1, h2 Hog) bool {
	match := bytes.Equal(h1.New(), h2.New())
	return match
}

// Private function
func generateSalt(secret string, saltSize uint) []byte {
	sc := []byte(secret)
	buf := make([]byte, saltSize, saltSize+sha1.Size)
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
