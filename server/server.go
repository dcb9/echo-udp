package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Start a echo UDP server, listen on 127.0.0.1:59999")

	echoUDPServer(context.Background(), "127.0.0.1:59999")
}

const defaultBuffMaxSize = 64 * 1024;
const defaultWriteTimeout = 10 * time.Second

func echoUDPServer(ctx context.Context, addr string) error {
	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer pc.Close()

	buff := make([]byte, defaultBuffMaxSize)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, rAddr, err := pc.ReadFrom(buff)
			if err != nil {
				return err
			}
			fmt.Printf("Packet-received: bytes=%d from=%s\n", n, rAddr.String() )

			time.Sleep(500*time.Millisecond)

			err = pc.SetWriteDeadline(time.Now().Add(defaultWriteTimeout))
			if err != nil {
				return err
			}

			n, err = pc.WriteTo(buff[:n], rAddr)
			if err != nil {
				return err
			}

			fmt.Printf("Packet-written: bytes=%d to=%s\n", n, rAddr.String())
		}
	}
}
