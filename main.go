package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type TableSchema struct {
	TableID uint64
	Schema  string
	Name    string

	Columns []*Column
}

type Column struct {
	Name     string
	TypeName string
}

func parseTableMapEvent(event *replication.TableMapEvent) {
	schema := &TableSchema{
		TableID: event.TableID,
		Schema:  string(event.Schema),
		Name:    string(event.Table),

		Columns: make([]*Column, 0),
	}

	for i, typeCode := range event.ColumnType {
		column := &Column{
			Name:     string(event.ColumnName[i]),
			TypeName: mysqlToTypeName(typeCode),
		}
		schema.Columns = append(schema.Columns, column)
	}

	fmt.Printf("Database: %s, Table %s(ID=%d)\n", schema.Schema, schema.Name, schema.TableID)
	fmt.Println("Column Schema:")
	for _, column := range schema.Columns {
		fmt.Printf("%s - %s\n", column.Name, column.TypeName)
	}
}

func mysqlToTypeName(t byte) string {
	switch t {
	case mysql.MYSQL_TYPE_TINY:
		return "TINYINT"
	case mysql.MYSQL_TYPE_SHORT:
		return "SMALLINT"
	case mysql.MYSQL_TYPE_LONG:
		return "INT"
	case mysql.MYSQL_TYPE_LONGLONG:
		return "BIGINT"
	case mysql.MYSQL_TYPE_FLOAT:
		return "FLOAT"
	case mysql.MYSQL_TYPE_DOUBLE:
		return "DOUBLE"
	case mysql.MYSQL_TYPE_NEWDECIMAL:
		return "DECIMAL"
	case mysql.MYSQL_TYPE_VARCHAR, mysql.MYSQL_TYPE_VAR_STRING:
		return "VARCHAR"
	case mysql.MYSQL_TYPE_STRING:
		return "CHAR"
	case mysql.MYSQL_TYPE_BLOB:
		return "BLOB/TEXT"
	case mysql.MYSQL_TYPE_JSON:
		return "JSON"
	case mysql.MYSQL_TYPE_DATE:
		return "DATE"
	case mysql.MYSQL_TYPE_TIME, mysql.MYSQL_TYPE_TIME2:
		return "TIME"
	case mysql.MYSQL_TYPE_DATETIME, mysql.MYSQL_TYPE_DATETIME2:
		return "DATETIME"
	case mysql.MYSQL_TYPE_TIMESTAMP, mysql.MYSQL_TYPE_TIMESTAMP2:
		return "TIMESTAMP"
	case mysql.MYSQL_TYPE_ENUM:
		return "ENUM"
	case mysql.MYSQL_TYPE_SET:
		return "SET"
	case mysql.MYSQL_TYPE_BIT:
		return "BIT"
	case mysql.MYSQL_TYPE_YEAR:
		return "YEAR"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", t)
	}
}

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

		switch event := ev.Event.(type) {
		case *replication.TableMapEvent:
			parseTableMapEvent(event)
		}
	}
}
