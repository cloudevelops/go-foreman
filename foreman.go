package foreman

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	//	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Foreman struct {
	Hostname  string
	Username  string
	password  string
	VerifySSL bool
	BaseURL   string
	client    *http.Client
	auth      string
}

func NewForeman(HostName string, UserName string, Password string) *Foreman {
	var foreman *Foreman
	var tr *http.Transport

	foreman = new(Foreman)
	foreman.Hostname = HostName
	foreman.Username = UserName
	foreman.password = Password
	foreman.VerifySSL = false
	foreman.BaseURL = "https://" + foreman.Hostname + "/api/"
	foreman.auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(UserName+":"+Password))

	if foreman.VerifySSL {
		tr = &http.Transport{}
	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	foreman.client = &http.Client{Transport: tr}
	return foreman
}

func (foreman *Foreman) Post(endpoint string, jsonData []byte) (map[string]interface{}, error) {
	var target string
	var data interface{}
	var req *http.Request

	target = foreman.BaseURL + endpoint
	//fmt.Println("POST form " + target)
	req, err := http.NewRequest("POST", target, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(jsonData)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", foreman.auth)
	r, err := foreman.client.Do(req)
	defer r.Body.Close()
	if err != nil {
		fmt.Println("Error while posting")
		fmt.Println(err)
		return nil, err
	}
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return nil, errors.New("HTTP Error " + r.Status)
	}

	response, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error while reading body")
		fmt.Println(err)
		return nil, err
	}
	err = json.Unmarshal(response, &data)
	if err != nil {
		fmt.Println("Error while processing JSON")
		fmt.Println(err)
		return nil, err
	}
	m := data.(map[string]interface{})
	return m, nil
}

func (foreman *Foreman) Put(endpoint string, jsonData []byte) (map[string]interface{}, error) {
	var target string
	var data interface{}
	var req *http.Request

	target = foreman.BaseURL + endpoint
	//fmt.Println("POST form " + target)
	req, err := http.NewRequest("PUT", target, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(jsonData)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", foreman.auth)
	r, err := foreman.client.Do(req)
	defer r.Body.Close()
	if err != nil {
		fmt.Println("Error while posting")
		fmt.Println(err)
		return nil, err
	}
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return nil, errors.New("HTTP Error " + r.Status)
	}

	response, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error while reading body")
		fmt.Println(err)
		return nil, err
	}
	err = json.Unmarshal(response, &data)
	if err != nil {
		fmt.Println("Error while processing JSON")
		fmt.Println(err)
		return nil, err
	}
	m := data.(map[string]interface{})
	return m, nil
}

func (foreman *Foreman) Get(endpoint string) (map[string]interface{}, error) {
	var target string
	var data interface{}

	target = foreman.BaseURL + endpoint
	req, err := http.NewRequest("GET", target, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", foreman.auth)
	r, err := foreman.client.Do(req)
	defer r.Body.Close()
	if err != nil {
		fmt.Println("Error while getting")
		fmt.Println(err)
		return nil, err
	}
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return nil, errors.New("HTTP Error " + r.Status)
	}

	response, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error while reading body")
		fmt.Println(err)
		return nil, err
	}
	err = json.Unmarshal(response, &data)
	if err != nil {
		fmt.Println("Error while processing JSON")
		fmt.Println(err)
		return nil, err
	}
	m := data.(map[string]interface{})
	return m, nil
}

func (foreman *Foreman) Delete(endpoint string) (map[string]interface{}, error) {
	var target string
	var data interface{}
	var req *http.Request

	target = foreman.BaseURL + endpoint
	req, err := http.NewRequest("DELETE", target, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", foreman.auth)
	r, err := foreman.client.Do(req)
	if r != nil {
		defer r.Body.Close()
	}
	if err != nil {
		fmt.Println("Error while deleting")
		fmt.Println(err)
		return nil, err
	}
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return nil, errors.New("HTTP Error " + r.Status)
	}
	response, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error while reading body")
		fmt.Println(err)
		return nil, err
	}
	err = json.Unmarshal(response, &data)
	if err != nil {
		fmt.Println("Error while processing JSON")
		fmt.Println(err)
		return nil, err
	}
	m := data.(map[string]interface{})
	return m, nil
}

type Host struct {
	HostGroupId string `json:"hostgroup_id"`
	Name        string `json:"name"`
	Mac         string `json:"mac"`
	Build       bool   `json:"build"`
}

type HostMap map[string]Host

func (foreman *Foreman) CreateHost(HostGroupId int, Name string, Mac string) (string, error) {
	var hostMap map[string]Host
	var err error

	hostMap = make(HostMap)
	hostMap["host"] = Host{
		HostGroupId: strconv.Itoa(HostGroupId),
		Name:        Name,
		Mac:         Mac,
		Build:       true,
	}
	jsonText, err := json.Marshal(hostMap)
	data, err := foreman.Post("hosts", jsonText)
	if err != nil {
		fmt.Print("Error ")
		fmt.Println(err)
		return "", err
	}
	return strconv.FormatFloat(data["id"].(float64), 'f', 0, 64), nil
}

func (foreman *Foreman) DeleteHost(HostID string) error {
	var err error

	_, err = foreman.Delete("hosts/" + HostID)
	return err
}

