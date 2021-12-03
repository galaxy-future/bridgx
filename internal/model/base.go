package model

import (
	"time"

	"github.com/galaxy-future/BridgX/internal/clients"
	"github.com/galaxy-future/BridgX/internal/constants"
	"github.com/galaxy-future/BridgX/internal/logs"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Base struct {
	Id       int64      `json:"id" gorm:"primary_key"`
	CreateAt *time.Time `json:"-"`
	UpdateAt *time.Time `json:"-"`
}

const (
	BATCH_SIZE = 5
)

// Create insert the value into database
func Create(model interface{}) error {
	err := clients.WriteDBCli.Create(model).Error
	if err != nil {
		logErr("create data to write db", err)
		return err
	}
	return nil
}

func CreateIgnoreDuplicate(model interface{}) error {
	if err := clients.WriteDBCli.Clauses(clause.OnConflict{DoNothing: true}).Create(model).Error; err != nil {
		logErr("create data to write db", err)
		return err
	}
	return nil
}

func BatchCreate(values interface{}) error {
	err := clients.WriteDBCli.CreateInBatches(values, BATCH_SIZE).Error
	if err != nil {
		logErr("create data to write db", err)
		return err
	}
	return nil
}

// Save update value in database, if the value doesn't have primary key, will insert it
func Save(model interface{}) error {
	err := clients.WriteDBCli.Save(model).Error
	if err != nil {
		logErr("save data to write db", err)
		return err
	}
	return nil
}

func Delete(model interface{}) error {
	err := clients.WriteDBCli.Delete(model).Error
	if err != nil {
		logErr("delete from write db", err)
		return err
	}
	return nil
}

func QueryFirst(where map[string]interface{}, out interface{}) error {
	err := clients.ReadDBCli.Where(where).First(out).Error
	if err != nil {
		logErr("query data from read db", err)
		return err
	}
	return nil
}

// Query records
func Query(where map[string]interface{}, page int, pageSize int, models interface{}, order string, withCount bool) (count int64, err error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > constants.DefaultPageSize {
		pageSize = constants.DefaultPageSize
	}
	offset := (page - 1) * pageSize
	query := clients.ReadDBCli.Where(where).Model(models).Order(order).Offset(offset).Limit(pageSize).Find(models)
	if err := query.Error; err != nil {
		logErr("query data from read db", err)
		return 0, err
	}
	if withCount {
		if err = query.Offset(-1).Limit(-1).Count(&count).Error; err != nil {
			logErr("query data from read db", err)
			return 0, err
		}
		return count, nil
	}
	return 0, nil
}

// Query records
func QueryWhere(where *gorm.DB, page int, pageSize int, models interface{}, order string, withCount bool) (count int64, err error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > constants.DefaultPageSize {
		pageSize = constants.DefaultPageSize
	}
	offset := (page - 1) * pageSize
	query := clients.ReadDBCli.Where(where).Order(order).Offset(offset).Limit(pageSize).Find(models)
	if err := query.Error; err != nil {
		logErr("query data from read db", err)
		return 0, err
	}
	if withCount {
		if err = query.Offset(-1).Limit(-1).Count(&count).Error; err != nil {
			logErr("query data from read db", err)
			return 0, err
		}
		return count, nil
	}
	return 0, nil
}

// QueryAll records
func QueryAll(where map[string]interface{}, models interface{}, order string) (err error) {
	if err = clients.ReadDBCli.Where(where).Order(order).Find(models).Error; err != nil {
		logErr("query all from read db", err)
		return err
	}
	return nil
}

// Updates update attributes with callbacks
func Updates(model interface{}, ids []int64, updates map[string]interface{}) error {
	if err := clients.WriteDBCli.Model(model).Where("id IN (?)", ids).Updates(updates).Error; err != nil {
		logErr("update data list to write db", err)
		return err
	}
	return nil
}

// Get find first record that match given conditions, order by primary key
func Get(id int64, out interface{}) error {
	if err := clients.ReadDBCli.Where("id = ?", id).First(out).Error; err != nil {
		logErr("get data from read db", err)
		return err
	}
	return nil
}

// Gets find records that match given conditions
func Gets(ids []int64, out interface{}) error {
	if err := clients.ReadDBCli.Where(ids).Find(out).Error; err != nil {
		logErr("get data list from read db", err)
		return err
	}
	return nil
}

// Count records
func Count(where map[string]interface{}, model interface{}) (int64, error) {
	var cnt int64
	if err := clients.ReadDBCli.Model(model).Where(where).Count(&cnt).Error; err != nil {
		logErr("query all from read db", err)
		return 0, err
	}
	return cnt, nil
}

func logErr(errType string, err error) {
	logs.Logger.Error(errType, err)
}
