# loser

Something to help you track down network problems, sort of- only works on Linux, it's probably buggy.

## What does it do?

- Spins up an echo server on TCP 6943
- Spins up an echo server on TCP 6943
- Spins up a echo client and echo server for all IPs given on the commandline
- Spins up a Prometheus exporter on TCP 6942
  - Exposes interface counters (from `sysfs`) as well as some sent / received / out-of-order / lost metrics for the TCP
    and UDP streams

### TODO

- Change it up a bit so that its always recording metrics (rather than just while a connection is active)
    - Right now, you have to know that if `sent` isn't incrementing, then the connection is down

## Usage

```shell
# build
GOOS=linux GOARCH=amd64 go build -o loser

# transfer it to server 1 (192.168.100.101) and server 2 (192.168.100.102) and ...

# server 1
loser 192.168.100.101

# server 2
loser 192.168.100.102
```

Now you can hit the following:

- [http://192.168.100.101:6942/metrics](http://192.168.100.101:6942/metrics)
- [http://192.168.100.101:6943/metrics](http://192.168.100.101:6943/metrics)

You should have some metrics like this:

```
eth0_collisions 0
eth0_if_index 244188
eth0_mtu 85500
eth0_multicast 0
eth0_rx_bytes 4.145843806e+09
eth0_rx_compressed 0
eth0_rx_crc_errors 0
eth0_rx_dropped 0
eth0_rx_errors 0
eth0_rx_fifo_errors 0
eth0_rx_frame_errors 0
eth0_rx_length_errors 0
eth0_rx_missed_errors 0
eth0_rx_nohandler 0
eth0_rx_over_errors 0
eth0_rx_packets 2.5475855e+07
eth0_speed 570000
eth0_tx_aborted_errors 0
eth0_tx_bytes 1.708747998e+09
eth0_tx_carrier_errors 0
eth0_tx_compressed 0
eth0_tx_dropped 0
eth0_tx_errors 0
eth0_tx_fifo_errors 0
eth0_tx_heartbeat_errors 0
eth0_tx_packets 2.5319479e+07
eth0_tx_window_errors 0
go_gc_duration_seconds{quantile="0"} 6.5292e-05
go_gc_duration_seconds{quantile="0.25"} 7.3459e-05
go_gc_duration_seconds{quantile="0.5"} 0.000138835
go_gc_duration_seconds{quantile="0.75"} 0.00020946
go_gc_duration_seconds{quantile="1"} 0.000505255
go_gc_duration_seconds_sum 0.001324678
go_gc_duration_seconds_count 8
go_gc_gogc_percent 100
go_gc_gomemlimit_bytes 9.223372036854776e+18
go_goroutines 18
go_info{version="go1.23.2"} 1
go_memstats_alloc_bytes 1.813296e+06
go_memstats_alloc_bytes_total 2.1060264e+07
go_memstats_buck_hash_sys_bytes 4114
go_memstats_frees_total 84660
go_memstats_gc_sys_bytes 2.57984e+06
go_memstats_heap_alloc_bytes 1.813296e+06
go_memstats_heap_idle_bytes 3.39968e+06
go_memstats_heap_inuse_bytes 3.547136e+06
go_memstats_heap_objects 4604
go_memstats_heap_released_bytes 1.851392e+06
go_memstats_heap_sys_bytes 6.946816e+06
go_memstats_last_gc_time_seconds 1.7326379732145946e+09
go_memstats_mallocs_total 89264
go_memstats_mcache_inuse_bytes 7200
go_memstats_mcache_sys_bytes 15600
go_memstats_mspan_inuse_bytes 106080
go_memstats_mspan_sys_bytes 114240
go_memstats_next_gc_bytes 4.194304e+06
go_memstats_other_sys_bytes 1.438542e+06
go_memstats_stack_inuse_bytes 1.441792e+06
go_memstats_stack_sys_bytes 1.441792e+06
go_memstats_sys_bytes 1.2540944e+07
go_sched_gomaxprocs_threads 6
go_threads 10
lo_collisions 0
lo_if_index 57
lo_mtu 3.735552e+06
lo_multicast 0
lo_rx_bytes 0
lo_rx_compressed 0
lo_rx_crc_errors 0
lo_rx_dropped 0
lo_rx_errors 0
lo_rx_fifo_errors 0
lo_rx_frame_errors 0
lo_rx_length_errors 0
lo_rx_missed_errors 0
lo_rx_nohandler 0
lo_rx_over_errors 0
lo_rx_packets 0
lo_speed 0
lo_tx_aborted_errors 0
lo_tx_bytes 0
lo_tx_carrier_errors 0
lo_tx_compressed 0
lo_tx_dropped 0
lo_tx_errors 0
lo_tx_fifo_errors 0
lo_tx_heartbeat_errors 0
lo_tx_packets 0
lo_tx_window_errors 0
process_cpu_seconds_total 3.28
process_max_fds 1.048576e+06
process_network_receive_bytes_total 6.6979134e+07
process_network_transmit_bytes_total 2.4327134e+07
process_open_fds 12
process_resident_memory_bytes 1.4897152e+07
process_start_time_seconds 1.7326379198e+09
process_virtual_memory_bytes 1.869778944e+09
process_virtual_memory_max_bytes 1.8446744073709552e+19
promhttp_metric_handler_requests_in_flight 1
promhttp_metric_handler_requests_total{code="200"} 11
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
tcp_172_17_0_2_6943_lost 0
tcp_172_17_0_2_6943_out_of_order 0
tcp_172_17_0_2_6943_received 165241
tcp_172_17_0_2_6943_sent 165272
udp_172_17_0_2_6943_lost 0
udp_172_17_0_2_6943_out_of_order 0
udp_172_17_0_2_6943_received 165240
udp_172_17_0_2_6943_sent 165272
```
