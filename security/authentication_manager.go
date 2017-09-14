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
	"github.com/uber/jaeger/model"
	"go.uber.org/zap"
	"github.com/uber/jaeger/pkg/cache"
	"time"
)

// Authenticator authenticates inbound spans
type Authenticator interface {
	Authenticate(token *AuthenticationToken) bool
	TokenFromSpan(span *model.Span) *AuthenticationToken
}

type AuthenticationManager struct {
	store  AuthenticationStore
	logger *zap.Logger
	cache  cache.Cache
	key    string
}

func NewAuthenticationManager(
	authStore AuthenticationStore,
	logger *zap.Logger,
	key string,
	cacheSize int,
	cacheTTL time.Duration,
) AuthenticationManager {
	return AuthenticationManager{
		store: authStore,
		logger: logger,
		cache: cache.NewLRUWithOptions(
			cacheSize,
			&cache.Options{
				TTL: time.Second * cacheTTL,
			}),
		key: key,
	}
}

// TokenFromSpan builds an authentication token from span / process tags.
// The same tag value is used for both user name and password fields.
func (am AuthenticationManager) TokenFromSpan(span *model.Span) *AuthenticationToken {
	if kv, ok := span.Tags.FindByKey(am.key); ok {
		return &AuthenticationToken{
			Username: kv.VStr,
			Password: kv.VStr,
		}
	}
	if kv, ok := span.Process.Tags.FindByKey(am.key); ok {
		return &AuthenticationToken{
			Username: kv.VStr,
			Password: kv.VStr,
		}
	}
	return nil
}

// Authenticate authenticates the inbound span. It collaborates with underlying authentication store
// to obtain principal info. Upon successful authentication, the context is cached to avoid subsequent
// round trips to the authentication store, and thus reducing the overall I/O.
func (am AuthenticationManager) Authenticate(token *AuthenticationToken) bool {
	ctx := am.ctxFromCache(token.Username)
	if ctx == nil {
		var err error
		ctx, err = am.store.FindPrincipal(*token)
		if err != nil {
			am.logger.Warn("Failed to load principal", zap.Error(err))
			return false
		}
		if ctx.PasswordEquals(token.Password) && !ctx.Locked {
			am.ctxToCache(token.Username, ctx)
			return true
		}
		return false
	}
	return !ctx.Locked
}

func (am *AuthenticationManager) ctxFromCache(username string) *AuthenticationContext {
	ctx := am.cache.Get(username)
	if c, ok := ctx.(*AuthenticationContext); ok {
		return c
	}
	return nil
}

func (am *AuthenticationManager) ctxToCache(username string, ctx *AuthenticationContext) {
	am.cache.Put(username, ctx)
}