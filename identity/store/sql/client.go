//
// Copyright (c) Sematext International
// All Rights Reserved
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

type Client struct {
	db *sql.DB
}

// NewClient builds a new SQL client. It's necessary to import the specific database driver for
// each relational database engine in order to manipulate and access the data.
func NewClient(
	driver string,
	datasource string,
) (*Client, error) {
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

// QueryForRow executes a given query with named parameters. The query has to return one field
// per row. If no rows are fetched into the result set, an error is returned
func (c Client) QueryForRow(query string, args ...interface{}) (interface{}, error) {
	var result interface{}
	// TODO: check query structure with SQL parser
	err := c.db.QueryRow(query, args...).Scan(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}