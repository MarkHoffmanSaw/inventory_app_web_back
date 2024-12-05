package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

type IncomingMaterialJSON struct {
	CustomerID   string `json:"customerId"`
	StockID      string `json:"stockId"`
	MaterialType string `json:"type"`
	Qty          string `json:"quantity"`
	Cost         string `json:"cost"`
	MinQty       string `json:"minQuantity"`
	MaxQty       string `json:"maxQuantity"`
	Description  string `json:"description"`
	Owner        string `json:"owner"`
	IsActive     bool   `json:"isActive"`
}

type IncomingMaterialDB struct {
	ShippingID   string  `field:"shipping_id"`
	CustomerName string  `field:"customer_name"`
	CustomerID   int     `field:"customer_id"`
	StockID      string  `field:"stock_id"`
	Cost         float64 `field:"cost"`
	Quantity     int     `field:"quantity"`
	MinQty       int     `field:"min_required_quantity"`
	MaxQty       int     `field:"max_required_quantity"`
	Notes        string  `field:"notes"`
	IsActive     bool    `field:"is_active"`
	MaterialType string  `field:"material_type"`
	Owner        string  `field:"owner"`
}

type MaterialJSON struct {
	IncomingMaterialID string `json:"incomingMaterialId"`
	LocationID         string `json:"locationId"`
	Notes              string `json:"notes"`
	Qty                string `json:"quantity"`
}

type MaterialDB struct {
	MaterialId   int       `field:"material_id"`
	StockId      string    `field:"stock_id"`
	LocationId   int       `field:"location_id"`
	CustomerId   int       `field:"customer_id"`
	MaterialType string    `field:"material_type"`
	Description  string    `field:"description"`
	Notes        string    `field:"notes"`
	Quantity     int       `field:"quantity"`
	UpdatedAt    time.Time `field:"updated_at"`
	IsActive     bool      `field:"is_active"`
	Cost         float64   `field:"cost"`
	MinQty       int       `field:"min_required_quantity"`
	MaxQty       int       `field:"max_required_quantity"`
	Owner        string    `field:"onwer"`
}

type TransactionInfo struct {
	materialId    int       `field:"material_id"`
	stockId       string    `field:"stock_id"`
	quantity      int       `field:"quantity_change"`
	notes         string    `field:"notes"`
	cost          float64   `field:"cost"`
	updatedAt     time.Time `field:"updated_at"`
	jobTicket     string    `field:"job_ticket"`
	isMove        bool      // opts
	newMaterialId int       // opts
}

func fetchMaterialTypes() []string {
	return []string{"Card", "Envelope"}
}

