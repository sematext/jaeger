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

package identity

import (
	"github.com/uber/jaeger/model"
	"go.uber.org/zap"
)

// Authenticator authenticates inbound spans
type Authenticator interface {
	Authenticate(span *model.Span) bool
}

type SpanAuthenticator struct {
	store 	 TokenStore
	logger   *zap.Logger
	tokenKey string
}

func NewSpanAuthenticator(
	tokenStore TokenStore,
	logger *zap.Logger,
	tokenKey string,
) SpanAuthenticator {
	return SpanAuthenticator{
		store: tokenStore,
		logger: logger,
		tokenKey: tokenKey,
	}
}

// Authenticate accept the incoming span if it carries a tag whose key name
// is specified by tokenKey, and the token associated with former key is present
// in the token store
func (ath SpanAuthenticator) Authenticate(span *model.Span) bool {
	if kv, ok := span.Tags.FindByKey(ath.tokenKey); ok != false {
		return ath.store.TokenExists(kv.VStr)
	} else {
		ath.logger.Warn("Token not found in tags", zap.String("token-key", ath.tokenKey))
		return false
	}
}