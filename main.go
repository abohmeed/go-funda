package main

import "github.com/abohmeed/go-funda/funda"
import "github.com/abohmeed/go-funda/SQLite"

func main() {
	listings := funda.GetListings("funda.nl", "amersfoort")
	SQLite.Save(listings, "amersfoort")
}
