package rullerlist

import (
	"io/ioutil"
	"os"
	"path"
	"encoding/xml"
	"fmt"
)

type RullerList struct {
	list []string
}

type Root struct {
	XMLName   xml.Name   `xml:"rules"`
	Rules []Rule `xml:"rule>key"`
}

type Rule struct {
	Key string `xml:",chardata"`
}


func (this *RullerList) IsInList(elem string) bool {
	for _, val := range this.list {
		if val == elem {
			return true
		}
	}
	return false
}

func (this *RullerList) _addRules(path string) {
	xmlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}

	var dict Root
	xml.Unmarshal([]byte(xmlFile), &dict)

	for _, rule := range dict.Rules {
		if this.IsInList(rule.Key) == false {
			this.list = append(this.list, rule.Key)
		}
	}

}

func (this *RullerList) GetSonarQubeRulesList(dir_to_find string) {
	if info, err := os.Stat(dir_to_find); err == nil {
		if info.IsDir() == true {
			files, err := ioutil.ReadDir(dir_to_find)
			if err != nil {
				fmt.Printf("Error %s in reading of dir %s\n", err, dir_to_find)
				return
			}

			for _, file := range files {
				if path.Ext(file.Name()) == ".xml" {
					file_name := dir_to_find + "/" + file.Name()
					this._addRules(file_name)
				}
			}
		}
	}
	return
}

