stadis 
============
<sub>redis + statsd</sub>

Golang implementation of sending redis metrics to statsd

**inspiration from:** https://github.com/keenlabs/redis-statsd

```
NAME:
   stadis - get redis info and submit to statsd

USAGE:
   stadis [global options] command [command options] [arguments...]

VERSION:
   0.0.2

AUTHOR(S):

COMMANDS:
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --redis-host, -r "localhost:6379"	host:port of redis servier
   --statsd-host, -s "localhost:8125"	host:port of statsd servier
   --prefix, -p 			host:port of redis servier [$HOSTNAME]
   --interval, -i "5000"		time in milliseconds to periodically check redis
   --debug, -d				enable debug logging
   --help, -h				show help
   --version, -v			print the version
```

Example Output:
```
-------------------
GAUGES:
instantaneous_ops_per_sec: 944
latest_fork_usec: 0
WARN: migrate_cached_socketsis not displayed in redis info
pubsub_channels: 0
pubsub_patterns: 0
uptime_in_seconds: 4502605
blocked_clients: 0
connected_clients: 12
used_memory_peak: 2694777624
used_memory: 811414544
used_memory_lua: 31744
mem_fragmentation_ratio: 0
used_memory_rss: 679784448
-------------------
COUNTERS:
rejected_connections: 0
WARN: sync_fullis not displayed in redis info
WARN: sync_partial_okis not displayed in redis info
total_commands_processed: 2590146137
expired_keys: 159267474
keyspace_hits: 1610246628
WARN: sync_partial_erris not displayed in redis info
total_connections_received: 3841
evicted_keys: 0
keyspace_misses: 330335235
```

#PROPS

thanks to https://github.com/keenlabs/redis-statsd for inspiration
