package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	xerrors "github.com/pkg/errors"
)


////////////////////////////////////////////////////////////////
// dao层

type User struct {
	Id			int32	`db:"id"`
	Name		string	`db:"name"`
	Password 	string	`db:"password"`
	Phone		string	`db:"phone"`
	Email		string	`db:"email"`
	CreateTime	string	`db:"create_time"`
	UpdateTime	string	`db:"update_time"`
}

func DbConnect(address, dbname,  username, password string) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, address, dbname)
	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		return nil, xerrors.Wrap(err, "db connect fail")
	}

	return db, nil
}

func DbGetUserInfo(db *sqlx.DB, id int32) (*User, error) {
	user := &User{}
	query := "select * from user where id = ?"
	err := db.Get(user, query, id)
	if err != nil {
		return nil, xerrors.Wrap(err, "user not found by id")
	}

	return user, nil
}

////////////////////////////////////////////////////////////////
// 业务处理层

type SysInfo struct {
	dbConf	*DbConfig
	db		*sqlx.DB
	// other info...
}

type DbConfig struct {
	address		string
	dbname		string
	username	string
	password	string
}

type SysProc struct {
	SysInfo
}

func (sysProc *SysProc) init(dbConf *DbConfig) error {
	db, err := DbConnect(dbConf.address,
		dbConf.dbname,
		dbConf.username,
		dbConf.password)
	if err != nil {
		return xerrors.WithMessage(err, "sys init fail")
	}

	sysProc.dbConf = dbConf
	sysProc.db = db
	return nil
}

func (sysProc *SysProc) userLogin(userid int32) error {
	user, err := DbGetUserInfo(sysProc.db, userid)
	if err != nil {
		return xerrors.WithMessage(err, "user login fail")
	}

	// 其他业务处理...
	fmt.Printf("User Info : \n\t%v\n", user)

	return nil
}

////////////////////////////////////////////////////////////////
// 业务逻辑测试
func testSucc() {
	dbConf := &DbConfig{"localhost:3306", "geektime", "geek", "geek"}
	sys := &SysProc{}
	err := sys.init(dbConf)
	if err != nil {
		fmt.Println("Init")
		fmt.Printf("original error: %T %v\n",
			xerrors.Cause(err), xerrors.Cause(err))
		fmt.Printf("stack trace:\n%+v", err)
	}

	err = sys.userLogin(1)
	if err != nil {
		fmt.Println("Login")
		fmt.Printf("original error: %T %v\n",
			xerrors.Cause(err), xerrors.Cause(err))
		fmt.Printf("stack trace:\n%+v", err)
	}
}


func testInitFail() {
	dbConf := &DbConfig{"localhost:3306", "geektime", "geek", "geek1"}
	sys := &SysProc{}
	err := sys.init(dbConf)
	if err != nil {
		fmt.Println("Init")
		fmt.Printf("original error: %T %v\n",
			xerrors.Cause(err), xerrors.Cause(err))
		fmt.Printf("stack trace:\n%+v", err)
	}
}

func testLoginFail() {
	dbConf := &DbConfig{"localhost:3306", "geektime", "geek", "geek"}
	sys := &SysProc{}
	err := sys.init(dbConf)
	if err != nil {
		fmt.Println("Init")
		fmt.Printf("original error: %T %v\n",
			xerrors.Cause(err), xerrors.Cause(err))
		fmt.Printf("stack trace:\n%+v", err)
	}

	err = sys.userLogin(2)
	if err != nil {
		fmt.Println("Login")
		fmt.Printf("original error: %T %v\n",
			xerrors.Cause(err), xerrors.Cause(err))
		fmt.Printf("stack trace:\n%+v", err)
	}
}

func main() {
	fmt.Printf("\n\nTEST SUCCESS:\n")
	testSucc()

	fmt.Printf("\n\nTEST INIT FAIL:\n")
	testInitFail()

	fmt.Printf("\n\nTEST LOGIN FAIL:\n")
	testLoginFail()
}


func main1() {
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
