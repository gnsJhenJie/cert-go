package certgo

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/Alonza0314/cert-go/util"
	logger "github.com/Alonza0314/logger-go"
)

func CreatePrivateKey(keyPath string, overwrite bool) (*ecdsa.PrivateKey, error) {
	logger.Info("CreatePrivateKey", "creating private key")

	// check if private key exists
	if util.FileExists(keyPath) {
		if !overwrite {
			logger.Error("CreatePrivateKey", fmt.Sprintf("private key already exists at %s.", keyPath))
			return nil, errors.New("private key already exists")
		}
		logger.Warn("CreatePrivateKey", "private key already exists. Overwrite it")
		if err := util.FileDelete(keyPath); err != nil {
			logger.Error("CreatePrivateKey", "failed to remove existing private key: "+err.Error())
			return nil, err
		}
	}

	// generate private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logger.Error("CreatePrivateKey", err.Error())
		return nil, err
	}

	// encode private key
	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		logger.Error("CreatePrivateKey", err.Error())
		return nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	})

	// check directory exists
	if !util.FileDirExists(keyPath) {
		logger.Warn("CreatePrivateKey", util.FileDir(keyPath)+" directory not exists, creating...")
		if err := util.FileDirCreate(keyPath); err != nil {
			return nil, err
		}
		logger.Info("CreatePrivateKey", util.FileDir(keyPath)+" directory created")
	}

	// save private key
	if err := util.FileWrite(keyPath, keyPEM, 0644); err != nil {
		return nil, err
	}

	logger.Info("CreatePrivateKey", "private key created")
	return privateKey, nil
}
