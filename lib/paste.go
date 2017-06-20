package lib

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"gopkg.in/redis.v3"
	"hash"
	"time"
)

// Size of the per-hash seed
const SEED_SIZE int = 8

// Base key under which we store hashes
const REDIS_KEY string = "net.as65342.go-paste"

func (p *Paste) Save() (hash_s string, err error) {
	var (
		hash     hash.Hash
		seed     []byte
		rkey     string
		rcontent string
		result   *redis.StatusCmd
	)

	if p.Hash != "" {
		err = errors.New("Paste already has a hash: " + p.Hash)
		return
	}

	if p.Expiration < 1*time.Minute {
		err = errors.New("Expiration needs to be larger then 1 minute")
		return
	}

	if p.Expiration > 1440*time.Hour {
		err = errors.New("Expiration needs to be smaller then 60 days")
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
	rcontent = base64.StdEncoding.EncodeToString(p.Content)

	if result = Redis.Set(rkey, rcontent, p.Expiration); result.Err() != nil {
		err = errors.New("Failed to write new hash: " + result.Err().Error())
		return
	}

	hash_s = p.Hash

	return
}

func LoadPaste(hash string) (p *Paste, err error) {
	var (
		rkey    string
		content []byte
		result  *redis.StringCmd
	)

	rkey = REDIS_KEY + "." + hash

	if result = Redis.Get(rkey); result.Err() != nil {
		err = errors.New("Failed to retrieve paste: " + result.Err().Error())
		return
	}

	if content, err = base64.StdEncoding.DecodeString(result.Val()); err != nil {
		err = errors.New("Failed to decode base64: " + err.Error())
		return
	}

	p = &Paste{
		Hash:    hash,
		Content: content,
	}

	return
}
