//usr/bin/env go run "$0" "$@" ; exit "$?"

//go:build ignore
// +build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/maxxant/alor4go"

	"github.com/oklog/ulid/v2"
	"nhooyr.io/websocket"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ctx2, _ := context.WithCancel(ctx)

	jwt := alor4go.GetJWT()

	sock, resp, err := websocket.Dial(ctx, "wss://api.alor.ru/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		fmt.Println(resp.Status)
	}
	defer sock.Close(websocket.StatusInternalError, "the sky is falling")

	go func() {
		for {
			// delay measurement
			time.Sleep(1 * time.Second)
			t := time.Now()
			if err = sock.Ping(ctx); err != nil {
				break
			} else {
				fmt.Println("ping-pong RTT:", time.Since(t))
			}
		}
	}()

	moex := alor4go.WsReqOrderBookGetAndSubscribeExchangeMOEX
	sber := "SBER"
	depth := 20
	guid := ulid.Make().String()

	opcode := alor4go.OrderBookGetAndSubscribe
	// opcode := "OrderBookGetAndSubscribe" // fix for old yaml spec

	obReq := alor4go.WSProdOrderBookJSONRequestBody{
		Token:    &jwt,
		Guid:     &guid,
		Exchange: &moex,
		Code:     &sber,
		Depth:    &depth,
		Opcode:   &opcode,
	}

	buf, err := json.Marshal(obReq)
	if err != nil {
		log.Fatalln(err)
		return
	}

	sock.Write(ctx, websocket.MessageText, buf)
	fmt.Println("sub ok")

	// example for resp:
	// MessageText { "data": {"snapshot":true,"bids":[],"asks":[],"timestamp":1671223802,"ms_timestamp":1671223802393,"existing":true}, "guid":"01GMM8RJ59N9AY05F8FGQEF3JC" }
	// MessageText {"requestGuid":"01GMM8RJ59N9AY05F8FGQEF3JC","httpCode":200,"message":"Handled successfully"}

	for {
		typ, dat, err := sock.Read(ctx)
		if err != nil {
			break
		}
		fmt.Println(typ, string(dat))
	}

	<-ctx2.Done()
	sock.Close(websocket.StatusNormalClosure, "")
}
