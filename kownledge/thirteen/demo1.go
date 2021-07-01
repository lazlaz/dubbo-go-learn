package main

import (
	"fmt"

	"log"
)

import (
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/sql/sem/tree"
)

func main() {
	sql := `INSERT INTO so_master (sysno, so_id, buyer_user_sysno, seller_company_code,
		receive_division_sysno, receive_address, receive_zip, receive_contact, receive_contact_phone, stock_sysno,
		payment_type, so_amt, status, order_date, appid, memo) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,now(),$14,$15)`

	if stmts, err := parser.Parse(sql); err != nil {
		log.Fatal(err)
	} else {
		for _, stmt := range stmts {
			fmt.Printf("Pretty SQL:\n%s", tree.Pretty(stmt.AST))
		}
	}

}
