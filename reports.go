package main

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/leekchan/accounting"
)

type Transaction struct {
	StockID      string    `field:"stock_id"`
	LocationName string    `field:"location_name"`
	MaterialType string    `field:"material_type"`
	Qty          int       `field:"quantity"`
	UnitCost     float64   `field:"unit_cost"`
	Cost         float64   `field:"cost"`
	UpdatedAt    time.Time `field:"updated_at"`
	TotalValue   float64   `field:"total_value"`
}

type SearchQuery struct {
	customerId   int
	materialType string
	dateFrom     string
	dateTo       string
	dateAsOf     string
}

type Report struct {
	db *sql.DB
}

type TransactionReport struct {
	Report
	trxFilter SearchQuery
}

type BalanceReport struct {
	Report
	blcFilter SearchQuery
}

type TransactionRep struct {
	StockID      string
	MaterialType string
	Qty          string
	UnitCost     string
	Cost         string
	Date         string
}

type BalanceRep struct {
	StockID      string
	LocationName string
	MaterialType string
	Qty          string
	TotalValue   string
}

var accLib accounting.Accounting = accounting.Accounting{Symbol: "$", Precision: 2}

func (t TransactionReport) getReportList() ([]TransactionRep, error) {
	rows, err := t.db.Query(`SELECT tl.stock_id, m.material_type,
								tl.quantity_change as "quantity",
								tl.cost as "unit_cost",
								(tl.quantity_change * tl.cost) as "cost",
								tl.updated_at
							 FROM transactions_log tl
							 LEFT JOIN materials m ON m.material_id = tl.material_id
							 LEFT JOIN customers c ON m.customer_id = c.customer_id
							 WHERE 
								($1 = 0 OR m.customer_id = $1) AND
								($2 = '' OR m.material_type::TEXT = $2) AND
								($3 = '' OR tl.updated_at::TEXT >= $3) AND
								($4 = '' OR tl.updated_at::TEXT <= $4)
							 ORDER BY transaction_id;`,
		t.trxFilter.customerId, t.trxFilter.materialType, t.trxFilter.dateFrom, t.trxFilter.dateTo)
	if err != nil {
		return []TransactionRep{}, err
	}

	trxList := []TransactionRep{}

	for rows.Next() {
		trx := Transaction{}

		err := rows.Scan(
			&trx.StockID,
			&trx.MaterialType,
			&trx.Qty,
			&trx.UnitCost,
			&trx.Cost,
			&trx.UpdatedAt,
		)
		if err != nil {
			return []TransactionRep{}, err
		}

		year, month, day := trx.UpdatedAt.Date()
		strDate := strconv.Itoa(int(month)) + "/" +
			strconv.Itoa(day) + "/" +
			strconv.Itoa(year)
		unitCost := accLib.FormatMoney(trx.UnitCost)
		cost := accLib.FormatMoney(trx.Cost)

		trxList = append(trxList, TransactionRep{
			StockID:      trx.StockID,
			MaterialType: trx.MaterialType,
			Qty:          strconv.Itoa(trx.Qty),
			UnitCost:     unitCost,
			Cost:         cost,
			Date:         strDate,
		})
	}

	return trxList, nil
}

func (b BalanceReport) getReportList() ([]BalanceRep, error) {
	rows, err := b.db.Query(`
	SELECT m.stock_id,
		   l.name as "location_name",
		   m.material_type,
		   SUM(tl.quantity_change) AS "quantity",
		   SUM(tl.quantity_change * tl.cost) AS "total_value"
	FROM transactions_log tl
	LEFT JOIN materials m ON m.material_id = tl.material_id
	LEFT JOIN locations l ON l.location_id = m.location_id
	WHERE
		($1 = 0 OR m.customer_id = $1) AND
		($2 = '' OR m.material_type::TEXT = $2) AND
		($3 = '' OR tl.updated_at::TEXT <= $3)
	GROUP BY m.stock_id, l.name, m.material_type
`,
		b.blcFilter.customerId, b.blcFilter.materialType, b.blcFilter.dateAsOf,
	)
	if err != nil {
		return []BalanceRep{}, err
	}

	blcList := []BalanceRep{}

	for rows.Next() {
		balance := Transaction{}

		err := rows.Scan(
			&balance.StockID,
			&balance.LocationName,
			&balance.MaterialType,
			&balance.Qty,
			&balance.TotalValue,
		)
		if err != nil {
			return []BalanceRep{}, err
		}

		totalValue := accLib.FormatMoney(balance.TotalValue)
		blcList = append(blcList, BalanceRep{
			StockID:      balance.StockID,
			LocationName: balance.LocationName,
			MaterialType: balance.MaterialType,
			Qty:          strconv.Itoa(balance.Qty),
			TotalValue:   totalValue,
		})
	}

	return blcList, err

}
