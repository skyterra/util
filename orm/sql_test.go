package orm_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
	"util/orm"
)

type UserInfo struct {
	UserID   int    `db:"user_id"`
	UserName string `db:"user_name"`
	City     string `db:"city"`
}

var _ = Describe("Sql", func() {

	Context("insert data", func() {
		It("should be succeed", func() {
			cols, err := orm.GetColNames(&UserInfo{}, "db")
			Expect(err).Should(Succeed())

			sql := fmt.Sprintf("SELECT %s FROM sample WHERE user_id = 1086", strings.Join(cols, ","))
			records, err := orm.Query(context.TODO(), db, sql, &UserInfo{})
			Expect(err).Should(Succeed())
			Expect(len(records) == 1).Should(BeTrue())

			r, ok := records[0].(*UserInfo)
			Expect(ok).Should(BeTrue())
			Expect(r.UserID == 1086).Should(BeTrue())
		})
	})
})
