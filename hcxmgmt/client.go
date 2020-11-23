package hcxmgmt

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Username   string
	Password   string
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type updateConfigurationModule struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}

type updateConfigurationModuleBody struct {
	ID     int                       `json:"id"`
	Update updateConfigurationModule `json:"_update"`
}

type enterprise_get_object_groups struct {
	Type string `json:"type"`
}

// NewClient -
func NewClient(hcx *string, username *string, password *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 300 * time.Second},
		// Default Hashicups URL
		HostURL: *hcx,
	}

	if (&hcx != nil) && (&username != nil) && (&password != nil) {

		c.Username = *username
		c.Password = *password

	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) (*http.Response, []byte, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(c.Username, c.Password)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	//b, err := json.MarshalIndent(req, "", "  ")
	//log.Printf("%s", b)
	//return nil, nil, fmt.Errorf("req: %s", req)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode != http.StatusOK {
		if res.StatusCode != http.StatusNoContent {
			if res.StatusCode != http.StatusAccepted {
				return nil, nil, fmt.Errorf("status: %d, response: %s", res.StatusCode, body)
			}
		}
	}

	return res, body, err
}

// DeepCopy ...
func DeepCopy(src map[string]interface{}) (map[string]interface{}, error) {

	var dst map[string]interface{}

	if src == nil {
		fmt.Println("Error src")
		return nil, fmt.Errorf("src cannot be nil")
	}

	bytes, err := json.Marshal(src)

	if err != nil {
		return nil, fmt.Errorf("Unable to marshal src: %s", err)
	}
	err = json.Unmarshal(bytes, &dst)

	if err != nil {
		fmt.Println("Error unmarshal")
		fmt.Println(err)
		return nil, fmt.Errorf("Unable to unmarshal into dst: %s", err)
	}
	return dst, nil

}
