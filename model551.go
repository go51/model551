package model551

import (
	"database/sql"
	"errors"
	"github.com/go51/string551"
	"reflect"
	"strconv"
)

type Model struct {
	models map[string]*ModelInformation
}

type ModelInformation struct {
	NewFunc          NewModelFunc
	NewFuncPointer   NewModelFunc
	ModelType        reflect.Type
	ModelName        string
	TableInformation *TableInformation
	SqlInformation   *SqlCache
}

type TableInformation struct {
	TableName      string
	PrimaryKey     string
	Fields         []string
	DeleteTable    bool
	DeletedAtField string
}

type SqlCache struct {
	Insert        string
	Select        string
	Update        string
	Delete        string
	LogicalDelete string
}

var modelInstance *Model

type NewModelFunc func() interface{}
type SqlType int

const (
	SQL_INSERT SqlType = iota
	SQL_UPDATE
	SQL_LOGICAL_DELETE
)

type PrimaryInterface interface {
	SetId(int64)
	GetId() int64
}
type ScanInterface interface {
	Scan(rows *sql.Rows) error
}
type ValuesInterface interface {
	SqlValues(sqlType SqlType) []interface{}
}

func Load() *Model {
	if modelInstance != nil {
		return modelInstance
	}

	modelInstance = &Model{
		models: map[string]*ModelInformation{},
	}

	return modelInstance
}

func (m *Model) Add(newFunc NewModelFunc, newPointerFunc NewModelFunc) {
	model := newFunc()

	mType := reflect.TypeOf(model)
	mName := mType.Name()

	if m.models[mName] != nil {
		panic(errors.New("追加されたモデルは既に登録されています。"))
	}

	mInfo := &ModelInformation{
		NewFunc:        newFunc,
		NewFuncPointer: newPointerFunc,
		ModelType:      mType,
		ModelName:      mName,
	}

	mInfo.TableInformation = loadTableSetting(mType)

	mInfo.SqlInformation = cacheSql(mInfo.TableInformation)

	m.models[mName] = mInfo
}

func (m *Model) Get(modelName string) *ModelInformation {
	return m.models[modelName]
}

func loadTableSetting(mType reflect.Type) *TableInformation {
	tInfo := &TableInformation{
		TableName:      string551.SnakeCase(mType.Name()),
		PrimaryKey:     "id",
		Fields:         []string{},
		DeleteTable:    false,
		DeletedAtField: "",
	}

	if name := loadTableName(mType); name != "" {
		tInfo.TableName = name
	}

	if primaryKey := loadPrimaryKey(mType); primaryKey != "" {
		tInfo.PrimaryKey = primaryKey
	}

	if del, name := loadDeleteAt(mType); del {
		tInfo.DeleteTable = true
		tInfo.DeletedAtField = name
	}

	tInfo.Fields = loadFields(mType)

	return tInfo
}

func loadTableName(mType reflect.Type) string {
	for i := 0; i < mType.NumField(); i++ {
		sField := mType.Field(i)
		if name := sField.Tag.Get("db_table"); name != "" {
			return name
		}
	}

	return ""
}

func loadPrimaryKey(mType reflect.Type) string {
	for i := 0; i < mType.NumField(); i++ {
		sField := mType.Field(i)
		pk, err := strconv.ParseBool(sField.Tag.Get("db_pk"))
		if err == nil && pk {
			db := sField.Tag.Get("db")
			if db == "" {
				return string551.SnakeCase(sField.Name)
			} else {
				return string551.SnakeCase(db)
			}
		}
	}

	return ""

}

func loadDeleteAt(mType reflect.Type) (bool, string) {

	for i := 0; i < mType.NumField(); i++ {
		sField := mType.Field(i)
		del, err := strconv.ParseBool(sField.Tag.Get("db_delete"))
		if err == nil && del {
			return true, string551.SnakeCase(sField.Name)
		}
	}

	return false, ""
}

func loadFields(mType reflect.Type) []string {
	fields := make([]string, 0)

	for i := 0; i < mType.NumField(); i++ {
		sField := mType.Field(i)
		db := sField.Tag.Get("db")
		del, err := strconv.ParseBool(sField.Tag.Get("db_delete"))
		if err != nil || !del {
			if db == "" {
				fields = append(fields, string551.SnakeCase(sField.Name))
			} else if db != "-" {
				fields = append(fields, string551.SnakeCase(db))
			}
		}
	}

	return fields

}

