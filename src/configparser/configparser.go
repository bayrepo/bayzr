//    BayZR - utility for managing set of static analysis tools
//    Copyright (C) 2016  Alexey Berezhok
//    e-mail: bayrepo.info@gmail.com
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with this program.  If not, see <http://www.gnu.org/licenses/>.

package configparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ConfigparserContainer struct {
	compilatorsList []string            //default compilator
	extList         []string            //default extention
	addParameters   map[string][]string //extended filename=params
	bash            []string            //default bashcmd
	showerroroutput bool                //default stderr = on/off
	checkby         []string            //default empty
	output          string              //default custom
	template        string              //"FILE"|"LINE"|"SEV"|"ID"|"MESSAGE"
	wrapstrings     int64               //default 10
	ignore          []string            //default empty
	glovaldefs      []string            //default empty
	outputfile      string              //default report.log
	html_template   string              //default empty
	fnd_path        []string            //only global option not set in config file
	cc_replacer     bool                //use CC, CXX replace or just analyze output
	list_of_files   []string            //list of files to output
}

/*
* Внутрення функция, для разбивки строки на подстроки на основе множества разделителей
 */
func SplitOwnLongSep(data string, sep []string) []string {
	var result []string
	var data_copy []string
	for i, sep_one := range sep {
		result = []string{}
		if i == 0 {
			data_copy = strings.Split(data, strings.Trim(sep_one, " \n\t"))
		}
		for _, value := range data_copy {
			for _, sep_sec := range strings.Split(value, strings.Trim(sep_one, " \n\t")) {
				if sep_sec != "" {
					result = append(result, strings.Trim(sep_sec, " \n\t"))
				}
			}
		}
		data_copy = result
	}
	return result
}

/*
* Внутрення функция, для разбивки строки на подстроки на основе множества разделителей
* из разделителей не удаляются ведущие и конечные экстра символы
 */
func SplitOwnLongSepNoTrimSep(data string, sep []string) []string {
	var result []string
	var data_copy []string
	for i, sep_one := range sep {
		result = []string{}
		if i == 0 {
			data_copy = strings.Split(data, sep_one)
		}
		for _, value := range data_copy {
			for _, sep_sec := range strings.Split(value, sep_one) {
				if sep_sec != "" {
					result = append(result, strings.Trim(sep_sec, " \n\t"))
				}
			}
		}
		data_copy = result
	}
	return result
}

/*
* Укороченная и адаптированная версия Split заточенная под нужды конфигуратора
 */
func SplitOwn(data string) []string {

	var result []string
	for _, main_value := range strings.Fields(data) {
		trimmed_main_value := strings.Trim(main_value, " \n\t")
		if trimmed_main_value != "" {
			for _, value := range SplitOwnLongSep(trimmed_main_value, []string{",", ";", "|"}) {
				value = strings.Trim(value, " \n\t")
				if value != "" {
					result = append(result, value)
				}
			}
		}
	}
	return result
}

/*
* Функция удаления дубликатов из среза
 */
func RemoveDuplicate(data []string) []string {
	var result []string

	for _, value := range data {
		if value == "" {
			continue
		}
		found := false
		for _, fnd := range result {
			if value == fnd {
				found = true
			}
		}
		if found == false {
			result = append(result, value)
		}
	}
	return result
}

/*
* Дефолтовый конструктор. Создан на случай,
* если понадобится инициализировать конфигурационный файл не пустыми параметрами
 */
func CreateDefaultConfig() *ConfigparserContainer {
	return &ConfigparserContainer{[]string{},
		[]string{},
		map[string][]string{},
		[]string{"/usr/bin/bash", "-c"},
		true, []string{}, string("custom"),
		string("FILE|LINE|SEV|ID|MESSAGE"), 10, []string{}, []string{}, string("report.log"),
		"", []string{}, false, []string{"*"}}
}

/*
* Функция чтония конфигруационного файла
* входной параметр: configName - имя файла
* результат: описание ошибки или nil
 */
