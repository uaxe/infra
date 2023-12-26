package crypto

import (
	"testing"
)

func TestCfbContentEncryptCipherError(t *testing.T) {
	// crypto bucket
	masterRsaCipher, _ := CreateMasterRsa(matDesc, rsaPublicKey, rsaPrivateKey)
	contentProvider := CreateAesCfbCipher(masterRsaCipher, AES_TYPE_128)
	cc, err := contentProvider.ContentCipher()
	if err != nil {
		t.Fatal(err)
		return
	}

	var cipherData CipherData
	cipherData.RandomKeyIv(31, 15)

	_, err = cc.Clone(cipherData)
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestCfbCreateCipherDataError(t *testing.T) {
	// crypto bucket
	masterRsaCipher, _ := CreateMasterRsa(matDesc, "", "")
	contentProvider := CreateAesCfbCipher(masterRsaCipher, AES_TYPE_128)

	v := contentProvider.(aesCtrCipherBuilder)
	_, err := v.createCipherData()
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestCfbContentCipherCDError(t *testing.T) {
	var cd CipherData

	// crypto bucket
	masterRsaCipher, _ := CreateMasterRsa(matDesc, "", "")
	contentProvider := CreateAesCfbCipher(masterRsaCipher, AES_TYPE_128)

	v := contentProvider.(aesCtrCipherBuilder)
	_, err := v.contentCipherCD(cd)
	if err != nil {
		t.Fatal(err)
		return
	}
}
