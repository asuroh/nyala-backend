package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// OauthGoogleURLAPI ...
const OauthGoogleURLAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

// AccessTokenGoogleURLAPI ...
const AccessTokenGoogleURLAPI = "https://oauth2.googleapis.com/token"

// GetGoogleProfile ...
func GetGoogleProfile(token string) (res map[string]interface{}, err error) {
	response, err := http.Get(OauthGoogleURLAPI + token)
	if err != nil {
		return res, err
	}
	if response.StatusCode >= 400 {
		fmt.Println(response)
		return res, errors.New("invalid_google_access_token")
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(responseBody, &res)
	if err != nil {
		return res, err
	}

	return res, err
}

// GetUser ...
func GetUser(token string) (res []byte, err error) {
	response, err := http.Get(OauthGoogleURLAPI + token)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	res, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetAccessToken ...
func GetAccessToken(authorizationCode string) (res map[string]interface{}, err error) {
	data := url.Values{}
	data.Set("code", authorizationCode)
	data.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("GOOGLE_SECRET_ID"))
	data.Set("redirect_uri", os.Getenv("GOOGLE_REDIRECT_URL"))
	data.Set("grant_type", "authorization_code")

	response, err := http.PostForm(AccessTokenGoogleURLAPI, data)
	if err != nil {
		return res, err
	}
	if response.StatusCode != 200 {
		return res, errors.New("Can't get access token")
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return res, errors.New("error_read_body")
	}
	err = json.Unmarshal(responseBody, &res)
	if err != nil {
		return res, err
	}

	return res, err
}
