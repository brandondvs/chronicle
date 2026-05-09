# Chronicle

Chronicle is a Go binary that connects to MySQL replication and sends the replication events to Kafka. It is intended to be a simple and lightweight Change Data Capture (CDC) pipeline.


# Development

Ensure you have Docker installed and run:

```
docker compose up -d
```

Get GTID from the MySQL instance

```
docker exec -it chronicle-mysql /bin/bash
mysql -u root -p
```

Once logged into the MySQL instance, run:
```
SHOW BINARY LOG STATUS
```

Note the value stored in the column: `Executed_Gtid_Set`.


```
Chronicle is simple CDC application for sending MySQL events to Kafka

Usage:
  chronicle [flags]
  chronicle [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  start       Connect to MySQL and start reading replication log events

Flags:
  -h, --help   help for chronicle

Use "chronicle [command] --help" for more information about a command.
```
