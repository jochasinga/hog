# hog

**hog** is a wrapper around Go `crypto` package for quick hashing of passwords
with salts appended.

Currently only supporting `md5`, `sha1`, and `sha256` packages for both the
salt and the combination.

## Examples

```Go

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
```
Without `Salt` field provided to `Combination`, it will auto-generate one based
on the password string provided using `SHA1` and `Salt.Size = 16`

```Go
// Creating a combination and corresponding salt
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
```

Or use helper functions to create hash

```Go

func main() {

	// Same variables as previous

	bsalt := hog.CreateHash(salt)
	hogged := hog.CreateHash(pw)

	// or just create default combination from password string
	fin := hog.CreateHashFromString(password)

}

```
