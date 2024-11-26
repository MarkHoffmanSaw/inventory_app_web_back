package main

import (
	"database/sql"
	"strconv"
)

type MaterialTemplate struct {
	Customer     string  `json:"customer"`
	StockID      string  `json:"stockId"`
	MaterialType string  `json:"type"`
	Qty          string  `json:"quantity"`
	Cost         float64 `json:"cost"`
	MinQty       string  `json:"minQuantity"`
	MaxQty       string  `json:"maxQuantity"`
	Description  string  `json:"description"`
	Owner        string  `json:"owner"`
	IsActive     bool    `json:"isActive"`
}

func sendMaterial(material MaterialTemplate, db *sql.DB) error {
	qty, _ := strconv.Atoi(material.Qty)
	minQty, _ := strconv.Atoi(material.MinQty)
	maxQty, _ := strconv.Atoi(material.MaxQty)

	_, err := db.Query(`
				INSERT INTO incoming_materials
					(customer_name, stock_id, cost, quantity,
					max_required_quantity, min_required_quantity,
					notes, is_active, type, owner)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		material.Customer, material.StockID, material.Cost,
		qty, maxQty, minQty,
		material.Description, material.IsActive, material.MaterialType,
		material.Owner,
	)

	if err != nil {
		return err
	}
	return nil
}
