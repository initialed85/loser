package packets

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

func RunTCPClient(ctx context.Context, rawDialAddr string, reportFn func(time.Time, int64, int64, int64, int64)) error {
	dialAddr, err := net.ResolveTCPAddr("tcp4", rawDialAddr)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp4", nil, dialAddr)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	log.Printf("connected to TCP %s", conn.RemoteAddr())
	defer func() {
		log.Printf("lost connection to TCP %s", conn.RemoteAddr())
	}()

	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	sendTicker := time.NewTicker(time.Millisecond * 10)
	defer func() {
		sendTicker.Stop()
	}()

	mu := new(sync.Mutex)

	sent := int64(0)
	received := int64(0)
	outOfOrder := int64(0)
	lost := int64(0)

	buf := make([]byte, 65536)

	reportTicker := time.NewTicker(time.Second * 1)
	defer func() {
		reportTicker.Stop()
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-reportTicker.C:
			}

			mu.Lock()
			sent := sent
			received := received
			outOfOrder := outOfOrder
			lost := lost
			mu.Unlock()

			log.Printf("TCP %s sent: %d, received: %d, outOfOrder: %d, lost: %d", conn.RemoteAddr(), sent, received, outOfOrder, lost)

			reportFn(time.Now(), sent, received, outOfOrder, lost)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-sendTicker.C:
		}

		now := time.Now()
		expiry := now.Add(time.Millisecond * 10)

		err = conn.SetWriteDeadline(expiry)
		if err != nil {
			return err
		}

		err = conn.SetReadDeadline(expiry)
		if err != nil {
			return err
		}

		sent++

		_, err = conn.Write([]byte(fmt.Sprintf("%d", sent)))
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			mu.Lock()
			lost++
			mu.Unlock()

			continue
		}

		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			mu.Lock()
			lost++
			mu.Unlock()

			return err
		}

		b := buf[:n]

		ack, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			return err
		}

		mu.Lock()
		if ack == sent {
			received++
		} else {
			outOfOrder++
		}
		mu.Unlock()
	}
}
