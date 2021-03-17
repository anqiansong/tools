package sqlx_test

import (
	"context"
	"fmt"
	"log"

	"tools/sqlx"

	"github.com/DATA-DOG/go-sqlmock"
)

func ExampleUnmarshalRow() {
	db, mock, err := sqlmock.New()
	must(err)

	rows := sqlmock.NewRows([]string{"id"}).FromCSVString("1")
	mock.ExpectQuery("select (.+) from foo where id = ?").WithArgs(1).WillReturnRows(rows)
	rs, err := db.QueryContext(context.Background(), `select id from foo where id = ?`, 1)
	must(err)
	var i int
	must(sqlx.UnmarshalRow(rs, &i))
	fmt.Println(i)
	// Output:
	// 1
}

func ExampleUnmarshalRows() {
	db, mock, err := sqlmock.New()
	must(err)

	rows := sqlmock.NewRows([]string{"id", "name", "age", "score", "gender", "graduate"}).FromCSVString("1,test,20,89.5,男,1\n2,test2,20,90.5,女,0")
	mock.ExpectQuery("select (.+) from foo where id = ?").WithArgs(1).WillReturnRows(rows)
	rs, err := db.QueryContext(context.Background(), `select * from foo where id = ?`, 1)
	must(err)

	type User struct {
		Id       string  `db:"id"`
		Name     string  `db:"name"`
		Age      int     `db:"age"`
		Score    float32 `db:"score"`
		Gender   string  `db:"gender"`
		Graduate int     `db:"graduate"`
	}

	var user []*User
	must(sqlx.UnmarshalRows(rs, &user))
	for _, e := range user {
		fmt.Println(e.Id, e.Name, e.Age, e.Score, e.Gender, e.Graduate)
	}

	// Output:
	// 1 test 20 89.5 男 1
	// 2 test2 20 90.5 女 0
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
