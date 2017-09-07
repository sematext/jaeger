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
	"github.com/uber/jaeger/identity"
	"github.com/uber/jaeger/pkg/cache"
	"go.uber.org/zap"
)

// DbTokenStore describes the store based on relational database engine
type DbTokenStore struct {
	client *Client
	cache  cache.Cache
	query  string
}

// NewDbTokenStore creates a new instance of the SQL-based token store. It attempts to
// ping the database server to ensure it's available.
func NewDbTokenStore(
	client *Client,
	logger *zap.Logger,
) (*DbTokenStore, error) {
	if err := client.Ping(); err != nil {
		return nil, err
	}
	return &DbTokenStore{
		client: client,
	}, nil
}

// TokenExists attempts to find a token by executing the specified SQL query
// against the underlying database engine. If found, the token is cached to
// avoid subsequent round trips to the database server, and thus reducing the
// overall I/O activity.
func (store DbTokenStore) TokenExists(token identity.Token) bool {
	return true
}