func (foreman *Foreman) SearchResource(Resource string, Query string) (map[string]interface{}, error) {
	escapedQuery := strings.Replace(Query, " ", "+", -1)
	data, err := foreman.Get(Resource + "?search=" + escapedQuery + "&per_page=10000")
	if err != nil {
		fmt.Print("Error searching for resource, retry in 5s:")
		fmt.Println(err)
		time.Sleep(5 * time.Second)
		data, err = foreman.Get(Resource + "?search=" + escapedQuery + "&per_page=10000")
		if err != nil {
			fmt.Print("Error searching for resource, retry in 15s:")
			fmt.Println(err)
			time.Sleep(15 * time.Second)
			data, err = foreman.Get(Resource + "?search=" + escapedQuery + "&per_page=10000")
			if err != nil {
				fmt.Print("Error searching for resource, retry in 60s:")
				fmt.Println(err)
				time.Sleep(60 * time.Second)
				data, err = foreman.Get(Resource + "?search=" + escapedQuery + "&per_page=10000")
				if err != nil {
					fmt.Print("Error searching for resource, failing")
					fmt.Println(err)
					return nil, err
				}
			}
		}
	}
	resultSlice := data["results"].([]interface{})
	if len(resultSlice) < 1 {
		fmt.Print("Resource not found")
		return nil, errors.New("Resource not found")
	}
	for _, resultItem := range resultSlice {
		resultData := resultItem.(map[string]interface{})
		//resultData := resultItem
		if title, ok := resultData["title"]; ok {
			if title == Query {
				return resultData, err
			}
		}
	}
	for _, resultItem := range resultSlice {
		resultData := resultItem.(map[string]interface{})
		if title, ok := resultData["name"]; ok {
			if title == Query {
				return resultData, err
			}
		}
	}
	//spew.Dump(data)
	fmt.Print("Resource not found")
	return nil, errors.New("Resource not found")
}

func (foreman *Foreman) SearchResourceName(Resource string, Query string) (map[string]interface{}, error) {
	escapedQuery := strings.Replace(Query, " ", "+", -1)
	data, err := foreman.Get(Resource + "?search=name~" + escapedQuery + "&per_page=10000")
	if err != nil {
		fmt.Print("Error searching for resource, retry in 5s:")
		fmt.Println(err)
		time.Sleep(5 * time.Second)
		data, err = foreman.Get(Resource + "?search=name~" + escapedQuery + "&per_page=10000")
		if err != nil {
			fmt.Print("Error searching for resource, retry in 15s:")
			fmt.Println(err)
			time.Sleep(15 * time.Second)
			data, err = foreman.Get(Resource + "?search=name~" + escapedQuery + "&per_page=10000")
			if err != nil {
				fmt.Print("Error searching for resource, retry in 60s:")
				fmt.Println(err)
				time.Sleep(60 * time.Second)
				data, err = foreman.Get(Resource + "?search=name~" + escapedQuery + "&per_page=10000")
				if err != nil {
					fmt.Print("Error searching for resource, failing")
					fmt.Println(err)
					return nil, err
				}
			}
		}
	}
	resultSlice := data["results"].([]interface{})
	if len(resultSlice) < 1 {
		fmt.Print("Resource not found")
		return nil, errors.New("Resource not found")
	}
	for _, resultItem := range resultSlice {
		resultData := resultItem.(map[string]interface{})
		//resultData := resultItem
		if title, ok := resultData["title"]; ok {
			if title == Query {
				return resultData, err
			}
		}
	}
	for _, resultItem := range resultSlice {
		resultData := resultItem.(map[string]interface{})
		if title, ok := resultData["name"]; ok {
			if title == Query {
				return resultData, err
			}
		}
	}
	//spew.Dump(data)
	fmt.Print("Resource not found")
	return nil, errors.New("Resource not found")
}

func (foreman *Foreman) SearchAnyResource(Resource string, Query string) (map[string]interface{}, error) {
	escapedQuery := strings.Replace(Query, " ", "+", -1)
	data, err := foreman.Get(Resource + "?search=" + escapedQuery + "&per_page=10000")
	if err != nil {
		fmt.Print("Error searching for resource, retry in 5s:")
		fmt.Println(err)
		time.Sleep(5 * time.Second)
		data, err = foreman.Get(Resource + "?search=" + escapedQuery + "&per_page=10000")
		if err != nil {
			fmt.Print("Error searching for resource, retry in 15s:")
			fmt.Println(err)
			time.Sleep(15 * time.Second)
			data, err = foreman.Get(Resource + "?search=" + escapedQuery + "&per_page=10000")
			if err != nil {
				fmt.Print("Error searching for resource, retry in 60s:")
				fmt.Println(err)
				time.Sleep(60 * time.Second)
				data, err = foreman.Get(Resource + "?search=" + escapedQuery + "&per_page=10000")
				if err != nil {
					fmt.Print("Error searching for resource, failing")
					fmt.Println(err)
					return nil, err
				}
			}
		}
	}
	//spew.Dump(data)
	resultSlice := data["results"].([]interface{})
	if len(resultSlice) < 1 {
		fmt.Print("Resource not found")
		return nil, errors.New("Resource not found")
	}
	return data, nil
}
