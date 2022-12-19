//usr/bin/env go run "$0" "$@" ; exit "$?"

//go:build ignore
// +build ignore

package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/maxxant/alor4go"

	"github.com/deepmap/oapi-codegen/pkg/types"
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

	client, _ := alor4go.NewClient("https://api.alor.ru", opt...)

	var bearer = "Bearer " + alor4go.GetJWT()

	requestEdit := func(ctx context.Context, req *http.Request) error {
		if req == nil {
			return errors.New("req==nil")
		}

		req.Header.Add("Authorization", bearer)
		return nil
	}

	username := os.Getenv("ALOR_USER")
	forts := os.Getenv("ALOR_FORTS")

	{
		resp, err := client.DevUserPortfolio(context.Background(), username, requestEdit)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Println(resp.Status)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error while reading the response bytes:", err)
			}
			fmt.Println("resp:", string([]byte(body)))
		}
	}

	{
		param := &alor4go.TradeStatsParams{
			DateFrom: &types.Date{Time: time.Now().Add(-time.Hour * 240)},
		}
		resp, err := client.TradeStats(context.Background(), alor4go.TradeStatsParamsExchangeMOEX, forts, param, requestEdit)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Println(resp.Status)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error while reading the response bytes:", err)
			}
			fmt.Println("resp:", string([]byte(body)))
		}
	}
}
