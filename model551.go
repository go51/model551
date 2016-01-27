package model551

type Model struct {}

var modelInstance *Model

func Load() *Model {
	if modelInstance != nil {
		return modelInstance
	}

	modelInstance = &Model{}

	return modelInstance
}