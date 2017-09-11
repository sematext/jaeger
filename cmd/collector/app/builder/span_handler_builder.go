// Copyright (c) 2017 Uber Technologies, Inc.
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

package builder

import (
	"errors"
	"os"

	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	basicB "github.com/uber/jaeger/cmd/builder"
	"github.com/uber/jaeger/cmd/collector/app"
	zs "github.com/uber/jaeger/cmd/collector/app/sanitizer/zipkin"
	"github.com/uber/jaeger/cmd/flags"
	"github.com/uber/jaeger/model"
	i"github.com/uber/jaeger/identity"
	sqlsc"github.com/uber/jaeger/identity/store/sql/config"
	dbTokenStore"github.com/uber/jaeger/identity/store/sql"
	cascfg "github.com/uber/jaeger/pkg/cassandra/config"
	escfg "github.com/uber/jaeger/pkg/es/config"
	casSpanstore "github.com/uber/jaeger/plugin/storage/cassandra/spanstore"
	esSpanstore "github.com/uber/jaeger/plugin/storage/es/spanstore"
	"github.com/uber/jaeger/storage/spanstore"
)

var (
	errMissingCassandraConfig     = errors.New("Cassandra not configured")
	errMissingMemoryStore         = errors.New("MemoryStore is not provided")
	errMissingElasticSearchConfig = errors.New("ElasticSearch not configured")
)

// SpanHandlerBuilder holds configuration required for handlers
type SpanHandlerBuilder struct {
	logger         	  *zap.Logger
	metricsFactory    metrics.Factory
	collectorOpts     *CollectorOptions
	spanWriter        spanstore.Writer
	spanAuthenticator i.Authenticator
}

// NewSpanHandlerBuilder returns new SpanHandlerBuilder with configured span storage.
func NewSpanHandlerBuilder(cOpts *CollectorOptions, sFlags *flags.SharedFlags, opts ...basicB.Option) (*SpanHandlerBuilder, error) {
	options := basicB.ApplyOptions(opts...)

	spanHb := &SpanHandlerBuilder{
		collectorOpts:  cOpts,
		logger:         options.Logger,
		metricsFactory: options.MetricsFactory,
	}

	var err error
	if sFlags.SpanStorage.Type == flags.CassandraStorageType {
		if options.CassandraSessionBuilder == nil {
			return nil, errMissingCassandraConfig
		}
		spanHb.spanWriter, err = spanHb.initCassStore(options.CassandraSessionBuilder)
	} else if sFlags.SpanStorage.Type == flags.MemoryStorageType {
		if options.MemoryStore == nil {
			return nil, errMissingMemoryStore
		}
		spanHb.spanWriter = options.MemoryStore
	} else if sFlags.SpanStorage.Type == flags.ESStorageType {
		if options.ElasticClientBuilder == nil {
			return nil, errMissingElasticSearchConfig
		}
		spanHb.spanWriter, err = spanHb.initElasticStore(options.ElasticClientBuilder)
	} else {
		return nil, flags.ErrUnsupportedStorageType
	}

	var tstore i.TokenStore
	if cOpts.AuthSpan {
		switch sFlags.TokenStore.Type {
		case flags.SQLTokenStoreType:
			tstore, err = spanHb.initDbTokenStore(options.DbTokenStoreClientBuilder)
		default:
			err = flags.ErrUnupportedTokenStoreType
		}
		if err == nil {
			spanHb.spanAuthenticator = i.NewSpanAuthenticator(
				tstore,
				options.Logger,
				cOpts.AuthTokenKey,
			)
		}
	}

	if err != nil {
		return nil, err
	}

	return spanHb, nil
}

func (spanHb *SpanHandlerBuilder) initCassStore(builder cascfg.SessionBuilder) (spanstore.Writer, error) {
	session, err := builder.NewSession()
	if err != nil {
		return nil, err
	}

	return casSpanstore.NewSpanWriter(
		session,
		spanHb.collectorOpts.WriteCacheTTL,
		spanHb.metricsFactory,
		spanHb.logger,
	), nil
}

func (spanHb *SpanHandlerBuilder) initElasticStore(esBuilder escfg.ClientBuilder) (spanstore.Writer, error) {
	client, err := esBuilder.NewClient()
	if err != nil {
		return nil, err
	}

	return esSpanstore.NewSpanWriter(
		client,
		spanHb.logger,
		spanHb.metricsFactory,
		esBuilder.GetNumShards(),
		esBuilder.GetNumReplicas(),
	), nil
}

func (spanHb *SpanHandlerBuilder) initDbTokenStore(dbClientBuilder sqlsc.DbClientBuilder) (i.TokenStore, error) {
	client, err := dbClientBuilder.NewDbClient()
	if err != nil {
		return nil, err
	}
	return dbTokenStore.NewDbTokenStore(
		client,
		spanHb.logger,
		dbClientBuilder.GetQuery(),
		dbClientBuilder.GetMaxCacheSize(),
		dbClientBuilder.GetCacheEviction(),
	)
}

// BuildHandlers builds span handlers (Zipkin, Jaeger)
func (spanHb *SpanHandlerBuilder) BuildHandlers() (app.ZipkinSpansHandler, app.JaegerBatchesHandler) {
	hostname, _ := os.Hostname()
	hostMetrics := spanHb.metricsFactory.Namespace(hostname, nil)

	zSanitizer := zs.NewChainedSanitizer(
		zs.NewSpanDurationSanitizer(),
		zs.NewSpanStartTimeSanitizer(),
		zs.NewParentIDSanitizer(),
		zs.NewErrorTagSanitizer(),
	)

	spanProcessor := app.NewSpanProcessor(
		spanHb.spanWriter,
		app.Options.ServiceMetrics(spanHb.metricsFactory),
		app.Options.HostMetrics(hostMetrics),
		app.Options.Logger(spanHb.logger),
		app.Options.SpanFilter(spanHb.defaultSpanFilter),
		app.Options.NumWorkers(spanHb.collectorOpts.NumWorkers),
		app.Options.QueueSize(spanHb.collectorOpts.QueueSize),
	)

	return app.NewZipkinSpanHandler(spanHb.logger, spanProcessor, zSanitizer),
		app.NewJaegerSpanHandler(spanHb.logger, spanProcessor)
}

func (spanHb *SpanHandlerBuilder) defaultSpanFilter(span *model.Span) bool {
	if spanHb.collectorOpts.AuthSpan {
		authenticator := spanHb.spanAuthenticator
		return authenticator.Authenticate(span)
	}
	return true
}
