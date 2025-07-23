package database;

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/Moody0101-X/Go_Api/models"
	"strings"
)

var DATABASE *sql.DB

func InitializeDb(dbPath string) (error, string) {
	
	var err error
	DATABASE, err = sql.Open("sqlite3", dbPath); if err != nil {
		return err, ""
	}
	return nil, dbPath
}

func isEmpty(s string) bool { return len(s) == 0 }

func GetNextUID(Table string) int {

	var id int;
	
	row, err := DATABASE.Query("select MAX(ID) from " + Table);
	defer row.Close()
	if err != nil { return 0 }
	for row.Next() {
		row.Scan(&id);
	}
	return id + 1;
}

func GetNewPostID() int {
	return GetLastID("Posts");
}

func GetLastID(Table string) int {

	var id int;
	
	row, err := DATABASE.Query("select MAX(ID) from " + Table);
	
	defer row.Close()

	if err != nil { return 0 }

	for row.Next() {
		row.Scan(&id);
	}

	return id;
}

func CheckCdnLink(link string) string {
	return strings.ReplaceAll(link, "http://localhost:8500", models.CDN_API);
}
