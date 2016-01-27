package model551
import (
	"reflect"
	"errors"
)

type Model struct {
	models map[string]*Detail
}

type Detail struct {
	NewFunc   NewModelFunc
	ModelType reflect.Type
	ModelName string
}

var modelInstance *Model

func Load() *Model {
	if modelInstance != nil {
		return modelInstance
	}

	modelInstance = &Model{
		models:map[string]*Detail{},
	}

	return modelInstance
}

type NewModelFunc func() interface{}

func (m *Model) Add(newFunc NewModelFunc) {
	model := newFunc()

	mType := reflect.TypeOf(model)
	mName := mType.Name()

	if m.models[mName] != nil {
		panic(errors.New("追加されたモデルは既に登録されています。"))
	}

	detail := &Detail{
		NewFunc:newFunc,
		ModelType:mType,
		ModelName:mName,
	}

	m.models[mName] = detail
}

func (m *Model) Get(modelName string) *Detail {
	return m.models[modelName]
}