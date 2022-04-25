package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	IE "autocall/internal/elastic"
	model "autocall/internal/elastic/model"

	elastic "github.com/olivere/elastic/v7"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Username              string
	Password              string
	Host                  []string
	Nodename              string
	ClusterName           string
	TimeoutConnect        int
	MaxIdleConnsPerHost   int
	PoolSize              int
	RetryOnStatus         []int
	DisableRetry          bool
	EnableRetryOnTimeout  bool
	MaxRetries            int
	RetriesTimeout        int
	ResponseHeaderTimeout int
	Index                 string
}

type ElasticClient struct {
	user                  string
	password              string
	host                  []string
	nodename              string
	clustername           string
	timeoutconnect        int
	maxIdleConnsPerHost   int
	poolSize              int
	retryOnStatus         []int
	disableRetry          bool
	enableRetryOnTimeout  bool
	maxRetries            int
	retriesTimeout        int
	responseHeaderTimeout int
	Client                *elastic.Client
	Index                 string
}

var (
	ElasticPool ElasticClient
)

func NewElasticClient(config Config) IE.IElasticConnector {
	elasticPool := &ElasticClient{
		host:                  config.Host,
		user:                  config.Username,
		password:              config.Password,
		nodename:              config.Nodename,
		clustername:           config.ClusterName,
		timeoutconnect:        config.TimeoutConnect,
		maxIdleConnsPerHost:   config.MaxIdleConnsPerHost,
		poolSize:              config.PoolSize,
		retryOnStatus:         config.RetryOnStatus,
		disableRetry:          config.DisableRetry,
		enableRetryOnTimeout:  config.EnableRetryOnTimeout,
		maxRetries:            config.MaxRetries,
		retriesTimeout:        config.RetriesTimeout,
		responseHeaderTimeout: config.ResponseHeaderTimeout,
	}
	err := elasticPool.GetConn()
	if err != nil {
		panic(err)
	}
	return elasticPool
}

func (e *ElasticClient) GetConn() error {

	client, err := elastic.NewClient(
		elastic.SetBasicAuth(e.user, e.password),
		elastic.SetURL(strings.Join(e.host[:], ",")),
		elastic.SetSniff(false),
		elastic.SetMaxRetries(e.maxRetries),
		elastic.SetHealthcheckInterval(time.Duration(e.responseHeaderTimeout)*time.Second),
		elastic.SetHealthcheckTimeout(time.Duration(e.responseHeaderTimeout)*time.Second),
	)

	if err != nil {
		// Handle error
		log.Error("ElasticClient - GetConn => Error : ", err.Error())
		return err
	}
	e.Client = client
	return nil
}

func (e *ElasticClient) GetClient() *elastic.Client {
	return e.Client
}

//Ping Kiểm tra kết nối
func (e *ElasticClient) Ping() {
	for _, hostUrl := range e.host {
		info, code, err := e.Client.Ping(hostUrl).Do(context.Background())
		if err != nil {
			log.Error("Error getting response: ", err)
			panic(err.Error())
		}
		if code != 200 {
			log.Error("Connect Elasticsearch fail, code : ", code)
			panic(info)
		}
		log.Info("Elasticsearch returned with code ", code, " and version ", info.Version.Number, "\n")
	}
}

