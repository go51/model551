package model551
import (
	"reflect"
	"errors"
	"github.com/go51/string551"
	"strconv"
)

type Model struct {
	models map[string]*ModelInformation
}

type ModelInformation struct {
	NewFunc          NewModelFunc
	ModelType        reflect.Type
	ModelName        string
	TableInformation *TableInformation
}

type TableInformation struct {
	TableName      string
	PrimaryKey     string
	Fields         []string
	DeleteTable    bool
	DeletedAtField string
}

var modelInstance *Model

type NewModelFunc func() interface{}

func Load() *Model {
	if modelInstance != nil {
		return modelInstance
	}

	modelInstance = &Model{
		models:map[string]*ModelInformation{},
	}

	return modelInstance
}

func (m *Model) Add(newFunc NewModelFunc) {
	model := newFunc()

	mType := reflect.TypeOf(model)
	mName := mType.Name()

	if m.models[mName] != nil {
		panic(errors.New("追加されたモデルは既に登録されています。"))
	}

	detail := &ModelInformation{
		NewFunc:newFunc,
		ModelType:mType,
		ModelName:mName,
		TableInformation:loadTableSetting(mType),
	}

	m.models[mName] = detail
}

func (m *Model) Get(modelName string) *ModelInformation {
	return m.models[modelName]
}

func loadTableSetting(mType reflect.Type) *TableInformation {
	tInfo := &TableInformation{
		TableName:string551.SnakeCase(mType.Name()),
		PrimaryKey:"id",
		Fields:[]string{},
		DeleteTable:false,
		DeletedAtField:"",
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
			return string551.SnakeCase(sField.Name)
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
		if err != nil || ! del {
			if db == "" {
				fields = append(fields, string551.SnakeCase(sField.Name))
			} else if db != "-" {
				fields = append(fields, string551.SnakeCase(db))
			}
		}
	}

	return fields

}
