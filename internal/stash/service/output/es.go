package output

import (
	"BlackHole/pkg/config"
	"context"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/rogpeppe/go-internal/semver"
	"github.com/zeromicro/go-zero/core/executors"
	"github.com/zeromicro/go-zero/core/logx"

	jsoniter "github.com/json-iterator/go"
)

const es8Version = "8.0.0"

type (
	EsWriter struct {
		docType   string
		esVersion string
		client    *elastic.Client
		inserter  *executors.ChunkExecutor
		indexer   *Index
	}

	valueWithIndex struct {
		index string
		val   string
	}
)

func NewElasticSearchWriter(c *config.ElasticSearchConf) (*EsWriter, error) {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(c.Hosts...),
		elastic.SetGzip(c.Compress),
		elastic.SetBasicAuth(c.Username, c.Password),
	)
	if err != nil {
		return nil, err
	}

	version, err := client.ElasticsearchVersion(c.Hosts[0])
	if err != nil {
		return nil, err
	}

	writer := EsWriter{
		docType:   c.DocType,
		client:    client,
		esVersion: version,
	}
	writer.inserter = executors.NewChunkExecutor(writer.execute, executors.WithChunkBytes(c.MaxChunkBytes))

	var loc *time.Location
	if len(c.TimeZone) > 0 {
		loc, err = time.LoadLocation(c.TimeZone)
		logx.Must(err)
	} else {
		loc = time.Local
	}
	writer.indexer = NewIndex(client, c.Index, loc)
	return &writer, nil
}

func (w *EsWriter) Write(val map[string]interface{}) error {
	index := w.indexer.GetIndex(val)

	bs, err := jsoniter.Marshal(val)
	if err != nil {
		return err
	}

	return w.inserter.Add(valueWithIndex{
		index: index,
		val:   string(bs),
	}, len(string(bs)))
}

func (w *EsWriter) execute(vals []interface{}) {
	var bulk = w.client.Bulk()
	for _, val := range vals {
		pair := val.(valueWithIndex)
		req := elastic.NewBulkIndexRequest().Index(pair.index)
		if isSupportType(w.esVersion) && len(w.docType) > 0 {
			req = req.Type(w.docType)
		}
		req = req.Doc(pair.val)
		bulk.Add(req)
	}
	resp, err := bulk.Do(context.Background())
	if err != nil {
		logx.Error(err)
		return
	}

	// bulk error in docs will report in response items
	if !resp.Errors {
		return
	}

	for _, imap := range resp.Items {
		for _, item := range imap {
			if item.Error == nil {
				continue
			}

			logx.Error(item.Error)
		}
	}
}

func isSupportType(version string) bool {
	// es8.x not support type field
	return semver.Compare(version, es8Version) < 0
}
