package elastic

import (
	model "autocall/internal/elastic/model"

	"github.com/olivere/elastic/v7"
)

// IElasticDriver - IModelRepository
type IElasticConnector interface {
	GetConn() error
	GetClient() *elastic.Client
	Ping()
	CreateIndex(name string, jsonBody string) (bool, error)
	CreateAlias(index string, alias string) (bool, error)
	CheckAliasExist(index, alias string) (bool, error)
	CreateDocs(name string, id string, body interface{}) (model.Response, error)
	CreateDocWithRouting(index, routing, id, index_type string, body interface{}) (model.Response, error)
	GetDoc(name string, id string) (model.Response, error)
	GetDocWithRouting(index, routing, id string) (model.Response, error)
	SearchDoc(index string, searchSource *elastic.SearchSource) (model.Response, error)
	SearchDocWithPaging(index string, searchSource *elastic.SearchSource, size, from int) (model.Response, error)
	GetDocBySearch(index string, searchSource *elastic.SearchSource) (model.Response, error)
	CountDocs(index string) (int64, error)
	CountDocsWithRouting(index, routing string) (int64, error)
}

var ElasticPool IElasticConnector
