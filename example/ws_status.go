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

	xchg := alor4go.WsReqInstrumentsGetAndSubscribeV2ExchangeMOEX
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
	opcode := alor4go.InstrumentsGetAndSubscribeV2

	obReq := alor4go.WsReqInstrumentsGetAndSubscribeV2{
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
				/*
					{
					  "data": {
					    "symbol": "SBER",
					    "exchange": "MOEX",
					    "board": "TQBR",
					    "tradingStatus": 17,
					    "tradingStatusInfo": "нормальный период торгов",
					    "priceMin": 228.78,
					    "priceMax": 279.96
					  },
					  "guid": "c328fcf1-e495-408a-a0ed-e20f95d6b813"
					}
				*/
				type InstrumentStatus struct {
					Symbol            string   `json:"symbol"`
					Exchange          string   `json:"exchange"`
					Board             string   `json:"board"`
					TradingStatus     *int     `json:"tradingStatus,omitempty"`     // Торговый статус инструмента
					TradingStatusInfo *string  `json:"tradingStatusInfo,omitempty"` // Описание торгового статуса инструмента
					PriceMin          *float64 `json:"priceMin,omitempty"`
					PriceMax          *float64 `json:"priceMax,omitempty"`
				}
				type InstrumentsGetAndSubscribeResponse struct {
					Data *InstrumentStatus `json:"data,omitempty"`
					Guid *string           `json:"guid,omitempty"`
				}
				statusResp := &InstrumentsGetAndSubscribeResponse{}

				if err := json.Unmarshal(dat, &statusResp); err != nil {
					fmt.Println("err json.Unmarshal", err)
					return
				}

				if statusResp.Data != nil {
					istatus := *statusResp.Data
					fmt.Println("Symbol:", istatus.Symbol, "Exchange:", istatus.Exchange, "Board:", istatus.Board)

					if istatus.TradingStatus != nil {
						fmt.Println("TradingStatus:", *istatus.TradingStatus)
					}
					if istatus.TradingStatusInfo != nil {
						fmt.Println("TradingStatusInfo:", *istatus.TradingStatusInfo)
					}
					if istatus.PriceMin != nil {
						fmt.Println("PriceMin:", *istatus.PriceMin)
					}
					if istatus.PriceMax != nil {
						fmt.Println("PriceMax:", *istatus.PriceMax)
					}
				}
			}()
		}
	}

	<-ctx.Done()
}
