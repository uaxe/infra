package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

type aesCfb struct {
	encrypter cipher.Stream
	decrypter cipher.Stream
}

func newAesCfb(cd CipherData) (Cipher, error) {
	block, err := aes.NewCipher(cd.Key)
	if err != nil {
		return nil, err
	}
	encrypter := cipher.NewCFBEncrypter(block, cd.IV)
	decrypter := cipher.NewCFBDecrypter(block, cd.IV)
	return &aesCfb{encrypter, decrypter}, nil
}

func (c *aesCfb) Encrypt(src io.Reader) io.Reader {
	reader := &cfbEncryptReader{
		encrypter: c.encrypter,
		src:       src,
	}
	return reader
}

type cfbEncryptReader struct {
	encrypter cipher.Stream
	src       io.Reader
}

func (reader *cfbEncryptReader) Read(data []byte) (int, error) {
	plainText := make([]byte, len(data), len(data))
	n, err := reader.src.Read(plainText)
	if n > 0 {
		plainText = plainText[0:n]
		reader.encrypter.XORKeyStream(data, plainText)
	}
	return n, err
}

func (c *aesCfb) Decrypt(src io.Reader) io.Reader {
	return &cfbDecryptReader{
		decrypter: c.decrypter,
		src:       src,
	}
}

type cfbDecryptReader struct {
	decrypter cipher.Stream
	src       io.Reader
}

func (reader *cfbDecryptReader) Read(data []byte) (int, error) {
	cryptoText := make([]byte, len(data), len(data))
	n, err := reader.src.Read(cryptoText)
	if n > 0 {
		cryptoText = cryptoText[0:n]
		reader.decrypter.XORKeyStream(data, cryptoText)
	}
	return n, err
}
