package main

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func databaseError(err error) {
	if err != nil {
		serverLogger("Database error", err.Error(), "Error")
	}
}

func initString() string {
	dbString := dbConf.username + ":" + dbConf.password + "@tcp(" + dbConf.addr + ":" + dbConf.port + ")/hive?charset=utf8&parseTime=true"
	return dbString
}

func checkToken(cliToken string) (string, string) {
	// Make connect string
	dbString := initString()
	// Connect
	db, err := sql.Open("mysql", dbString)
	databaseError(err)

	// Query
	res, err := db.Query("SELECT * FROM token")
	databaseError(err)

	db.Close()

	// Extract data
	for res.Next() {
		var tag string
		var admin string
		var token string
		var timestamp string
		err := res.Scan(&tag, &admin, &token, &timestamp)
		databaseError(err)
		// Compare
		if token == cliToken {
			return admin, tag
		}
	}

	// Invalid token
	return "", ""
}

func checkTagDuplicate(cliTag string) bool {
	// Make connect string
	dbString := initString()
	// Connect
	db, err := sql.Open("mysql", dbString)
	databaseError(err)

	// Query
	res, err := db.Query("SELECT * FROM token")
	databaseError(err)

	db.Close()

	// Extract data
	for res.Next() {
		var tag string
		var admin string
		var token string
		var timestamp string
		err := res.Scan(&tag, &admin, &token, &timestamp)
		databaseError(err)
		// Compare
		if tag == cliTag {
			return true
		}
	}

	// No duplicate tag
	return false
}

func checkAdminDuplicate(newAdmin string) bool {
	// Make connect string
	dbString := initString()
	// Connect
	db, err := sql.Open("mysql", dbString)
	databaseError(err)

	// Query
	res, err := db.Query("SELECT * FROM admin")
	databaseError(err)

	db.Close()

	// Extract data
	for res.Next() {
		var admin string
		var token string
		err := res.Scan(&admin, &token)
		databaseError(err)
		// Compare
		if admin == newAdmin {
			return true
		}
	}

	// No duplicate tag
	return false
}

func storeMessage(admin string, tag string, content string) bool {
	// Make connect string
	dbString := initString()
	// Connect
	db, err := sql.Open("mysql", dbString)
	databaseError(err)
	// Statement
	stmt, err := db.Prepare("INSERT message SET tag=?, admin=?, content=?, timestamp=?")
	databaseError(err)

	// Insert
	res, err := stmt.Exec(tag, admin, content, time.Now().Format("2006-01-02 15:04:05"))
	databaseError(err)

	db.Close()

	// Validate
	num, err := res.RowsAffected()
	databaseError(err)

	if num == 1 {
		return true
	}
	return false
}

func fetchMessages(admin string) []message {
	// Make connect string
	dbString := initString()
	// Connect
	db, err := sql.Open("mysql", dbString)
	databaseError(err)

	// Statement
	stmt, err := db.Prepare("SELECT * FROM message WHERE admin=? ORDER BY timestamp DESC")
	databaseError(err)

	// Query
	res, err := stmt.Query(admin)
	databaseError(err)

	db.Close()

	// Data array
	var messages []message
	// Extract data
	for res.Next() {
		var id int
		var tag string
		var admin string
		var content string
		var timestamp string
		err := res.Scan(&id, &tag, &admin, &content, &timestamp)
		databaseError(err)
		// Append data
		var m message
		m.ID = id
		m.Tag = tag
		m.Admin = admin
		m.Content = content
		m.Timestamp = timestamp
		messages = append(messages, m)
	}
	return messages
}

func storeToken(tag string, token string) bool {
	admin := cliConf.admin
	// Make connect string
	dbString := initString()
	// Connect
	db, err := sql.Open("mysql", dbString)
	databaseError(err)

	// Statement
	stmt, err := db.Prepare("INSERT token SET tag=?, token=?, admin=?, timestamp=?")
	databaseError(err)

	// Insert
	res, err := stmt.Exec(tag, token, admin, time.Now().Format("2006-01-02 15:04:05"))
	databaseError(err)

	db.Close()

	// Validate
	num, err := res.RowsAffected()
	databaseError(err)

	if num == 1 {
		return true
	}
	return false
}

func delToken(cliTag string) bool {
	admin := cliConf.admin
	// Make connect string
	dbString := initString()
	// Connect
	db, err := sql.Open("mysql", dbString)
	databaseError(err)

	// Statement
	stmt, err := db.Prepare("DELETE FROM token WHERE tag=? AND admin=?")
	databaseError(err)

	// Insert
	res, err := stmt.Exec(cliTag, admin)
	databaseError(err)

	db.Close()

	// Validate
	num, err := res.RowsAffected()
	databaseError(err)

	if num == 1 {
		return true
	}
	return false
}

func fetchToken() []string {
	admin := cliConf.admin
	// Make connect string
	dbString := initString()
	// Connect
	db, err := sql.Open("mysql", dbString)
	databaseError(err)

	// Statement
	stmt, err := db.Prepare("SELECT * FROM token WHERE admin=?")
	databaseError(err)

	// Query
	res, err := stmt.Query(admin)
	databaseError(err)

	db.Close()

	// Data array
	var tagNtoken []string
	// Extract data
	for res.Next() {
		var tag string
		var admin string
		var token string
		var timestamp string
		err := res.Scan(&tag, &admin, &token, &timestamp)
		databaseError(err)
		// Append data
		str := tag + ":" + token
		tagNtoken = append(tagNtoken, str)
	}
	return tagNtoken
}

func storeAdmin(admin string, token string) bool {
	// Make connect string
	dbString := initString()
	// Connect
	db, err := sql.Open("mysql", dbString)
	databaseError(err)

	// Statement
	stmt, err := db.Prepare("INSERT admin SET name=?, token=?, timestamp=?")
	databaseError(err)

	// Insert
	res, err := stmt.Exec(admin, token, time.Now().Format("2006-01-02 15:04:05"))
	databaseError(err)

	db.Close()

	// Validate
	num, err := res.RowsAffected()
	databaseError(err)

	if num == 1 {
		return true
	}
	return false
}

func fetchAdmin(token string) (string, bool) {
	// Make connect string
	dbString := initString()
	// Connect
	db, err := sql.Open("mysql", dbString)
	databaseError(err)

	// Statement
	stmt, err := db.Prepare("SELECT * FROM admin WHERE token=?")
	databaseError(err)

	// Query
	var name string
	var dbToken string
	var timestamp string
	err = stmt.QueryRow(token).Scan(&name, &dbToken, &timestamp)
	if err == sql.ErrNoRows {
		return "", false
	}
	databaseError(err)

	db.Close()

	return name, true
}
