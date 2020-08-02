package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/jessevdk/go-flags"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
)

type options struct {

	// DB Connections Details
	Host     string `long:"host" short:"H" default:"localhost" description:"Database server host or socket directory."`
	User     string `long:"username" short:"U" default:"postgres" description:"Database user name."`
	Password string `long:"password" short:"W" default:"postgres" description:"Password for the database user."`
	Database string `long:"database" short:"d" default:"postgres" description:"Database name to connect to"`
	Port     int    `long:"port" short:"p" default:"5432" description:"Port for database connection."`
	SslMode  string `long:"sslmode" short:"S" default:"disable" description:"sslMode to use to connect to database."`
	Timeout  int    `long:"timeout" short:"t" default:"5" description:"Database connection timeout value."`

	// Query Parameters
	TableSize int    `long:"tablesize" short:"s" default:"10737418240" description:"(in bytes) Only find xid_age for tables greater than this size."`
	Limit     string `long:"limit" short:"l" default:"10" description:"Limit the number of rows returned by the query. You can pass the value as 'ALL' to remove limits."`

	// Thresholds
	Warning  int `long:"warning" short:"w" default:"190000000" description:"Warning threshold for max xid_age. "`
	Critical int `long:"critical" short:"c" default:"195000000" description:"Critical threshold for max xid_age."`
}

func printEntries(entries []string) {
	for _, v := range entries {
		fmt.Println(v)
	}
	return
}

func main() {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		opts.Host, opts.Port, opts.User, opts.Password, opts.Database, opts.SslMode, opts.Timeout)
	conn, err := sql.Open(dbDriver, dbinfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database. %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	xidAgeQuery := fmt.Sprintf(`
			SELECT relname as tablename, age(relfrozenxid) as xid_age
			FROM pg_class
			WHERE relkind = 'r' and pg_table_size(oid) > %d
			ORDER BY age(relfrozenxid)
			DESC LIMIT %s;
			`, opts.TableSize, opts.Limit)

	rows, err := conn.Query(xidAgeQuery)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to query results. %v\n", err)
		os.Exit(1)
	}

	var warningEntries []string
	var criticalEntries []string
	var readableLine string

	defer rows.Close()
	for rows.Next() {
		var tableName string
		var xidAge int
		err = rows.Scan(&tableName, &xidAge)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed while scanning rows. %v\n", err)
			os.Exit(1)
		}

		if xidAge >= opts.Critical {
			readableLine = "Critical: Table " + tableName + " has max xid_age " + strconv.Itoa(xidAge)
			criticalEntries = append(criticalEntries, readableLine)
			continue
		}

		if xidAge >= opts.Warning {
			readableLine = "Warning: Table " + tableName + " has max xid_age " + strconv.Itoa(xidAge)
			warningEntries = append(warningEntries, readableLine)
		}

	}

	if len(criticalEntries) > 0 {
		printEntries(criticalEntries)
		if len(warningEntries) > 0 {
			printEntries(warningEntries)
		}
		os.Exit(2)
	}

	if len(warningEntries) > 0 {
		printEntries(warningEntries)
		os.Exit(1)
	}

	fmt.Println("All tables have xid_age below threshold")

}
