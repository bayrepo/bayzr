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

package reporter

import (
	"bufio"
	"configparser"
	"fmt"
	"io"
	"os"
	"resultanalyzer"
	"sort"
	"strconv"
	"strings"
)

const (
	NORMAL int = iota
	WARNING
	DANGER
	LINE
	ERRLINE
	ERRLINE_CONT
	BREAK
)

type ReporterContainerLineItemLink struct {
	Number int64
	Line   string
	FndStr bool
}

func (this *ReporterContainerLineItemLink) String() string {
	return fmt.Sprintf("[%d:%s:%s]", this.Number, this.Line, this.FndStr)
}

type ReporterContainerLineItem struct {
	Number int64
	Line   string
	Value  int
	Plugin string
	Link   []ReporterContainerLineItemLink
}

func (this *ReporterContainerLineItem) String() string {
	tp := " "
	switch {
	case this.Value == NORMAL:
		tp = "NORMAL"
	case this.Value == WARNING:
		tp = "WARNING"
	case this.Value == DANGER:
		tp = "DANGER"
	case this.Value == ERRLINE:
		tp = "*"
	case this.Value == LINE:
		tp = " "
	}
	buf := fmt.Sprintf("%d | (%s %s) %s", this.Number, tp, this.Plugin, this.Line)
	if len(this.Link) > 0 {
		for _, val := range this.Link {
			buf += fmt.Sprintf("\n%s", val.String())
		}
	}
	return buf
}

type ReporterContainerItem struct {
	Start_pos    int64
	End_pos      int64
	File         string
	List_strings []ReporterContainerLineItem
}

func (this ReporterContainerItem) String() string {
	buf := fmt.Sprintf("File: %s(%d-%d)\n", this.File, this.Start_pos, this.End_pos)
	for _, value := range this.List_strings {
		buf += fmt.Sprintf("%s\n", value.String())
	}
	return buf
}

type ReporterContainer struct {
	config           *configparser.ConfigparserContainer
	list             *[]*resultanalyzer.ResultAnalyzerConatiner
	err_list         []*ReporterContainerItem
	err_list_files   map[string][]*ReporterContainerItem
	err_list_names   []string
	list_of_commands map[string][]string
}

func (this *ReporterContainer) String() string {
	buf := ""
	for _, value := range this.err_list {
		buf += fmt.Sprintf("---------------------------------------\n%s\n", value.String())
	}
	return buf
}

func (this *ReporterContainer) rebuildListToMapList() {
	this.err_list_files = map[string][]*ReporterContainerItem{}
	for i := range this.err_list {
		if isThereDisableCommet(this.err_list[i]) == true {
			tmp, found := this.err_list_files[this.err_list[i].File]
			if found == false {
				tmp = []*ReporterContainerItem{}
			}
			tmp = append(tmp, this.err_list[i])
			this.err_list_files[this.err_list[i].File] = tmp
		}
		this.err_list_names = append(this.err_list_names, this.err_list[i].File)
	}
	this.err_list_names = configparser.RemoveDuplicate(this.err_list_names)
	sort.Strings(this.err_list_names)
}

func Make_ReporterContainer(conf *configparser.ConfigparserContainer, lst *[]*resultanalyzer.ResultAnalyzerConatiner, list_of_cmds map[string][]string) *ReporterContainer {
	return &ReporterContainer{conf, lst, []*ReporterContainerItem{}, map[string][]*ReporterContainerItem{}, []string{}, list_of_cmds}
}

func quickCommentAnalysis(fn string, line_in string) bool {
	ln := int64(0)
	i, err := strconv.ParseInt(strings.Trim(line_in, " \n\t"), 10, 64)
	if err == nil {
		ln = i
	}
	if ln > 0 {
		file_source, err := os.Open(fn)
		if err != nil {
			fmt.Printf("Can't open file %s\n", fn)
			return false
		}
		defer file_source.Close()
		reader := bufio.NewReader(file_source)
		counter := int64(1)

		for {
			line, err := reader.ReadString('\n')
			if (line != "") || (line == "" && err == nil) {
				if counter == (ln - 1) {
					if strings.Contains(line, "bzr:ignore") == false {
						return true
					} else {
						return false
					}
				}
			}
			if counter > ln {
				break
			}
			if err != nil {
				if err != io.EOF {
					fmt.Printf("Read %s file error: %s\n", fn, err)
				}
				break
			}
			counter++
		}
	}

	return true
}

