package orm_test

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var db *sql.DB

const (
	UseDatabase      = "USE util"
	CreateDatabase   = "CREATE DATABASE IF NOT EXISTS util"
	CreateTableSql   = "CREATE TABLE IF NOT EXISTS sample (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, user_id int NOT NULL, user_name VARCHAR(255) NOT NULL, city TEXT NOT NULL)"
	InsertRowsFormat = `INSERT INTO sample (user_id, user_name, city) VALUES (%d, "%s", "%s")`
	DropDatabase     = "DROP DATABASE util"
)

func TestOrm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Orm Suite")
}

var _ = BeforeSuite(func() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(CreateDatabase)
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(UseDatabase)
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(CreateTableSql)
	if err != nil {
		panic(err.Error())
	}

	city := []string{"beijing", "shanghai", "chengdu", "chongqing"}
	for i := 0; i < 100; i++ {
		script := fmt.Sprintf(InsertRowsFormat, i+1000, fmt.Sprintf("user_%d", i), city[i%len(city)])
		_, err = db.Exec(script)
		if err != nil {
			panic(err.Error())
		}
	}
})

var _ = AfterSuite(func() {
	db.Exec(DropDatabase)
	db.Close()
})
