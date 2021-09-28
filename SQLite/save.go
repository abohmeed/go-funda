package SQLite

import (
	"database/sql"
	"fmt"
	"github.com/abohmeed/go-funda/funda"
	"log"
	_ "github.com/mattn/go-sqlite3"

)

func Save(listings []funda.Listing, city string) {
	db, err := sql.Open("sqlite3", "./houses.db")
	checkErr("Could not open database:",err)
	for _, l := range listings {
		stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s(Title,Address,Price,ListedSince,Insulation,EnergyLabel,Facilities,ConstructionYear,LivingArea,Status,URL) values(?,?,?,?,?,?,?,?,?,?,?,?)", city))
		checkErr("Could not prepare SQL statement:",err)
		res, err := stmt.Exec(l.Title, l.Address, l.Price, l.ListedSince, l.Insulation, l.EnergyLabel, l.Facilities, l.ConstructionYear, l.Bedrooms, l.LivingArea, l.Status, l.URL)
		checkErr("Could not execute statement:",err)
		id, err := res.LastInsertId()
		checkErr("Could not get last inserted ID:",err)
		fmt.Println(fmt.Sprintf("Inserted row with id %d\n\n", id))
	}
}
func checkErr(msg string,err error) {
	if err != nil {
		log.Fatal(err)
	}
}
