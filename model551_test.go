package model551_test
import (
	"testing"
	"github.com/go51/model551"
)

func TestLoad(t *testing.T) {
	m1 := model551.Load()
	m2 := model551.Load()

	if m1 == nil {
		t.Errorf("インスタンスの生成に失敗しました。")
	}
	if m2 == nil {
		t.Errorf("インスタンスの生成に失敗しました。")
	}
}

func BenchmarkLoad(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model551.Load()
	}
}

type SampleModel struct {
	Name        string
	Description string
}

func NewSampleModel() interface{} {
	return SampleModel{}
}

func TestAdd(t *testing.T) {
	m := model551.Load()

	m.Add(NewSampleModel)

	detail := m.Get("SampleModel")

	if detail.ModelName != "SampleModel" {
		t.Errorf("モデルの保持に失敗しました。")
	}
	if detail.ModelType.Name() != "SampleModel" {
		t.Errorf("モデルの保持に失敗しました。")
	}
}