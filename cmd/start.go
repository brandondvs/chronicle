package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
)

var (
	gtidFlag string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Connect to MySQL and start reading replication log events",
	RunE: func(cmd *cobra.Command, args []string) error {
		startCmdMain()
		return nil
	},
}

func init() {
	startCmd.Flags().StringVar(&gtidFlag, "gtid", "", "GTID is the position to start reading replication events")

	rootCmd.AddCommand(startCmd)
}

type TableSchema struct {
	TableID uint64 `json:"table_id"`
	Schema  string `json:"schema"`
	Name    string `json:"name"`

	Columns []*Column `json:"columns"`
}

type Column struct {
	Name     string `json:"name"`
	TypeName string `json:"type_name"`
}

func parseTableMapEvent(event *replication.TableMapEvent) *TableSchema {
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

	return schema
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

func createKafkaWriter() *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP("localhost:9092"),
		Topic:        "chronicle",
		Balancer:     &kafka.Hash{},
		Async:        true,
		RequiredAcks: kafka.RequireAll,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
}

func startCmdMain() {
	cfg := replication.BinlogSyncerConfig{
		ServerID: 1,
		Flavor:   "mysql",
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "chronicle",
		Password: "chroniclepw",
	}

	syncer := replication.NewBinlogSyncer(cfg)

	gtidSet, err := mysql.ParseGTIDSet(mysql.MySQLFlavor, gtidFlag)
	if err != nil {
		log.Fatalln("Failed parsing GTID set", err)
	}

	streamer, err := syncer.StartSyncGTID(gtidSet)
	if err != nil {
		log.Fatalln("Failed to start GTID streamer", err)
	}

	kafkaWriter := createKafkaWriter()
	defer kafkaWriter.Close()

	var tableSchemaCache map[uint64]*TableSchema = make(map[uint64]*TableSchema)

	for {
		ev, err := streamer.GetEvent(context.Background())
		if err != nil {
			fmt.Println("Failed to get event", err)
			continue
		}

		switch event := ev.Event.(type) {
		case *replication.TableMapEvent:
			schema := parseTableMapEvent(event)
			tableSchemaCache[schema.TableID] = schema

		case *replication.RowsEvent:
			schema, ok := tableSchemaCache[event.TableID]
			if !ok {
				fmt.Printf("TableID (%d) is missing from the table schema cache", event.TableID)
			}

			for _, row := range event.Rows {
				for i, column := range schema.Columns {
					data := fmt.Sprintf("%s (%s) -> %v\n", column.Name, column.TypeName, row[i])
					msg := kafka.Message{Value: []byte(data)}
					if err := kafkaWriter.WriteMessages(context.Background(), msg); err != nil {
						fmt.Println("Failed to write kafka message", err)
					}
				}
			}
		}
	}
}
