package model

import (
	"FileStore-System/service/user/model"
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
	"strings"
)

var _ UserStoreModel = (*customUserStoreModel)(nil)

var userStoreRowsExpectAutoSetButId = strings.Join(stringx.Remove(userStoreFieldNames, "`update_at`", "`updated_at`", "`update_time`", "`create_at`", "`created_at`", "`create_time`"), ",")

type (
	// UserStoreModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserStoreModel.
	UserStoreModel interface {
		userStoreModel
		FindByFactors(ctx context.Context, parentId int64, userId int64, storeId int64) (*UserStore, error)
		FindAllInPage(ctx context.Context, parentId int64, userId int64, startIndex int64, pageSize int64) ([]*UserStore, error)
		CountStoreUsage(ctx context.Context, storeId int64) (int64, error)
		CountByIdAndParentId(ctx context.Context, id int64, userId int64, parentId int64) (int64, error)
		DeleteByParentId(ctx context.Context, parentId int64) error
		CountByParentIdAndName(ctx context.Context, parentId int64, userId int64, Name string) (int64, error)
		InsertByBatch(ctx context.Context, data []*UserStore) (sql.Result, error)
		InsertWithId(ctx context.Context, data *UserStore) (sql.Result, error)
	}

	customUserStoreModel struct {
		*defaultUserStoreModel
	}
)

// NewUserStoreModel returns a model for the database table.
func NewUserStoreModel(conn sqlx.SqlConn, c cache.CacheConf) UserStoreModel {
	return &customUserStoreModel{
		defaultUserStoreModel: newUserStoreModel(conn, c),
	}
}

func (m *defaultUserStoreModel) FindByFactors(ctx context.Context, parentId int64, userId int64, storeId int64) (*UserStore, error) {
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.Where("parent_id = ?", parentId).Where("user_id = ?", userId).Where("store_id = ?", storeId).ToSql()
	var resp UserStore
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, model.ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserStoreModel) FindAllInPage(ctx context.Context, parentId int64, userId int64, startIndex int64, pageSize int64) ([]*UserStore, error) {
	var resp []*UserStore
	rowBuilder := m.RowBuilder()
	rowBuilder = rowBuilder.Where("parent_id = ?", parentId).Where("user_id = ?", userId)
	if startIndex != 0 || pageSize != 0 {
		rowBuilder = rowBuilder.Offset(uint64(startIndex)).Limit(uint64(pageSize))
	}
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, model.ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserStoreModel) CountStoreUsage(ctx context.Context, storeId int64) (int64, error) {
	countBuilder := m.CountBuilder("store_id")
	query, values, err := countBuilder.Where("store_id = ?", storeId).ToSql()
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return 0, model.ErrNotFound
	default:
		return 0, err
	}
}

func (m *defaultUserStoreModel) CountByIdAndParentId(ctx context.Context, id int64, userId int64, parentId int64) (int64, error) {
	countBuilder := m.CountBuilder("id")
	query, values, err := countBuilder.Where("parent_id = ?", parentId).Where("user_id = ?", userId).Where("id = ?", id).ToSql()
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return 0, model.ErrNotFound
	default:
		return 0, err
	}
}

func (m *defaultUserStoreModel) CountByParentIdAndName(ctx context.Context, parentId int64, userId int64, Name string) (int64, error) {
	countBuilder := m.CountBuilder("id")
	query, values, err := countBuilder.Where("parent_id = ?", parentId).Where("user_id = ?", userId).Where("name = ?", Name).ToSql()
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return 0, model.ErrNotFound
	default:
		return 0, err
	}
}

func (m *defaultUserStoreModel) DeleteByParentId(ctx context.Context, parentId int64) error {
	userStoreParentIdKey := fmt.Sprintf("%s%v", cacheUserStoreParentIdPrefix, parentId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `parent_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, parentId)
	}, userStoreParentIdKey)
	return err
}

func (m *defaultUserStoreModel) InsertByBatch(ctx context.Context, data []*UserStore) (sql.Result, error) {
	// 存放 (?, ?) 的slice
	valueStrings := make([]string, 0, len(data))
	// 存放values的slice
	valueArgs := make([]interface{}, 0, len(data)*2)
	// 遍历users准备相关数据
	for _, u := range data {
		// 此处占位符要与插入值的个数对应
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, u.Id)
		valueArgs = append(valueArgs, u.UserId)
		valueArgs = append(valueArgs, u.ParentId)
		valueArgs = append(valueArgs, u.StoreId)
		valueArgs = append(valueArgs, u.Name)
	}
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values %s", m.table, storeRowsExpectAutoSetButId, strings.Join(valueStrings, ","))
		return conn.ExecCtx(ctx, query, data)
	})
	return ret, err
}

func (m *defaultUserStoreModel) InsertWithId(ctx context.Context, data *UserStore) (sql.Result, error) {
	userStoreIdKey := fmt.Sprintf("%s%v", cacheUserStoreIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, userStoreRowsExpectAutoSetButId)
		return conn.ExecCtx(ctx, query, data.Id, data.UserId, data.ParentId, data.StoreId, data.Name)
	}, userStoreIdKey)
	return ret, err
}

func (m *defaultUserStoreModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(userStoreRows).From(m.table)
}

func (m *defaultUserStoreModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("count(" + field + ")").From(m.table)
}
