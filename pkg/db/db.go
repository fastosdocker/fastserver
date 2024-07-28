package db

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var session *gorm.DB

type Query struct {
	Select string
	Where  string
	Order  string
	Group  string
	Offset int
	Limit  int
}

func Init(dsn string) error {
	err := initDB(dsn)
	if err != nil {
		return err
	}

	err = initTable()
	if err != nil {
		return err
	}

	return nil
}

func initDB(dsn string) error {
	log.Println("初始化数据库")
	var err error
	session, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %s", err)
	}

	//session = session.Debug()

	return nil
}

func initTable() error {
	err := CreateTable(
		Application{},
		ApplicationLog{},
		ContainerStats{},
		User{},
	)
	if err != nil {
		return fmt.Errorf("初始化表结构失败: %s", err)
	}

	return nil
}

func HasTable(model any) bool {
	return session.Migrator().HasTable(model)
}

func CreateTable(model ...any) error {
	return session.Migrator().AutoMigrate(model...)
}

func Create(model any) error {
	return session.Create(model).Error
}

func BatchCreate(model any) error {
	return session.Create(model).Error
}

func Delete(model any, cond string) error {
	if cond == "" {
		return fmt.Errorf("删除条件不能为空")
	}

	return session.Where(cond).Delete(model).Error
}

func Truncate(model schema.Tabler) error {
	return session.Exec("TRUNCATE TABLE ?", model.TableName()).Error
}

func Save(model any) error {
	return session.Save(model).Error
}

func Update(model any, val map[string]interface{}, conds ...string) error {
	s := session.Model(model)

	if len(conds) > 0 {
		s.Where(conds[0])
	}

	return s.Updates(val).Error
}

func Find(model any, query ...Query) error {
	s := session.Where("1=1")

	if len(query) > 0 {
		q := query[0]

		if q.Select != "" {
			s.Select(q.Select)
		}
		if q.Where != "" {
			s.Where(q.Where)
		}
		if q.Group != "" {
			s.Where(q.Group)
		}
		if q.Order != "" {
			s.Order(q.Order)
		}
		if q.Offset > 0 {
			s.Offset(q.Offset)
		}

		if q.Limit > 0 {
			s.Limit(q.Limit)
		}
	}

	return s.Find(model).Error
}

func FindOne(model any, query ...Query) error {
	if len(query) > 0 {
		query[0].Limit = 1
	} else {
		query = append(query, Query{Limit: 1})
	}

	return Find(model, query...)
}

func Count(model any, query ...Query) (int64, error) {
	s := session.Model(model)

	if len(query) > 0 {
		q := query[0]

		if q.Select != "" {
			s.Select(q.Select)
		}
		if q.Where != "" {
			s.Where(q.Where)
		}
		if q.Group != "" {
			s.Where(q.Group)
		}
		if q.Order != "" {
			s.Order(q.Order)
		}
		if q.Offset > 0 {
			s.Offset(q.Offset)
		}

		if q.Limit > 0 {
			s.Limit(q.Limit)
		}
	}

	var count int64
	err := s.Count(&count).Error

	return count, err
}

// Get get db session
func Get() *gorm.DB {
	return session
}