func (this *ReporterContainer) saveAnalyzisInfoDirect(file_name string, report_template string) {
	file, err := os.Create(file_name)
	if err != nil {
		fmt.Printf("Can't create file %s: %s", file_name, err)
		os.Exit(1)
	}
	defer file.Close()
	old_plugin_name := ""
	for _, value := range *this.list {
		array_list, plugin_name := value.GetListOfErrors()

		is_fnd_file := false

		for _, message := range array_list {
			if quickCommentAnalysis(message.File, message.Line) == true {
				if this.config.CheckFile(message.File) == true {
					is_fnd_file = true
					break
				}
			}
		}

		if is_fnd_file == true {
			if old_plugin_name != plugin_name {
				file.WriteString(plugin_name + "\n")
				old_plugin_name = plugin_name
			}
			for _, message := range array_list {
				if this.config.CheckFile(message.File) == false {
					continue
				}
				if this.config.CheckFileLine(message.File, message.Line) == false {
					continue
				}

				if quickCommentAnalysis(message.File, message.Line) == true {
					res := report_template
					res = strings.Replace(res, "FILE", message.File, -1)
					res = strings.Replace(res, "LINE", message.Line, -1)
					res = strings.Replace(res, "SEV", message.Sev, -1)
					res = strings.Replace(res, "ID", message.Id, -1)
					res = strings.Replace(res, "MESSAGE", message.Message, -1)
					file.WriteString(res + "\n")
				}
			}
		}
	}
}

func isEndOfStringC(line string) bool {
	res := strings.Trim(line, " \n\t")
	if len(res) > 0 {
		ch := res[len(res)-1]
		if ch == ';' || ch == '}' || ch == '{' {
			return true
		}
	} else {
		return true
	}
	return false
}

type MixedListItem struct {
	Plugin_name string
	Postion     int64
	Item        *resultanalyzer.ResultAnalyzerConatinerItem
	level       int
}

type tMixedListItem []MixedListItem

func (this tMixedListItem) Len() int {
	return len(this)
}

func (this tMixedListItem) Less(i, j int) bool {
	return this[i].Postion < this[j].Postion
}

func (this tMixedListItem) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type MixedList struct {
	From int64
	To   int64
	List tMixedListItem
}

func removeFromSlise(list []MixedList, key int) []MixedList {
	if key == (len(list) - 1) {
		return list[:key]
	}
	if key >= len(list) {
		return list
	}
	return append(list[:key], list[key+1:]...)
}

func makeMixedArray(this *ReporterContainer, wrap int64) *[]MixedList {
	list := []MixedList{}
	for _, value := range *this.list {
		array_list, plugin_name := value.GetListOfErrors()
		for _, message := range array_list {
			position, err := strconv.ParseInt(message.Line, 10, 64)
			if err != nil {
				fmt.Printf("Incorrect line number %s in file %s for message %s\n", message.Line, message.File, message.Message)
				continue
			}
			if this.config.CheckFile(message.File) == false {
				continue
			}
			if this.config.CheckFileLine(message.File, message.Line) == false {
				continue
			}

			fnd := false
			for key := range list {
				for val := range list[key].List {
					if (position > (list[key].List[val].Postion - wrap)) &&
						(position < (list[key].List[val].Postion + wrap)) {
						if list[key].List[val].Item.File == message.File {
							item := MixedListItem{plugin_name, position, message, getDangerLevel(value.GetListPlugin().GetResultLevels(message.Sev))}
							list[key].List = append(list[key].List, item)
							fnd = true
							break
						}
					}
				}
				if fnd == true {
					break
				}
			}
			if fnd == false {
				item_obj := MixedListItem{plugin_name, position, message, getDangerLevel(value.GetListPlugin().GetResultLevels(message.Sev))}
				item := MixedList{0, 0, []MixedListItem{}}
				item.List = append(item.List, item_obj)
				list = append(list, item)
			}
		}

	}
	

	for key := range list {
		if len(list[key].List) == 1 {
			list[key].From = list[key].List[0].Postion - wrap
			list[key].To = list[key].List[0].Postion + wrap
		} else if len(list[key].List) > 1 {
			max := list[key].List[0].Postion
			min := list[key].List[0].Postion
			for val := range list[key].List {
				if max < list[key].List[val].Postion {
					max = list[key].List[val].Postion
				}
				if min > list[key].List[val].Postion {
					min = list[key].List[val].Postion
				}
			}
			list[key].From = min - wrap
			list[key].To = max + wrap
		}
	}

	//counter := 0

	for {
		//	counter += 1
		//	for _, i := range list {
		//		fmt.Printf("%d F:%d, T:%d\n", counter, i.From, i.To)
		//	}

		is_spare := false
		for key := range list {
			fnd := false
			for key2 := range list {
			    if list[key2].List[0].Item.File != list[key].List[0].Item.File {
			        continue
			    }
				if (key != key2) && ((list[key].To >= list[key2].From && list[key].From <= list[key2].From) ||
					(list[key2].To >= list[key].From && list[key2].From <= list[key].From) ||
					(list[key2].From >= list[key].From && list[key2].To <= list[key].To) ||
					(list[key].From >= list[key2].From && list[key].To <= list[key2].To)) {

					//				fmt.Printf("key %d F:%d, T:%d vs key %d F:%d, T:%d \n", key, list[key].From, list[key].To, key2, list[key2].From, list[key2].To)

					list[key].List = append(list[key].List, list[key2].List...)
					if list[key].From > list[key2].From {
						list[key].From = list[key2].From
					}
					if list[key].To < list[key2].To {
						list[key].To = list[key2].To
					}
					list = removeFromSlise(list, key2)
					fnd = true
					break
				}
			}
			if fnd == true {
				is_spare = true
				break
			}
		}
		if is_spare == false {
			break
		}
	}

	for key := range list {
		if len(list[key].List) > 1 {
			sort.Sort(list[key].List)
		}
	}

	/*	for _, i := range list {
		fmt.Printf("F:%d, T:%d\n", i.From, i.To)
		for _, j := range i.List {
			fmt.Printf("PLG %s Position %d %s\n", j.Plugin_name, j.Postion, j.Item.String())
		}
	}*/

	return &list
}

