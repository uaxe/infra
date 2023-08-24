package crypto

import (
	"io"
	"math/rand"
	"strings"
	"time"

	. "gopkg.in/check.v1"
)

var (
	letters            = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	timeoutInOperation = 3 * time.Second
)

func RandStr(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}

func RandLowStr(n int) string {
	return strings.ToLower(RandStr(n))
}

type CryptoSuite struct {
}

var _ = Suite(&CryptoSuite{})

func (s *CryptoSuite) TestAesCtr(c *C) {
	var cipherData CipherData
	cipherData.RandomKeyIv(32, 16)
	cipher, _ := newAesCtr(cipherData)

	byteReader := strings.NewReader(RandLowStr(100))
	enReader := cipher.Encrypt(byteReader)
	encrypter := &CryptoEncrypter{Body: byteReader, Encrypter: enReader}
	encrypter.Close()
	buff := make([]byte, 10)
	n, err := encrypter.Read(buff)
	c.Assert(n, Equals, 0)
	c.Assert(err, Equals, io.EOF)

	deReader := cipher.Encrypt(byteReader)
	Decrypter := &CryptoDecrypter{Body: byteReader, Decrypter: deReader}
	Decrypter.Close()
	buff = make([]byte, 10)
	n, err = Decrypter.Read(buff)
	c.Assert(n, Equals, 0)
	c.Assert(err, Equals, io.EOF)
}

func (s *CryptoSuite) TestAesCfb(c *C) {
	var cipherData CipherData
	cipherData.RandomKeyIv(32, 16)
	cipher, _ := newAesCfb(cipherData)

	byteReader := strings.NewReader(RandLowStr(100))
	enReader := cipher.Encrypt(byteReader)
	encrypter := &CryptoEncrypter{Body: byteReader, Encrypter: enReader}
	encrypter.Close()
	buff := make([]byte, 10)
	n, err := encrypter.Read(buff)
	c.Assert(n, Equals, 0)
	c.Assert(err, Equals, io.EOF)

	deReader := cipher.Encrypt(byteReader)
	Decrypter := &CryptoDecrypter{Body: byteReader, Decrypter: deReader}
	Decrypter.Close()
	buff = make([]byte, 10)
	n, err = Decrypter.Read(buff)
	c.Assert(n, Equals, 0)
	c.Assert(err, Equals, io.EOF)
}
