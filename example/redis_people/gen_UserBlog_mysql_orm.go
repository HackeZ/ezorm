package test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ezbuy/ezorm/db"
)

var (
	_ time.Time
	_ bytes.Buffer
	_ = strings.Index
)

// -----------------------------------------------------------------------------

func (m *_UserBlogMgr) queryOne(ctx context.Context, query string, args ...interface{}) (*UserBlog, error) {
	ret, err := m.queryLimit(ctx, query, 1, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}
	return ret[0], nil
}

func (m *_UserBlogMgr) query(ctx context.Context, query string, args ...interface{}) (results []*UserBlog, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (m *_UserBlogMgr) Query(query string, args ...interface{}) (results []*UserBlog, err error) {
	return m.QueryContext(context.Background(), query, args...)
}

func (m *_UserBlogMgr) QueryContext(ctx context.Context, query string, args ...interface{}) (results []*UserBlog, err error) {
	return m.queryLimit(ctx, query, -1, args...)
}

func (*_UserBlogMgr) queryLimit(ctx context.Context, query string, limit int, args ...interface{}) (results []*UserBlog, err error) {
	rows, err := db.MysqlQueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ezorm.UserBlog query error: %v", err)
	}
	defer rows.Close()

	offset := 0
	for rows.Next() {
		if limit >= 0 && offset >= limit {
			break
		}
		offset++

		var result UserBlog
		err := rows.Scan(&(result.Key),
			&(result.Value),
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ezorm.UserBlog fetch result error: %v", err)
	}

	return
}