func makeListOfFiles(this *ReporterContainer) []string {
	list := map[string]int64{}
	for _, value := range *this.list {
		array_list, _ := value.GetListOfErrors()
		for _, message := range array_list {
			if _, found := list[message.File]; found == true {
				list[message.File]++
			} else {
				list[message.File] = 1
			}
		}
	}
	list_string := []string{}
	for key := range list {
		list_string = append(list_string, key)
	}
	sort.Strings(list_string)
	return list_string
}

func getDangerLevel(value string) int {
	if value == "MEDIUM" {
		return WARNING
	}
	if value == "HIGH" {
		return DANGER
	}
	return NORMAL
}

func (this *ReporterContainer) makeAddErrorInfo(lst []string, wrap_strings int64) []ReporterContainerLineItemLink {
	list := []ReporterContainerLineItemLink{}
	for _, val := range lst {
		tmp1 := strings.SplitN(val, ":", 2)
		if len(tmp1) > 1 {
			fn := strings.Trim(tmp1[0], " \n\t")
			ln := int64(0)
			i, err := strconv.ParseInt(strings.Trim(tmp1[1], " \n\t"), 10, 64)
			if err == nil {
				ln = i
			}
			if ln > 0 {
				file_source, err := os.Open(fn)
				if err != nil {
					fmt.Printf("Can't open file %s\n", fn)
					return []ReporterContainerLineItemLink{}
				}
				defer file_source.Close()
				reader := bufio.NewReader(file_source)
				counter := int64(1)

				for {
					line, err := reader.ReadString('\n')
					if (line != "") || (line == "" && err == nil) {
						if counter < (ln+wrap_strings) && counter > (ln-wrap_strings) {
							item := ReporterContainerLineItemLink{counter, line, false}
							if ln == counter {
								item.FndStr = true
							}
							list = append(list, item)
						}
					}
					if counter > ln+(wrap_strings) {
						break
					}
					if err != nil {
						if err != io.EOF {
							fmt.Printf("Read %s file error: %s\n", fn, err)
						}
						break
					}
					counter++
				}
				item := ReporterContainerLineItemLink{counter, "...", false}
				list = append(list, item)
			}
		}
	}
	return list
}

func (this *ReporterContainer) saveAnalyzisInfo(file_name string, wrap_strings int64) {
	file, err := os.Create(file_name)
	if err != nil {
		fmt.Printf("Can't create file %s: %s", file_name, err)
		os.Exit(1)
	}
	defer file.Close()
	list := makeMixedArray(this, wrap_strings)
	list_of_files := makeListOfFiles(this)
	for _, fn := range list_of_files {
		file_source, err := os.Open(fn)
		if err != nil {
			fmt.Printf("Can't open file %s\n", fn)
			continue
		}
		defer file_source.Close()
		reader := bufio.NewReader(file_source)
		counter := int64(1)

		mark_val := false

		var item_res *ReporterContainerItem = nil
		for {
			line, err := reader.ReadString('\n')
			if (line != "") || (line == "" && err == nil) {
				for _, items := range *list {
					if len(items.List) > 0 {
						if items.List[0].Item.File == fn && counter >= items.From && counter <= items.To {
							item_res = nil
							for i := range this.err_list {
								if this.err_list[i].File == fn && this.err_list[i].Start_pos == items.From && this.err_list[i].End_pos == items.To {
									item_res = this.err_list[i]
								}
							}
							if item_res == nil {
								item_res = &ReporterContainerItem{}
								item_res.Start_pos = items.From
								item_res.End_pos = items.To
								item_res.File = fn
								this.err_list = append(this.err_list, item_res)
							}
							fnd_str := false
							for _, items_val := range items.List {
								if items_val.Postion == counter {
									obj := ReporterContainerLineItem{}
									obj.Number = items_val.Postion
									obj.Plugin = items_val.Plugin_name
									obj.Value = items_val.level
									obj.Line = items_val.Item.Message
									obj.Link = this.makeAddErrorInfo(items_val.Item.Pos_link_file_line, wrap_strings)
									item_res.List_strings = append(item_res.List_strings, obj)
									fnd_str = true
								}
							}
							if fnd_str == true || mark_val == true {
								obj := ReporterContainerLineItem{}
								obj.Number = counter
								obj.Plugin = ""
								if mark_val == true {
									obj.Value = ERRLINE_CONT
								} else {
									obj.Value = ERRLINE
								}
								obj.Line = line
								item_res.List_strings = append(item_res.List_strings, obj)
								if isEndOfStringC(line) == true {
									mark_val = false
								} else {
									mark_val = true
								}
							} else {
								obj := ReporterContainerLineItem{}
								obj.Number = counter
								obj.Plugin = ""
								obj.Value = LINE
								obj.Line = line
								item_res.List_strings = append(item_res.List_strings, obj)
							}
						}
					}
				}
			}
			if err != nil {
				if err != io.EOF {
					fmt.Printf("Read %s file error: %s\n", fn, err)
				}
				break
			}
			counter++
		}
	}
}

