package packets

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func RunTCPServer(ctx context.Context, port int) error {
	listenAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp4", listenAddr)
	if err != nil {
		return err
	}
	defer func() {
		_ = listener.Close()
	}()

	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	handleConn := func(conn *net.TCPConn) {
		log.Printf("connection from TCP %s", conn.RemoteAddr())

		defer func() {
			_ = conn.Close()
			log.Printf("lost connection from TCP %s", conn.RemoteAddr())
		}()

		buf := make([]byte, 65536)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			n, err := conn.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
					return
				}

				log.Printf("failed conn.Read: %s", err)
				return
			}

			b := buf[:n]

			_, err = conn.Write(b)
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
					return
				}

				log.Printf("failed conn.Write: %s", err)
				return
			}
		}
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		conn, err := listener.AcceptTCP()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}

		go handleConn(conn)
	}
}
