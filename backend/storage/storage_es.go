package storage

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/baixiaozhou/perfmonitorscan/backend/conf"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"net/http"
	"strconv"
	"strings"
)

const (
	HTTP = "http://"
)

var (
	//go:embed monitoring_cpu_data_mapping.json
	mappingFile embed.FS
)

func InitEsClient(dbConf conf.DBConfig) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			HTTP + dbConf.Host + ":" + dbConf.Port,
		},
		Username: dbConf.Username,
		Password: dbConf.Password,
	}
	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		conf.Logger.Error("Failed to init es client:", err.Error())
		return nil, err
	}
	res, err := esClient.Info()
	if err != nil {
		conf.Logger.Error("Failed to get es info:", err.Error())
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		conf.Logger.Error(res.Status(), res.String())
	}

	return esClient, nil
}

func getDataNodeCount(esClient *elasticsearch.Client) (int, error) {
	req := esapi.CatNodesRequest{
		Format: "json",
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		conf.Logger.Error("Failed to get node count:", err.Error())
		return 0, err
	}
	defer res.Body.Close()
	if res.IsError() {
		conf.Logger.Error(res.Status(), res.String())
	}

	var (
		data      []types.NodesRecord
		nodeCount = 0
	)

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		conf.Logger.Error("Failed to get node count:", err.Error())
		return 0, err
	}
	for _, node := range data {
		// data node
		if strings.Contains(*node.NodeRole, "d") {
			nodeCount++
		}
	}

	return nodeCount, nil
}

func InitIndex(esClient *elasticsearch.Client, indexName string) error {
	exists, err := esClient.Indices.Exists([]string{indexName})
	if err != nil {
		conf.Logger.Error("Failed to get index:", err.Error())
	}
	if exists.StatusCode == http.StatusOK {
		conf.Logger.Info("The index is exists:", indexName)
		return nil
	}

	// index settings and mappings
	dataNodeCount, err := getDataNodeCount(esClient)
	if err != nil {
		dataNodeCount = 1
	}

	mappings, err := mappingFile.ReadFile("monitoring_cpu_data_mapping.json")
	if err != nil {
		conf.Logger.Error("Failed to read monitoring cpu data mapping:", err.Error())
		return err
	}

	var (
		indexConfig bytes.Buffer
	)

	indexConfig.Write([]byte(strings.Replace(string(mappings), "dataNodeCount", strconv.Itoa(dataNodeCount), -1)))

	fmt.Println(string(indexConfig.Bytes()))
	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  bytes.NewReader(indexConfig.Bytes()),
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		conf.Logger.Error("Failed to create index:", err.Error())
		return err
	}

	fmt.Println(res.String())
	if res.StatusCode != http.StatusOK {
		conf.Logger.Error("Failed to create index:", res.Status())
	}
	return nil
}

func SaveEsData(esClient *elasticsearch.Client, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		conf.Logger.Error("Failed to marshal data:", err.Error())
		return err
	}
	fmt.Println(string(jsonData))
	req := esapi.IndexRequest{
		Index:   INDEX_NAME,
		Body:    bytes.NewReader(jsonData),
		Refresh: "true",
	}
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		conf.Logger.Error("Failed to save data to es:", err.Error())
		return err
	}
	if res.StatusCode != http.StatusCreated {
		conf.Logger.Error("Failed to save data to es:", res.Status())
		return err
	}
	return nil
}
