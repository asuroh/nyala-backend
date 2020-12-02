package middleware

import (
	"context"
	"errors"
	"fmt"
	"kriyapeople/model"
	"strings"

	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	apiHandler "kriyapeople/server/handler"
	"kriyapeople/usecase"
)

type jwtClaims struct {
	jwt.StandardClaims
}

// VerifyMiddlewareInit ...
type VerifyMiddlewareInit struct {
	*usecase.ContractUC
}

// VerifyPermissionInit ...
type VerifyPermissionInit struct {
	*usecase.ContractUC
	Menu string
}

func userContextInterface(ctx context.Context, req *http.Request, subject string, body map[string]interface{}) context.Context {
	return context.WithValue(ctx, subject, body)
}

func (m VerifyMiddlewareInit) verifyJWT(r *http.Request, role string, singleLogin bool) (res map[string]interface{}, err error) {
	claims := &jwtClaims{}

	tokenAuthHeader := r.Header.Get("Authorization")
	if !strings.Contains(tokenAuthHeader, "Bearer") {
		return res, errors.New("Invalid token")
	}
	tokenAuth := strings.Replace(tokenAuthHeader, "Bearer ", "", -1)

	_, err = jwt.ParseWithClaims(tokenAuth, claims, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		secret := m.ContractUC.EnvConfig["TOKEN_SECRET"]
		return []byte(secret), nil
	})
	if err != nil {
		return res, errors.New("Invalid Token!")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return res, errors.New("Expired Token!")
	}

	// Decrypt payload
	res, err = m.ContractUC.Jwe.Rollback(claims.Id)
	if err != nil {
		return res, errors.New("Error when load the payload!")
	}

	// Check if the token provided has a valid role
	if res["role"] == nil {
		return res, errors.New("Invalid " + role + " token!")
	}
	if res["role"].(string) != role {
		return res, errors.New("Not an " + role + " token!")
	}

	if singleLogin && role == "user" {
		var deviceID string
		err = m.ContractUC.GetFromRedis("userDeviceID"+res["id"].(string), &deviceID)
		if err != nil {
			return res, errors.New("Invalid Device!")
		}
		if deviceID != res["device_id"].(string) {
			return res, errors.New("Expired Device Token!")
		}
	}

	return res, nil
}

func (m VerifyMiddlewareInit) verifyRefreshJWT(r *http.Request, role string) (res map[string]interface{}, err error) {
	claims := &jwtClaims{}

	tokenAuthHeader := r.Header.Get("Authorization")
	if !strings.Contains(tokenAuthHeader, "Bearer") {
		return res, errors.New("Invalid token")
	}
	tokenAuth := strings.Replace(tokenAuthHeader, "Bearer ", "", -1)

	_, err = jwt.ParseWithClaims(tokenAuth, claims, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		secret := m.ContractUC.EnvConfig["TOKEN_REFRESH_SECRET"]
		return []byte(secret), nil
	})
	if err != nil {
		return res, errors.New("Invalid Token!")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return res, errors.New("Expired Token!")
	}

	// Decrypt payload
	res, err = m.ContractUC.Jwe.Rollback(claims.Id)
	if err != nil {
		return res, errors.New("Error when load the payload!")
	}

	// Check if the token provided has a valid role
	if res["role"] == nil {
		return res, errors.New("Invalid " + role + " token!")
	}
	if res["role"].(string) != role {
		return res, errors.New("Not an " + role + " token!")
	}

	return res, nil
}

// VerifySuperadminTokenCredential ...
func (m VerifyMiddlewareInit) VerifySuperadminTokenCredential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jweRes, err := m.verifyJWT(r, "admin", false)
		if err != nil {
			apiHandler.RespondWithJSON(w, 401, 401, err.Error(), []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		// Check id in table
		adminUc := usecase.AdminUC{ContractUC: m.ContractUC}
		admin, err := adminUc.FindByID(jweRes["id"].(string), false)
		if admin.ID == "" {
			apiHandler.RespondWithJSON(w, 401, 401, "Not found!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}
		if admin.RoleName != model.RoleCodeSuperadmin {
			apiHandler.RespondWithJSON(w, 401, 401, "only admin can access", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		jweRes["userName"] = admin.Information.UserName
		jweRes["roleName"] = admin.RoleName

		ctx := userContextInterface(r.Context(), r, "user", jweRes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// VerifyAdminTokenCredential ...
func (m VerifyMiddlewareInit) VerifyAdminTokenCredential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jweRes, err := m.verifyJWT(r, "admin", false)
		if err != nil {
			apiHandler.RespondWithJSON(w, 401, 401, err.Error(), []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		// Check id in table
		adminUc := usecase.AdminUC{ContractUC: m.ContractUC}
		admin, err := adminUc.FindByID(jweRes["id"].(string), false)
		if admin.ID == "" {
			apiHandler.RespondWithJSON(w, 401, 401, "Not found!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		jweRes["userName"] = admin.Information.UserName
		jweRes["roleName"] = admin.RoleName

		ctx := userContextInterface(r.Context(), r, "user", jweRes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
