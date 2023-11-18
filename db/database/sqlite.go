package db

import (
	"database/sql"
	log "emailTest/logger"
	_ "github.com/mattn/go-sqlite3"
	"sync"
	"time"
)

var (
	SqliteDb *SqliteDbType
	once     sync.Once
)

type SqliteDbType struct {
	db *sql.DB
}
type User struct {
	uid      int
	name     string
	created  string
	birthday string
}
type UserList struct {
	list   []User
	second int
}

func InitDb(path string) {
	once.Do(func() {
		db, err := sql.Open("sqlite3", path)
		if err != nil {
			log.Logger.Error(err)
			return
		}
		SqliteDb = &SqliteDbType{
			db: db,
		}
		table := `
    CREATE TABLE IF NOT EXISTS user (
        uid INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(128) NULL,
        created DATE NULL,
        birthday VARCHAR(128) NULL
    );
    `
		_, err = SqliteDb.db.Exec(table)
		if err != nil {
			log.Logger.Error(err)
			return
		}

	})

}
func (t *SqliteDbType) AddPerson(name string, birthday string) {
	searched := t.Search(name, birthday)
	if searched {
		t.Modify(name, birthday)
		return
	}
	t.Add(name, birthday)
}

func (t *SqliteDbType) Search(name string, birthday string) bool {
	rows, err := t.db.Query("SELECT * FROM user where name=?", name)
	if err != nil {
		panic(err)
	}
	searched := false
	defer rows.Close()
	var user = &User{}
	for rows.Next() {
		err = rows.Scan(&user.uid, &user.name, &user.created, &user.birthday)
		if err != nil {
			panic(err)
		}
		searched = true

	}
	return searched
}
func (t *SqliteDbType) Add(name string, birthday string) {
	stmt, err := t.db.Prepare("INSERT INTO user(name,  created,birthday) values(?,?,?)")
	if err != nil {
		panic(err)
	}
	// res 为返回结果
	res, err := stmt.Exec(name, time.Now(), birthday)
	if err != nil {
		panic(err)
	}

	// 可以通过res取自动生成的id
	_, err = res.LastInsertId()
	if err != nil {
		panic(err)
	}
}
func (t *SqliteDbType) Delete(name string) {
	stmt, err := t.db.Prepare("delete from user where uid=?")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(name)
	if err != nil {
		panic(err)
	}
}
func (t *SqliteDbType) Modify(name string, birthday string) {
	stmt, err := t.db.Prepare("update user set birthday=? where name=?")
	if err != nil {
		panic(err)
	}

	res, err := stmt.Exec(birthday, name)
	if err != nil {
		panic(err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		panic(err)
	}
}
