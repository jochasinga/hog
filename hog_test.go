package hog_test

import (
	"github.com/jochasinga/hog"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"reflect"
)

var (
	salt              hog.Salt
	combination       hog.Combination
	password          hog.Hog
	badPassword       hog.Hog
	emptySaltPassword hog.Hog
	badSaltedPassword hog.Hog
	correctPassword   hog.Hog
)

var _ = Describe("Salt", func() {
	BeforeEach(func() {
		salt = hog.Salt{
			Func:   hog.SHA1,
			Size:   16,
			Secret: "mypassword",
		}
	})

	Describe("Create new Salt", func() {
		Context("With New() method", func() {
			It("should return `[]uint8`", func() {
				Expect(reflect.TypeOf(salt.New()).String()).To(Equal("[]uint8"))
			})
		})
	})

	Describe("Two Salts", func() {
		Context("When comparing two Salts", func() {
			It("should never be identical", func() {
				Expect(bytes.Equal(salt.New(), salt.New())).To(Equal(false))
			})
		})
	})
})

var _ = Describe("Combination", func() {

	BeforeEach(func() {
		password = hog.Combination{
			Func:   hog.SHA1,
			Secret: "mypassword",
			Salt:   []byte("@#!s_&"),
		}
		badPassword = hog.Combination{
			Func:   hog.SHA1,
			Secret: "badPassWord",
			Salt:   []byte("@#!s_&"),
		}
		badSaltedPassword = hog.Combination{
			Func:   hog.SHA1,
			Secret: "mypassword",
			Salt:   []byte("!_@$1aB"),
		}
		correctPassword = hog.Combination{
			Func:   hog.SHA1,
			Secret: "mypassword",
			Salt:   []byte("@#!s_&"),
		}
		emptySaltPassword = hog.Combination{
			Func:   hog.SHA1,
			Secret: "mypassword",
		}
	})

	Describe("Create new password combination", func() {
		Context("With New() method", func() {
			It("should return a `[]uint8`", func() {
				Expect(reflect.TypeOf(password.New()).String()).To(Equal("[]uint8"))
			})
		})
		Context("with New() method and no Salt provided", func() {
			It("should generate a salt based on the password", func() {

			})
		})
	})

	Describe("Two combinations", func() {
		Context("with identical passwords and salts", func() {
			It("should be identical", func() {
				Expect(bytes.Equal(password.New(), correctPassword.New())).To(Equal(true))
			})
		})
		Context("with identical salts but different passwords", func() {
			It("should not be identical", func() {
				Expect(bytes.Equal(password.New(), badPassword.New())).To(Equal(false))
			})
		})
		Context("with identical passwords but different salts", func() {
			It("should not be identical", func() {
				Expect(bytes.Equal(password.New(), badSaltedPassword.New())).To(Equal(false))
			})
		})
	})
})

var _ = Describe("Hog", func() {
	BeforeEach(func() {
		salt = hog.Salt{
			Func:   hog.SHA1,
			Size:   16,
			Secret: "mypassword",
		}
		password = hog.Combination{
			Func:   hog.SHA1,
			Secret: "mypassword",
			Salt:   []byte("@#!s_&"),
		}
	})

	Describe("Create a new salt", func() {
		Context("With CreateHash() function", func() {
			It("should return a `[]uint8`", func() {
				Expect(reflect.TypeOf(hog.CreateHash(salt)).String()).To(Equal("[]uint8"))
			})
		})
	})
	Describe("Create a new combination", func() {
		Context("With CreateHash() function", func() {
			It("should return a `[]uint8`", func() {
				Expect(reflect.TypeOf(hog.CreateHash(password)).String()).To(Equal("[]uint8"))
			})
		})
	})
})
