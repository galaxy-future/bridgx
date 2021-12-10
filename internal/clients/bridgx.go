package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/galaxy-future/BridgX/internal/logs"
	gf_cluster "github.com/galaxy-future/BridgX/pkg/gf-cluster"
	"go.uber.org/zap"
)

var bridgxClient *Client

type Client struct {
	ServerAddress string
	httpClient    *http.Client
}

func NewServerAdderss(serverAddress string) *Client {
	return &Client{
		ServerAddress: serverAddress,
		httpClient: &http.Client{
			Timeout: 5000 * time.Millisecond,
		},
	}
}

func InitializeClient(serverAddress string) {
	bridgxClient = NewServerAdderss(serverAddress)
}

func GetClient() *Client {
	return bridgxClient
}

func (client *Client) GetUnusedCluster(token string, pageNumber, pageSize int) (*gf_cluster.ListBirdgxClusterByClusterResponse, error) {
	request := gf_cluster.ListBridgxClusterByTagRequest{
		Tags:       map[string]string{gf_cluster.UsageKey: gf_cluster.UnusedValue},
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	data, err := json.Marshal(&request)
	if err != nil {
		logs.Logger.Error("marshall request failed", zap.Error(err))
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/cluster/list_by_tags", client.ServerAddress), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response gf_cluster.ListBirdgxClusterByClusterResponse
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (client *Client) UpdateBridgxClusterUsingTag(token string, clusterName string, using bool) (*gf_cluster.EditBridgxClusterTagResponse, error) {
	request := gf_cluster.EditBridgxClusterTagRequst{
		Tags:        map[string]string{gf_cluster.UsageKey: gf_cluster.GalaxyfutureCloudUsage},
		ClusterName: clusterName,
	}

	if !using {
		request.Tags[gf_cluster.UsageKey] = gf_cluster.UnusedValue
	}

	data, err := json.Marshal(&request)
	if err != nil {
		logs.Logger.Error("marshall request failed", zap.Error(err))
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/cluster/edit_tags", client.ServerAddress), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response gf_cluster.EditBridgxClusterTagResponse
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (client *Client) GetBridgxClusterInstances(token string, clusterName string, pageSize, pageNum int) (*gf_cluster.GetBridgxClusterInstanceResponse, error) {

	data := "{}"
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/instance/describe_all?cluster_name=%s&status=running", client.ServerAddress, clusterName), strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response gf_cluster.GetBridgxClusterInstanceResponse
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (client *Client) GetBridgxClusterAllInstances(token string, clusterName string) ([]*gf_cluster.BridgxInstance, error) {
	response, err := client.GetBridgxClusterInstances(token, clusterName, 0, 50)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := response.Data.InstanceList
	pager := 1
	for len(result) < response.Data.Pager.Total {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			response, err = client.GetBridgxClusterInstances(token, clusterName, pager, 50)
			if err != nil {
				return nil, err
			}
			if response.Code != 200 {
				return nil, fmt.Errorf(response.Msg)
			}
			result = append(result, response.Data.InstanceList...)
			pager += 1
		}
	}
	return result, nil
}

func (client *Client) GetBriodgxClusterDetails(token string, clusterName string) (*gf_cluster.BridgxClusterDetailsResponse, error) {
	data := "{}"
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/cluster/name/%s", client.ServerAddress, clusterName), strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response gf_cluster.BridgxClusterDetailsResponse
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (client *Client) GetAKSKClusterDetails(token string, clusterName string) (*gf_cluster.GetAKSKResponse, error) {
	data := "{}"
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/cloud_account/info?cluster_name=%s", client.ServerAddress, clusterName), strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response gf_cluster.GetAKSKResponse
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (client *Client) Login(username string, password string) (token string, err error) {

	request := struct {
		Username string
		Password string
	}{Username: username, Password: password}

	data, _ := json.Marshal(&request)

	resp, err := client.httpClient.Post(fmt.Sprintf("%s/user/login", client.ServerAddress), "application/json", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	response := struct {
		Code int
		Msg  string
		Data string
	}{}
	err = json.Unmarshal(respData, &response)
	if err != nil {
		return "", err
	}
	if response.Code != 200 {
		return "", fmt.Errorf("Wrong code: %v, response : %s", err, respData)
	}
	return response.Data, nil
}
