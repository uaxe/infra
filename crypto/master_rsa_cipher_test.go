package crypto

import (
	"strings"
	"testing"
)

func TestMasterRsaError(t *testing.T) {
	masterRsaCipher, _ := CreateMasterRsa(matDesc, RandLowStr(100), rsaPrivateKey)
	_, err := masterRsaCipher.Encrypt([]byte("123"))
	if err != nil {
		t.Fatal(err)
		return
	}
	masterRsaCipher, _ = CreateMasterRsa(matDesc, rsaPublicKey, RandLowStr(100))
	_, err = masterRsaCipher.Decrypt([]byte("123"))
	if err != nil {
		t.Fatal(err)
		return
	}

	testPrivateKey := rsaPrivateKey
	[]byte(testPrivateKey)[100] = testPrivateKey[90]
	masterRsaCipher, _ = CreateMasterRsa(matDesc, rsaPublicKey, testPrivateKey)
	_, err = masterRsaCipher.Decrypt([]byte("123"))
	if err != nil {
		t.Fatal(err)
		return
	}

	masterRsaCipher, _ = CreateMasterRsa(matDesc, rsaPublicKey, rsaPrivateKey)

	var cipherData CipherData
	err = cipherData.RandomKeyIv(16/2, ivSize/4)
	if err != nil {
		t.Fatal(err)
		return
	}

	masterRsaCipher, _ = CreateMasterRsa(matDesc, rsaPublicKey, rsaPrivateKey)
	v := masterRsaCipher.(MasterRsaCipher)

	v.PublicKey = strings.Replace(rsaPublicKey, "PUBLIC KEY", "CERTIFICATE", -1)
	_, err = v.Encrypt([]byte("HELLOW"))
	if err != nil {
		t.Fatal(err)
		return
	}

	v.PrivateKey = strings.Replace(rsaPrivateKey, "PRIVATE KEY", "CERTIFICATE", -1)
	_, err = v.Decrypt([]byte("HELLOW"))
	if err != nil {
		t.Fatal(err)
		return
	}
}
