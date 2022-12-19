//usr/bin/env go run "$0" "$@" ; exit "$?"

//go:build ignore
// +build ignore

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/maxxant/alor4go"
)

func main() {
	/*
		Чтобы выполнить авторизованный запрос, добавьте в запрос заголовок с именем "Authorization"
		и значением, состоящим из префикса "Bearer " и валидного JWT токена.

		Срок жизни JWT короткий: это сделано для безопасности. Для большинства вариантов использования
		API мы рекоммендуем использовать механизм refresh token.

		Тестовый контур: https://apidev.alor.ru
		Боевой контур: https://api.alor.ru
	*/

	opt := []alor4go.ClientOption{}

	client, _ := alor4go.NewClientWithResponses("https://api.alor.ru", opt...)

	// anonymous req.  15 min delay

	{
		edi := []alor4go.RequestEditorFn{}
		resp, err := client.LocalTimeWithResponse(context.Background(), edi...)
		if err != nil {
			panic(err)
		}

		defer resp.HTTPResponse.Body.Close()
		if resp.HTTPResponse.StatusCode != 200 {
			fmt.Println(resp.HTTPResponse.Status)
		} else {
			//fmt.Println("resp:", *resp.JSON200)
			fmt.Println("resp:", time.Unix(int64(*resp.JSON200), 0))
		}
	}

	//

	var bearer = "Bearer " + alor4go.GetJWT()

	requestEdit := func(ctx context.Context, req *http.Request) error {
		if req == nil {
			return errors.New("req==nil")
		}

		req.Header.Add("Authorization", bearer)
		return nil
	}

	// auth req. real time

	{
		resp, err := client.LocalTimeWithResponse(context.Background(), requestEdit)
		if err != nil {
			panic(err)
		}

		defer resp.HTTPResponse.Body.Close()
		if resp.HTTPResponse.StatusCode != 200 {
			fmt.Println(resp.HTTPResponse.Status)
		} else {
			//fmt.Println("resp:", *resp.JSON200)
			fmt.Println("resp:", time.Unix(int64(*resp.JSON200), 0))
		}
	}
}
