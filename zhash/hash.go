package zhash

import (
	"encoding/hex"
	"hash"
	"io"

	"golang.org/x/crypto/bcrypt"
)

func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashData(h hash.Hash, data []byte) string {
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func HashReader(h hash.Hash, r io.Reader) (string, error) {
	_, err := io.Copy(h, r)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func HashFile(h hash.Hash, file io.ReadSeeker) (string, error) {
	str, err := HashReader(h, file)
	if err != nil {
		return "", err
	}
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return str, err
	}
	return str, nil
}
