package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id			int32	`db:"id"`
	Name		string	`db:"name"`
	Password 	string	`db:"password"`
	Phone		string	`db:"phone"`
	Email		string	`db:"email"`
	CreateTime	string	`db:"create_time"`
	UpdateTime	string	`db:"update_time"`
}

func main() {
	db, err :=sqlx.Connect("mysql", "geek:geek@tcp(localhost:3306)/geektime")
	if err != nil {
		fmt.Println(err)
		return
	}

	query := "select * from user where id = ?"

	user := User{}
	id := 1
	err = db.Get(&user, query, id)
	fmt.Printf("db.Get(%v): \n\t" +
		"user = %v\n\t" +
		"err : %v\n\n",
		id, user, err)

	user = User{}
	id = 2
	err = db.Get(&user, query, id)
	fmt.Printf("db.Get(%v): \n\t" +
		"user = %v\n\t" +
		"err : %v\n\n",
		id, user, err)

}
