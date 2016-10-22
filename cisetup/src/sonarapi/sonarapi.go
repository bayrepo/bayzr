package sonarapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type SonarApi struct {
	auth_ok  bool
	address  string
	login    string
	password string
}

func MakeSonarApi(addr string, lg string, pass string) *SonarApi {
	return &SonarApi{false, addr, lg, pass}
}

func readBody(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return []byte{}, fmt.Errorf("Unexpected status code", resp.StatusCode)
	}

	return body, nil

}

func (this *SonarApi) Connect() error {
	client := &http.Client{}
	url_str := this.address + "/api/authentication/login"
	form := url.Values{}
	form.Set("login", this.login)
	form.Set("password", this.password)
	req, err := http.NewRequest("POST", url_str, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := readBody(resp)
	if err != nil {
		return err
	}

	resp, err = http.Get(this.address + "/api/authentication/validate")
	if err != nil {
		return err
	}

	body, err = readBody(resp)
	if err != nil {
		return err
	}

	var result map[string]interface{}

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if _, ok := result["valid"]; ok == false {
		return fmt.Errorf("Incorrect auth responce")
	}

	if result["valid"].(bool) != true {
		return fmt.Errorf("Login and password not accepted")
	}

	this.auth_ok = true

	return nil
}

func (this *SonarApi) GetPluginList(tp bool) ([]string, error) {
	plagins_list := []string{}
	if this.auth_ok != true {
		return plagins_list, fmt.Errorf("Auth not completed")
	}

	client := &http.Client{}
	url_str := ""
	if tp {
		url_str = this.address + "/api/plugins/available"
	} else {
		url_str = this.address + "/api/plugins/installed"
	}
	form := url.Values{}
	req, err := http.NewRequest("POST", url_str, strings.NewReader(form.Encode()))
	if err != nil {
		return plagins_list, err
	}
	req.SetBasicAuth(this.login, this.password)

	resp, err := client.Do(req)
	if err != nil {
		return plagins_list, err
	}

	body, err := readBody(resp)
	if err != nil {
		return plagins_list, err
	}

	var result map[string]interface{}

	if err := json.Unmarshal(body, &result); err != nil {
		return plagins_list, err
	}

	if _, ok := result["plugins"]; ok == false {
		return plagins_list, nil
	}

	for _, item := range result["plugins"].([]interface{}) {
		res := item.(map[string]interface{})
		if _, ok := res["key"]; ok == false {
			continue
		}
		if _, ok := res["name"]; ok == false {
			continue
		}
		plagins_list = append(plagins_list, res["key"].(string))
	}

	return plagins_list, nil
}

func (this *SonarApi) InstallPlugin(plugin string, avail []string, installed []string, tp bool) error {
	if this.auth_ok != true {
		return fmt.Errorf("Auth not completed")
	}
	if tp {
		for _, val := range installed {
			if val == plugin {
				return nil
			}
		}
		is_avail := false
		for _, val := range avail {
			if val == plugin {
				is_avail = true
				break
			}
		}
		if is_avail == false {
			return fmt.Errorf("Plugin does not exists %s", plugin)
		}
	} else {
		is_avail := false
		for _, val := range installed {
			if val == plugin {
				is_avail = true
				break
			}
		}
		if is_avail == false {
			return nil
		}
	}
	fmt.Println("Ready to install ", plugin)

	client := &http.Client{}

	url_str := ""
	if tp {
		url_str = this.address + "/api/plugins/install"
	} else {
		url_str = this.address + "/api/plugins/uninstall"
	}

	form := url.Values{}
	form.Set("key", plugin)
	req, err := http.NewRequest("POST", url_str, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(this.login, this.password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	_, err = readBody(resp)
	if err != nil {
		return err
	}

	return nil

}

func (this *SonarApi) SetSonarOption(id string, value string) error {
	if this.auth_ok != true {
		return fmt.Errorf("Auth not completed")
	}
	
	client := &http.Client{}

	url_str := this.address + "/api/properties"


	form := url.Values{}
	form.Set("id", id)
	form.Set("value", value)
	req, err := http.NewRequest("POST", url_str, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(this.login, this.password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	_, err = readBody(resp)
	if err != nil {
		return err
	}

	return nil

}
