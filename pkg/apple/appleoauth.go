package apple

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"kriyapeople/pkg/interfacepkg"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthKeysAPI ...
const AuthKeysAPI = "https://appleid.apple.com/auth/keys"

// GetJWK ...
func GetJWK() (res map[string]interface{}, err error) {
	response, err := http.Get(AuthKeysAPI)
	if err != nil {
		return res, errors.New("invalid_apple_api")
	}
	if response.StatusCode >= 400 {
		return res, errors.New("invalid_apple_api")
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return res, errors.New("error_read_body")
	}
	err = json.Unmarshal(responseBody, &res)
	fmt.Println(res["error"])
	if err != nil {
		fmt.Println(err)
		return res, err
	}

	return res, err
}

// GetPublicKey ...
func GetPublicKey(kid string) (res string, err error) {
	JwkArrayInterface, err := GetJWK()
	if err != nil {
		return res, err
	}

	JwkArray := JwkArrayInterface["keys"].([]interface{})
	for _, jwk := range JwkArray {
		JwkKid := interfacepkg.InterfaceStringToString(jwk, "kid")
		if JwkKid == kid {
			kty := interfacepkg.InterfaceStringToString(jwk, "kty")
			if kty != "RSA" {
				return res, errors.New("invalid key type")
			}

			// decode the base64 bytes for n
			n := interfacepkg.InterfaceStringToString(jwk, "n")
			nb, err := base64.RawURLEncoding.DecodeString(n)
			if err != nil {
				return res, err
			}

			e := 0
			// The default exponent is usually 65537, so just compare the
			// base64 for [1,0,1] or [0,1,0,1]
			JwkE := interfacepkg.InterfaceStringToString(jwk, "e")
			if JwkE == "AQAB" || JwkE == "AAEAAQ" {
				e = 65537
			} else {
				// need to decode "e" as a big-endian int
				return res, errors.New("need to deocde e")
			}

			pk := &rsa.PublicKey{
				N: new(big.Int).SetBytes(nb),
				E: e,
			}

			der, err := x509.MarshalPKIXPublicKey(pk)
			if err != nil {
				return res, err
			}

			block := &pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: der,
			}

			var out bytes.Buffer
			pem.Encode(&out, block)
			return out.String(), err
		}
	}
	return res, err
}

// VerifyJWT ...
func VerifyJWT(token, email string) (err error) {
	claims := jwt.MapClaims{}

	// Parse token
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return nil, err
	})
	if err != nil {
		return err
	}

	// Generate Public Key
	pubKey, err := GetPublicKey(jwtToken.Header["kid"].(string))
	if err != nil {
		return err
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubKey))
	if err != nil {
		return err
	}

	// Verify token with Public Key
	parts := strings.Split(token, ".")
	err = jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), parts[2], key)
	if err != nil {
		return err
	}

	exp, _ := claims["exp"].(int64)
	if exp < time.Now().Unix() {
		return errors.New("Apple token already expired")
	}
	appleEmail := claims["email"].(string)
	if appleEmail != email {
		return errors.New("Invalid Apple ID")
	}
	emailVerified, _ := claims["email_verified"].(bool)
	if emailVerified != true {
		return errors.New("Email is not verified")
	}

	return err
}

// generate private key from pem encoded string
func getPrivKey(pemEncoded []byte) (*ecdsa.PrivateKey, error) {
	var block *pem.Block
	var x509Encoded []byte
	var err error
	var privateKeyI interface{}
	var privateKey *ecdsa.PrivateKey
	var ok bool

	// decode the pem format & check if its is private key
	block, _ = pem.Decode(pemEncoded)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("Failed to decode PEM block containing private key")
	}

	// get the encoded bytes
	x509Encoded = block.Bytes

	// generate the private key object
	privateKeyI, err = x509.ParsePKCS8PrivateKey(x509Encoded)
	if err != nil {
		return nil, errors.New("Private key decoding failed. " + err.Error())
	}
	// cast into ecdsa.PrivateKey object
	privateKey, ok = privateKeyI.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("Private key is not ecdsa key")
	}

	return privateKey, nil
}
