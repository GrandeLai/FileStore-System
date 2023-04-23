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

var _ ShareModel = (*customShareModel)(nil)

var shareRowsExpectAutoSetButId = strings.Join(stringx.Remove(shareFieldNames, "`create_time`", "`update_at`", "`updated_at`", "`update_time`", "`create_at`", "`created_at`"), ",")

type (
	// ShareModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShareModel.
	ShareModel interface {
		shareModel
		InsertWithId(ctx context.Context, data *Share) (sql.Result, error)
		InsertByBatch(ctx context.Context, data []*Share) (sql.Result, error)
		FindByShareURL(ctx context.Context, shareUrl string) ([]*Share, error)
	}

	customShareModel struct {
		*defaultShareModel
	}
)

// NewShareModel returns a model for the database table.
func NewShareModel(conn sqlx.SqlConn, c cache.CacheConf) ShareModel {
	return &customShareModel{
		defaultShareModel: newShareModel(conn, c),
	}
}

func (m *defaultShareModel) InsertWithId(ctx context.Context, data *Share) (sql.Result, error) {
	shareIdKey := fmt.Sprintf("%s%v", cacheShareIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ? ,? )", m.table, shareRowsExpectAutoSetButId)
		return conn.ExecCtx(ctx, query, data.Id, data.UserId, data.StoreId, data.ExpiredTime, data.ShareUrl, data.ExtractionCode)
	}, shareIdKey)
	return ret, err
}

func (m *defaultShareModel) InsertByBatch(ctx context.Context, data []*Share) (sql.Result, error) {
	// 存放 (?, ?) 的slice
	valueStrings := make([]string, 0, len(data))
	// 存放values的slice
	valueArgs := make([]interface{}, 0, len(data)*2)
	// 遍历users准备相关数据
	for _, u := range data {
		// 此处占位符要与插入值的个数对应
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, u.Id)
		valueArgs = append(valueArgs, u.UserId)
		valueArgs = append(valueArgs, u.StoreId)
		valueArgs = append(valueArgs, u.ExpiredTime)
		valueArgs = append(valueArgs, u.ShareUrl)
		valueArgs = append(valueArgs, u.ExtractionCode)
	}
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values %s", m.table, shareRowsExpectAutoSetButId, strings.Join(valueStrings, ","))
		return conn.ExecCtx(ctx, query, valueArgs...)
	})
	return ret, err
}

func (m *defaultShareModel) FindByShareURL(ctx context.Context, shareUrl string) ([]*Share, error) {
	var resp []*Share
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.Where("share_url = ?", shareUrl).ToSql()
	if err != nil {
		return nil, err
	}
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultShareModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(shareRows).From(m.table)
}
