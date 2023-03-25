package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type Cryptoer interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

// CryptJWT signing Key
type CryptJWT struct {
	signingKey []byte
	cryptoer   Cryptoer
}

type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewJWT creates a new JWT instance.
// The signing key is used to sign the token.
// The cryptoer is used to encrypt and decrypt the token.It can be nil.
func NewJWT(SigningKey []byte, cryptoer Cryptoer) *CryptJWT {
	return &CryptJWT{signingKey: SigningKey, cryptoer: cryptoer}
}

// CreateToken creates a new token
func (j *CryptJWT) CreateToken(claims CustomClaims) (string, error) {
	jwTtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwTtoken.SignedString(j.signingKey)
	if err != nil {
		return "", err
	}
	if j.cryptoer != nil {
		return j.cryptoer.Encrypt(token)
	}
	return token, nil
}

// ParseToken parses the token.
func (j *CryptJWT) ParseToken(tokenString string) (*CustomClaims, error) {

	// 解密token
	if j.cryptoer != nil {
		var err error
		tokenString, err = j.cryptoer.Decrypt(tokenString)
		if err != nil {
			return nil, err
		}
	}

	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return j.signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	claims, ok := token.Claims.(*CustomClaims)

	if ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, jwt.ErrInvalidType
}