func (storage *ConfigparserContainer) ReadConfig(configName string) error {
	file, err := os.Open(configName)
	if err != nil {
		return fmt.Errorf("ReadConfig open config file error: %s\n", err)
	}
	defer file.Close()

	section := "default"
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if line != "" {
			section_reg := regexp.MustCompile("^\\s*\\[([\\S\\s\t]+)\\]\\s*$")
			if matches := section_reg.FindStringSubmatch(line); matches != nil {
				section = strings.Trim(matches[1], " \n\t")
			} else {
				section_comment := regexp.MustCompile("^\\s*;")
				if result := section_comment.FindString(line); result == "" {
					if section == "ignore" {
						storage.ignore = append(storage.ignore, strings.Trim(line, " \n\t"))
					} else {
						section_result := regexp.MustCompile("^\\s*([a-zA-z\\. ,/\\\\0-9\t]+)=([\\S\\s\t]+)$")
						if matches = section_result.FindStringSubmatch(line); matches != nil {
							section_key := strings.Trim(matches[1], " \n\t")
							if section == "default" {
								if section_key == "compilator" {
									storage.compilatorsList = append(storage.compilatorsList,
										SplitOwn(strings.Trim(matches[2], " \n\t"))...)
								}
								if section_key == "extention" {
									storage.extList = append(storage.extList,
										SplitOwn(strings.Trim(matches[2], " \n\t"))...)
								}
								if section_key == "bashcmd" {
									storage.bash = SplitOwn(strings.Trim(matches[2], " \n\t"))
								}
								if section_key == "stderr" {
									storage.showerroroutput = (strings.ToLower(strings.Trim(matches[2], " \n\t")) == "on")
								}
								if section_key == "replace" {
									storage.cc_replacer = (strings.ToLower(strings.Trim(matches[2], " \n\t")) == "on")
								}
								if section_key == "globaldefs" {
									storage.glovaldefs = SplitOwn(strings.Trim(matches[2], " \n\t"))
								}
								if section_key == "outputfile" {
									storage.outputfile = strings.Trim(matches[2], " \n\t")
								}
							}
							if section == "extended" && section_key != "" {
								storage.addParameters[section_key] = append(storage.addParameters[section_key],
									SplitOwn(strings.Trim(matches[2], " \n\t"))...)
							}
							if section == "plugins" {
								if section_key == "checkby" {
									storage.checkby = SplitOwn(strings.Trim(matches[2], " \n\t"))
								}
								if section_key == "output" {
									storage.output = strings.ToLower(strings.Trim(matches[2], " \n\t"))
								}
								if section_key == "template" {
									storage.template = strings.Trim(matches[2], " \n\t")
								}
								if section_key == "html_template" {
									storage.html_template = strings.Trim(matches[2], " \n\t")
								}
								if section_key == "wrapstrings" {
									i, err := strconv.ParseInt(strings.Trim(matches[2], " \n\t"), 10, 64)
									if err == nil {
										storage.wrapstrings = i
									}
								}
							}
						}
					}
				}
			}

		}
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("ReadConfig read config file error: %s\n", err)
			}
			break
		}
	}
	storage.compilatorsList = RemoveDuplicate(storage.compilatorsList)
	storage.extList = RemoveDuplicate(storage.extList)
	storage.checkby = RemoveDuplicate(storage.checkby)
	for key := range storage.addParameters {
		storage.addParameters[key] = RemoveDuplicate(storage.addParameters[key])
	}

	for key := range storage.extList {
		if []rune(storage.extList[key])[0] != rune('.') {
			storage.extList[key] = "." + storage.extList[key]
		}
	}

	return nil
}

/*
* Проверка, установлены ли параметры compilatorsList и extList
* если нет, то установить значения по умолчанию
 */
func (storage *ConfigparserContainer) DefaultPropogate() {
	if len(storage.compilatorsList) == 0 {
		storage.compilatorsList = []string{"gcc", "cc", "c++", "gcc++"}
	}
	if len(storage.extList) == 0 {
		storage.extList = []string{".c", ".cpp", ".cc", ".h", ".hh", ".hpp"}
	}
}

func (storage ConfigparserContainer) String() string {
	var tmp_result string = ""
	for key, value := range storage.addParameters {
		tmp_result += fmt.Sprintf("%s=%s\n", key, strings.Join(value, ", "))
	}
	var errstr string = "off"
	if storage.showerroroutput == true {
		errstr = "on"
	}
	return fmt.Sprintf("StdErr: %s\nCompilators: %s\nFile extentions: %s\nExtend:\n%s\nGlobalDefs: %s\nPlugins: %s\nOutput: %s\nTemplate: %s\nWrapStrings: %d\nIgnore: %s\nOutFile: %s\nHtmlTpl: %s\nFnd Path:%s\n Replacer: %s\n",
		errstr,
		strings.Join(storage.compilatorsList, ", "),
		strings.Join(storage.extList, ", "),
		tmp_result,
		strings.Join(storage.glovaldefs, ", "),
		strings.Join(storage.checkby, ", "),
		storage.output, storage.template, storage.wrapstrings,
		strings.Join(storage.ignore, ", "), storage.outputfile,
		storage.html_template, strings.Join(storage.fnd_path, ", "),
		storage.cc_replacer)
}

func (storage ConfigparserContainer) GetBash() []string {
	return storage.bash
}

func (storage ConfigparserContainer) GetStdError() bool {
	return storage.showerroroutput
}

func (storage ConfigparserContainer) GetFilesList() []string {
	return storage.extList
}

func (storage ConfigparserContainer) GetCompillatorsList() []string {
	return storage.compilatorsList
}

func (storage ConfigparserContainer) GetFileAddInfo(name string) string {
	var res_ret string = ""
	res, found := storage.addParameters[name]
	if found == true {
		res_ret = strings.Join(res, " ")
	}
	if len(storage.glovaldefs) > 0 {
		res_ret += " " + strings.Join(storage.glovaldefs, " ")
	}
	return res_ret
}

func (storage ConfigparserContainer) IsFileIgnored(name string) bool {
	for _, val := range storage.ignore {
		if val == name {
			return true
		}
	}
	return false
}

func (storage ConfigparserContainer) GetListOfPlugins() []string {
	return storage.checkby
}

func (storage ConfigparserContainer) GetReport() (string, string, int64, string) {
	return storage.output, storage.template, storage.wrapstrings, storage.outputfile
}

func (storage ConfigparserContainer) GetHtmlTemplate() string {
	return storage.html_template
}

func (storage *ConfigparserContainer) AddFndPath(path string) {
	storage.fnd_path = append(storage.fnd_path, path)
}

func (storage ConfigparserContainer) RetFndPath() []string {
	return storage.fnd_path
}

func (storage ConfigparserContainer) Replacer() bool {
	return storage.cc_replacer
}

func (storage *ConfigparserContainer) SetFilesList(list string) {
	storage.list_of_files = append(storage.list_of_files,
		SplitOwn(strings.Trim(list, " \n\t"))...)
}

func (storage ConfigparserContainer) CheckFile(file string) bool {
	for _, f_item := range storage.list_of_files {
		if f_item == "*" {
			return true
		} else {
			if strings.Contains(file, f_item) == true {
				return true
			}
		}
	}
	return false
}
