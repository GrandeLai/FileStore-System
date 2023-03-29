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

var _ UserModel = (*customUserModel)(nil)

var userRowsExpectAutoSetButIdAndTotalVolume = strings.Join(stringx.Remove(userFieldNames, "`total_volume`", "`create_at`", "`created_at`", "`create_time`", "`update_at`", "`updated_at`", "`update_time`"), ",")

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		InsertWithId(ctx context.Context, data *User) (sql.Result, error)
		JudgeUserExist(ctx context.Context, name string, password string) (*User, error)
		FindByName(ctx context.Context, name string) (*User, error)
		UpdateVolume(ctx context.Context, id int64, size int64) (result sql.Result, err error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c),
	}
}

func (m *defaultUserModel) InsertWithId(ctx context.Context, data *User) (sql.Result, error) {
	userEmailKey := fmt.Sprintf("%s%v", cacheUserEmailPrefix, data.Email)
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	userNameKey := fmt.Sprintf("%s%v", cacheUserUserNamePrefix, data.UserName)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, userRowsExpectAutoSetButIdAndTotalVolume)
		return conn.ExecCtx(ctx, query, data.Id, data.UserName, data.UserPwd, data.Email, data.NowVolume)
	}, userEmailKey, userIdKey, userNameKey)
	return ret, err
}

func (m *defaultUserModel) JudgeUserExist(ctx context.Context, name string, password string) (*User, error) {
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.Where("name = ?", name).Where("password = ?", password).ToSql()
	var resp User
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

func (m *defaultUserModel) FindByName(ctx context.Context, name string) (*User, error) {
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.Where("name = ?", name).ToSql()
	var resp User
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

func (m *defaultUserModel) UpdateVolume(ctx context.Context, id int64, size int64) (result sql.Result, err error) {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	userEmailKey := fmt.Sprintf("%s%v", cacheUserEmailPrefix, data.Email)
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	userNameKey := fmt.Sprintf("%s%v", cacheUserUserNamePrefix, data.UserName)
	res, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set now_volume = now_volume + ? where `id` = ? and `now_volume` + ? <= `total_volume`", m.table)
		return conn.ExecCtx(ctx, query, size, id, size)
	}, userEmailKey, userIdKey, userNameKey)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *defaultUserModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(userRows).From(m.table)
}

func (m *defaultUserModel) UpdateBuilder() squirrel.UpdateBuilder {
	return squirrel.Update(m.table)
}