func cacheSql(tInfo *TableInformation) *SqlCache {
	sqlCache := &SqlCache{
		Insert:        cacheSqlInsert(tInfo),
		Select:        cacheSqlSelect(tInfo),
		Update:        cacheSqlUpdate(tInfo),
		Delete:        cacheSqlDelete(tInfo),
		LogicalDelete: cacheSqlLogicalDelete(tInfo),
	}

	return sqlCache
}

func cacheSqlInsert(tInfo *TableInformation) string {
	sql := ""
	var append int = 0

	sql = string551.Join(sql, "INSERT INTO `"+tInfo.TableName+"` ")
	sql = string551.Join(sql, "(")
	for i := 0; i < len(tInfo.Fields); i++ {
		if tInfo.Fields[i] == tInfo.PrimaryKey {
			continue
		}
		if append == 0 {
			sql = string551.Join(sql, "`"+tInfo.Fields[i]+"`")
		} else {
			sql = string551.Join(sql, ", `"+tInfo.Fields[i]+"`")
		}
		append++
	}

	append = 0
	sql = string551.Join(sql, ") VALUES (")
	for i := 0; i < len(tInfo.Fields); i++ {
		if tInfo.Fields[i] == tInfo.PrimaryKey {
			continue
		}
		if append == 0 {
			sql = string551.Join(sql, "?")
		} else {
			sql = string551.Join(sql, ", ?")
		}
		append++
	}
	sql = string551.Join(sql, ")")

	return sql
}

func cacheSqlSelect(tInfo *TableInformation) string {
	sql := ""

	sql = string551.Join(sql, "SELECT ")
	for i := 0; i < len(tInfo.Fields); i++ {
		if i == 0 {
			sql = string551.Join(sql, "`"+tInfo.Fields[i]+"`")
		} else {
			sql = string551.Join(sql, ", `"+tInfo.Fields[i]+"`")
		}
	}
	sql = string551.Join(sql, " FROM `"+tInfo.TableName+"` WHERE 1 = 1")

	return sql
}

func cacheSqlUpdate(tInfo *TableInformation) string {
	sql := ""
	var append int = 0

	sql = string551.Join(sql, "UPDATE `"+tInfo.TableName+"` SET ")
	for i := 0; i < len(tInfo.Fields); i++ {
		if tInfo.Fields[i] == tInfo.PrimaryKey {
			continue
		}
		if append == 0 {
			sql = string551.Join(sql, "`"+tInfo.Fields[i]+"` = ?")
		} else {
			sql = string551.Join(sql, ", `"+tInfo.Fields[i]+"` = ?")
		}
		append++
	}
	sql = string551.Join(sql, " WHERE `"+tInfo.PrimaryKey+"` = ?")

	return sql
}

func cacheSqlDelete(tInfo *TableInformation) string {
	sql := ""

	sql = string551.Join(sql, "DELETE FROM `"+tInfo.TableName+"` WHERE `"+tInfo.PrimaryKey+"` = ?")

	return sql
}

func cacheSqlLogicalDelete(tInfo *TableInformation) string {
	sql := ""

	if tInfo.DeleteTable == false {
		return sql
	}

	sql = string551.Join(sql, "INSERT INTO `"+tInfo.TableName+"_delete` ")
	sql = string551.Join(sql, "(")
	for i := 0; i < len(tInfo.Fields); i++ {
		if i == 0 {
			sql = string551.Join(sql, "`"+tInfo.Fields[i]+"`")
		} else {
			sql = string551.Join(sql, ", `"+tInfo.Fields[i]+"`")
		}
	}
	sql = string551.Join(sql, ", `"+tInfo.DeletedAtField+"`) VALUES (")
	for i := 0; i < len(tInfo.Fields); i++ {
		if i == 0 {
			sql = string551.Join(sql, "?")
		} else {
			sql = string551.Join(sql, ", ?")
		}
	}
	sql = string551.Join(sql, ", ?)")

	return sql
}
