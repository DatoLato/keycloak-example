package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func generateSynapseToken(userID string, macKey []byte) string {
	encodedUser := base64.RawURLEncoding.EncodeToString([]byte(userID))

	mac := hmac.New(sha1.New, macKey)
	mac.Write([]byte(userID))
	signature := mac.Sum(nil)

	encodedSig := base64.RawURLEncoding.EncodeToString(signature)

	return fmt.Sprintf("syt_%s_%s", encodedUser, encodedSig)
}

func main() {

	config := initConfiguration()
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("username", config.Username)
	data.Set("password", config.Password)
	data.Set("scope", "openid profile email offline_access")

	resp, err := http.Post(config.KeycloakURL, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Ошибка запроса к Keycloak: %s", body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		panic(err)
	}

	fmt.Println("Access Token:", tokenResp.AccessToken)
	fmt.Println("-------------------------------------")
	fmt.Println("-------------------------------------")
	fmt.Println("Refresh Token:", tokenResp.RefreshToken)

	// TODO: Я ПОКА НЕ НАШЕЛ ЗАКОННОГО СПОСОБА  ОТПРАВИТЬ ЗАПРОС НА ПОЛУЧЕНИЕ ТОКЕНА ОТ SYNAPSE, ПОЭТОМУ ПОКА ВЫДЕРНУЛ ИЗ БАЗЫ
	// Основная проблема в том, что это не bearer в нашем понимании, состоит из syt - обазательный префикс, dGVzdDEyMz00MG1haWwucnU - зашифрованный userID в base64, оставшаяся часть - шифрованный макарун секрет
	req, _ := http.NewRequest("GET", config.MatrixURL, nil)
	//token := generateSynapseToken("@admin:vm-196f9ca4.na4u.ru", []byte(config.MacaroonSercetKey))
	req.Header.Set("Authorization", "Bearer syt_dGVzdDEyMz00MG1haWwucnU_KcBgvuBBhYWgTylxDLVX_1MkgYb")

	client := &http.Client{}
	matrixResp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer matrixResp.Body.Close()

	matrixBody, _ := io.ReadAll(matrixResp.Body)
	fmt.Println("-------------------------------------")
	fmt.Println("-------------------------------------")

	fmt.Println("Ответ от Synapse:", string(matrixBody))
}
