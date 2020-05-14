package pks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry/metric-store-release/src/pkg/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

// Client handles interaction with the PKS API
type Client struct {
	url        string
	httpClient *http.Client
	log        *logger.Logger
}

func NewClient(addr string, httpClient *http.Client, log *logger.Logger) *Client {
	//TODO use an interface for HTTPClient
	return &Client{
		url:        addr,
		httpClient: httpClient,
		log:        log,
	}
}

func (client *Client) GetClusters(authorization string) ([]string, error) {
	url := fmt.Sprintf("%s/v1/clusters", client.url)
	client.log.Debug("cluster request", logger.String("url", url))
	pksRequest, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}
	pksRequest.Header.Add("Authorization", authorization)

	responseBody, err := client.doRequest(pksRequest, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var pksClusters []PksClustersResponse
	err = json.Unmarshal(responseBody, &pksClusters)
	if err != nil {
		return nil, err
	}

	var clusterNames []string
	for _, response := range pksClusters {
		clusterNames = append(clusterNames, response.Name)
	}
	return clusterNames, nil
}

func (client *Client) GetCredentials(clusterName string, authorization string) (*Credentials, error) {
	pksRequest, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/clusters/%s/binds", client.url, clusterName), bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}
	pksRequest.Header.Add("Authorization", authorization)
	pksRequest.Header.Add("Content-Type", "application/json")
	pksRequest.Header.Add("Media-Type", "application/json")

	responseBody, err := client.doRequest(pksRequest, http.StatusCreated)
	if err != nil {
		return nil, err
	}

	pksCredentials := &CredentialsResponse{}
	client.log.Debug("Credential Response", zap.ByteString("body", responseBody))
	err = json.Unmarshal(responseBody, pksCredentials)
	if err != nil {
		return nil, err
	}

	return &Credentials{
		CaData:    pksCredentials.Clusters[0].Cluster.CertificateAuthorityData,
		UserToken: pksCredentials.Users[0].User.Token,
		Server:    pksCredentials.Clusters[0].Cluster.Server}, nil
}

func (client *Client) doRequest(req *http.Request, expectedStatus int) ([]byte, error) {
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("unexpected response code %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

type Credentials struct {
	CaData    string
	UserToken string
	Server    string
}

type PksClustersResponse struct {
	Name string `json:"name"`
}

type CredentialsResponse struct {
	Clusters []clusters `json:"clusters"`
	Users    []users    `json:"users"`
}

type clusters struct {
	Name    string  `json:"name"`
	Cluster cluster `json:"cluster"`
}

type cluster struct {
	CertificateAuthorityData string `json:"certificate-authority-data"`
	Server                   string `json:"server"`
}

type users struct {
	Name string `json:"name"`
	User user   `json:"user"`
}

type user struct {
	Token string `json:"token"`
}