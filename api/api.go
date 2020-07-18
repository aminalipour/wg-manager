package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// API is a utility for communicating with the Mullvad API
type API struct {
	Username string
	Password string
	BaseURL  string
	AdminBaseURL  string
	Client   *http.Client
}

// WireguardPeerList is a list of Wireguard peers
type WireguardPeerList []WireguardPeer

// WireguardPeer is a wireguard peer
type WireguardPeer struct {
	IPv4   string `json:"ipv4"`
	IPv6   string `json:"ipv6"`
	Ports  []int  `json:"ports"`
	Pubkey string `json:"pubkey"`
}
type PeerUsages struct{
	Receive int64 `json:"receive"`
	Transmit int64 `json:"transmit"`
}

type PeerUsagesData map[string][]PeerUsages

// GetWireguardPeers fetches a list of wireguard peers from the API and returns it
func (a *API) GetWireguardPeers() (WireguardPeerList, error) {
	req, err := http.NewRequest("GET", a.BaseURL+"/wg/active-pubkeys/v2/", nil)
	if err != nil {
		return WireguardPeerList{}, err
	}

	if a.Username != "" && a.Password != "" {
		req.SetBasicAuth(a.Username, a.Password)
	}

	response, err := a.Client.Do(req)
	if err != nil {
		return WireguardPeerList{}, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return WireguardPeerList{}, err
	}

	var decodedResponse WireguardPeerList
	err = json.Unmarshal(body, &decodedResponse)
	if err != nil {
		return WireguardPeerList{}, fmt.Errorf("error decoding wireguard peers")
	}

	return decodedResponse, nil
}

// GetWireguardPeers fetches a list of wireguard peers from the API and returns it
func (a *API) UpdateServerData(connectedPeers int, CPUUsage float64, receive uint64,transfer uint64) {
	values := map[string]string{
		"connected_peers": strconv.Itoa(connectedPeers),
		"cpu_usage": strconv.FormatFloat(CPUUsage,'f', 6, 64),
		"receive": strconv.FormatUint(receive,10),
		"transfer": strconv.FormatUint(transfer,10),
	}

	jsonValue, _ := json.Marshal(values)

	req, _ := http.NewRequest("POST", a.AdminBaseURL+"/update-server-data/",bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	if a.Username != "" && a.Password != "" {
		req.SetBasicAuth(a.Username, a.Password)
	}

	response, err := a.Client.Do(req)
	if err!=nil {
		panic(err.Error())
	}

	defer response.Body.Close()
}


func (a *API) UpdatePeersBandwidthUsages(peersUsages PeerUsagesData) {
	values := map[string]PeerUsagesData{
		"peers": peersUsages,
	}

	jsonValue, _ := json.Marshal(values)

	req, _ := http.NewRequest("POST", a.AdminBaseURL+"/update-peers-bandwidth-usages/",bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	if a.Username != "" && a.Password != "" {
		req.SetBasicAuth(a.Username, a.Password)
	}

	response, err := a.Client.Do(req)
	if err!=nil {
		fmt.Println(err.Error())
	}

	defer response.Body.Close()
}
