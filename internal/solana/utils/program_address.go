package utils

import (
	"crypto/ed25519"
	"crypto/sha256"
	"errors"

	"filippo.io/edwards25519"
	"github.com/mr-tron/base58"
)

const (
	MaxSeedLength = 32
	MaxSeeds      = 16
	PDA_MARKER    = "ProgramDerivedAddress"
)

func FindProgramAddressSync(seeds [][]byte, programId []byte) (string, uint8, error) {
	var nonce uint8 = 255

	for nonce > 0 {
		seedsWithNonce := append(seeds, []byte{nonce})
		address, isValidAddress := CreateProgramAddressSync(seedsWithNonce, programId)
		if !isValidAddress {
			nonce--
			continue
		}
		return address, nonce, nil
	}

	return "", 0, errors.New("unable to find a viable program address nonce")
}

func CreateProgramAddressSync(seeds [][]byte, programID []byte) (string, bool) {
	buf := []byte{}

	for _, seed := range seeds {
		if len(seed) > MaxSeedLength {
			return "", false
		}
		buf = append(buf, seed...)
	}

	buf = append(buf, programID[:]...)
	buf = append(buf, []byte(PDA_MARKER)...)
	hash := sha256.Sum256(buf)

	if IsOnCurve(hash[:]) {
		return "", false
	}

	return base58.Encode(hash[:]), true
}

func IsOnCurve(b []byte) bool {
	if len(b) != ed25519.PublicKeySize {
		return false
	}
	_, err := new(edwards25519.Point).SetBytes(b)
	isOnCurve := err == nil
	return isOnCurve
}
