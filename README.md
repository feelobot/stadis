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
   0.0.0

COMMANDS:
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --redis-host, -r "localhost:6379"	host:port of redis servier
   --statsd-host, -s "localhost:8125"	host:port of statsd servier
   --prefix, -p 			host:port of redis servier [$HOSTNAME]
   --help, -h				show help
   --version, -v			print the version
```
