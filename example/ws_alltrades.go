//usr/bin/env go run "$0" "$@" ; exit "$?"

//go:build ignore
// +build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/maxxant/alor4go"

	"github.com/oklog/ulid/v2"
	"nhooyr.io/websocket"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jwt := alor4go.GetJWT()

	xchg := alor4go.WsReqAllTradesGetAndSubscribeExchangeMOEX
	tiker := "SBER"

	sock, resp, err := websocket.Dial(ctx, "wss://api.alor.ru/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 101 {
		fmt.Println(resp.Status)
	}

	var ws_ping_age_ms int32

	go func() {

		go func() {
			for {
				// control ping age
				time.Sleep(100 * time.Millisecond)
				atomic.AddInt32(&ws_ping_age_ms, 100)
				if atomic.LoadInt32(&ws_ping_age_ms) > 3000 {
					cancel() // disconnect ws via context
					break
				}
			}
		}()

		for {
			// ping delay measurement
			time.Sleep(1 * time.Second)
			t := time.Now()
			if err = sock.Ping(ctx); err != nil {
				cancel() // disconnect ws via context
				break
			} else {
				_ = t
				//fmt.Println("ping-pong RTT:", time.Since(t))
				atomic.StoreInt32(&ws_ping_age_ms, 0)
			}
		}
	}()

	guid := ulid.Make().String()
	opcode := alor4go.AllTradesGetAndSubscribe

	obReq := alor4go.WSProdAllTradesGetAndSubscribeJSONRequestBody{
		Token:    &jwt,
		Guid:     &guid,
		Exchange: &xchg,
		Code:     &tiker,
		Opcode:   &opcode,
	}

	buf, err := json.Marshal(obReq)
	if err != nil {
		log.Fatalln(err)
		return
	}

	sock.Write(ctx, websocket.MessageText, buf)
	fmt.Println("sub ok")

	for {
		typ, dat, err := sock.Read(ctx)
		if err != nil {
			cancel()
			break
		}

		//fmt.Println(typ, string(dat))

		if typ == websocket.MessageText {
			func() {
				// example for resp:
				// { "data": {"id":2022426556961456082,"orderno":0,"symbol":"CNY-3.23","qty":39,"price":10.24500,"time":"2023-01-17T14:11:58.8950000","timestamp":1673964718895,"side":"sell","oi":9183872,"existing":false}, "guid":"01GQ00XMT11ZN2SB2KWQH53Y0W" }

				trades := &alor4go.WsResTradesGetAndSubscribe{}
				if err := json.Unmarshal(dat, &trades); err != nil {
					fmt.Println("err json.Unmarshal", err)
					return
				}

				if trades.Data != nil && trades.Data.Price != nil && trades.Data.Id != nil {
					id, _ := trades.Data.Id.Int64()
					fmt.Println("price:", *trades.Data.Price, " id:", id)
				}
			}()
		}
	}

	<-ctx.Done()
}
