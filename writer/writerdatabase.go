package writer

import (
	"database/sql"
	"fmt"
	"log"
	"recordfilelist/raicredentials"
	"strings"
	"time"
)

// database writer
type DatabaseWriter struct {
	database    *sql.DB
	files       []ReferencedFile
	DbName      string
	DbTable     string
	Credentials string
}

func (db *DatabaseWriter) WriteInformation(file ReferencedFile) {
	db.files = append(db.files, file)
}
func (db *DatabaseWriter) NewDirectory(directory ReferencedDirectory) {
	db.files = nil
}
func (db *DatabaseWriter) EndDirectory(directory ReferencedDirectory) {
	sqlParts := []string{}
	vals := []any{}

	if len(db.files) == 0 {
		return
	}
	for _, file := range db.files {
		sqlParts = append(sqlParts, "(?,?,?,?,?)")
		vals = append(vals, file.Disk, file.Directory, file.FileName, file.FileSize, time.Now().UTC())
	}

	sqlQuery := fmt.Sprintf(
		"INSERT INTO `%s`.`%s` (`disk`,`directory`,`file`,`size`,`created`) VALUES %s",
		db.DbName, db.DbTable, strings.Join(sqlParts, ","))
	fmt.Println(sqlQuery)
	result, err := db.database.Exec(sqlQuery, vals...)
	if err == nil {
		rows, _ := result.RowsAffected()
		fmt.Printf("records added: %d \n", rows)
	} else {
		fmt.Println(err)
	}
}
func (db *DatabaseWriter) InitiateStore(disk string) {
	credentials, err := raicredentials.Credentials(db.Credentials)
	if err != nil {
		panic(err)
	}
	replacer := strings.NewReplacer(
		"{host}", credentials.Host,
		"{port}", credentials.Port,
		"{user}", credentials.User,
		"{password}", credentials.Password,
	)
	dbConnection, err := sql.Open("mysql", replacer.Replace("{user}:{password}@tcp({host}:{port})/forgogo?parseTime=true"))
	if err != nil {
		log.Fatal(err)
	}
	db.database = dbConnection
	db.database.Exec(fmt.Sprintf("DELETE FROM `%s`.`%s` WHERE `disk`=?", db.DbName, db.DbTable), disk)
}
func (db *DatabaseWriter) CloseStore() {
	db.database.Close()
}

// query := "SELECT `id`,`file`,`size`,`created` FROM `forgogo`.`aenlever` ORDER BY `id` DESC"
// result, _ := dbConnection.Query(query)
// defer result.Close()

// type record struct {
// 	id      int
// 	file    string
// 	size    int
// 	created time.Time
// }
// for result.Next() {
// 	var record record
// 	_ = result.Scan(&record.id, &record.file, &record.size, &record.created)
// 	fmt.Printf("%d - %s - %d (%v) \n", record.id, record.file, record.size, record.created)
// }
