//usr/bin/env go run "$0" "$@" ; exit "$?"

//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {

	refreshToken := os.Getenv("ALOR_REFRESH")

	if refreshToken == "" {
		panic("empty refresh token")
	}

	fmt.Println("refreshToken:", refreshToken)

	/*
		Чтобы выполнить авторизованный запрос, добавьте в запрос заголовок с именем "Authorization"
		и значением, состоящим из префикса "Bearer " и валидного JWT токена.

		Срок жизни JWT короткий: это сделано для безопасности. Для большинства вариантов использования
		API мы рекоммендуем использовать механизм refresh token.

		Механизм refresh token позволяет получать JWT с новым сроком жизни.
		Для этого отправьте POST запрос на адрес
		  - https://oauthdev.alor.ru/refresh?token={refreshToken} (тестовый контур)
		  - или https://oauth.alor.ru/refresh?token={refreshToken} (боевой контур).
		Если у refresh token не истек срок жизни и не он не был отозван, то в теле ответа в поле AccessToken вернётся свежий JWT токен.

		Срок жизни refresh token, получаемого обычным способом — 1 месяц.
		Срок жизни refresh token, получаемого самостоятельным выписыванием — год.
	*/

	apiUrl := "https://oauth.alor.ru"
	resource := "/refresh"

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodPost, urlStr, nil)

	query := req.URL.Query()
	query.Add("token", refreshToken)
	req.URL.RawQuery = query.Encode()

	if false {
		dump, _ := httputil.DumpRequestOut(req, true)
		fmt.Println(string(dump))
	}

	resp, _ := client.Do(req)
	// fmt.Println(resp.Status)

	defer resp.Body.Close()

	type JSResp struct {
		AccessToken string
	}

	decoder := json.NewDecoder(resp.Body)
	var jsresp JSResp
	{
		err := decoder.Decode(&jsresp)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("jsresp:", jsresp)
}
