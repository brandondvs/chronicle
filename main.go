package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

func main() {
	cfg := replication.BinlogSyncerConfig{
		ServerID: 1,
		Flavor:   "mysql",
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "chronicle",
		Password: "chroniclepw",
	}

	syncer := replication.NewBinlogSyncer(cfg)

	gtidSet, err := mysql.ParseGTIDSet(mysql.MySQLFlavor, "1ec671c5-4184-11f1-8f8c-0e8b9f5e7e44:1-15")
	if err != nil {
		log.Fatalln("Failed parsing GTID set", err)
	}

	streamer, err := syncer.StartSyncGTID(gtidSet)
	if err != nil {
		log.Fatalln("Failed to start GTID streamer", err)
	}

	for {
		ev, err := streamer.GetEvent(context.Background())
		if err != nil {
			fmt.Println("Failed to get event", err)
			continue
		}

		ev.Dump(os.Stdout)
	}
}
