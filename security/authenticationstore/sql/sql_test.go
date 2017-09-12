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
	"testing"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"github.com/pkg/errors"
	"github.com/uber/jaeger/security"
)

type clientMock struct {
	mock.Mock
}

func (c clientMock) Ping() error {
	args := c.Called()
	return args.Error(0)
}

func (c clientMock) QueryForRow(query string, mapper RowMapper, args ...interface{}) (interface{}, error) {
	arguments := c.Called(query, "func(*sql.Row) (interface{}, error)", args)
	return arguments.Get(0), arguments.Error(1)
}

func TestNewDbAuthenticationStore(t *testing.T) {
	c := new(clientMock)

	c.On("Ping").Return(nil)

	logger := zap.NewNop()
	store, err := NewDbAuthenticationStore(
		c,
		logger,
		"SELECT token, token, active FROM system WHERE token = ?",
	)
	require.NoError(t, err)
	assert.NotNil(t, store)
}

func TestNewDbAuthenticationStorePingError(t *testing.T) {
	c := new(clientMock)

	c.On("Ping").Return(errors.New("Connection timeout"))

	logger := zap.NewNop()
	_, err := NewDbAuthenticationStore(
		c,
		logger,
		"SELECT token, token, active FROM system WHERE token = ?",
	)
	require.Error(t, err)
}

func TestFindValidPrincipal(t *testing.T) {
	c := new(clientMock)

	ctx := security.AuthenticationContext{
		Principal: "c15a1793-71b7-46a5-88c5-bc76f9c772a0",
		Password: "c15a1793-71b7-46a5-88c5-bc76f9c772a0",
		Locked: false,
	}
	c.On("Ping").Return(nil)
	c.On("QueryForRow",
		"SELECT token, token, active FROM system WHERE token = ?",
		"func(*sql.Row) (interface{}, error)",
		[]interface{}{"c15a1793-71b7-46a5-88c5-bc76f9c772a0"},
	).Return(ctx, nil)
	logger := zap.NewNop()
	store, _ := NewDbAuthenticationStore(
		c,
		logger,
		"SELECT token, token, active FROM system WHERE token = ?",
	)
	token := security.AuthenticationToken{
		Username: "c15a1793-71b7-46a5-88c5-bc76f9c772a0",
		Password: "c15a1793-71b7-46a5-88c5-bc76f9c772a0",
	}
	context, err := store.FindPrincipal(token)
	require.NoError(t, err)
	assert.NotNil(t, context)
	assert.Equal(t, token.Username, ctx.Principal)
	assert.Equal(t, token.Password, ctx.Password)
	assert.False(t, context.Locked)
}

func TestFindNonExistentPrincipal(t *testing.T) {
	c := new(clientMock)

	c.On("Ping").Return(nil)
	c.On("QueryForRow",
		"SELECT token, token, active FROM system WHERE token = ?",
		"func(*sql.Row) (interface{}, error)",
		[]interface{}{"c15a1793-71b7-46a5-88c5-bc76f9c772a0"},
	).Return(nil, errors.New("sql: no rows in result set"))
	logger := zap.NewNop()
	store, _ := NewDbAuthenticationStore(
		c,
		logger,
		"SELECT token, token, active FROM system WHERE token = ?",
	)
	token := security.AuthenticationToken{
		Username: "c15a1793-71b7-46a5-88c5-bc76f9c772a0",
		Password: "c15a1793-71b7-46a5-88c5-bc76f9c772a0",
	}
	context, err := store.FindPrincipal(token)
	require.Error(t, err)
	assert.Nil(t, context)
}

