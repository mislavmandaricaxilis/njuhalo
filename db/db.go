package db

import (
	"database/sql"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/lpredova/njuhalo/model"
	_ "github.com/mattn/go-sqlite3" // SQLlite db
)

const dbName = "./njuhalo.db"

var usr, _ = user.Current()
var dbPath = usr.HomeDir + "/.njuhalo/" + "njuhalo.db"

// InsertItem method inserts new offer into database
func InsertItem(offers []model.Offer) bool {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO items(itemID, url, name, image, price, description, createdAt) values(?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	for _, offer := range offers {
		_, err := stmt.Exec(offer.ID, offer.URL, offer.Name, offer.Image, offer.Price, offer.Description, int32(time.Now().Unix()))
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}

	return true
}

// GetItem method that checks if there is alreay offer with that ID
func GetItem(itemID string) bool {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM items where itemID = %s", itemID))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	for rows.Next() {
		return true
	}

	return false
}

// CreateDatabase creates sqllite db file in user home dir
func CreateDatabase() bool {

	usr, err := user.Current()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	err = os.MkdirAll(usr.HomeDir+"/.njuhalo/", 0755)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if _, err = os.Stat(usr.HomeDir + "/.njuhalo"); os.IsNotExist(err) {
		os.Mkdir(usr.HomeDir+"/.njuhalo", 0755)
	}
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	f, err := os.Create(dbPath)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer f.Close()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	stmt, err := db.Prepare("CREATE TABLE items (id integer PRIMARY KEY AUTOINCREMENT, itemID integer, url text, name text, image text, price text, description text, createdAt integer)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	stmt, err = db.Prepare("CREATE INDEX index_itemID ON items (itemID)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}
