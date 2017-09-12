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
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uber/jaeger/security"
)

func TestNewInMemoryAuthenticationStoreEmptyPrincipals(t *testing.T) {
	_, err := NewInMemoryAuthenticationStore([]string{})
	require.Error(t, err)
}

func TestFindPrincipal(t *testing.T) {
	store, _ := NewInMemoryAuthenticationStore(
		[]string{
			"38094e7a-96f8-11e7-bc87-9bae86d05b5b",
			"45f4b11e-96f8-11e7-9e14-a7be8cf5fadf",
			"4c788060-96f8-11e7-99db-5f748d21718d",
		},
	)
	ctx1, _ := store.FindPrincipal(security.AuthenticationToken{Username: "38094e7a-96f8-11e7-bc87-9bae86d05b5b", Password: "38094e7a-96f8-11e7-bc87-9bae86d05b5b"})
	assert.NotNil(t, ctx1)
	ctx2, _ := store.FindPrincipal(security.AuthenticationToken{Username: "44f4b11e-96f8-11e7-9e14-a7be8cf5fadf", Password: "c15a1793-71b7-46a5-88c5-bc76f9c772a0"})
	assert.Nil(t, ctx2)
}
