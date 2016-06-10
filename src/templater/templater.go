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

package templater

import (
	"configparser"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"reporter"
	"time"
	"sort"
)

type ListOfErrorsShort struct {
	Position      int64
	Error_message string
}

type ListOfErrorLong struct {
	Critical int64
	Warning  int64
	Normal   int64
	List     []reporter.ReporterContainerLineItem
}

type PreparedToOutput struct {
	ReportName     string
	ListOfCheckers map[string]int64
	ListOfShort    map[string][]ListOfErrorsShort
	ListOfLong     map[string]ListOfErrorLong
	ListOfFiles    []string
	Consts         []int
}

func MakeTemplater() *PreparedToOutput {
	return &PreparedToOutput{"", map[string]int64{}, map[string][]ListOfErrorsShort{},
		map[string]ListOfErrorLong{}, []string{}, []int{}}
}

func (this *PreparedToOutput) PropogateData(obj *reporter.ReporterContainer, path string, cfg *configparser.ConfigparserContainer) {
	t := time.Now()
	this.ReportName = t.Format("Report from 2006-01-02 15:04:05")
	list := obj.GetErrList()
	for key, value := range *list {
		this.ListOfFiles = append(this.ListOfFiles, key)
		for ind, item_cont := range value {
			for _, item := range item_cont.List_strings {
				val, fnd := this.ListOfLong[key]
				if fnd == false {
					val = ListOfErrorLong{}
				}
				val.List = append(val.List, item)
				fnd_err := false
				if item.Value == reporter.NORMAL {
					val.Normal++
					fnd_err = true
				}
				if item.Value == reporter.WARNING {
					val.Warning++
					fnd_err = true
				}
				if item.Value == reporter.DANGER {
					val.Critical++
					fnd_err = true
				}
				if fnd_err == true {
					this.ListOfCheckers[item.Plugin]++
					this.ListOfShort[key] = append(this.ListOfShort[key],
						ListOfErrorsShort{item.Number, item.Line})
				}
				this.ListOfLong[key] = val
			}
			val, fnd := this.ListOfLong[key]
			if fnd == true {
				if len(val.List) > 0 && len(value)-1 != ind {
					val.List = append(val.List, reporter.ReporterContainerLineItem{0, "", reporter.BREAK, "", []reporter.ReporterContainerLineItemLink{}})
					this.ListOfLong[key] = val
				}
			}
		}
	}
	this.ListOfFiles = configparser.RemoveDuplicate(this.ListOfFiles)
	sort.Strings(this.ListOfFiles)
	this.Consts = []int{reporter.NORMAL, reporter.WARNING, reporter.DANGER, reporter.LINE, reporter.ERRLINE, reporter.ERRLINE_CONT, reporter.BREAK}
	tpl := template.New(filepath.Base(path))
	tpl, err := tpl.ParseFiles(path)
	if err != nil {
		fmt.Printf("Template %s parsing error %s\n", path, err)
		os.Exit(1)
	}

	_, _, _, output_file := cfg.GetReport()

	file, err := os.Create(output_file)
	if err != nil {
		fmt.Printf("Can't create file %s: %s", output_file, err)
		os.Exit(1)
	}
	defer file.Close()

	err = tpl.Execute(file, this)
	if err != nil {
		fmt.Printf("Template %s parsing error %s\n", path, err)
		os.Exit(1)
	}
}
