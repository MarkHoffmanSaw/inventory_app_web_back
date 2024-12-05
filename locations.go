package main

import (
	"database/sql"
	"log"
)

type LocationFilter struct {
	query1 string
}

type LocationDB struct {
	ID          int    `field:"location_id"`
	Name        string `field:"name"`
	WarehouseID int    `field:"warehouse_id"`
}

func fetchLocations(db *sql.DB, options LocationFilter) ([]LocationDB, error) {
	rows, err := db.Query("SELECT * FROM locations;")
	if err != nil {
		log.Println("Error fetchLocations1: ", err)
		return nil, err
	}
	defer rows.Close()

	var locations []LocationDB

	for rows.Next() {
		var location LocationDB
		if err := rows.Scan(&location.ID, &location.Name, &location.WarehouseID); err != nil {
			log.Println("Error fetchLocations2: ", err)
			return locations, err
		}
		locations = append(locations, location)
	}
	if err = rows.Err(); err != nil {
		return locations, err
	}

	return locations, nil
}
