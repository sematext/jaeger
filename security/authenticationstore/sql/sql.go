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
	sec"github.com/uber/jaeger/security"
	"go.uber.org/zap"
	"database/sql"
)

// DbAuthenticationStore describes the store based on relational database engine
type DbAuthenticationStore struct {
	client DbClient
	logger *zap.Logger
	query  string
}

// NewDbAuthenticationStore creates a new instance of the SQL-based authentication store. It attempts to
// ping the database server to ensure it's available.
func NewDbAuthenticationStore(
	client DbClient,
	logger *zap.Logger,
	query string,
) (*DbAuthenticationStore, error) {
	if err := client.Ping(); err != nil {
		return nil, err
	}
	return &DbAuthenticationStore{
		client: client,
		logger: logger,
		query: query,
	}, nil
}

// FindPrincipal attempts to find the principal associated with an authentication token. The result of the
// SQL select sentence has to contain three fields describing the user name, the password and the status of the
// account (active|locked).
func (db *DbAuthenticationStore) FindPrincipal(token sec.AuthenticationToken) (*sec.AuthenticationContext, error) {
	ctx, err := db.client.QueryForRow(db.query, func(row *sql.Row) (interface{}, error) {
		var ctx sec.AuthenticationContext
		err := row.Scan(&ctx.Principal, &ctx.Password, &ctx.Locked)
		if err != nil {
			return nil, err
		}
		return ctx, nil
	}, token.Username)
	if err != nil {
		return nil, err
	}
	c, _ := ctx.(sec.AuthenticationContext)
	return &c, nil
}

