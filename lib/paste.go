package lib

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash"
)

// Size of the per-hash seed
const SEED_SIZE int = 8

// Base key under which we store hashes
const REDIS_KEY string = "net.as65342.go-paste"

func (p *Paste) Save() (hash_s string, err error) {
	var hash hash.Hash
	var seed []byte
	var rkey string

	if p.Hash != "" {
		err = errors.New("Paste already has a hash: " + p.Hash)
		return
	}

	if p.Expiration.String() == "" {
		err = errors.New("Expiration cannot be empty")
		return
	}

	if Redis == nil {
		err = errors.New("No redis client found")
		return
	}

	// Generate random seed for this paste
	seed = make([]byte, SEED_SIZE)
	if _, err = rand.Read(seed); err != nil {
		err = errors.New("Unable to generate random seed: " + err.Error())
		return
	}

	hash = sha1.New()
	hash.Write(seed)
	hash.Write(p.Content)

	p.Hash = hex.EncodeToString(hash.Sum(nil))

	rkey = REDIS_KEY + "." + p.Hash

	Redis.Set(rkey, base64.StdEncoding.EncodeToString(p.Content), p.Expiration)

	hash_s = p.Hash

	return
}

func LoadPaste(hash string) (p *Paste, err error) {
	var rkey string
	var value string
	var content []byte

	rkey = REDIS_KEY + "." + hash

	value = Redis.Get(rkey).Val()
	if value == "" {
		err = errors.New("No such id: " + hash)
		return
	}

	if content, err = base64.StdEncoding.DecodeString(value); err != nil {
		err = errors.New("Failed to decode base64: " + err.Error())
		return
	}

	p = &Paste{
		Hash:    hash,
		Content: content,
	}

	return
}
