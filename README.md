# Chronicle

Chronicle is a Go binary that connects to MySQL replication and sends the replication events to Kafka. It is intented to be a simple and lightweight Change Data Capture (CDC) pipeline.


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