func sendMaterial(material IncomingMaterialJSON, db *sql.DB) error {
	qty, _ := strconv.Atoi(material.Qty)
	minQty, _ := strconv.Atoi(material.MinQty)
	maxQty, _ := strconv.Atoi(material.MaxQty)

	_, err := db.Query(`
				INSERT INTO incoming_materials
					(customer_id, stock_id, cost, quantity,
					max_required_quantity, min_required_quantity,
					notes, is_active, type, owner)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		material.CustomerID, material.StockID, material.Cost,
		qty, maxQty, minQty,
		material.Description, material.IsActive, material.MaterialType,
		material.Owner,
	)

	log.Println(err)

	if err != nil {
		return err
	}
	return nil
}

func getIncomingMaterials(db *sql.DB) ([]IncomingMaterialDB, error) {
	rows, err := db.Query(`
		SELECT shipping_id, c.name, c.customer_id, stock_id, cost, quantity,
		min_required_quantity, max_required_quantity, notes, is_active, type, owner
		FROM incoming_materials im
		LEFT JOIN customers c ON c.customer_id = im.customer_id
		`)
	if err != nil {
		return nil, fmt.Errorf("Error querying incoming materials: %w", err)
	}
	defer rows.Close()

	var materials []IncomingMaterialDB
	for rows.Next() {
		var material IncomingMaterialDB
		if err := rows.Scan(
			&material.ShippingID,
			&material.CustomerName,
			&material.CustomerID,
			&material.StockID,
			&material.Cost,
			&material.Quantity,
			&material.MinQty,
			&material.MaxQty,
			&material.Notes,
			&material.IsActive,
			&material.MaterialType,
			&material.Owner,
		); err != nil {
			return nil, fmt.Errorf("Error scanning row: %w", err)
		}
		materials = append(materials, material)
	}
	return materials, nil
}

func createMaterial(material MaterialJSON, db *sql.DB) error {
	var incomingMaterial IncomingMaterialDB

	err := db.QueryRow(`
		SELECT customer_id, stock_id, cost, min_required_quantity,
		max_required_quantity, notes, is_active, type, owner
		FROM incoming_materials
		WHERE shipping_id = $1`, material.IncomingMaterialID).
		Scan(
			&incomingMaterial.CustomerID,
			&incomingMaterial.StockID,
			&incomingMaterial.Cost,
			&incomingMaterial.MinQty,
			&incomingMaterial.MaxQty,
			&incomingMaterial.Notes,
			&incomingMaterial.IsActive,
			&incomingMaterial.MaterialType,
			&incomingMaterial.Owner,
		)

	if err != nil {
		return err
	}

	// Update material in the current location
	rows, err := db.Query(`
					UPDATE materials
					SET quantity = (quantity + $1),
						notes = $2
					WHERE stock_id = $3
						AND location_id = $4
						AND owner = $5
					RETURNING material_id;
					`, material.Qty, material.Notes, incomingMaterial.StockID, material.LocationID, incomingMaterial.Owner,
	)

	if err != nil {
		return err
	}

	var materialId int

	for rows.Next() {
		err := rows.Scan(&materialId)
		if err != nil {
			return err
		}
	}

	// If there is no the same material in the current location
	// Then add the material in the chosen one
	if materialId == 0 {
		err := db.QueryRow(`
						INSERT INTO materials
						(
							stock_id,
							location_id,
							customer_id,
							material_type,
							description,
							notes,
							quantity,
							updated_at,
							min_required_quantity,
							max_required_quantity,
							is_active,
							cost,
							owner
						)
						VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING material_id;`,
			incomingMaterial.StockID,
			material.LocationID,
			incomingMaterial.CustomerID,
			incomingMaterial.MaterialType,
			incomingMaterial.Notes,
			material.Notes,
			material.Qty,
			time.Now(),
			incomingMaterial.MinQty,
			incomingMaterial.MaxQty,
			incomingMaterial.IsActive,
			incomingMaterial.Cost,
			incomingMaterial.Owner,
		).Scan(&materialId)

		if err != nil {
			return err
		}
	}

	// Remove the material from incoming
	shippingId, _ := strconv.Atoi(material.IncomingMaterialID)
	err = deleteIncomingMaterial(db, shippingId)

	if err != nil {
		return err
	}

	qty, _ := strconv.Atoi(material.Qty)
	err = addTranscation(&TransactionInfo{
		materialId: materialId,
		stockId:    incomingMaterial.StockID,
		quantity:   qty,
		notes:      material.Notes,
		updatedAt:  time.Now(),
		cost:       incomingMaterial.Cost,
	}, db)

	if err != nil {
		return err
	}
	return nil
}

func deleteIncomingMaterial(db *sql.DB, shippingId int) error {
	if _, err := db.Exec(`
			DELETE FROM incoming_materials WHERE shipping_id = $1;`,
		shippingId); err != nil {
		return err
	}

	return nil
}

func addTranscation(trx *TransactionInfo, db *sql.DB) error {
	if trx.quantity < 0 {
		removingQty := int(math.Abs(float64(trx.quantity)))

		emptyCost := []string{"0"}

		for removingQty > 0 {
			var transactionId int
			var cost float64
			var remainingQty int

			// Find a last deduction
			db.QueryRow(`
				SELECT transaction_id, cost, remaining_quantity FROM transactions_log
				WHERE material_id = $1 AND stock_id = $2 AND quantity_change < 0
					AND cost NOT IN (`+strings.Join(emptyCost, ",")+`)
				ORDER BY transaction_id DESC LIMIT 1;
						`,
				trx.materialId,
				trx.stockId).Scan(&transactionId, &cost, &remainingQty)

			// First deduction is NOT found
			if transactionId == 0 {
				db.QueryRow(`
					SELECT transaction_id, cost, remaining_quantity FROM transactions_log
					WHERE material_id = $1 AND stock_id = $2  AND quantity_change > 0
						AND cost NOT IN (`+strings.Join(emptyCost, ",")+`)
					ORDER BY transaction_id LIMIT 1;
							`,
					trx.materialId,
					trx.stockId,
				).Scan(&transactionId, &cost, &remainingQty)

				// When neither positive nor negative calculations found
				if transactionId == 0 {
					return errors.New("no remains found")
				}

				// First deduction is found, but remains are zero
			} else if transactionId != 0 && remainingQty == 0 {
				emptyCost = append(emptyCost, strconv.FormatFloat(cost, 'f', -1, 64))
				continue
			}

			// Deduct from the balance
			if remainingQty < removingQty {
				removingQty -= remainingQty

				_, errInsert := db.Exec(
					`INSERT INTO transactions_log
							(material_id, stock_id, quantity_change, notes,
							cost, job_ticket, updated_at, remaining_quantity)
							 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
							 `, trx.materialId, trx.stockId, -remainingQty, trx.notes,
					cost, trx.jobTicket, trx.updatedAt, 0)

				if errInsert != nil {
					log.Println("err1", errInsert)
					return errInsert
				}

				emptyCost = append(emptyCost, strconv.FormatFloat(cost, 'f', -1, 64))

				if trx.isMove {
					addTranscation(&TransactionInfo{
						materialId: trx.newMaterialId,
						stockId:    trx.stockId,
						quantity:   remainingQty,
						notes:      trx.notes,
						cost:       cost,
						updatedAt:  trx.updatedAt,
						jobTicket:  trx.jobTicket,
					}, db)
				}
			} else if remainingQty >= removingQty {
				remainingQty -= removingQty

				_, errInsert := db.Exec(
					`INSERT INTO transactions_log
							(material_id, stock_id, quantity_change, notes,
							cost, job_ticket, updated_at, remaining_quantity)
							 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
							 `, trx.materialId, trx.stockId, -removingQty, trx.notes,
					cost, trx.jobTicket, trx.updatedAt, remainingQty)

				if errInsert != nil {
					log.Println("err2", errInsert)
					return errInsert
				}

				if trx.isMove {
					addTranscation(&TransactionInfo{
						materialId: trx.newMaterialId,
						stockId:    trx.stockId,
						quantity:   removingQty,
						notes:      trx.notes,
						cost:       cost,
						updatedAt:  trx.updatedAt,
						jobTicket:  trx.jobTicket,
					}, db)
				}

				removingQty = 0
			}
		}
	} else {
		// Check if an ID with the same cost exists
		var transactionId int
		db.QueryRow(`
				SELECT transaction_id FROM transactions_log
				WHERE
					material_id = $1 AND
					stock_id = $2 AND
					quantity_change > 0 AND
					cost = $3
				ORDER BY transaction_id DESC LIMIT 1;
						`,
			trx.materialId, trx.stockId, trx.cost).Scan(&transactionId)

		// If the ID exists then update it
		if transactionId > 0 {
			_, e := db.Query(`
				UPDATE transactions_log
				SET quantity_change = quantity_change + $2,
					remaining_quantity = remaining_quantity + $2,
					updated_at = NOW()
				WHERE transaction_id = $1

		`, transactionId, trx.quantity)

			if e != nil {
				return e
			}
		} else {
			// If an ID doesn't exist then add a new one
			_, e := db.Exec(
				`INSERT INTO transactions_log
			(material_id, stock_id, quantity_change, notes,
			cost, job_ticket, updated_at, remaining_quantity)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			 `, trx.materialId, trx.stockId, trx.quantity, trx.notes,
				trx.cost, trx.jobTicket, trx.updatedAt, trx.quantity)

			if e != nil {
				return e
			}
		}
	}
	return nil
}
