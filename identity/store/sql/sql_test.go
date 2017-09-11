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
)

type clientMock struct {
	mock.Mock
}

func (c clientMock) Ping() error {
	args := c.Called()
	return args.Error(0)
}

func (c clientMock) QueryForRow(query string, args ...interface{}) (interface{}, error) {
	arguments := c.Called(query, args)
	return arguments.String(0), arguments.Error(1)
}

func TestNewDbTokenStore(t *testing.T) {
	c := new(clientMock)

	c.On("Ping").Return(nil)

	logger := zap.NewNop()
	store, err := NewDbTokenStore(
		c,
		logger,
		"SELECT token FROM system WHERE token = ?",
		100,
		3600,
	)
	require.NoError(t, err)
	assert.NotNil(t, store)
	assert.Equal(t, 0, store.cache.Size())
}

func TestNewDbTokenStorePingError(t *testing.T) {
	c := new(clientMock)

	c.On("Ping").Return(errors.New("Connection timeout"))

	logger := zap.NewNop()
	_, err := NewDbTokenStore(
		c,
		logger,
		"SELECT token FROM system WHERE token = ?",
		100,
		3600,
	)
	require.Error(t, err)
}

func TestFindToken(t *testing.T) {
	c := new(clientMock)

	c.On("Ping").Return(nil)
	c.On("QueryForRow",
		"SELECT token FROM system WHERE token = ?",
		[]interface{}{"c15a1793-71b7-46a5-88c5-bc76f9c772a0"},
	).Return("c15a1793-71b7-46a5-88c5-bc76f9c772a0", nil)

	logger := zap.NewNop()
	store, _ := NewDbTokenStore(
		c,
		logger,
		"SELECT token FROM system WHERE token = ?",
		100,
		3600,
	)
	assert.Nil(t, store.cache.Get("c15a1793-71b7-46a5-88c5-bc76f9c772a0"))
	found, err := store.FindToken("c15a1793-71b7-46a5-88c5-bc76f9c772a0")
	require.NoError(t, err)
	assert.True(t, found)
	assert.NotNil(t, store.cache.Get("c15a1793-71b7-46a5-88c5-bc76f9c772a0"))
}

func TestFindTokenQueryError(t *testing.T) {
	c := new(clientMock)

	c.On("Ping").Return(nil)
	c.On("QueryForRow",
		"SELECT token FROM system WHERE token = ?",
		[]interface{}{"c15a1793-71b7-46a5-88c5-bc76f9c772a0"},
	).Return("c15a1793-71b7-46a5-88c5-bc76f9c772a0", errors.New("sql: no rows in result set"))

	logger := zap.NewNop()
	store, _ := NewDbTokenStore(
		c,
		logger,
		"SELECT token FROM system WHERE token = ?",
		100,
		3600,
	)
	assert.Nil(t, store.cache.Get("c15a1793-71b7-46a5-88c5-bc76f9c772a0"))
	_, err := store.FindToken("c15a1793-71b7-46a5-88c5-bc76f9c772a0")
	require.Error(t, err)
	assert.Nil(t, store.cache.Get("c15a1793-71b7-46a5-88c5-bc76f9c772a0"))
}