// Кузьмин Павел - август 2020
// Этот исходный код подготовлен в рамках исполнения тестового задания.

package incrementor

import (
    "testing"
    "database/sql"

    "github.com/DATA-DOG/go-sqlmock"
)

const (
	NO_INSERT_ID = 0
)

// Создание базовго мока БД
// документация по работе с макетами БД
// библиотеки sqlmock - https://godoc.org/github.com/DATA-DOG/go-sqlmock
func getMockupDb() (*sql.DB, sqlmock.Sqlmock) {
	const qselect string = "^SELECT (.+) FROM `incrementor`*"

	db, mock, err := sqlmock.New()

	idrows  := sqlmock.NewRows([]string{"id"}).AddRow(1)

	valrows := sqlmock.NewRows([]string{"val"}).
	        AddRow(0).
	        AddRow(0).
	        AddRow(1)

	if err != nil {
		panic(err.Error())
	}

	// SELECT `id`, `key`, `val` FROM `incrementor` LIMIT 1
	// запрос выполняется при создании экземпляра dbStore (newDbStore)
	mock.ExpectQuery("^SELECT `id`, `key`, `val` FROM `incrementor`*").WillReturnRows(idrows)

	// SELECT `id` FROM `incrementor` WHERE `key`=?
	// запрос выполняется при получении значений из БД при создании
	// экземпляра Incrementor (NewIncrementor)
	mock.ExpectQuery(qselect).WithArgs(KEY_MAX_VAL).WillReturnRows(valrows)
	mock.ExpectQuery(qselect).WithArgs(KEY_CUR_VAL).WillReturnRows(valrows)
	mock.ExpectQuery(qselect).WithArgs(KEY_INC_VAL).WillReturnRows(valrows)

	return db, mock
}

// Тестирование создания экземпляра Incrementor
func TestCreating(t *testing.T) {
	db, _ := getMockupDb()
	defer db.Close()

	incr := NewIncrementor(db)

	if incr == nil {
		t.Error("cannot create an instance of Incrementor")
	}
}

// Тестирование значения счетчика по умолчанию
func TestGetNumberDefault(t *testing.T) {
	db, _ := getMockupDb()
	defer db.Close()

	incr := NewIncrementor(db)

	if incr == nil {
		t.Error("cannot create an instance of Incrementor")
		return
	}

	if val := incr.GetNumber(); val != 0 {
		t.Errorf("invalid increment default value, %d want 0", val)
	}
}

// Тестирование значения приращения по умолчанию
func TestGetIncValDefault(t *testing.T) {
	db, _ := getMockupDb()
	defer db.Close()

	incr := NewIncrementor(db)

	if incr == nil {
		t.Error("cannot create an instance of Incrementor")
		return
	}

	if val := incr.GetIncVal(); val != 1 {
		t.Errorf("invalid increment default value, %d want 1", val)
	}
}

// Тестирование максимального значения по умолчанию
func TestGetMaxValDefault(t *testing.T) {
	db, _ := getMockupDb()
	defer db.Close()

	incr := NewIncrementor(db)

	if incr == nil {
		t.Error("cannot create an instance of Incrementor")
		return
	}

	if incr.GetMaxVal() != MAX_INT {
		t.Error("invalid maximum default value")
	}
}

// Тестирование установки максимального значения
func TestSetGetMaxVal(t *testing.T) {
	var want int = 1000

	db, mock := getMockupDb()
	defer db.Close()

	// ожидание возникновения одного запроса на обновление
	// максимального значения
	mock.ExpectExec("^UPDATE `incrementor`*").
		 WithArgs(want, KEY_MAX_VAL).
		 WillReturnResult( sqlmock.NewResult(NO_INSERT_ID, 1) )

	incr := NewIncrementor(db)

	if incr == nil {
		t.Error("cannot create an instance of Incrementor")
		return
	}

	if err := incr.SetMaximumValue(want); err != nil {
		t.Errorf("error set maximum value: %s", err.Error())
		return
	}

	if val := incr.GetMaxVal(); val != want {
		t.Errorf("invalid maximum value %d want %d", val, want)
		return
	}
}

