package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/ThingsIXFoundation/data-aggregator/mapping/store/clouddatastore/models"
	"github.com/ThingsIXFoundation/data-aggregator/utils"
	"github.com/ThingsIXFoundation/http-utils/encoding"
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/ThingsIXFoundation/types"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
)

type ChallengeRequest struct {
	Owner common.Address `json:"owner"`
}

type ChallengeResponse struct {
	Owner     common.Address `json:"owner"`
	Challenge string         `json:"challenge"`
}

type SignatureRequest struct {
	Owner     common.Address `json:"owner"`
	Challenge string         `json:"challenge"`
	Signature string         `json:"signature"`
}

type SignatureResponse struct {
	Owner common.Address `json:"owner"`
	Code  string         `json:"code"`
}

type CodeCheckRequest struct {
	MapperID types.ID `json:"mapperId"`
	Code     string   `json:"code"`
}

func generateCode() (string, error) {
	length := 8
	charset := "abcdefghijkmnpqrstuvwxyz23456789" // Excludes confusing characters
	charsetLength := big.NewInt(int64(len(charset)))
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		password[i] = charset[randomIndex.Int64()]
	}

	return string(password), nil
}

func (mapi *MappingAPI) CreateChallenge(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 1*time.Minute)
	)
	defer cancel()
	defer r.Body.Close()

	challengeRequest := ChallengeRequest{}

	err := json.NewDecoder(r.Body).Decode(&challengeRequest)
	if err != nil {
		log.Warnf("invalid request body")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	challenge := make([]byte, 32)
	_, err = rand.Read(challenge)
	if err != nil {
		log.WithError(err).Error("cannot get random string")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	authToken := models.DBMappingAuthToken{
		Owner:      utils.AddressToString(challengeRequest.Owner),
		Expiration: time.Now().Add(6 * 30 * 24 * time.Hour),
		Code:       "",
		Challenge:  hex.Dump(challenge),
	}

	err = mapi.store.StoreMappingAuthToken(ctx, &authToken)
	if err != nil {
		log.WithError(err).Error("cannot get random string")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := ChallengeResponse{
		Owner:     challengeRequest.Owner,
		Challenge: hex.Dump(challenge),
	}

	encoding.ReplyJSON(w, r, http.StatusOK, resp)
}

func (mapi *MappingAPI) SubmitSignature(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 1*time.Minute)
	)
	defer cancel()
	defer r.Body.Close()

	signatureRequest := SignatureRequest{}

	err := json.NewDecoder(r.Body).Decode(&signatureRequest)
	if err != nil {
		log.WithError(err).Warnf("invalid request body")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	hash, err := challengeHash(signatureRequest.Owner, signatureRequest.Challenge)
	if err != nil {
		log.WithError(err).Warnf("invalid request")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	authToken, err := mapi.store.GetMappingAuthTokenByChallenge(ctx, signatureRequest.Challenge)
	if err != nil {
		log.WithError(err).Error("cannot get auth token for challenge")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if authToken == nil {
		log.Warnf("invalid request")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	signature := common.FromHex(signatureRequest.Signature)
	if len(signature) > 0 && signature[len(signature)-1] >= 27 {
		signature[len(signature)-1] = signature[len(signature)-1] - 27
	}

	pub, err := crypto.SigToPub(hash, signature)
	if err != nil {
		log.WithError(err).Error("invalid signature")
		http.Error(w, "invalid signature", http.StatusBadRequest)
		return
	}

	signer := crypto.PubkeyToAddress(*pub)
	if signer == (common.Address{}) || signer != signatureRequest.Owner || signer != common.HexToAddress(authToken.Owner) {
		log.WithFields(logrus.Fields{
			"signer": signer,
			"owner":  signatureRequest.Owner,
		}).Error("signature not created by owner")
		http.Error(w, "invalid signature", http.StatusBadRequest)
		return
	}

	code, err := generateCode()
	if err != nil {
		log.WithError(err).Error("cannot get random code")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	authToken.Expiration = time.Now().Add(6 * 30 * 24 * time.Hour)
	authToken.Code = code

	err = mapi.store.DeleteAllMappingAuthTokens(ctx, authToken.Owner)
	if err != nil {
		log.WithError(err).Error("cannot clean auth token")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = mapi.store.StoreMappingAuthToken(ctx, authToken)
	if err != nil {
		log.WithError(err).Error("cannot store auth token")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := SignatureResponse{
		Owner: signatureRequest.Owner,
		Code:  code,
	}

	encoding.ReplyJSON(w, r, http.StatusOK, resp)
}

func challengeHash(owner common.Address, challenge string) ([]byte, error) {
	stringTy, _ := abi.NewType("string", "", nil)
	addressTy, _ := abi.NewType("address", "", nil)

	args := abi.Arguments{
		{Type: stringTy},
		{Type: stringTy},
		{Type: addressTy},
		{Type: stringTy},
		{Type: stringTy},
	}

	packed, err := args.Pack("MAPPINGAUTH", "|", owner, "|", challenge)
	if err != nil {
		return nil, err
	}

	// simulate wallet signing routine
	// https://docs.ethers.org/v5/api/signer/#Signer-signMessage
	// the dashboard uses a hex encoded form of packed (including the 0x prefix)
	// so wallets can show hex instead of a chunck of raw binary data.
	return accounts.TextHash([]byte("0x" + hex.EncodeToString(packed))), nil
}

func (mapi *MappingAPI) CheckCode(w http.ResponseWriter, r *http.Request) {
	var (
		log         = logging.WithContext(r.Context())
		ctx, cancel = context.WithTimeout(r.Context(), 1*time.Minute)
	)
	defer cancel()
	defer r.Body.Close()

	codeCheckRequest := CodeCheckRequest{}

	err := json.NewDecoder(r.Body).Decode(&codeCheckRequest)
	if err != nil {
		log.Warnf("invalid request body")
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	authToken, err := mapi.store.GetMappingAuthTokenByCode(ctx, codeCheckRequest.Code)
	if err != nil {
		log.WithError(err).Error("cannot get auth token for challenge")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if authToken == nil {
		log.Warnf("invalid code")
		http.Error(w, "invalid code", http.StatusUnauthorized)
		return
	}

	mapper, err := mapi.mapperStore.Get(ctx, codeCheckRequest.MapperID)
	if err != nil {
		log.WithError(err).Error("cannot get mapper")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if mapper == nil {
		log.Warnf("unknown mapper")
		http.Error(w, "unknown mapper", http.StatusNotFound)
		return
	}

	if mapper.Owner != nil && *mapper.Owner == common.HexToAddress(authToken.Owner) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
