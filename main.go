package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/initialed85/loser/pkg/network_interfaces"
	"github.com/initialed85/loser/pkg/packets"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.Printf("starting loser...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//
	// network interface ticker and handler
	//

	networkInterfacesBodyMu := new(sync.Mutex)
	networkInterfacesBody := []byte("{}")
	networkInterfacesTicker := time.NewTicker(time.Second * 1)

	countersMu := new(sync.Mutex)
	counters := make(map[string]prometheus.Counter)

	go func() {
		log.Printf("starting network interface ticker...")

		lastNetworkInterfaces := make(map[string]network_interfaces.NetworkInterface)

		for {
			select {
			case <-ctx.Done():
				return
			case <-networkInterfacesTicker.C:
			}

			err := func() error {
				rawNetworkInterfaces, err := network_interfaces.GetNetworkInterfaces()
				if err != nil {
					return fmt.Errorf("warning: failed GetNetworkInterfaces(): %s", err)
				}

				body, err := json.MarshalIndent(rawNetworkInterfaces, "", "  ")
				if err != nil {
					return fmt.Errorf("warning: failed json.Marshal() for networkInterfaces: %s", err)
				}

				networkInterfacesBodyMu.Lock()
				networkInterfacesBody = body
				networkInterfacesBodyMu.Unlock()

				networkInterfaces := make(map[string]network_interfaces.NetworkInterface)

				for _, networkInterface := range rawNetworkInterfaces {
					networkInterfaces[networkInterface.Name] = networkInterface
				}

				// handle new interfaces we've only just seen
				for networkInterfaceName := range networkInterfaces {
					_, exists := lastNetworkInterfaces[networkInterfaceName]
					if !exists {
						log.Printf("adding interface %s...", networkInterfaceName)

						countersMu.Lock()
						counters[fmt.Sprintf("%s_IFIndex", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_if_index", networkInterfaceName)})
						counters[fmt.Sprintf("%s_MTU", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_mtu", networkInterfaceName)})
						counters[fmt.Sprintf("%s_Speed", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_speed", networkInterfaceName)})
						counters[fmt.Sprintf("%s_Collisions", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_collisions", networkInterfaceName)})
						counters[fmt.Sprintf("%s_Multicast", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_multicast", networkInterfaceName)})
						counters[fmt.Sprintf("%s_RxBytes", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_rx_bytes", networkInterfaceName)})
						counters[fmt.Sprintf("%s_RxCompressed", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_rx_compressed", networkInterfaceName)})
						counters[fmt.Sprintf("%s_RxCrcErrors", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_rx_crc_errors", networkInterfaceName)})
						counters[fmt.Sprintf("%s_RxDropped", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_rx_dropped", networkInterfaceName)})
						counters[fmt.Sprintf("%s_RxErrors", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_rx_errors", networkInterfaceName)})
						counters[fmt.Sprintf("%s_RxFifoErrors", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_rx_fifo_errors", networkInterfaceName)})
						counters[fmt.Sprintf("%s_RxFrameErrors", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_rx_frame_errors", networkInterfaceName)})
						counters[fmt.Sprintf("%s_RxLengthErrors", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_rx_length_errors", networkInterfaceName)})
						counters[fmt.Sprintf("%s_RxMissedErrors", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_rx_missed_errors", networkInterfaceName)})
						counters[fmt.Sprintf("%s_RxNohandler", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_rx_nohandler", networkInterfaceName)})
						counters[fmt.Sprintf("%s_RxOverErrors", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_rx_over_errors", networkInterfaceName)})
						counters[fmt.Sprintf("%s_RxPackets", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_rx_packets", networkInterfaceName)})
						counters[fmt.Sprintf("%s_TxAbortedErrors", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_tx_aborted_errors", networkInterfaceName)})
						counters[fmt.Sprintf("%s_TxBytes", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_tx_bytes", networkInterfaceName)})
						counters[fmt.Sprintf("%s_TxCarrierErrors", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_tx_carrier_errors", networkInterfaceName)})
						counters[fmt.Sprintf("%s_TxCompressed", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_tx_compressed", networkInterfaceName)})
						counters[fmt.Sprintf("%s_TxDropped", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_tx_dropped", networkInterfaceName)})
						counters[fmt.Sprintf("%s_TxErrors", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_tx_errors", networkInterfaceName)})
						counters[fmt.Sprintf("%s_TxFifoErrors", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_tx_fifo_errors", networkInterfaceName)})
						counters[fmt.Sprintf("%s_TxHeartbeatErrors", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_tx_heartbeat_errors", networkInterfaceName)})
						counters[fmt.Sprintf("%s_TxPackets", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_tx_packets", networkInterfaceName)})
						counters[fmt.Sprintf("%s_TxWindowErrors", networkInterfaceName)] = promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("%s_tx_window_errors", networkInterfaceName)})
						countersMu.Unlock()
					}
				}

				// handle the actual metrics
				countersMu.Lock()
				for networkInterfaceName, networkInterface := range networkInterfaces {
					counters[fmt.Sprintf("%s_IFIndex", networkInterfaceName)].Add(float64(networkInterface.IFIndex))
					counters[fmt.Sprintf("%s_MTU", networkInterfaceName)].Add(float64(networkInterface.MTU))
					counters[fmt.Sprintf("%s_Speed", networkInterfaceName)].Add(float64(networkInterface.Speed))
					counters[fmt.Sprintf("%s_Collisions", networkInterfaceName)].Add(float64(networkInterface.Collisions))
					counters[fmt.Sprintf("%s_Multicast", networkInterfaceName)].Add(float64(networkInterface.Multicast))
					counters[fmt.Sprintf("%s_RxBytes", networkInterfaceName)].Add(float64(networkInterface.RxBytes))
					counters[fmt.Sprintf("%s_RxCompressed", networkInterfaceName)].Add(float64(networkInterface.RxCompressed))
					counters[fmt.Sprintf("%s_RxCrcErrors", networkInterfaceName)].Add(float64(networkInterface.RxCrcErrors))
					counters[fmt.Sprintf("%s_RxDropped", networkInterfaceName)].Add(float64(networkInterface.RxDropped))
					counters[fmt.Sprintf("%s_RxErrors", networkInterfaceName)].Add(float64(networkInterface.RxErrors))
					counters[fmt.Sprintf("%s_RxFifoErrors", networkInterfaceName)].Add(float64(networkInterface.RxFifoErrors))
					counters[fmt.Sprintf("%s_RxFrameErrors", networkInterfaceName)].Add(float64(networkInterface.RxFrameErrors))
					counters[fmt.Sprintf("%s_RxLengthErrors", networkInterfaceName)].Add(float64(networkInterface.RxLengthErrors))
					counters[fmt.Sprintf("%s_RxMissedErrors", networkInterfaceName)].Add(float64(networkInterface.RxMissedErrors))
					counters[fmt.Sprintf("%s_RxNohandler", networkInterfaceName)].Add(float64(networkInterface.RxNohandler))
					counters[fmt.Sprintf("%s_RxOverErrors", networkInterfaceName)].Add(float64(networkInterface.RxOverErrors))
					counters[fmt.Sprintf("%s_RxPackets", networkInterfaceName)].Add(float64(networkInterface.RxPackets))
					counters[fmt.Sprintf("%s_TxAbortedErrors", networkInterfaceName)].Add(float64(networkInterface.TxAbortedErrors))
					counters[fmt.Sprintf("%s_TxBytes", networkInterfaceName)].Add(float64(networkInterface.TxBytes))
					counters[fmt.Sprintf("%s_TxCarrierErrors", networkInterfaceName)].Add(float64(networkInterface.TxCarrierErrors))
					counters[fmt.Sprintf("%s_TxCompressed", networkInterfaceName)].Add(float64(networkInterface.TxCompressed))
					counters[fmt.Sprintf("%s_TxDropped", networkInterfaceName)].Add(float64(networkInterface.TxDropped))
					counters[fmt.Sprintf("%s_TxErrors", networkInterfaceName)].Add(float64(networkInterface.TxErrors))
					counters[fmt.Sprintf("%s_TxFifoErrors", networkInterfaceName)].Add(float64(networkInterface.TxFifoErrors))
					counters[fmt.Sprintf("%s_TxHeartbeatErrors", networkInterfaceName)].Add(float64(networkInterface.TxHeartbeatErrors))
					counters[fmt.Sprintf("%s_TxPackets", networkInterfaceName)].Add(float64(networkInterface.TxPackets))
					counters[fmt.Sprintf("%s_TxWindowErrors", networkInterfaceName)].Add(float64(networkInterface.TxWindowErrors))
				}
				countersMu.Unlock()

				// handle old interfaces we're no longer seeing
				for lastNetworkInterfaceName := range lastNetworkInterfaces {
					_, exists := networkInterfaces[lastNetworkInterfaceName]
					if !exists {
						log.Printf("removing interface %s...", lastNetworkInterfaceName)

						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_IFIndex", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_MTU", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_Speed", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_Collisions", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_Multicast", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_RxBytes", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_RxCompressed", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_RxCrcErrors", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_RxDropped", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_RxErrors", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_RxFifoErrors", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_RxFrameErrors", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_RxLengthErrors", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_RxMissedErrors", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_RxNohandler", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_RxOverErrors", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_RxPackets", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_TxAbortedErrors", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_TxBytes", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_TxCarrierErrors", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_TxCompressed", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_TxDropped", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_TxErrors", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_TxFifoErrors", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_TxHeartbeatErrors", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_TxPackets", lastNetworkInterfaceName)])
						_ = prometheus.DefaultRegisterer.Unregister(counters[fmt.Sprintf("%s_TxWindowErrors", lastNetworkInterfaceName)])

						countersMu.Lock()
						delete(counters, fmt.Sprintf("%s_IFIndex", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_MTU", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_Speed", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_Collisions", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_Multicast", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_RxBytes", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_RxCompressed", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_RxCrcErrors", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_RxDropped", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_RxErrors", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_RxFifoErrors", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_RxFrameErrors", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_RxLengthErrors", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_RxMissedErrors", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_RxNohandler", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_RxOverErrors", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_RxPackets", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_TxAbortedErrors", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_TxBytes", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_TxCarrierErrors", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_TxCompressed", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_TxDropped", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_TxErrors", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_TxFifoErrors", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_TxHeartbeatErrors", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_TxPackets", lastNetworkInterfaceName))
						delete(counters, fmt.Sprintf("%s_TxWindowErrors", lastNetworkInterfaceName))
						countersMu.Unlock()
					}
				}

				lastNetworkInterfaces = networkInterfaces

				return nil
			}()
			if err != nil {
				log.Fatal(err)
				continue
			}
		}
	}()

	log.Printf("registering /network-interfaces endpoint")
	http.Handle("/network-interfaces", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		networkInterfacesBodyMu.Lock()
		body := networkInterfacesBody
		networkInterfacesBodyMu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))

	//
	// tcp server
	//

	go func() {
		err := packets.RunTCPServer(ctx, 6943)
		if err != nil {
			log.Fatal(err)
		}
	}()

	//
	// tcp clients
	//

	for _, rawDialAddr := range os.Args[1:] {
		go func() {
			friendlyRawDialAddr := strings.ReplaceAll(rawDialAddr, ".", "_")
			friendlyRawDialAddr = strings.ReplaceAll(friendlyRawDialAddr, ":", "_")

			sentCounter := promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("tcp_%s_sent", friendlyRawDialAddr)})
			receivedCounter := promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("tcp_%s_received", friendlyRawDialAddr)})
			outOfOrderCounter := promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("tcp_%s_out_of_order", friendlyRawDialAddr)})
			lostCounter := promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("tcp_%s_lost", friendlyRawDialAddr)})

			reportFn := func(now time.Time, sent int64, received int64, outOfOrder int64, lost int64) {
				sentCounter.Add(float64(sent))
				receivedCounter.Add(float64(received))
				outOfOrderCounter.Add(float64(outOfOrder))
				lostCounter.Add(float64(lost))
			}

			for {
				err := packets.RunTCPClient(ctx, rawDialAddr, reportFn)
				if err != nil {
					log.Printf("warning: failed packets.RunTCPClient: %s", err)
					time.Sleep(time.Second * 1)
				}
			}
		}()
	}

	//
	// udp server
	//

	go func() {
		err := packets.RunUDPServer(ctx, 6943)
		if err != nil {
			log.Fatal(err)
		}
	}()

	//
	// udp clients
	//

	for _, rawDialAddr := range os.Args[1:] {
		go func() {
			friendlyRawDialAddr := strings.ReplaceAll(rawDialAddr, ".", "_")
			friendlyRawDialAddr = strings.ReplaceAll(friendlyRawDialAddr, ":", "_")

			sentCounter := promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("udp_%s_sent", friendlyRawDialAddr)})
			receivedCounter := promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("udp_%s_received", friendlyRawDialAddr)})
			outOfOrderCounter := promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("udp_%s_out_of_order", friendlyRawDialAddr)})
			lostCounter := promauto.NewCounter(prometheus.CounterOpts{Name: fmt.Sprintf("udp_%s_lost", friendlyRawDialAddr)})

			reportFn := func(now time.Time, sent int64, received int64, outOfOrder int64, lost int64) {
				sentCounter.Add(float64(sent))
				receivedCounter.Add(float64(received))
				outOfOrderCounter.Add(float64(outOfOrder))
				lostCounter.Add(float64(lost))
			}

			for {
				err := packets.RunUDPClient(ctx, rawDialAddr, reportFn)
				if err != nil {
					log.Printf("warning: failed packets.RunUDPClient: %s", err)
					time.Sleep(time.Second * 1)
				}
			}
		}()
	}

	//
	// general stuff
	//

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":6942", nil)
}
