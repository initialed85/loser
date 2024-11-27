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

func RunUDPClient(ctx context.Context, host string, actualReportFn func(time.Time, int64, int64, int64, int64)) error {
	mu := new(sync.Mutex)

	sent := int64(0)
	received := int64(0)
	outOfOrder := int64(0)
	lost := int64(0)

	lastSent := int64(0)
	lastReceived := int64(0)
	lastOutOfOrder := int64(0)
	lastLost := int64(0)

	dialAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:6943", host))
	if err != nil {
		time.Sleep(time.Second * 1)
		return err
	}

	conn, err := net.DialUDP("udp4", nil, dialAddr)
	if err != nil {
		time.Sleep(time.Second * 1)
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	log.Printf("connected to UDP %s", conn.RemoteAddr())
	defer func() {
		log.Printf("lost connection to UDP %s", conn.RemoteAddr())
	}()

	reportFn := func() {
		mu.Lock()

		thisSent := sent - lastSent
		thisReceived := received - lastReceived
		thisOutOfOrder := outOfOrder - lastOutOfOrder
		thisLost := lost - lastLost

		if lastSent == 0 {
			thisSent = 0
		}

		if lastReceived == 0 {
			thisReceived = 0
		}

		if lastOutOfOrder == 0 {
			thisOutOfOrder = 0
		}

		if lastLost == 0 {
			thisLost = 0
		}

		lastSent = sent
		lastReceived = received
		lastOutOfOrder = outOfOrder
		lastLost = lost

		mu.Unlock()

		log.Printf("UDP %s sent: %d, received: %d, outOfOrder: %d, lost: %d", conn.RemoteAddr(), thisSent, thisReceived, thisOutOfOrder, thisLost)
		actualReportFn(time.Now(), thisSent, thisReceived, thisOutOfOrder, thisLost)
	}

	reportFn()

	defer func() {
		reportFn()
	}()

	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	sendTicker := time.NewTicker(time.Millisecond * 10)
	defer func() {
		sendTicker.Stop()
	}()

	buf := make([]byte, 65536)

	reportTicker := time.NewTicker(time.Second * 5)
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

			reportFn()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-sendTicker.C:
		}

		now := time.Now()
		expiry := now.Add(time.Second * 1)

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