// Тестирование установки максимального значения
// с обнулением счетчика при возникновении условия,
// где значение значение счетчика больше максимального значения 
func TestSetMaximumValueOverflow(t *testing.T) {
	var maxval int = 25
	var iters int  = 30

	db, mock := getMockupDb()
	defer db.Close()

	incr := NewIncrementor(db)

	if incr == nil {
		t.Error("cannot create an instance of Incrementor")
		return
	}

	// ожидание множественных запросов на обновление счетчика
	for i := 1; i <= iters; i++ {
		mock.ExpectExec("^UPDATE `incrementor`*").
			WithArgs(i, KEY_CUR_VAL).
			WillReturnResult( sqlmock.NewResult(NO_INSERT_ID, 1) )
	}

	for i := 0; i < iters; i++ {
		if err := incr.IncrementNumber(); err != nil {
			t.Errorf("error increment value: %s", err.Error())
			return
		}
	}

	// описываем транзакцию, т.к. в этом случае, когда максимальное значение
	// больше текущего, будет произведен сброс счетчика в ноль и установка
	// максимального значения, сохранение в БД производится через транзакцию
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `incrementor`*").
		 WithArgs(0, KEY_CUR_VAL).
		 WillReturnResult( sqlmock.NewResult(NO_INSERT_ID, 1) )
	mock.ExpectExec("^UPDATE `incrementor`*").
		 WithArgs(maxval, KEY_MAX_VAL).
		 WillReturnResult( sqlmock.NewResult(NO_INSERT_ID, 1) )
	mock.ExpectCommit()

	if err := incr.SetMaximumValue(maxval); err != nil {
		t.Errorf("error set maximum value: %s", err.Error())
		return
	}

	if val := incr.GetMaxVal(); val != maxval {
		t.Errorf("invalid maximum value %d want %d", val, maxval)
		return
	}
}

// Тестирование увеличения счетчика
func TestIncrementNumber(t *testing.T) {
	var iters int = 137

	db, mock := getMockupDb()
	defer db.Close()

	incr := NewIncrementor(db)

	if incr == nil {
		t.Error("cannot create an instance of Incrementor")
		return
	}

	// ожидание множественных запросов на обновление счетчика
	for i := 1; i <= iters; i++ {
		mock.ExpectExec("^UPDATE `incrementor`*").
			WithArgs(i, KEY_CUR_VAL).
			WillReturnResult( sqlmock.NewResult(NO_INSERT_ID, 1) )
	}

	for i := 1; i <= iters; i++ {
		if err := incr.IncrementNumber(); err != nil {
			t.Errorf("error increment value: %s", err.Error())
			return
		}

		if val := incr.GetNumber(); val != i {
			t.Errorf("invalid value %d want %d", val, i)
			return
		}
	}
}

// Тестирование установки значения приращения
func TestSetIncrValue(t *testing.T) {
	var incval int = 50
	var iters  int = 100

	db, mock := getMockupDb()
	defer db.Close()

	incr := NewIncrementor(db)

	if incr == nil {
		t.Error("cannot create an instance of Incrementor")
		return
	}

	// ожидание запроса на обновление значения приращения
	mock.ExpectExec("^UPDATE `incrementor`*").
		 WithArgs(incval, KEY_INC_VAL).
		 WillReturnResult( sqlmock.NewResult(NO_INSERT_ID, 1) )

	// ожидание множественных запросов на обновление счетчика
	for i := 1; i <= iters; i++ {
		mock.ExpectExec("^UPDATE `incrementor`*").
			WithArgs(i * incval, KEY_CUR_VAL).
			WillReturnResult( sqlmock.NewResult(NO_INSERT_ID, 1) )
	}

	if err := incr.SetIncrValue(incval); err != nil {
		t.Errorf("error set increment value: %s", err.Error())
		return
	}

	for i := 1; i <= iters; i++ {
		if err := incr.IncrementNumber(); err != nil {
			t.Errorf("error increment value: %s", err.Error())
			return
		}

		if val := incr.GetNumber(); val != i * incval {
			t.Errorf("invalid value %d want %d", val, i * incval)
			return
		}
	}
}

// Тестирование установки настроек, максимального значения и значения
// приращения и проверки через инкрементацию счетчика с переполнением
func TestSetSettings(t *testing.T) {
	var iters   int = 100
	var incval  int = 5
	var maxval  int = 50
	var overval int

	db, mock := getMockupDb()
	defer db.Close()

	incr := NewIncrementor(db)

	if incr == nil {
		t.Error("cannot create an instance of Incrementor")
		return
	}

	// ожидание транзакции на обновление значений настроке
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `incrementor`*").
		 WithArgs(maxval, KEY_MAX_VAL).
		 WillReturnResult( sqlmock.NewResult(NO_INSERT_ID, 1) )
	mock.ExpectExec("^UPDATE `incrementor`*").
		 WithArgs(incval, KEY_INC_VAL).
		 WillReturnResult( sqlmock.NewResult(NO_INSERT_ID, 1) )
	mock.ExpectCommit()

	if err := incr.SetSettings(maxval, incval); err != nil {
		t.Errorf("error set settings values: %s", err.Error())
		return
	}

	overval = 0

	for i := 0; i < iters; i++ {
		overval += incval
		
		if (overval > maxval) {
			overval = 0
		}
		
		// ожидание множественных запросов на обновление счетчика
		mock.ExpectExec("^UPDATE `incrementor`*").
			WithArgs(overval, KEY_CUR_VAL).
			WillReturnResult( sqlmock.NewResult(NO_INSERT_ID, 1) )
		
	}

	for i := 0; i < iters; i++ {
		if err := incr.IncrementNumber(); err != nil {
			t.Errorf("error increment value: %s", err.Error())
			return
		}
	}

	if val := incr.GetNumber(); val != overval {
		t.Errorf("invalid value %d want %d", val, overval)
		return
	}
}