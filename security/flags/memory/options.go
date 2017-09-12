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
	"github.com/spf13/viper"
	"flag"
	"strings"
	"github.com/uber/jaeger/security/authenticationstore/memory/config"
)

const (
	suffixPrincipals = ".principals"
)

// Options describes various configuration for the in-memory token store
type Options struct {
	primary *namespaceConfig
}

type namespaceConfig struct {
	config.Configuration
	namespace string
}

func NewOptions (namespace string) *Options {
	return &Options {
		primary: &namespaceConfig{
			Configuration: config.Configuration{
				Principals: []string{},
			},
			namespace: namespace,
		},
	}
}

func (opt *Options) GetPrimary() *config.Configuration {
	return &opt.primary.Configuration
}

// AddFlags adds flags for Options
func (opt *Options) AddFlags(flagSet *flag.FlagSet) {
	addFlags(flagSet, opt.primary)
}

func addFlags(flagSet *flag.FlagSet, nsConfig *namespaceConfig) {
	flagSet.String(
		nsConfig.namespace + suffixPrincipals,
		"",
		"The comma-separated list of authentication principals")
}

// InitFromViper initializes Options with properties from viper
func (opt *Options) InitFromViper(v *viper.Viper) {
	initFromViper(opt.primary, v)
}

func initFromViper(cfg *namespaceConfig, v *viper.Viper) {
	principals := v.GetString(cfg.namespace + suffixPrincipals)
	if principals != "" {
		cfg.Principals = strings.Split(principals, ",")
	}
}
