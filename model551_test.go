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
	Id          int64
	Name        string
	Description string
}

type SampleModelTableInfo struct {
	Key         int64    `db_table:"table_information" db_pk:"true" db:"id"`
	Name        string   `db:"-"`
	Description string
	DeletedAt   string   `db_delete:"true"`
}

func NewSampleModel() interface{} {
	return SampleModel{}
}

func NewSampleModelTableInfo() interface{} {
	return SampleModelTableInfo{}
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

func TestTableInfo(t *testing.T) {
	m := model551.Load()

	m.Add(NewSampleModelTableInfo)

	mSampleModel := m.Get("SampleModel")
	mSampleModelTableInfo := m.Get("SampleModelTableInfo")

	// table name
	if mSampleModel.TableInformation.TableName != "sample_model" {
		t.Errorf("テーブル名の解析に失敗しました。\n\"sample_model\" => \"%s\"\n", mSampleModel.TableInformation.TableName)
	}
	if mSampleModelTableInfo.TableInformation.TableName != "table_information" {
		t.Errorf("テーブル名の解析に失敗しました。\n\"table_information\" => \"%s\"\n", mSampleModelTableInfo.TableInformation.TableName)
	}

	// primary key
	if mSampleModel.TableInformation.PrimaryKey != "id" {
		t.Errorf("プライマリーキーの解析に失敗しました。\n\"id\" => \"%s\"\n", mSampleModel.TableInformation.PrimaryKey)
	}
	if mSampleModelTableInfo.TableInformation.PrimaryKey != "id" {
		t.Errorf("プライマリーキーの解析に失敗しました。\n\"id\" => \"%s\"\n", mSampleModelTableInfo.TableInformation.PrimaryKey)
	}

	// delete table
	if mSampleModel.TableInformation.DeleteTable {
		t.Errorf("論理削除の解析に失敗しました。\n\"false\" => \"ture\"\n")
	}
	if !mSampleModelTableInfo.TableInformation.DeleteTable {
		t.Errorf("論理削除の解析に失敗しました。\n\"true\" => \"false\"\n")
	}
	if mSampleModelTableInfo.TableInformation.DeletedAtField != "deleted_at" {
		t.Errorf("論理削除の解析に失敗しました。\n\"deleted_at\" => \"%s\"\n", mSampleModelTableInfo.TableInformation.DeletedAtField)
	}

	// Field
	if len(mSampleModel.TableInformation.Fields) != 3 {
		t.Errorf("フィールドの解析に失敗しました。\n\"3\" => \"%d\"\nDump: %#v\n", len(mSampleModel.TableInformation.Fields), mSampleModel.TableInformation.Fields)
	}
	if len(mSampleModelTableInfo.TableInformation.Fields) != 2 {
		t.Errorf("フィールドの解析に失敗しました。\n\"2\" => \"%d\"\nDump: %#v\n", len(mSampleModelTableInfo.TableInformation.Fields), mSampleModelTableInfo.TableInformation.Fields)
	}
	if mSampleModelTableInfo.TableInformation.Fields[0] != "id" {
		t.Errorf("フィールドの解析に失敗しました。\n\"id\" => \"%s\"\nDump: %#v\n", mSampleModelTableInfo.TableInformation.Fields[0], mSampleModelTableInfo.TableInformation.Fields)
	}
}

func TestSql(t *testing.T) {
	m := model551.Load()

	mSampleModel := m.Get("SampleModel")
	mSampleModelTableInfo := m.Get("SampleModelTableInfo")

	// Insert
	sql := "INSERT INTO `sample_model` (`name`, `description`) VALUES (?, ?)"
	if mSampleModel.SqlInformation.Insert != sql {
		t.Errorf("SQLキャッシュが失敗しました。\nOK: %s\nNG: %s\n", sql, mSampleModel.SqlInformation.Insert)
	}
	sql = "INSERT INTO `table_information` (`description`) VALUES (?)"
	if mSampleModelTableInfo.SqlInformation.Insert != sql {
		t.Errorf("SQLキャッシュが失敗しました。\nOK: %s\nNG: %s\n", sql, mSampleModelTableInfo.SqlInformation.Insert)
	}

	// Select
	sql = "SELECT `id`, `name`, `description` FROM `sample_model` WHERE 1 = 1"
	if mSampleModel.SqlInformation.Select != sql {
		t.Errorf("SQLキャッシュが失敗しました。\nOK: %s\nNG: %s\n", sql, mSampleModel.SqlInformation.Select)
	}
	sql = "SELECT `id`, `description` FROM `table_information` WHERE 1 = 1"
	if mSampleModelTableInfo.SqlInformation.Select != sql {
		t.Errorf("SQLキャッシュが失敗しました。\nOK: %s\nNG: %s\n", sql, mSampleModelTableInfo.SqlInformation.Select)
	}

	// Update
	sql = "UPDATE `sample_model` SET `name` = ?, `description` = ? WHERE `id` = ?"
	if mSampleModel.SqlInformation.Update != sql {
		t.Errorf("SQLキャッシュが失敗しました。\nOK: %s\nNG: %s\n", sql, mSampleModel.SqlInformation.Update)
	}
	sql = "UPDATE `table_information` SET `description` = ? WHERE `id` = ?"
	if mSampleModelTableInfo.SqlInformation.Update != sql {
		t.Errorf("SQLキャッシュが失敗しました。\nOK: %s\nNG: %s\n", sql, mSampleModelTableInfo.SqlInformation.Update)
	}

	// Delete
	sql = "DELETE FROM `sample_model` WHERE `id` = ?"
	if mSampleModel.SqlInformation.Delete != sql {
		t.Errorf("SQLキャッシュが失敗しました。\nOK: %s\nNG: %s\n", sql, mSampleModel.SqlInformation.Delete)
	}
	sql = "DELETE FROM `table_information` WHERE `id` = ?"
	if mSampleModelTableInfo.SqlInformation.Delete != sql {
		t.Errorf("SQLキャッシュが失敗しました。\nOK: %s\nNG: %s\n", sql, mSampleModelTableInfo.SqlInformation.Delete)
	}

	// Delete
	sql = ""
	if mSampleModel.SqlInformation.LogicalDelete != sql {
		t.Errorf("SQLキャッシュが失敗しました。\nOK: %s\nNG: %s\n", sql, mSampleModel.SqlInformation.LogicalDelete)
	}
	sql = "INSERT INTO `table_information_delete` (`id`, `description`, `deleted_at`) VALUES (?, ?, ?)"
	if mSampleModelTableInfo.SqlInformation.LogicalDelete != sql {
		t.Errorf("SQLキャッシュが失敗しました。\nOK: %s\nNG: %s\n", sql, mSampleModelTableInfo.SqlInformation.LogicalDelete)
	}

}

