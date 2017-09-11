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

package identity

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"github.com/uber/jaeger/model"
	"github.com/stretchr/testify/assert"
)

type tokenStoreMock struct {
	mock.Mock
}

func (store *tokenStoreMock) FindToken(token string, parameters ...TokenParameters) (bool, error) {
	args := store.Called(token, parameters)
	return args.Bool(0), args.Error(1)
}

func TestAuthenticate(t *testing.T) {

	store := new(tokenStoreMock)
	logger := zap.NewNop()

	store.On(
		"FindToken",
		"n15a1793-a1b7-16a5-88c5-bc76f9c772a1",
		[]TokenParameters(nil),).Return(true, nil)

	spanAuthenticator := SpanAuthenticator{
		store,
		logger,
		"token",
	}
	span := &model.Span {
		Tags: model.KeyValues{
			model.String("token", "n15a1793-a1b7-16a5-88c5-bc76f9c772a1"),
		},
	}
	success := spanAuthenticator.Authenticate(span)
	assert.True(t, success)
}

func TestAuthenticateKeyNotFound(t *testing.T) {

	store := new(tokenStoreMock)
	logger := zap.NewNop()

	store.On(
		"FindToken",
		"n15a1793-a1b7-16a5-88c5-bc76f9c772a1",
		[]TokenParameters(nil),).Return(true, nil)

	spanAuthenticator := SpanAuthenticator{
		store,
		logger,
		"token",
	}
	span := &model.Span {
		Tags: model.KeyValues{
			model.String("apptoken", "n15a1793-a1b7-16a5-88c5-bc76f9c772a1"),
		},
	}
	success := spanAuthenticator.Authenticate(span)
	assert.False(t, success)
}

func TestAuthenticateTokenNotFound(t *testing.T) {

	store := new(tokenStoreMock)
	logger := zap.NewNop()

	store.On(
		"FindToken",
		"115a1793-a1b7-16a5-88c5-bc76f9c772a1",
		[]TokenParameters(nil),).Return(false, nil)

	spanAuthenticator := SpanAuthenticator{
		store,
		logger,
		"token",
	}
	span := &model.Span {
		Tags: model.KeyValues{
			model.String("token", "115a1793-a1b7-16a5-88c5-bc76f9c772a1"),
		},
	}
	success := spanAuthenticator.Authenticate(span)
	assert.False(t, success)
}