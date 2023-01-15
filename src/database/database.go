package database;

import (
	// "fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	// "sync"
)

// struct DB type {
// 	sync.RWMutex
// 	dataBase *sql.DB
// }


var dataBase *sql.DB

// db initializer: Opens the db, then evluates a global conn variable.
func InitializeDb(dbPath string) (error, string) {
	
	var err error

	dataBase, err = sql.Open("sqlite3", dbPath); if err != nil {
		return err, ""
	}

	return nil, dbPath
}

func isEmpty(s string) bool { return len(s) == 0 }

func GetNextUID(Table string) int {

	var id int;
	row, err := dataBase.Query("select MAX(ID) from " + Table);
	
	defer row.Close()

	if err != nil { return 0 }

	for row.Next() {
		row.Scan(&id);
	}

	return id + 1;
}
