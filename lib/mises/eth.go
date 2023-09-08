package mises

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"

	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/btcutil/bech32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ripemd160"
)

func AuthWithEthSignature(auth string) (misesid string, misPubkey string, err error) {

	v, err := url.ParseQuery(auth)
	if err != nil {
		return
	}
	address := v.Get("address")
	sigStr := v.Get("sig")
	nonce := v.Get("nonce")
	ethPubkey := v.Get("pubkey")

	// verify signature
	sinMsg := "address=" + address + "&nonce=" + nonce
	if !VerifySignature(sinMsg, ethPubkey, sigStr) {
		err = errors.New("Invalid auth signature")
		return
	}
	// get misesid by eth public key
	misesid, misPubkey, err = getMisesidByEthPubkey(ethPubkey)
	if err != nil {
		return
	}

	return
}

func VerifySignature(msg, pubkey, sigStr string) bool {

	hash := toHash(msg)
	sig, err := hex.DecodeString(sigStr)
	if err != nil {
		fmt.Println("sig DecodeString error: ", err.Error())
		return false
	}
	publicKeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		fmt.Println("publicKey DecodeString error: ", err.Error())
		return false
	}
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), sig)
	if err != nil {
		fmt.Println("Ecrecover error: ", err.Error())
		return false
	}
	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	if !matches {
		return false
	}
	signatureNoRecoverID := sig[:len(sig)-1] // remove recovery ID

	return crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
}

func toHash(msg string) (hash common.Hash) {

	data := []byte(msg)
	hash = crypto.Keccak256Hash(data)

	return
}

var AddressPrefix = "mises"
var MisesIDPrefix = "did:mises:"

func getMisesidByEthPubkey(ethPubkey string) (misesid string, misesPubkey string, err error) {

	misesPubkey, err = ethPubkeyToMisesPubkey(ethPubkey)
	if err != nil {
		return
	}
	publicKeyBytes, err := hex.DecodeString(misesPubkey)
	if err != nil {
		return
	}
	mid, err := ConvertAndEncode(
		AddressPrefix,
		PubKeyAddrBytes(publicKeyBytes),
	)
	if err != nil {
		return
	}
	misesid = MisesIDPrefix + mid

	return
}

func ethPubkeyToMisesPubkey(pubkey string) (mpub string, err error) {

	pubBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		return
	}
	pub, err := crypto.UnmarshalPubkey(pubBytes)
	if err != nil {
		return
	}
	misesPubkey := (*btcec.PublicKey)(pub)
	mpub = hex.EncodeToString(misesPubkey.SerializeCompressed())

	return
}

func PubKeyAddrBytes(pubkey []byte) []byte {
	sha := sha256.Sum256(pubkey)
	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha[:]) // does not error
	pubKeyAddrBytes := hasherRIPEMD160.Sum(nil)
	return pubKeyAddrBytes
}

func ConvertAndEncode(hrp string, data []byte) (string, error) {
	converted, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("encoding bech32 failed: %w", err)
	}

	return bech32.Encode(hrp, converted)
}
