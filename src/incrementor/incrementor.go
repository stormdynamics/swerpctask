// Кузьмин Павел - август 2020
// Этот исходный код подготовлен в рамках исполнения тестового задания.

package incrementor

import (
	"log"
	"sync"
	"errors"
	"database/sql"
)

const (
	// Максимальное значение знакового целого int
	MAX_INT int = int(^uint(0) >> 1) 

	// Ключи для значений в таблице БД
	KEY_MAX_VAL = "max_val"
	KEY_CUR_VAL = "cur_val"
	KEY_INC_VAL = "inc_val"
)

type Incrementor struct {
	// Мьютекс для потокобезопасного чтения
	// и обновления значений
	mux sync.RWMutex

	// Экземпляр структуры хелпера
	// для манипуляций данными
	db  *dbStore

	// Максимальное значение для счетчика
	maxval  int

	// Текущее значение, счетчик
	counter int

	// Значение приращения
	incval  int
}

// Возвращает максимальное значение для счетчика
func (self *Incrementor) GetMaxVal() int {
	self.mux.RLock()
	defer self.mux.RUnlock()

	return self.maxval
}

// Возвращает текущее значение для счетчика
func (self *Incrementor) GetNumber() int {
	self.mux.RLock()
	defer self.mux.RUnlock()

	return self.counter
}

// Возвращает значение приращения
func (self *Incrementor) GetIncVal() int {
	self.mux.RLock()
	defer self.mux.RUnlock()

	return self.incval
}

// Увеличение значение счетчика на значение приращения.
// Обновляет новое значение счетчика в БД
func (self *Incrementor) IncrementNumber() error {
	var err error

	self.mux.Lock()
	defer self.mux.Unlock()

	// безопасная проверка переполнения счетчика
	// если значение переполняется оно обнуляется
	if self.maxval - self.counter - self.incval < 0 {
		if err = self.db.setInt(KEY_CUR_VAL, 0); err == nil {
			self.counter = 0
		}
	} else {
		if err = self.db.setInt(KEY_CUR_VAL, self.counter + self.incval); err == nil {
			self.counter += self.incval
		}
	}

	return err
}

// Устанавливает максимальное значение, до которого можно увеличивать
// значение счетчика.
func (self *Incrementor) SetMaximumValue(max int) error {
	var tdx *dbStore = self.db
	var err error
	var counter int

	if max <= 0 {
		return errors.New("maximum value must be greater than zero")
	}

	self.mux.Lock()
	defer self.mux.Unlock()

	counter = self.counter

	// если при установке нового максимального значения
	// значение счетчика больше, то обнуляем значение счетчика
	if counter > max {

		// т.к. необходимо обновить сразу два значение, то
		// обновление происходит в режиме транзакции
		if tdx, err = self.db.txBegin(); err != nil {
			return err
		}
		defer tdx.txRollback()

		if err = tdx.setInt(KEY_CUR_VAL, 0); err != nil {
			return err
		}

		counter = 0
	}

	if err = tdx.setInt(KEY_MAX_VAL, max); err != nil {
		return err
	}
	
	if tdx != self.db {
		if err = tdx.txCommit(); err != nil {
			return err
		}
	}

	self.maxval  = max
	self.counter = counter

	return nil
}

// Установка значения приращения.
// Обновляет новое значение приращения в БД
func (self *Incrementor) SetIncrValue(incval int) error {

	if incval <= 0 {
		return errors.New("increment value must be greater than zero")
	}

	if incval == MAX_INT {
		return errors.New("increment value must be less than maximum integer value")
	}

	self.mux.Lock()
	defer self.mux.Unlock()

	err := self.db.setInt(KEY_INC_VAL, incval)

	if err == nil {
		self.incval = incval
	}

	return err
}

// Обновление максимального значения счетчика и значения приращения.
// Сохранение в БД происходит через транзакцию
func (self *Incrementor) SetSettings(maxval, incval int) error {
	var tdx *dbStore
	var err error
	var counter int

	if maxval <= 0 {
		return errors.New("maximum value must be greater than zero")
	}

	if incval <= 0 {
		return errors.New("increment value must be greater than zero")
	}

	if incval == MAX_INT {
		return errors.New("increment value must be less than maximum integer value")
	}

	self.mux.Lock()
	defer self.mux.Unlock()

	counter = self.counter;

	if tdx, err = self.db.txBegin(); err != nil {
		return err
	}
	defer tdx.txRollback()

	if err = tdx.setInt(KEY_MAX_VAL, maxval); err != nil {
		return err
	}

	if err = tdx.setInt(KEY_INC_VAL, incval); err != nil {
		return err
	}

	if self.counter > maxval {
		if err = tdx.setInt(KEY_CUR_VAL, 0); err != nil {
			return err
		}
		counter = 0;
	}

	if err = tdx.txCommit(); err != nil {
		return err;
	}

	self.maxval  = maxval
	self.incval  = incval
	self.counter = counter

	return nil
}

// Получение значений счетчика, максимального значения и
// значения приращения из БД c проверкой.
// Если значение получено, но оно не проходит проверку,
// то будет использовано значение по умолчания
// и выполнена запись в лог
func (self *Incrementor) loadValsFromDb() error {
	var err error
	var val int

	if val, err = self.db.getInt(KEY_MAX_VAL); err != nil {
		return err
	} else if val > 0 {
		self.maxval = val
	} else if val != 0 {
		log.Printf("ivalid value loaded from DB `max_val`=%d, set by default\n", val)
	}

	if val, err = self.db.getInt(KEY_CUR_VAL); err != nil {
		return err
	} else if val >= 0 {
		self.counter = val
	} else {
		log.Printf("ivalid value loaded from DB `cur_val`=%d, set by default\n", val)
	}

	if val, err = self.db.getInt(KEY_INC_VAL); err != nil {
		return err
	} else if val > 0 && val != MAX_INT {
		self.incval = val
	} else {
		log.Printf("ivalid value loaded from DB `inc_val`=%d, set by default\n", val)
	}

	return nil
}

// Создает новый экземпляр Incrementor
func NewIncrementor(sql_db *sql.DB) *Incrementor {

	db, err := newDbStore(sql_db)

	if err != nil {
		log.Fatalf("error `incrementor` database: %s", err.Error())
	}

	incr := &Incrementor{
		db:      db,
		maxval:  MAX_INT,
		counter: 0,
		incval:  1,
	}

	// Получаем значения из БД
	err = incr.loadValsFromDb()

	if err != nil {
		log.Fatalf("load values from db error: %s", err.Error())
	}

	return incr
}