func isThereDisableCommet(cont *ReporterContainerItem) bool {
	err_nums := 0
	for key, value := range cont.List_strings {
		if value.Value == ERRLINE {
			if key > 0 {
				fnd_bzr := false
				for key_i := key - 1; key_i >= 0; key_i-- {
					if cont.List_strings[key_i].Value == LINE && strings.Contains(cont.List_strings[key_i].Line, "bzr:ignore") == true {
						fnd_bzr = true
						break
					}
				}
				if fnd_bzr == true {
					cont.List_strings[key].Value = LINE
					for key_i := key + 1; key_i < len(cont.List_strings); key_i++ {
						if cont.List_strings[key_i].Value == ERRLINE_CONT {
							cont.List_strings[key_i].Value = LINE
						} else {
							break
						}

					}
				} else {
					err_nums++
				}
			} else {
				err_nums++
			}
		}
	}
	return err_nums != 0
}

func (this *ReporterContainer) saveAnalyzisInfoDirectTxt(file_name string) {
	file, err := os.Create(file_name)
	if err != nil {
		fmt.Printf("Can't create file %s: %s", file_name, err)
		os.Exit(1)
	}
	defer file.Close()
	for plugin_name, plugin_commands := range this.list_of_commands {
		file.WriteString("Commands for plugin " + plugin_name + "\n")
		file.WriteString("--------------------------------------------------------------------------\n")
		for _, cmd_string := range plugin_commands {
			file.WriteString(cmd_string + "\n")
		}
		file.WriteString("\n\n")
	}
	this.rebuildListToMapList()
	for _, file_nm := range this.err_list_names {
		file.WriteString(file_nm + "\n")
		file.WriteString("\n")
		for _, list := range this.err_list_files[file_nm] {
			for _, value := range list.List_strings {
				buf := ""
				if value.Value <= DANGER {
					buf = fmt.Sprintf("//**DETECT** %s:      %s\n", value.Plugin, value.Line)
				} else {
					buf = fmt.Sprintf("%10d|  %s", value.Number, value.Line)
				}
				file.WriteString(buf)
			}
			file.WriteString(".......\n")
		}
	}
}

func (this *ReporterContainer) getAnalyzisInfoDirectHtml(file_name string) string {
	this.rebuildListToMapList()
	path_list := this.config.RetFndPath()
	tpl_file := this.config.GetHtmlTemplate()
	if tpl_file != "" {
		path := ""
		for _, val := range path_list {
			if _, err := os.Stat(val + tpl_file); err == nil {
				path = val + tpl_file
			}
		}
		if path != "" {
			return path
		} else {
			fmt.Println("Template html file does not exist")
			os.Exit(1)
		}
	} else {
		fmt.Println("Empty template html file")
		os.Exit(1)
	}
	return ""
}

func (this *ReporterContainer) CreateReport() (string, bool) {
	report_type, report_template, report_number, report_file := this.config.GetReport()
	switch {
	case report_type == "custom":
		this.saveAnalyzisInfoDirect(report_file, report_template)
		return "", false
	case report_type == "html" || report_type == "txt":
		this.saveAnalyzisInfo(report_file, report_number)
		if report_type == "txt" {
			this.saveAnalyzisInfoDirectTxt(report_file)
			return "", false
		} else {
			return this.getAnalyzisInfoDirectHtml(report_file), true
		}
	}
	return "", false
}

func (this *ReporterContainer) GetErrList() *map[string][]*ReporterContainerItem {
	return &this.err_list_files
}

func (this *ReporterContainer) GetCmdList() map[string][]string {
	return this.list_of_commands
}
