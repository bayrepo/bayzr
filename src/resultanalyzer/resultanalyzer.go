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

package resultanalyzer

import (
	"bufio"
	"checker"
	"configparser"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type ResultAnalyzerConatinerItem struct {
	File               string
	Line               string
	Sev                string
	Id                 string
	Message            string
	Pos_link_file_line []string
}

type ResultAnalyzerConatiner struct {
	chk          *checker.PluginInfoDataContainer
	result_array []*ResultAnalyzerConatinerItem
	file_name    string
	home_dir     string
	garbade_list []string
	mtx          sync.Mutex
	config       configparser.ConfigparserContainer
}

func (this ResultAnalyzerConatinerItem) String() string {
	return fmt.Sprintf("[\nFile:%s\nLine:%s\nSeverity:%s\nId:%s\nMessage:%s\nList of links: %s\n]\n",
		this.File, this.Line, this.Sev, this.Id, this.Message, strings.Join(this.Pos_link_file_line, ", "))
}

func Make_ResultAnalyzerConatinerItem() *ResultAnalyzerConatinerItem {
	return &ResultAnalyzerConatinerItem{}
}

func Make_ResultAnalyzerConatiner(fname string, chk_in *checker.PluginInfoDataContainer, home_d string, cfg configparser.ConfigparserContainer) *ResultAnalyzerConatiner {
	return &ResultAnalyzerConatiner{chk_in, []*ResultAnalyzerConatinerItem{}, fname, home_d, []string{}, sync.Mutex{}, cfg}
}

var tmp_array []*ResultAnalyzerConatinerItem = []*ResultAnalyzerConatinerItem{}
var tmp_descr string = ""
var fnd_end = false
var pos_link_ln = ""
var pos_link_fl = ""

func isWtSpaceBg(line string, numb int) bool {
	if numb > 0 {
		tab := strings.Repeat(" ", numb)
		if strings.HasPrefix(line, tab) == true {
			return true
		}
		return false
	}
	return false
}

func (this *ResultAnalyzerConatiner) makeGarbageList(line string) {
	if this.chk.GetClean() == "" {
		return
	}
	if currentPath, err := os.Getwd(); err == nil {
		section_reg := regexp.MustCompile(this.chk.GetClean())
		if matches := section_reg.FindStringSubmatch(line); matches != nil {
			file_object := strings.Trim(matches[1], " \n\t")
			f_to_delete := currentPath + "/" + file_object
			if _, err := os.Stat(f_to_delete); err == nil {
				if f_to_delete != "/" {
					this.mtx.Lock()
					defer this.mtx.Unlock()
					this.garbade_list = append(this.garbade_list, f_to_delete)
				}
			}
		}
	}
}

var last_File_Name string = ""

func (this *ResultAnalyzerConatiner) makeStringAnalysis(line string) {
	if len(line) > 0 {
		fields := this.chk.GetResult()
		result := strings.SplitN(line, string(this.chk.GetDelim()), len(fields))
		if len(result) == len(fields) && isWtSpaceBg(line, this.chk.GetSpaces()) == false {
			if fnd_end == true {
				for key := range tmp_array {
				    if tmp_array[key].Id=="" {
				        tmp_array[key].Id = tmp_array[key].Message
				    }
					tmp_array[key].Message += "-->" + strings.Trim(tmp_descr, " \r\n")
					tmp_array[key].Pos_link_file_line = append(tmp_array[key].Pos_link_file_line, fmt.Sprintf("%s:%s", pos_link_fl, pos_link_ln))
				}
				this.result_array = append(this.result_array, tmp_array...)
				tmp_descr = ""
				pos_link_ln = ""
				pos_link_fl = ""
				tmp_array = []*ResultAnalyzerConatinerItem{}
			}
			fnd_end = false
			item := Make_ResultAnalyzerConatinerItem()
			for i := range fields {
				switch fields[i] {
				case ":FILE":
					item.File = strings.Trim(result[i], " \n\t")
				case ":LINE":
					item.Line = strings.Trim(result[i], " \n\t")
				case ":SEV":
					item.Sev = strings.Trim(result[i], " \n\t")
				case ":ID":
					item.Id = strings.Trim(result[i], " \n\t")
				case ":MESSAGE":
					item.Message = result[i]
				}
			}
			if item.File == "" && last_File_Name != "" {
				item.File = last_File_Name
			}
			if item.File != "" && item.File != last_File_Name {
				last_File_Name = item.File
			}

			if item.Message != "" || item.Id != "" || item.Sev != "" {
				path_to_file_f_, err := filepath.Abs(item.File)
				if err == nil {
					item.File = path_to_file_f_
				}
				if strings.HasPrefix(path_to_file_f_, this.home_dir) == true && err == nil && this.config.CheckFile(item.File) == true {
					if this.chk.GetType() == "1" {
						this.result_array = append(this.result_array, item)
					} else if this.chk.GetType() == "2" {
						tmp_array = append(tmp_array, item)
					}
				}
			}

		} else {
			if this.chk.GetType() == "2" {
				if len(result) == len(fields) && isWtSpaceBg(line, this.chk.GetSpaces()) == true {
					for i := range fields {
						switch fields[i] {
						case ":FILE":
							pos_link_fl = strings.Trim(result[i], " \n\t")
						case ":LINE":
							pos_link_ln = result[i]
						}
					}
				}
				res := strings.Trim(line, " \n\t")
				if res != "" {
					tmp_descr += " " + res
					fnd_end = true
				}
			}
		}
	}
}

func (this *ResultAnalyzerConatiner) ParseResultOfCommand(cmd_in string, config *configparser.ConfigparserContainer) error {
	var start_cmd []string
	var cmd *exec.Cmd
	if len(cmd_in) > 0 {
		clean_cmd := cmd_in
		start_cmd = append(config.GetBash(), clean_cmd)
		if len(start_cmd) == 1 {
			cmd = exec.Command(start_cmd[0])
		} else {
			cmd = exec.Command(start_cmd[0], start_cmd[1:]...)
		}
		var err error
		var stdout io.ReadCloser
		var stderr io.ReadCloser
		if stdout, err = cmd.StdoutPipe(); err != nil {
			return fmt.Errorf("Result analyzer error: open stdout pipe error %s\n", err)
		}
		if stderr, err = cmd.StderrPipe(); err != nil {
			return fmt.Errorf("Result analyzer error: open stdout pipe error %s\n", err)
		}

		var scanner_out *bufio.Scanner
		var scanner_err *bufio.Scanner

		if this.chk.GetStream() == "stdout" {
			scanner_out = bufio.NewScanner(stdout)
			scanner_err = bufio.NewScanner(stderr)
		} else {
			scanner_err = bufio.NewScanner(stdout)
			scanner_out = bufio.NewScanner(stderr)
		}

		if config.GetStdError() == true {
			go func() {
				for scanner_err.Scan() {
					fmt.Println(scanner_err.Text())
					this.makeGarbageList(scanner_err.Text())
				}
			}()
		}

		if err = cmd.Start(); err != nil {
			return fmt.Errorf("Result analyzer error: start command error %s\n", err)
		}

		for scanner_out.Scan() {
			fmt.Println(scanner_out.Text())
			this.makeGarbageList(scanner_out.Text())
			this.makeStringAnalysis(scanner_out.Text())
		}
		if config.GetStdError() == true {
			if err = scanner_out.Err(); err != nil {
				fmt.Printf("Result analyzer error: stdout scanner error %s\n", err)
			}
		}

		if tmp_descr != "" && len(tmp_array) > 0 {
			for key := range tmp_array {
				tmp_array[key].Message += "-->" + strings.Trim(tmp_descr, " \r\n")
			}
			this.result_array = append(this.result_array, tmp_array...)
		}
		tmp_descr = ""
		tmp_array = []*ResultAnalyzerConatinerItem{}

		if err = cmd.Wait(); err != nil {
			if this.chk.GetResultStop() == false {
				return fmt.Errorf("Result analyzer error: wait command error %s\n", err)
			} else {
				fmt.Printf("Result analyzer error: wait command error %s\n", err)
			}
		}
		return nil
	}
	return fmt.Errorf("Result analyzer error: empty check command\n")
}

func (this *ResultAnalyzerConatiner) String() string {
	buf := ""
	for _, value := range this.result_array {
		buf += fmt.Sprintf("%s", value)
	}
	return buf
}

func (this *ResultAnalyzerConatiner) GetListOfErrors() ([]*ResultAnalyzerConatinerItem, string) {
	return this.result_array, this.chk.GetName()
}

func (this *ResultAnalyzerConatiner) GetListPlugin() *checker.PluginInfoDataContainer {
	return this.chk
}

func (this *ResultAnalyzerConatiner) MakePreCommand(config *configparser.ConfigparserContainer) error {
	var start_cmd []string
	var cmd *exec.Cmd
	res := strings.Trim(this.chk.GetPreCommand(), " \n\t")
	if res != "" {
		clean_cmd := res
		start_cmd = append(config.GetBash(), clean_cmd)
		if len(start_cmd) == 1 {
			cmd = exec.Command(start_cmd[0])
		} else {
			cmd = exec.Command(start_cmd[0], start_cmd[1:]...)
		}

		var err error
		var stdout io.ReadCloser
		var stderr io.ReadCloser

		var scanner_out *bufio.Scanner
		var scanner_err *bufio.Scanner

		if stdout, err = cmd.StdoutPipe(); err != nil {
			return fmt.Errorf("Result analyzer error: open stdout pipe error %s\n", err)
		}
		if stderr, err = cmd.StderrPipe(); err != nil {
			return fmt.Errorf("Result analyzer error: open stdout pipe error %s\n", err)
		}

		scanner_out = bufio.NewScanner(stdout)
		scanner_err = bufio.NewScanner(stderr)

		go func() {
			for scanner_err.Scan() {
				fmt.Println(scanner_err.Text())
			}
		}()

		if err = cmd.Start(); err != nil {
			return fmt.Errorf("Result analyzer error: start command error %s\n", err)
		}

		for scanner_out.Scan() {
			fmt.Println(scanner_out.Text())
		}
		if config.GetStdError() == true {
			if err = scanner_out.Err(); err != nil {
				fmt.Printf("Result analyzer error: stdout scanner error %s\n", err)
			}
		}

		if err = cmd.Wait(); err != nil {
			fmt.Printf("Result analyzer error: wait command error %s\n", err)
		}

	}
	return nil
}

func (this *ResultAnalyzerConatiner) RemoveGarbage(config *configparser.ConfigparserContainer) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	for _, val := range this.garbade_list {
		err := os.RemoveAll(val)
		if err != nil {
			fmt.Println("Can't delete Granbade file or directory: ", err)
		}
	}
}
