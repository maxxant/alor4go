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
	"time"

	"github.com/maxxant/alor4go"
)

func main() {
	opt := []alor4go.ClientOption{}

	client, _ := alor4go.NewClientWithResponses("https://api.alor.ru", opt...)

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
		param := &alor4go.DevHistoryParams{
			Symbol:   "SBER",
			Exchange: "MOEX",
			Tf:       "60",
			From:     int32(time.Date(2022, time.December, 8, 13, 0, 0, 0, time.UTC).Unix()),
			To:       int32(time.Date(2022, time.December, 8, 13, 3, 0, 0, time.UTC).Unix()),
		}
		resp, err := client.DevHistory(context.Background(), param, requestEdit)
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

			//fmt.Println("resp:", *resp.JSON200)
			// fmt.Println("resp:", time.Unix(int64(*resp.JSON200), 0))
		}
	}
}
