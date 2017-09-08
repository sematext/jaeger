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

package config

import (
	"fmt"
	"github.com/uber/jaeger/identity/store/sql"
)

// DbClientBuilder creates a new SQL client
type DbClientBuilder interface {
	NewDbClient() (*sql.Client, error)
	GetQuery() string
}

// Configuration describes the config properties needed to connect to a SQL database
type Configuration struct {
	// Driver specifies the type of database engine (mysql, postgres, etc.)
	Driver string
	// Database represents the name of the database where token table is stored
	Database string
	// Host the hostname where database server is located
	Host string
	// Port the port where database server is waiting for requests
	Port int16
	// Username to authenticate to an instance of the database
	Username string
	// Password to authenticate to an instance of the database
	Password string
	// Query specifies the SQL query that's used to obtain the token
	Query string

	CacheEviction int
	UnixSocket string
}

func (c *Configuration) NewDbClient() (*sql.Client, error) {
	client, err := sql.NewClient(c.Driver, c.buildDataSource())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Configuration) GetQuery() string {
	return c.Query
}

func (c *Configuration) buildDataSource() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
}
