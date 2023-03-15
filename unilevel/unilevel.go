package unilevel

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Member struct {
	ID             int
	Ref_code       string
	Child_ref_code string
	Firstname      string
}

func main() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_ziipfund")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM tbl_users")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var members []Member

	for rows.Next() {
		var member Member
		err := rows.Scan(&member.ID, &member.Ref_code, &member.Child_ref_code, &member.Firstname)
		if err != nil {
			panic(err.Error())
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		panic(err.Error())
	}

	fmt.Printf("%+v\n", members)
}
