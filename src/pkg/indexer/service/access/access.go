// Package access User document keys index
package access

import (
	"encoding/hex"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/keystore"
)

type AccessIndex struct {
	index    indexer.Indexer
	keystore *keystore.KeyStore
}

func New() *AccessIndex {
	return &AccessIndex{
		index:    indexer.Init("access"),
		keystore: keystore.New(),
	}
}

// Add user's document key
func (u *AccessIndex) Add(userId string, docStorageId *[32]byte, docKey []byte) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	// Getting user privateKey
	userPubKey, _, err := u.keystore.Get(userId)
	if err != nil {
		return err
	}

	// Document key encryption
	keyEncrypted, err := keybox.SealAnonymous(docKey, userPubKey)
	if err != nil {
		return err
	}

	// Index doc_id -> encrypted_doc_key
	indexKey := sha3.Sum256(append(docStorageId[:], userUUID[:]...))
	indexKeyStr := hex.EncodeToString(indexKey[:])

	if err = u.index.Add(indexKeyStr, keyEncrypted); err != nil {
		return err
	}

	return nil
}

// Get user key
func (u *AccessIndex) Get(userId string) (*[]byte, error) {
	var keyEncrypted []byte
	err := u.index.GetById(userId, &keyEncrypted)
	return &keyEncrypted, err
}
