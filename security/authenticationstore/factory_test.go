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
	"testing"
	"go.uber.org/zap"
	"github.com/uber/jaeger/cmd/builder"
	"github.com/stretchr/testify/mock"
	"github.com/uber/jaeger/security/authenticationstore/sql"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger/security/authenticationstore/memory"
)

type clientMock struct {
	mock.Mock
}

func (c clientMock) Ping() error {
	args := c.Called()
	return args.Error(0)
}

func (c clientMock) QueryForRow(query string, mapper sql.RowMapper, args ...interface{}) (interface{}, error) {
	arguments := c.Called(query, "func(*sql.Row) (interface{}, error)", args)
	return arguments.Get(0), arguments.Error(1)
}

type dbAuthenticationStoreClientBuilderMock struct {
	mock.Mock
}

type inMemAuthenticationStoreBuilderMock struct {
	mock.Mock
}

func (m dbAuthenticationStoreClientBuilderMock) NewDbClient() (sql.DbClient, error) {
	args := m.Called()
	return args.Get(0).(sql.DbClient), args.Error(1)
}

func (m dbAuthenticationStoreClientBuilderMock) GetQuery() string {
	args := m.Called()
	return args.String(0)
}

func (m inMemAuthenticationStoreBuilderMock) GetPrincipals() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func TestNewSqlAuthenticationStore(t *testing.T) {
	dbAuthStoreClientBuilderMock := new(dbAuthenticationStoreClientBuilderMock)

	clientMock := new(clientMock)
	clientMock.On("Ping").Return(nil)

	dbAuthStoreClientBuilderMock.On("NewDbClient").Return(clientMock, nil)
	dbAuthStoreClientBuilderMock.On("GetQuery").Return("SELECT token, token, active FROM System WHERE token = ?")

	basicOpts := builder.BasicOptions{
		DbAuthenticationStoreClientBuilder: dbAuthStoreClientBuilderMock,
	}
	logger := zap.NewNop()

	as, err := NewAuthenticationStore("sql", logger, basicOpts)

	require.NoError(t, err)
	assert.IsType(t, &sql.DbAuthenticationStore{}, as)
}

func TestNewMemoryAuthenticationStore(t *testing.T) {
	inMemoryAuthStoreBuilderMock := new(inMemAuthenticationStoreBuilderMock)

	inMemoryAuthStoreBuilderMock.On("GetPrincipals").Return([]string{"c15a1393-61b7-46a5-88c5-bc77f9c772a0"})
	basicOpts := builder.BasicOptions{
		InMemoryAuthenticationStoreBuilder: inMemoryAuthStoreBuilderMock,
	}
	logger := zap.NewNop()

	as, err := NewAuthenticationStore("memory", logger, basicOpts)

	require.NoError(t, err)
	assert.IsType(t, &memory.InMemoryAuthenticationStore{}, as)
}

func TestNewUnsupportedAuthenticationStore(t *testing.T) {
	inMemoryAuthStoreBuilderMock := new(inMemAuthenticationStoreBuilderMock)

	basicOpts := builder.BasicOptions{
		InMemoryAuthenticationStoreBuilder: inMemoryAuthStoreBuilderMock,
	}
	logger := zap.NewNop()

	as, err := NewAuthenticationStore("ldap", logger, basicOpts)

	require.Error(t, err)
	require.Nil(t, as)
}