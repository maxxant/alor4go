//usr/bin/env go run "$0" "$@" ; exit "$?"

//go:build ignore
// +build ignore

package main

// вебсокет. пинг понг

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/maxxant/alor4go"

	"nhooyr.io/websocket"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ctx2, _ := context.WithCancel(ctx)

	var bearer = "Bearer " + alor4go.GetJWT()
	head := http.Header{}
	head.Add("Authorization", bearer)
	opt := &websocket.DialOptions{
		HTTPHeader: head,
	}

	sock, resp, err := websocket.Dial(ctx, "wss://api.alor.ru/ws", opt)
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
			fmt.Println("ping..")
			t := time.Now()
			if err = sock.Ping(ctx); err != nil {
				fmt.Println("no pong")
				break
			} else {
				fmt.Println("..pong", time.Since(t))
			}
		}
	}()

	for {
		typ, dat, err := sock.Read(ctx)
		if err != nil {
			break
		}
		fmt.Println(typ, dat)
	}

	<-ctx2.Done()
	sock.Close(websocket.StatusNormalClosure, "")
}
