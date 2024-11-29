package main

import (
	"database/sql"
	"log"
)

type CustomerJSON struct {
	Name string `json:"customerName"`
	Code string `json:"customerCode"`
}

type Customer struct {
	ID   int    `field:"id"`
	Name string `field:"name"`
	Code string `field:"customer_code"`
}

func addCustomer(customer CustomerJSON, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO customers (name, customer_code) VALUES ($1,$2)",
		customer.Name, customer.Code)

	if err != nil {
		return err
	}
	return nil
}

func fetchCustomers(db *sql.DB) ([]Customer, error) {
	rows, err := db.Query("SELECT * FROM customers;")
	if err != nil {
		log.Println("Error fetchCustomers1: ", err)
		return nil, err
	}
	defer rows.Close()

	var customers []Customer

	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Code); err != nil {
			log.Println("Error fetchCustomers2: ", err)
			return customers, err
		}
		customers = append(customers, customer)
	}
	if err = rows.Err(); err != nil {
		return customers, err
	}

	return customers, nil
}
