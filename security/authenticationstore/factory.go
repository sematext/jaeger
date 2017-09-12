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

package authenticationstore

import (
	"github.com/uber/jaeger/security"
	dbAuthStore"github.com/uber/jaeger/security/authenticationstore/sql"
	memAuthStore"github.com/uber/jaeger/security/authenticationstore/memory"
	"go.uber.org/zap"
	"github.com/uber/jaeger/cmd/builder"
	"github.com/uber/jaeger/cmd/flags"
	"fmt"
)

// NewAuthenticationStore yields a specific authentication store or returns an error if authentication store
// type is not supported.
func NewAuthenticationStore(authStoreType string, logger *zap.Logger, opts builder.BasicOptions) (security.AuthenticationStore, error) {
	switch authStoreType {
	case flags.SQLAuthenticationStoreType:
		dbClientBuilder := opts.DbAuthenticationStoreClientBuilder
		client, err := dbClientBuilder.NewDbClient()
		if err != nil {
			return nil, err
		}
		return dbAuthStore.NewDbAuthenticationStore(
			client,
			logger,
			dbClientBuilder.GetQuery(),
		)
	case flags.InMemoryAuthenticationStoreType:
		memAuthStoreBuilder := opts.InMemoryAuthenticationStoreBuilder
		return memAuthStore.NewInMemoryAuthenticationStore(memAuthStoreBuilder.GetPrincipals())
	default:
		return nil, fmt.Errorf("%s is unsupported authentication store", authStoreType)
	}
}