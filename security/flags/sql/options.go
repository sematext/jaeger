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
	"github.com/uber/jaeger/security/authenticationstore/sql/config"
	"github.com/spf13/viper"
	"flag"
)

const (
	suffixUsername    	= ".username"
	suffixPassword    	= ".password"
	suffixDriver      	= ".driver"
	suffixHost		  	= ".host"
	suffixPort		  	= ".port"
	suffixDatabase	    = ".database"
	suffixQuery 	  	= ".query"
)

// Options describes various configuration for the SQL authentication store
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
				Driver: "mysql",
				Host: "localhost",
				Port: 3306,
				Username: "root",
				Password: "",
				Database: "",
				Query: "SELECT token, token, active FROM System WHERE token = ?",
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
		nsConfig.namespace + suffixUsername,
		nsConfig.Username,
		"The username required by SQL authentication store")
	flagSet.String(
		nsConfig.namespace + suffixPassword,
		nsConfig.Password,
		"The password required by SQL authentication store")
	flagSet.String(
		nsConfig.namespace + suffixDriver,
		nsConfig.Driver,
		"The name of the SQL driver used to connect with the underlying database server")
	flagSet.String(
		nsConfig.namespace + suffixHost,
		nsConfig.Host,
		"The name / IP address of the host where database server is listening")
	flagSet.Uint(
		nsConfig.namespace + suffixPort,
		nsConfig.Port,
		"The port where database server is listening for requests")
	flagSet.String(
		nsConfig.namespace + suffixDatabase,
		nsConfig.Database,
		"The name of the database where principal data can be find")
	flagSet.String(
		nsConfig.namespace + suffixQuery,
		nsConfig.Query,
		"The SQL query that retrieves principal info from the database")
}

// InitFromViper initializes Options with properties from viper
func (opt *Options) InitFromViper(v *viper.Viper) {
	initFromViper(opt.primary, v)
}

func initFromViper(cfg *namespaceConfig, v *viper.Viper) {
	cfg.Username = v.GetString(cfg.namespace + suffixUsername)
	cfg.Password = v.GetString(cfg.namespace + suffixPassword)
	cfg.Driver = v.GetString(cfg.namespace + suffixDriver)
	cfg.Host = v.GetString(cfg.namespace + suffixHost)
	cfg.Port = uint(v.GetInt(cfg.namespace + suffixPort))
	cfg.Database = v.GetString(cfg.namespace + suffixDatabase)
	cfg.Query = v.GetString(cfg.namespace + suffixQuery)
}
