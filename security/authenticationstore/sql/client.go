//
// Copyright (c) Sematext International
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type RowMapper func(*sql.Row) (interface{}, error)

type DbClient interface {
	Ping() error
	QueryForRow(query string, mapper RowMapper, args ...interface{}) (interface{}, error)
}

type Client struct {
	db *sql.DB
}

// NewClient builds a new SQL client. The driver for a specific database engine needs to be
// registered in order to manipulate and access the data.
func NewClient(
	driver string,
	datasource string,
) (DbClient, error) {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		return nil, err
	}
	return &Client{
		db: db,
	}, nil
}

func (c Client) Ping() error {
	return c.db.Ping()
}

// QueryForRow executes a given query. RowMapper function is called on the fetched row to scan
// the acquired fields and bind them to variables or to fields of the structure.
func (c Client) QueryForRow(query string, mapper RowMapper, args ...interface{}) (interface{}, error) {
	row := c.db.QueryRow(query, args...)
	result, err := mapper(row)
	if err != nil {
		return nil, err
	}
	return result, nil
}