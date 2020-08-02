# go-check-xid-age

## Description:
Helps protect againt Transaction ID Exhaustion (Wraparound) in PostgreSQL by alerting at given thresholds.
Performs Scan on all tables in a given database where Table Size is greater than `tableSize`.


## Usage:
```
  go-check-xid-age [OPTIONS]

Application Options:
  -H, --host=      Database server host or socket directory. (default: localhost)
  -U, --username=  Database user name. (default: postgres)
  -W, --password=  Password for the database user. (default: postgres)
  -d, --database=  Database name to connect to (default: postgres)
  -p, --port=      Port for database connection. (default: 5432)
  -S, --sslmode=   sslMode to use to connect to database. (default: disable)
  -t, --timeout=   Database connection timeout value. (default: 5)
  -s, --tablesize= (in bytes) Only find xid_age for tables greater than this size. (default: 10737418240)
  -l, --limit=     Limit the number of rows returned by the query. You can pass the value as 'ALL' to remove limits. (default: 10)
  -w, --warning=   Warning threshold for max xid_age.  (default: 190000000)
  -c, --critical=  Critical threshold for max xid_age. (default: 195000000)

Help Options:
  -h, --help       Show this help message
```


## Example: 
```
go-check-xid-age -H "localhost" -d "test" -U "postgres" -W ${PGPASSWORD} -w 18800 -c 18900
Warning: Table order_transition has max xid_age 18809
Warning: Table order_no_pool has max xid_age 18808
```
