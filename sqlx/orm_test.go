package sqlx

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestOrm(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	rs := mock.NewRows([]string{"id"}).FromCSVString("1")
	t.Run("basic", func(t *testing.T) {
		mock.ExpectQuery("select (.+) from user where id = ?").WithArgs(1).WillReturnRows(rs)
		row, err := db.Query("select id from user where id = ?", 1)
		assert.Nil(t, err)

		var data int
		err = UnmarshalRow(row, &data)
		assert.Nil(t, err)
		assert.Equal(t, 1, data)
	})

	t.Run("struct", func(t *testing.T) {
		rs := mock.NewRows([]string{"id", "name", "Age"}).FromCSVString("1,test,20")
		mock.ExpectQuery("select (.+) from user where id = ?").WithArgs(1).WillReturnRows(rs)
		type Foo struct {
			Id   int64  `db:"id"`
			Name string `db:"name"`
			Age  int
		}
		var foo Foo
		rows, err := db.Query("select id,name,Age from user where id = ?", 1)
		assert.Nil(t, err)

		err = UnmarshalRow(rows, &foo)
		assert.Nil(t, err)
		assert.True(t, foo.Id == 1 && foo.Name == "test" && foo.Age == 20)
	})

	t.Run("struct", func(t *testing.T) {
		rs := mock.NewRows([]string{"id", "name", "age"}).FromCSVString("1,test,20")
		mock.ExpectQuery("select (.+) from user where id = ?").WithArgs(1).WillReturnRows(rs)
		type Foo struct {
			Id   int64  `db:"id"`
			Name string `db:"name"`
			Age  int
		}
		var foo Foo
		rows, err := db.Query("select id,name,age from user where id = ?", 1)
		assert.Nil(t, err)

		err = UnmarshalRow(rows, &foo)
		assert.Nil(t, err)
		assert.True(t, foo.Id == 1 && foo.Name == "test" && foo.Age == 0)
	})

	t.Run("struct anonymous", func(t *testing.T) {
		rs := mock.NewRows([]string{"id", "name", "Age", "id_number", "gender", "nickname"}).FromCSVString("1,test,20,1001,男,test_alias")
		mock.ExpectQuery("select (.+) from user where id = ?").WithArgs(1).WillReturnRows(rs)
		type Bar struct {
			IdNumber string `db:"id_number"`
			Gender   string `db:"gender"`
		}

		type Foo struct {
			Id   int64  `db:"id"`
			Name string `db:"name"`
			Age  int
			Bar
		}

		var foo Foo
		rows, err := db.Query("select id,name,age from user where id = ?", 1)
		assert.Nil(t, err)

		err = UnmarshalRow(rows, &foo)
		assert.Nil(t, err)
		assert.True(t, foo.Id == 1 && foo.Name == "test" && foo.Age == 20 && foo.IdNumber == "1001" && foo.Gender == "男")
	})

	t.Run("slice pointer base", func(t *testing.T) {
		rs := mock.NewRows([]string{"id"}).FromCSVString("1\n2")
		mock.ExpectQuery("select (.+) from user where id = ?").WithArgs(1).WillReturnRows(rs)

		var foo []*int
		rows, err := db.Query("select id from user where id = ?", 1)
		assert.Nil(t, err)

		err = UnmarshalRows(rows, &foo)
		assert.Nil(t, err)
		var ret []int
		for _, i := range foo {
			ret = append(ret, *i)
		}
		assert.Equal(t, []int{1, 2}, ret)
	})

	t.Run("slice base", func(t *testing.T) {
		rs := mock.NewRows([]string{"id"}).FromCSVString("1\n2")
		mock.ExpectQuery("select (.+) from user where id = ?").WithArgs(1).WillReturnRows(rs)

		var foo []int
		rows, err := db.Query("select id from user where id = ?", 1)
		assert.Nil(t, err)

		err = UnmarshalRows(rows, &foo)
		assert.Nil(t, err)
		assert.Equal(t, []int{1, 2}, foo)
	})

	t.Run("slice struct", func(t *testing.T) {
		rs := mock.NewRows([]string{"id", "name"}).FromCSVString("1,test1\n2,test2")
		mock.ExpectQuery("select (.+) from user where id = ?").WithArgs(1).WillReturnRows(rs)

		type Foo struct {
			Id   int64  `db:"id"`
			Name string `db:"name"`
		}
		var foo []*Foo
		rows, err := db.Query("select id from user where id = ?", 1)
		assert.Nil(t, err)

		err = UnmarshalRows(rows, &foo)
		assert.Nil(t, err)
		assert.Equal(t, []*Foo{
			{
				Id:   1,
				Name: "test1",
			},
			{
				Id:   2,
				Name: "test2",
			},
		}, foo)
	})

	t.Run("slice pointer struct", func(t *testing.T) {
		rs := mock.NewRows([]string{"id", "name"}).FromCSVString("1,test1\n2,test2")
		mock.ExpectQuery("select (.+) from user where id = ?").WithArgs(1).WillReturnRows(rs)

		type Foo struct {
			Id   int64  `db:"id"`
			Name string `db:"name"`
		}
		var foo []*Foo
		rows, err := db.Query("select id from user where id = ?", 1)
		assert.Nil(t, err)

		err = UnmarshalRows(rows, &foo)
		assert.Nil(t, err)
		assert.Equal(t, []*Foo{
			{
				Id:   1,
				Name: "test1",
			},
			{
				Id:   2,
				Name: "test2",
			},
		}, foo)
	})
}
