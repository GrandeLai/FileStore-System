package model

import (
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

var _ StoreModel = (*customStoreModel)(nil)

var storeRowsExpectAutoSetButId = strings.Join(stringx.Remove(storeFieldNames, "`updated_at`", "`update_time`", "`create_at`", "`created_at`", "`create_time`", "`update_at`"), ",")

type (
	// StoreModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStoreModel.
	StoreModel interface {
		storeModel
		FindFolderByParentId(ctx context.Context, parentId int64) (*Store, error)
		InsertWithId(ctx context.Context, data *Store) (sql.Result, error)
		CountByHash(ctx context.Context, hash string) (int64, error)
		FindStoreByHash(ctx context.Context, hash string) (*Store, error)
	}

	customStoreModel struct {
		*defaultStoreModel
	}
)

// NewStoreModel returns a model for the database table.
func NewStoreModel(conn sqlx.SqlConn, c cache.CacheConf) StoreModel {
	return &customStoreModel{
		defaultStoreModel: newStoreModel(conn, c),
	}
}

func (m *defaultStoreModel) InsertWithId(ctx context.Context, data *Store) (sql.Result, error) {
	storeHashKey := fmt.Sprintf("%s%v", cacheStoreHashPrefix, data.Hash)
	storeIdKey := fmt.Sprintf("%s%v", cacheStoreIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, storeRowsExpectAutoSetButId)
		return conn.ExecCtx(ctx, query, data.Id, data.Hash, data.Size, data.Ext, data.Path, data.Name, data.Status, data.IsFolder)
	}, storeHashKey, storeIdKey)
	return ret, err
}

func (m *defaultStoreModel) FindFolderByParentId(ctx context.Context, parentId int64) (*Store, error) {
	storeIdKey := fmt.Sprintf("%s%v", cacheStoreIdPrefix, parentId)
	var resp Store
	err := m.QueryRowCtx(ctx, &resp, storeIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", storeRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, parentId)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultStoreModel) CountByHash(ctx context.Context, hash string) (int64, error) {
	countBuilder := m.CountBuilder("id")
	query, values, err := countBuilder.Where("hash = ?", hash).ToSql()
	if err != nil {
		return 0, err
	}
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

func (m *defaultStoreModel) FindStoreByHash(ctx context.Context, hash string) (*Store, error) {
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.Where("hash = ?", hash).ToSql()
	if err != nil {
		return nil, err
	}
	var resp Store
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultStoreModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(storeRows).From(m.table)
}

func (m *defaultStoreModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("count(" + field + ")").From(m.table)
}
