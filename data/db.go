package data

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "postgres"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "postgres"
)

var db *sql.DB

func init() {
	sqlcon := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, _ = sql.Open("postgres", sqlcon)
	db.Ping()
}
func addUser(a, b, c string) error {
	qry := `insert into users(name,email,password) values($1,$2,$3)`
	_, er := db.Exec(qry, a, b, c)
	if er != nil {
		return errors.New("user already exists")
	}
	return nil
}

func verifyUser(b, c string) bool {
	qry := `select * from users where email=$1 and password=$2`
	r, _ := db.Exec(qry, b, c)
	i, _ := r.RowsAffected()
	if i != 1 {
		return false
	} else {
		return true
	}
}

func addTodo(a, b, c string) (Todo, error) {
	qry := `insert into todo(title,description,email) values ($1,$2,$3)`
	_, er := db.Exec(qry, a, b, c)
	var td Todo
	if er != nil {
		return td, er
	}
	row := db.QueryRow("select * from todo where id=(select max(id) from todo)")
	row.Scan(&td.Id, &td.Title, &td.Desc, &td.email)
	return td, nil

}

func updateTodo(a, b, c, d string) (Todo, error) {
	var td Todo
	qry := `update todo set title=$1,description=$2 where id=$3 and email=$4`
	r, er := db.Exec(qry, a, b, c, d)
	if er != nil {
		return td, errors.New("forbidden")
	}
	ind, _ := r.RowsAffected()
	if ind == 0 {
		return td, errors.New("forbidden")
	}
	row := db.QueryRow("select * from todo where id=$1", c)
	row.Scan(&td.Id, &td.Title, &td.Desc, &td.email)
	return td, nil
}

func deleteTodo(a, b string) error {
	r, _ := db.Exec(`delete from todo where id=$1 and email=$2`, a, b)
	row, _ := r.RowsAffected()
	if row == 1 {
		return nil
	} else {
		return errors.New("forbidden")
	}
}

func fetchtodo(email string, limit, offset int) ([]Todo, error) {
	var items []Todo
	var dt Todo
	rows, _ := db.Query(`select id,title,description from todo where email=$1 order by id limit $2 offset $3`, email, limit, offset)
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("no data found")
	}
	for rows.Next() {
		rows.Scan(&dt.Id, &dt.Title, &dt.Desc)
		items = append(items, dt)
	}
	return items, nil
}
