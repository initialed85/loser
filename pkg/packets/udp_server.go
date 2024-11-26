package packets

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func RunUDPServer(ctx context.Context, port int) error {
	listenAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}

	listener, err := net.ListenUDP("udp", listenAddr)
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

	handleConn := func(conn *net.UDPConn) {
		buf := make([]byte, 65536)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
					return
				}

				log.Printf("failed conn.Read: %s", err)
				return
			}

			b := buf[:n]

			_, err = conn.WriteToUDP(b, addr)
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
					return
				}

				log.Printf("failed conn.Write: %s", err)
				return
			}
		}
	}

	handleConn(listener)

	return nil
}
