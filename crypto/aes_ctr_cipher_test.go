package crypto

import "testing"

func TestCreateCipherDataError(t *testing.T) {
	// crypto bucket
	masterRsaCipher, _ := CreateMasterRsa(matDesc, "", "")
	contentProvider := CreateAesCtrCipher(masterRsaCipher)

	v := contentProvider.(aesCtrCipherBuilder)
	_, err := v.createCipherData()
	if err == nil {
		t.Fatal(err)
	}
}
