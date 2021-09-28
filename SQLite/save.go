package SQLite

import (
	"database/sql"
	"fmt"
	"github.com/abohmeed/go-funda/funda"
	"log"
)

func Save(listings []funda.Listing, city string) {
	db, err := sql.Open("sqlite3", "./houses.db")
	checkErr(err)
	for _, l := range listings {
		stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s(Title,Address,Price,ListedSince,Insulation,EnergyLabel,Facilities,ConstructionYear,LivingArea,Status,URL) values(?,?,?,?,?,?,?,?,?,?,?,?)", city))
		checkErr(err)
		res, err := stmt.Exec(l.Title, l.Address, l.Price, l.ListedSince, l.Insulation, l.EnergyLabel, l.Facilities, l.ConstructionYear, l.Bedrooms, l.LivingArea, l.Status, l.URL)
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)
		fmt.Println(fmt.Sprintf("Inserted row with id %d\n\n", id))
	}
}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