func (e *ElasticClient) CreateIndex(name string, jsonBody string) (bool, error) {
	log.Debug("CreateIndex: name: ", name, "jsonBody: ", jsonBody)
	isExist, err := e.Client.IndexExists(name).Do(context.Background())
	if err != nil {
		log.Error("ElasticClient - CreateIndex => Error : ", err.Error())
		return false, err
	}
	if !isExist {
		// Create a new index.
		createIndex, err := e.Client.CreateIndex(name).BodyString(jsonBody).Do(context.Background())
		if err != nil {
			// Handle error
			log.Error("ElasticClient - CreateIndex => Error : ", err.Error())
			return false, err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
		if err != nil {
			log.Error("CreateIndex: Error getting response: ", err)
			return false, err
		}
		log.Info("ElasticClient - CreateIndex => Create : ", name)
		return true, nil
	} else {
		log.Info("ElasticClient - CreateIndex => Index is Exist : ", name)
		return true, nil
	}
}
func (e *ElasticClient) CreateAlias(index string, alias string) (bool, error) {
	res, err := e.Client.Alias().Action().Add(index, alias).Do(context.Background())
	if err != nil {
		log.Error("ElasticClient - CreateAlias => Error : ", err.Error())
		return false, err
	}
	log.Info("ElasticClient - CreateAlias => Create : ", alias, " res: ", res)
	return true, nil
}

func (e *ElasticClient) CheckAliasExist(index, alias string) (bool, error) {
	res, err := e.Client.Aliases().Index(index).Do(context.Background())
	if err != nil {
		log.Error("ElasticClient - CheckAliasExist => Error : ", err.Error())
		return false, err
	}
	if len(res.Indices) > 0 {
		indices := res.Indices[index]
		if indices.HasAlias(alias) {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, nil
}

func (e *ElasticClient) CreateDocs(name string, id string, body interface{}) (model.Response, error) {

	docId := fmt.Sprintf("%v", id)
	res, err := e.Client.Index().Index(name).Id(docId).BodyJson(body).Refresh("true").Do(context.Background())
	if err != nil {
		log.Error("CreateDocs: Error when request : ", err)
		return model.Response{}, err
	}
	j, err := json.Marshal(res)
	if err != nil {
		log.Error("SearchDoc: Error parse to json : ", err)
		return model.Response{}, err
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(j), &jsonMap)
	if err != nil {
		log.Error("SearchDoc: Error parse to map : ", err)
		return model.Response{}, err
	}
	response := model.Response{
		StatusCode: res.Status,
		Body:       jsonMap,
	}
	return response, err
}

func (e *ElasticClient) CreateDocWithRouting(index, routing, id, index_type string, body interface{}) (model.Response, error) {

	docId := fmt.Sprintf("%v", id)
	res, err := e.Client.Index().Index(index).Routing(routing).Id(docId).Type(index_type).BodyJson(body).Refresh("true").Do(context.Background())
	if err != nil {
		log.Error("CreateDocs: Error when request : ", err)
		return model.Response{}, err
	}
	j, err := json.Marshal(res)
	if err != nil {
		log.Error("SearchDoc: Error parse to json : ", err)
		return model.Response{}, err
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(j), &jsonMap)
	if err != nil {
		log.Error("SearchDoc: Error parse to map : ", err)
		return model.Response{}, err
	}
	response := model.Response{
		StatusCode: res.Status,
		Body:       jsonMap,
	}
	return response, err
}
func (e *ElasticClient) GetDoc(name string, id string) (model.Response, error) {
	res, err := e.Client.Get().Index(name).Id(id).Type("cdr").Pretty(true).Routing(name).Do(context.Background())
	if err != nil {
		log.Error("GetDoc: Error getting response: ", err)
		return model.Response{}, err
	}
	j, err := json.Marshal(res)
	if err != nil {
		log.Error("GetDoc: Error parse to json : ", err)
		return model.Response{}, err
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(j), &jsonMap)
	if err != nil {
		log.Error("GetDoc: Error parse to map : ", err)
		return model.Response{}, err
	}
	response := model.Response{
		Body: jsonMap,
	}
	return response, err
}

func (e *ElasticClient) GetDocWithRouting(index, routing, id string) (model.Response, error) {
	res, err := e.Client.Get().Index(index).Routing(routing).Id(id).Pretty(true).Do(context.Background())
	if err != nil {
		log.Error("GetDoc: Error getting response: ", err)
		return model.Response{}, err
	}
	j, err := json.Marshal(res)
	if err != nil {
		log.Error("GetDoc: Error parse to json : ", err)
		return model.Response{}, err
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(j), &jsonMap)
	if err != nil {
		log.Error("GetDoc: Error parse to map : ", err)
		return model.Response{}, err
	}
	response := model.Response{
		Body: jsonMap,
	}
	return response, err
}

func (e *ElasticClient) SearchDoc(index string, searchSource *elastic.SearchSource) (model.Response, error) {

	res, err := e.Client.Search().Index(index).SearchSource(searchSource).Do(context.Background())
	if err != nil {
		log.Error("SearchDoc: Error getting response: ", err)
		return model.Response{}, err
	}
	j, err := json.Marshal(res)
	if err != nil {
		log.Error("SearchDoc: Error parse to json : ", err)
		return model.Response{}, err
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(j), &jsonMap)
	if err != nil {
		log.Error("SearchDoc: Error parse to map : ", err)
		return model.Response{}, err
	}
	response := model.Response{
		StatusCode: res.Status, //int
		// Header:     res.Header,
		Body: jsonMap,
	}
	return response, err

}

func (e *ElasticClient) SearchDocWithPaging(index string, searchSource *elastic.SearchSource, size, from int) (model.Response, error) {

	res, err := e.Client.Search().Index(index).SearchSource(searchSource).Size(size).From(from).Do(context.Background())
	if err != nil {
		log.Error("SearchDoc: Error getting response: ", err)
		return model.Response{}, err
	}
	j, err := json.Marshal(res)
	if err != nil {
		log.Error("SearchDoc: Error parse to json : ", err)
		return model.Response{}, err
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(j), &jsonMap)
	if err != nil {
		log.Error("SearchDoc: Error parse to map : ", err)
		return model.Response{}, err
	}
	response := model.Response{
		StatusCode: res.Status,
		Header:     res.Header,
		Body:       jsonMap,
	}
	return response, err

}

func (e *ElasticClient) GetDocBySearch(index string, searchSource *elastic.SearchSource) (model.Response, error) {

	res, err := e.Client.Search().Index(index).Routing(index).SearchSource(searchSource).Do(context.Background())
	if err != nil {
		log.Error("SearchDoc: Error getting response: ", err)
		return model.Response{}, err
	}
	j, err := json.Marshal(res)
	if err != nil {
		log.Error("SearchDoc: Error parse to json : ", err)
		return model.Response{}, err
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(j), &jsonMap)
	if err != nil {
		log.Error("SearchDoc: Error parse to map : ", err)
		return model.Response{}, err
	}
	body, _ := jsonMap["hits"].(map[string]interface{})
	listLog, _ := body["hits"].([]interface{})
	var result interface{}
	if len(listLog) > 0 {
		item, _ := listLog[0].(map[string]interface{})
		result = item["_source"]
	}
	resBody := make(map[string]interface{})
	resBody["cdr"] = result
	response := model.Response{
		StatusCode: res.Status, //int
		Header:     res.Header,
		Body:       resBody,
	}
	return response, err

}

func (e *ElasticClient) CountDocs(index string) (int64, error) {
	res, err := e.Client.Count(index).Do(context.Background())
	if err != nil {
		log.Error("CountDoc: Error getting response: ", err)
		return 0, err
	}
	return res, err
}

func (e *ElasticClient) CountDocsWithRouting(index, routing string) (int64, error) {
	res, err := e.Client.Count(index).Routing(routing).Do(context.Background())
	if err != nil {
		log.Error("CountDoc: Error getting response: ", err)
		return 0, err
	}
	return res, err
}
