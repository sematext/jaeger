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

package security

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"github.com/uber/jaeger/model"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/pkg/errors"
)

type authenticationStoreMock struct {
	mock.Mock
}

func (store *authenticationStoreMock) FindPrincipal(token AuthenticationToken) (*AuthenticationContext, error) {
	args := store.Called(token)
	ctx, _ := args.Get(0).(AuthenticationContext)
	return &ctx, args.Error(1)
}

func TestAuthenticate(t *testing.T) {

	store := new(authenticationStoreMock)
	logger := zap.NewNop()

	ctx := AuthenticationContext{
		Principal: "315a1793-a1b7-16a5-88c5-bc76f9c772a1",
		Password: "315a1793-a1b7-16a5-88c5-bc76f9c772a1",
		Locked: false,
	}
	authenticationManager := NewAuthenticationManager(
		store,
		logger,
		"api-token",
		100,
		time.Second * 60,
	)
	span := &model.Span {
		Tags: model.KeyValues{
			model.String("api-token", "315a1793-a1b7-16a5-88c5-bc76f9c772a1"),
		},
	}
	token := authenticationManager.TokenFromSpan(span)
	store.On("FindPrincipal", *token).Return(ctx, nil)

	success := authenticationManager.Authenticate(token)
	assert.True(t, success)
	assert.Equal(t, authenticationManager.cache.Size(), 1)

	authenticationManager.Authenticate(token)
	store.AssertNumberOfCalls(t, "FindPrincipal", 1)
}

func TestAuthenticateIncorrectPassword(t *testing.T) {

	store := new(authenticationStoreMock)
	logger := zap.NewNop()

	ctx := AuthenticationContext{
		Principal: "315a1793-a1b7-16a5-88c5-bc76f9c772a1",
		Password: "315a1793-a1b7-16a5-88c5-bc76f9c772a1",
		Locked: false,
	}
	authenticationManager := NewAuthenticationManager(
		store,
		logger,
		"api-token",
		100,
		time.Second * 60,
	)
	span := &model.Span {
		Tags: model.KeyValues{
			model.String("api-token", "315a2793-a1b7-16a5-88c5-bc76f9c772a1"),
		},
	}
	token := authenticationManager.TokenFromSpan(span)
	store.On("FindPrincipal", *token).Return(ctx, nil)

	success := authenticationManager.Authenticate(token)
	assert.False(t, success)
	assert.Equal(t, authenticationManager.cache.Size(), 0)
}

func TestAuthenticatePrincipalLocked(t *testing.T) {

	store := new(authenticationStoreMock)
	logger := zap.NewNop()

	ctx := AuthenticationContext{
		Principal: "315a1793-a1b7-16a5-88c5-bc76f9c772a1",
		Password: "315a1793-a1b7-16a5-88c5-bc76f9c772a1",
		Locked: true,
	}
	authenticationManager := NewAuthenticationManager(
		store,
		logger,
		"api-token",
		100,
		time.Second * 60,
	)
	span := &model.Span {
		Tags: model.KeyValues{
			model.String("api-token", "315a1793-a1b7-16a5-88c5-bc76f9c772a1"),
		},
	}
	token := authenticationManager.TokenFromSpan(span)
	store.On("FindPrincipal", *token).Return(ctx, nil)

	success := authenticationManager.Authenticate(token)
	assert.False(t, success)
	assert.Equal(t, authenticationManager.cache.Size(), 0)
}

func TestAuthenticatePrincipalNotFound(t *testing.T) {

	store := new(authenticationStoreMock)
	logger := zap.NewNop()

	authenticationManager := NewAuthenticationManager(
		store,
		logger,
		"api-token",
		100,
		time.Second * 60,
	)
	span := &model.Span {
		Tags: model.KeyValues{
			model.String("api-token", "315a1793-a1b7-16a5-88c5-bc76f9c772a1"),
		},
	}
	token := authenticationManager.TokenFromSpan(span)
	store.On("FindPrincipal", *token).Return(nil, errors.New("sql: no rows in result set"))

	success := authenticationManager.Authenticate(token)
	assert.False(t, success)
	assert.Equal(t, authenticationManager.cache.Size(), 0)
}

