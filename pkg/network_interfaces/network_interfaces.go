package network_interfaces

import (
	"fmt"
	_log "log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var log = _log.New(
	os.Stdout,
	"loser: ",
	_log.Ldate|
		_log.Ltime|
		_log.Lmicroseconds|
		_log.LUTC|
		_log.Lmsgprefix|
		_log.LstdFlags,
)

type NetworkInterface struct {
	Timestamp         time.Time `json:"timestamp"`
	Name              string    `json:"name"`
	MAC               string    `json:"mac"`
	IFIndex           int       `json:"if_index"`
	MTU               int       `json:"mtu"`
	Speed             int       `json:"speed"`
	Collisions        int64     `json:"collisions"`
	Multicast         int64     `json:"multicast"`
	RxBytes           int64     `json:"rx_bytes"`
	RxCompressed      int64     `json:"rx_compressed"`
	RxCrcErrors       int64     `json:"rx_crc_errors"`
	RxDropped         int64     `json:"rx_dropped"`
	RxErrors          int64     `json:"rx_errors"`
	RxFifoErrors      int64     `json:"rx_fifo_errors"`
	RxFrameErrors     int64     `json:"rx_frame_errors"`
	RxLengthErrors    int64     `json:"rx_length_errors"`
	RxMissedErrors    int64     `json:"rx_missed_errors"`
	RxNohandler       int64     `json:"rx_nohandler"`
	RxOverErrors      int64     `json:"rx_over_errors"`
	RxPackets         int64     `json:"rx_packets"`
	TxAbortedErrors   int64     `json:"tx_aborted_errors"`
	TxBytes           int64     `json:"tx_bytes"`
	TxCarrierErrors   int64     `json:"tx_carrier_errors"`
	TxCompressed      int64     `json:"tx_compressed"`
	TxDropped         int64     `json:"tx_dropped"`
	TxErrors          int64     `json:"tx_errors"`
	TxFifoErrors      int64     `json:"tx_fifo_errors"`
	TxHeartbeatErrors int64     `json:"tx_heartbeat_errors"`
	TxPackets         int64     `json:"tx_packets"`
	TxWindowErrors    int64     `json:"tx_window_errors"`
}

var relevantSysClassNetItems = []string{
	"address",
	"ifindex",
	"mtu",
	// TODO: speed not reliable on all drivers it seems
	// "speed",
}

func GetNetworkInterfaces() ([]NetworkInterface, error) {
	sysClassNetDirEntries, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return nil, fmt.Errorf("failed os.ReadDir for sysClassNetDirEntries: %s", err)
	}

	networkInterfaces := make([]NetworkInterface, 0)

	for _, sysClassNetDirEntry := range sysClassNetDirEntries {
		now := time.Now()

		statsDirEntries, err := os.ReadDir(filepath.Join("/sys/class/net", sysClassNetDirEntry.Name(), "statistics"))
		if err != nil {
			// TODO: not everything has stats apparently
			// return nil, fmt.Errorf("failed os.ReadDir for statsDirEntries: %s", err)

			// TODO: keeping the noise down for this common and unimportant failure
			// log.Printf("warning: %s", fmt.Errorf("failed os.ReadDir for statsDirEntries: %s", err))
			continue
		}

		items := make(map[string]string)

		for _, relevantSysClassNetItem := range relevantSysClassNetItems {
			// attempts to get speed for lo return an "invalid argument" error; so let's not do that
			if relevantSysClassNetItem == "speed" && sysClassNetDirEntry.Name() == "lo" {
				items[relevantSysClassNetItem] = "0"
				continue
			}

			itemRaw, err := os.ReadFile(filepath.Join("/sys/class/net", sysClassNetDirEntry.Name(), relevantSysClassNetItem))
			if err != nil {
				return nil, fmt.Errorf("failed os.ReadFile for itemRaw: %#+v: %s", itemRaw, err)
			}

			items[relevantSysClassNetItem] = strings.TrimSpace(string(itemRaw))
		}

		ifIndex, err := strconv.ParseInt(items["ifindex"], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed strconv.ParseInt for ifIndex: %#+v: %s", items["ifindex"], err)
		}

		mtu, err := strconv.ParseInt(items["mtu"], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed strconv.ParseInt for mtu: %#+v: %s", items["mtu"], err)
		}

		// TODO: speed not reliable on all drivers it seems
		// speed, err := strconv.ParseInt(items["speed"], 10, 64)
		// if err != nil {
		// 	return nil, fmt.Errorf("failed strconv.ParseInt for speed: %#+v: %s", items["speed"], err)
		// }

		networkInterface := NetworkInterface{
			Timestamp: now,
			Name:      sysClassNetDirEntry.Name(),
			MAC:       items["address"],
			IFIndex:   int(ifIndex),
			MTU:       int(mtu),
			// TODO: speed not reliable on all drivers it seems
			// Speed:     int(speed),
		}

		stats := make(map[string]int64)

		for _, statsDirEntry := range statsDirEntries {
			statRaw, err := os.ReadFile(filepath.Join("/sys/class/net", sysClassNetDirEntry.Name(), "statistics", statsDirEntry.Name()))
			if err != nil {
				return nil, err
			}

			stat, err := strconv.ParseInt(strings.TrimSpace(string(statRaw)), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed strconv.ParseInt for %v: %#+v: %s", statsDirEntry, string(statRaw), err)
			}

			stats[statsDirEntry.Name()] = stat
		}

		networkInterface.Collisions = stats["collisions"]
		networkInterface.Multicast = stats["multicast"]
		networkInterface.RxBytes = stats["rx_bytes"]
		networkInterface.RxCompressed = stats["rx_compressed"]
		networkInterface.RxCrcErrors = stats["rx_crc_errors"]
		networkInterface.RxDropped = stats["rx_dropped"]
		networkInterface.RxErrors = stats["rx_errors"]
		networkInterface.RxFifoErrors = stats["rx_fifo_errors"]
		networkInterface.RxFrameErrors = stats["rx_frame_errors"]
		networkInterface.RxLengthErrors = stats["rx_length_errors"]
		networkInterface.RxMissedErrors = stats["rx_missed_errors"]
		networkInterface.RxNohandler = stats["rx_nohandler"]
		networkInterface.RxOverErrors = stats["rx_over_errors"]
		networkInterface.RxPackets = stats["rx_packets"]
		networkInterface.TxAbortedErrors = stats["tx_aborted_errors"]
		networkInterface.TxBytes = stats["tx_bytes"]
		networkInterface.TxCarrierErrors = stats["tx_carrier_errors"]
		networkInterface.TxCompressed = stats["tx_compressed"]
		networkInterface.TxDropped = stats["tx_dropped"]
		networkInterface.TxErrors = stats["tx_errors"]
		networkInterface.TxFifoErrors = stats["tx_fifo_errors"]
		networkInterface.TxHeartbeatErrors = stats["tx_heartbeat_errors"]
		networkInterface.TxPackets = stats["tx_packets"]
		networkInterface.TxWindowErrors = stats["tx_window_errors"]

		networkInterfaces = append(networkInterfaces, networkInterface)
	}

	return networkInterfaces, nil
}
