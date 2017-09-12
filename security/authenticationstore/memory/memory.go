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

package memory

import (
	sec"github.com/uber/jaeger/security"
	"github.com/pkg/errors"
)

// InMemoryAuthenticationStore is in-memory implementation of the authentication store
type InMemoryAuthenticationStore struct {
	ctxs map[string]*sec.AuthenticationContext
}

func NewInMemoryAuthenticationStore(
	principals []string,
) (*InMemoryAuthenticationStore, error) {
	if len(principals) < 1 {
		return nil, errors.New("No principals provided for in-memory store")
	}
	store := &InMemoryAuthenticationStore{
		ctxs: make(map[string]*sec.AuthenticationContext),
	}
	for _, principal := range principals {
		store.ctxs[principal] = &sec.AuthenticationContext{
			Principal: principal,
			Password: principal,
			Locked: false,
		}
	}
	return store, nil
}

func (mem *InMemoryAuthenticationStore) FindPrincipal(token sec.AuthenticationToken) (*sec.AuthenticationContext, error) {
	if ctx, ok := mem.ctxs[token.Username]; ok {
		return ctx, nil
	}
	return nil, nil
}