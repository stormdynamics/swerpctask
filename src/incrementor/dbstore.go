// Кузьмин Павел - август 2020
// Этот исходный код подготовлен в рамках исполнения тестового задания.

package incrementor

import (
	"errors"
	"database/sql"
)

var (
	errTxNotOpened = errors.New("transaction not opened")
	errTxOpened    = errors.New("transaction already opened")
	errMultiAff    = errors.New("expected to affect 1 row, affected")
)

// Предоставляет возможность обновлять и получать
// данные из таблицы `incrementor`
type dbStore struct {
	// Указатель на основной экземпляр БД
	db  *sql.DB

	// Указатель на экземпляр транзакции
	tx  *sql.Tx
}

// Открытие транзакции и возврат нового экземпляра dbStore
// если у текущего экземпляра транзакция уже открыта то вернет ошибку
func (self *dbStore) txBegin() (*dbStore, error) {
	
	if self.tx == nil {
		tx, err := self.db.Begin()

		if err != nil {
			return nil, err
		}

		return &dbStore{self.db, tx}, nil
	}

	return nil, errTxOpened
}

// Осуществление транзакции в случае открытой транзакции у текущего экземпляра
func (self *dbStore) txCommit() error {
	if self.tx != nil {
		return self.tx.Commit()
	}

	return errTxNotOpened
}

// Отмена транзакции в случае открытой транзакции у текущего экземпляра
func (self *dbStore) txRollback() error {
	if self.tx != nil {
		return self.tx.Rollback()
	}

	return errTxNotOpened
}

// Обновляет целое значение по ключу в таблице `incrementor`
// может вернуть ошибку в случае плохого запроса или
// если запрос возвращает результат количества затронутых
// строк больше 1
func (self *dbStore) setInt(key string, val int) error {
	var res sql.Result = nil
	var err error      = nil
	const qstr string  = "UPDATE `incrementor` SET `val` = ? WHERE `key` = ?"

	// в зависимости от того, начата или нет транзакция
	// запрос выполяется на экземпляре sql.DB или sql.Tx
	if self.tx == nil {
		res, err = self.db.Exec(qstr, val, key)
	} else {
		res, err = self.tx.Exec(qstr, val, key)
	}

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()

	if err != nil {
		return err
	}

	// количество затронутых записей должно быть равным единице
	if count != 1 { 
		return errMultiAff
	}

	return nil
}

// Возвращает из таблицы `incrementor` знаковое целое 
// по параметру ключа
func (self *dbStore) getInt(key string) (int, error) {
	var val int       = 0
	var err error     = nil
	const qstr string = "SELECT `val` FROM `incrementor` WHERE `key`=?"
	
	// в зависимости от того, начата или нет транзакция
	// запрос выполяется на экземпляре sql.DB или sql.Tx
	if self.tx == nil {
		err = self.db.QueryRow(qstr, key).Scan(&val)
	} else {
		err = self.tx.QueryRow(qstr, key).Scan(&val)
	}

	return val, err
}

// Создает новый экземпляр dbStore
// т.к. функция делает некоторые проверки входных параметров,
// тестовую выборку из таблицы incrementor, то возвращает два значения
// ссылку на экземпляр dbStore и ошибку
func newDbStore(sql_db *sql.DB) (*dbStore, error) {

	if sql_db == nil {
		return nil, errors.New("error: sql.DB is nil")
	}

	// пустой запрос к БД для тестирования подключения
	// и проверки существования таблицы `incrementor`
	q, err := sql_db.Query("SELECT `id`, `key`, `val` FROM `incrementor` LIMIT 1")
	
	if err != nil {
		return nil, err
	}

	q.Close()

	return &dbStore{ sql_db, nil }, nil
}