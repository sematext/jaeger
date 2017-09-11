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
	i"github.com/uber/jaeger/identity"
	"github.com/uber/jaeger/pkg/cache"
	"go.uber.org/zap"
	"time"
)

// DbTokenStore describes the store based on relational database engine
type DbTokenStore struct {
	client DbClient
	logger *zap.Logger
	cache  cache.Cache
	query  string
}

// NewDbTokenStore creates a new instance of the SQL-based token store. It attempts to
// ping the database server to ensure it's available.
func NewDbTokenStore(
	client DbClient,
	logger *zap.Logger,
	query string,
	maxCacheSize int,
	cacheEviction time.Duration,
) (*DbTokenStore, error) {
	if err := client.Ping(); err != nil {
		return nil, err
	}
	return &DbTokenStore{
		client: client,
		logger: logger,
		cache: cache.NewLRUWithOptions(
			maxCacheSize,
			&cache.Options{
				TTL: time.Second * cacheEviction,
			},
		),
		query: query,
	}, nil
}

// FindToken attempts to find a token in the cache. If not found, the specified SQL query
// is executed against the underlying database engine. The token is cached to
// avoid subsequent round trips to the database server, and thus reducing the
// overall I/O activity.
func (store *DbTokenStore) FindToken(token string, parameters ...i.TokenParameters) (bool, error) {
	if store.cache.Get(token) != nil {
		return true, nil
	} else {
 		if err := store.findToken(token, parameters...); err != nil {
			return false, err
		}
		store.cache.Put(token, token)
		return true, nil
	}
}

func (store *DbTokenStore) findToken(token string, parameters ...i.TokenParameters) error {
	args := make([]interface{}, len(parameters))
	for j, param := range parameters {
		args[j] = param
	}
	args = append([]interface{}{token}, args...)
	_, err := store.client.QueryForRow(store.query, args...)
	return err
}

