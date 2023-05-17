package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	err := GenerateKeyPair(2048)
	if err != nil {
		t.Errorf("Failed to generate key pair: %v", err)
	}

	// Read and parse the private key
	privateKeyData, err := ioutil.ReadFile(PrivateKeyPath)
	if err != nil {
		t.Errorf("Failed to read private key file: %v", err)
	}

	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		t.Errorf("Failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Errorf("Failed to parse private key: %v", err)
	}

	// Read and parse the public key
	publicKeyData, err := ioutil.ReadFile(PublicKeyPath)
	if err != nil {
		t.Errorf("Failed to read public key file: %v", err)
	}

	block, _ = pem.Decode(publicKeyData)
	if block == nil {
		t.Errorf("Failed to parse PEM block containing the public key")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		t.Errorf("Failed to parse public key: %v", err)
	}

	publicKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		t.Errorf("Failed to cast to RSA public key")
	}

	// Perform some basic tests on the keys
	if privateKey.N.Cmp(publicKey.N) != 0 {
		t.Errorf("Private and public key do not have the same modulus")
	}

	// Clean up
	defer os.Remove(PrivateKeyPath)
	defer os.Remove(PublicKeyPath)
}

func TestLoadPrivateKey(t *testing.T) {
	err := GenerateKeyPair(2048)
	if err != nil {
		t.Errorf("Failed to generate key pair: %v", err)
	}

	privateKey, err := LoadPrivateKey(PrivateKeyPath)
	if err != nil {
		t.Fatalf("Failed to load private key: %v", err)
	}

	if privateKey == nil {
		t.Fatalf("Private key is nil")
	}

	// Validate the private key by trying to sign and verify some data
	data := []byte("test data")
	hashed := sha256.Sum256(data)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		t.Fatalf("Failed to sign data with private key: %v", err)
	}

	publicKey, err := LoadPublicKey(PublicKeyPath)
	if err != nil {
		t.Fatalf("Failed to load public key: %v", err)
	}

	if publicKey == nil {
		t.Fatalf("Public key is nil")
	}

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		t.Fatalf("Failed to verify data with public key: %v", err)
	}

	// Clean up
	defer os.Remove(PrivateKeyPath)
	defer os.Remove(PublicKeyPath)
}

func TestLoadPublicKey(t *testing.T) {
	err := GenerateKeyPair(2048)
	if err != nil {
		t.Errorf("Failed to generate key pair: %v", err)
	}

	publicKey, err := LoadPublicKey(PublicKeyPath)
	if err != nil {
		t.Errorf("Failed to load public key: %v", err)
		return
	}
	if publicKey == nil {
		t.Fatalf("Public key is nil")
	}

	// Validate the public key by trying to encrypt and decrypt some data
	data := []byte("test data")

	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		t.Errorf("Failed to encrypt data with public key: %v", err)
	}

	privateKey, err := LoadPrivateKey(PrivateKeyPath)
	if err != nil {
		t.Errorf("Failed to load private key: %v", err)
	}

	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		t.Errorf("Failed to decrypt data with private key: %v", err)
	}

	if string(plainText) != string(data) {
		t.Errorf("Decrypted data does not match original data")
	}

	// Clean up
	defer os.Remove(PrivateKeyPath)
	defer os.Remove(PublicKeyPath)
}

func TestLoadPrivateKeyFail(t *testing.T) {
	// Ensure no private key file exists
	os.Remove(PrivateKeyPath)

	nonExistingPath := "path/to/non/existing/file"
	_, err := LoadPrivateKey(nonExistingPath)
	if err == nil {
		t.Errorf("Expected an error when loading non-existent private key, got nil")
	}
}

func TestLoadPublicKeyFail(t *testing.T) {
	// Ensure no public key file exists
	os.Remove(PublicKeyPath)

	nonExistingPath := "path/to/non/existing/file"
	_, err := LoadPublicKey(nonExistingPath)
	if err == nil {
		t.Errorf("Expected an error when loading non-existent public key, got nil")
	}
}
