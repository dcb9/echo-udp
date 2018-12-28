package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	msg := "hi"
	reader := bytes.NewReader([]byte(msg))
	writer := new(bytes.Buffer)
	err := udpClient(context.Background(), "127.0.0.1:59999", reader, writer)
	if err != nil {
		fmt.Println("got error: ", err)
	}

	fmt.Printf("Sent=%s Received=%s\n", msg, writer.String())
}

const defaultReadTimeout = 10 * time.Second
const defaultBuffMaxSize = 64 * 1024;

func udpClient(ctx context.Context, addr string, reader io.Reader, writer io.Writer)  error {
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return err
	}

	done := make(chan error)
	go func() {
		n, err := io.Copy(conn, reader)
		if err != nil {
			done <- err
			return
		}

		fmt.Printf("Packet-written: bytes=%d \n", n)

		err = conn.SetReadDeadline(time.Now().Add(defaultReadTimeout))
		if err != nil {
			done <- err
			return
		}

		buff := make([]byte, defaultBuffMaxSize)
		num, from, err := conn.ReadFrom(buff)
		if err != nil {
			done <- err
			return
		}
		fmt.Printf("Packet-received: bytes=%d from=%s\n", num, from.String())

		io.Copy(writer, bytes.NewReader(buff[:num]))

		done <- nil
	}()

	select {
		case err = <- done:
		case <-ctx.Done():
			err = ctx.Err()
	}

	return err
}
