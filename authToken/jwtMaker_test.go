package authToken

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/saalikmubeen/go-grpc-implementation/utils"
	"github.com/stretchr/testify/require"
)

// Happy path test
func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	role := utils.DepositorRole
	duration := time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
}

// Expired token case
func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(utils.RandomOwner(), utils.DepositorRole, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

// Invalid token case where "none" algorithm header is used
// This is a well-known vulnerability in JWT used by attackers to bypass
// the signature verification
func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(utils.RandomOwner(), utils.DepositorRole, time.Minute)
	require.NoError(t, err)

	// Create a new token with "none" algorithm
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)

	// jwt-go doesn't allow creating a token with "none" algorithm
	// as it is a well-known vulnerability in JWT.
	// So we need to use UnsafeAllowNoneSignatureType to bypass this check
	// by telling jwt-go that we know what we are doing and we are using "none" algorithm
	// intentionally by passing jwt.UnsafeAllowNoneSignatureType as the secret key
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
