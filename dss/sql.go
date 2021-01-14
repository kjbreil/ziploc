package dss

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	_ "github.com/denisenkom/go-mssqldb"
)

func SqlDSS() *DSS {
	query := url.Values{}
	query.Add("database", "DB_STORE")

	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword("RA", "bigBall5"),
		Host:   fmt.Sprintf("%s:%d", "NCBP", 1433),
		// Path:  instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}

	log.Println(u.String())

	db, err := sql.Open("sqlserver", u.String())

	if err != nil {
		log.Panic(err)
	}

	sqlStatement := `SELECT TOP 1 F2727,F2728,F2729,F2730,F2731,F2732,F253,F2733 FROM DSS_TAB`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println("Cannot query: ", err.Error())
		return nil
	}
	for rows.Next() {
		var t Table
		err = rows.Scan(&t.Priority, &t.Author, &t.Option, &t.Destination, &t.Script, &t.FileDate, &t.LastChangeDate, &t.Signature)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(t)
	}
	rows.Close()
	db.Close()

	return &DSS{
		Name:     "",
		SIL:      nil,
		Data:     nil,
		priority: 0,
	}
}
