package mises

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/btcutil/bech32"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mises-id/sns-socialsvc/lib/utils"
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

	if sigStr == "" {
		return "", "", errors.New("Signature cannot be empty")
	}
	// verify signature
	sigMsg := "address=" + address + "&nonce=" + nonce
	//isValid := VerifySignature(sinMsg, ethPubkey, sigStr)
	//isValid := VerifyEIP191Signature(sinMsg, ethPubkey, sigStr)
	isValid, ethPubkey := VerifyEIP191SignatureByAddress(sigMsg, address, sigStr)
	if !isValid {
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

func VerifyEIP191SignatureByAddress(msg, address, sigStr string) (ok bool, ethPubkey string) {

	hash := eip191Hash(msg)
	sigBytes, err := hexutil.Decode(utils.AddHexPrefix(sigStr))
	if err != nil {
		fmt.Println("sig DecodeString error: ", err.Error())
		return
	}
	if len(sigBytes) < 64 {
		fmt.Println("Invalid signature")
		return
	}
	sigBytes[64] %= 27
	if sigBytes[64] != 0 && sigBytes[64] != 1 {
		fmt.Println("Invalid signature recovery byte")
		return
	}
	pkey, err := crypto.SigToPub(hash.Bytes(), sigBytes)
	if err != nil {
		fmt.Println("Failed to recover public key from signature")
		return
	}

	sigAddress := crypto.PubkeyToAddress(*pkey)
	ethPubkey = hex.EncodeToString(crypto.FromECDSAPub(pkey))
	if err != nil {
		fmt.Println("publicKey DecodeString", ethPubkey, "error:", err.Error())
		return
	}
	if strings.ToLower(sigAddress.Hex()) != strings.ToLower(address) {
		fmt.Println("Signer address must match message address", sigAddress.Hex(), address)
		return
	}

	return true, ethPubkey
}

func VerifyEIP191Signature(msg, pubkey, sigStr string) (ok bool) {

	hash := eip191Hash(msg)
	sigBytes, err := hexutil.Decode(utils.AddHexPrefix(sigStr))
	if err != nil {
		fmt.Println("sig DecodeString error: ", err.Error())
		return
	}
	sigBytes[64] %= 27
	if sigBytes[64] != 0 && sigBytes[64] != 1 {
		fmt.Println("Invalid signature recovery byte")
		return
	}
	publicKeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		fmt.Println("publicKey DecodeString error: ", err.Error())
		return
	}
	signatureNoRecoverID := sigBytes[:len(sigBytes)-1] // remove recovery ID

	return crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
}

func VerifySignature(msg, pubkey, sigStr string) (ok bool) {

	hash := toHash(msg)

	sig, err := hex.DecodeString(utils.RemoveHexPrefix(sigStr))
	if err != nil {
		fmt.Println("sig DecodeString error: ", err.Error())
		return
	}
	publicKeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		fmt.Println("publicKey DecodeString error: ", err.Error())
		return
	}
	signatureNoRecoverID := sig[:len(sig)-1] // remove recovery ID

	return crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
}

func eip191Hash(data string) common.Hash {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256Hash([]byte(msg))
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
