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

package outputanalyzer

import (
	"bufio"
	"configparser"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

var compillatorCached string = ""
var dir_name string = "."

type OutPutAnalyzerContainerItem struct {
	IncludesList     []string
	DefList          []string
	Dir              string
	File             string
	Raw              []string
	is_first_gcc     bool
	number_of_params int
}

func (item OutPutAnalyzerContainerItem) String() string {
	return fmt.Sprintf("[\nIncludesList: %s\nDefList: %s\nDir: %s\nFile: %s\n]", strings.Join(item.IncludesList, ", "),
		strings.Join(item.DefList, ", "), item.Dir, item.File)
}

type OutPutAnalyzerContainer struct {
	filesList map[string]*OutPutAnalyzerContainerItem
	config    *configparser.ConfigparserContainer
}

func chekExistsNext(slice []string, index int) bool {
	if index >= len(slice) || index < 0 {
		return false
	}
	return true
}

func checkCompillator(lines []string, storage *OutPutAnalyzerContainer) (bool, bool) {
	if len(lines) <= 0 {
		return false, false
	}
	if compillatorCached != "" {
		if strings.Contains(lines[0], compillatorCached) == true {
			return true, true
		}

		for _, compil := range lines {
			if strings.Contains(compil, compillatorCached) == true {
				return true, false
			}
		}
	}

	for _, value := range storage.config.GetCompillatorsList() {
		if strings.Contains(lines[0], value) == true {
			compillatorCached = lines[0]
			return true, true
		}
	}

	for _, compil := range lines {
		for _, value := range storage.config.GetCompillatorsList() {
			if strings.Contains(compil, value) == true {
				compillatorCached = compil
				return true, false
			}
		}
	}
	return false, false
}

func findFileNameInSlice(file_name string, slice []string) int {
	for i := range slice {
		if slice[i] == file_name {
			return i
		}
	}
	return -1
}

func checkFile(lines []string, storage *OutPutAnalyzerContainer, found_files []string) string {
	if len(lines) <= 0 {
		return ""
	}
	for _, param := range lines {
		utf8_param := []rune(param)

		if len(utf8_param) > 2 && utf8_param[:2][0] == rune('-') && utf8_param[:2][1] == rune('c') {
			file := param[2:]
			for _, ext := range storage.config.GetFilesList() {
				utf8_ext := []rune(ext)
				if utf8_ext[0] != rune('.') {
					ext = "." + ext
				}
				if strings.Contains(file, ext) == true {
					if findFileNameInSlice(file, found_files) == -1 {
						return file
					}
					break
				}
			}
		}
	}
	for _, param_no_c := range lines {
		for _, ext := range storage.config.GetFilesList() {
			utf8_ext := []rune(ext)
			if utf8_ext[0] != rune('.') {
				ext = "." + ext
			}
			if len(ext) > len(param_no_c) {
				continue
			}
			file_ext := param_no_c[len(param_no_c)-len(ext):]
			if strings.Compare(file_ext, ext) == 0 {
				if findFileNameInSlice(param_no_c, found_files) == -1 {
					utf8_param := []rune(param_no_c)

					if len(utf8_param) > 2 && utf8_param[:2][0] == rune('-') && utf8_param[:2][1] == rune('c') {
						continue
					} else if len(utf8_param) > 2 && utf8_param[:2][0] == rune('-') && utf8_param[:2][1] == rune('o') {
						continue
					} else {
						return param_no_c
					}
				}
			}
		}
	}
	return ""
}

func getDefinitionList(lines []string) []string {
	var def_list []string
	if len(lines) <= 0 {
		return def_list
	}
	var skip int = 0
	for i, param := range lines {
		if skip == 1 {
			skip = 0
			continue
		}
		utf8_param := []rune(param)
		if len(utf8_param) > 2 && utf8_param[:2][0] == rune('-') && (utf8_param[:2][1] == rune('D') || utf8_param[:2][1] == rune('I')) {
			file := param[2:]
			if file == "" {
				skip = 1
				if chekExistsNext(lines, i+1) == true {
					def_list = append(def_list, param+lines[i+1])
				}
			} else {
				def_list = append(def_list, param)
			}

		}
	}
	return def_list
}

func getDefOnly(lines []string) []string {
	var def_list []string
	for _, param := range lines {
		utf8_param := []rune(param)
		if len(utf8_param) > 2 && utf8_param[:2][0] == rune('-') && utf8_param[:2][1] == rune('D') {
			tmp_param := param[2:]
			if tmp_param != "" {
				def_list = append(def_list, tmp_param)
			}
		}
	}
	return def_list
}

func getIncOnly(lines []string) []string {
	var def_list []string
	for _, param := range lines {
		utf8_param := []rune(param)
		if len(utf8_param) > 2 && utf8_param[:2][0] == rune('-') && utf8_param[:2][1] == rune('I') {
			tmp_param := param[2:]
			if tmp_param != "" {
				def_list = append(def_list, tmp_param)
			}
		}
	}
	return def_list
}

func getRawWithoutC(lines []string) []string {
	var _list []string
	for _, param := range lines {
		utf8_param := []rune(param)
		if len(utf8_param) > 2 && utf8_param[:2][0] == rune('-') && utf8_param[:2][1] == rune('c') {
			continue
		}
		_list = append(_list, param)
	}
	return _list
}

func checkForDirSwitching(lines []string) (string, bool) {
	result := strings.Join(lines, " ")
	section_reg := regexp.MustCompile("Entering directory `([\\s\\S]+)'$")
	if matches := section_reg.FindStringSubmatch(result); matches != nil {
		current_path, _ := os.Getwd()
		f_name := matches[1]
		if path.IsAbs(f_name) == false {
			f_name = current_path + "/" + f_name
		}
		if _, err := os.Stat(f_name); os.IsNotExist(err) {
			f_name = current_path + "/" + f_name
			if _, err := os.Stat(f_name); os.IsNotExist(err) {
				return ".", false
			}
		}
		return strings.Trim(f_name, " \n\t"), true
	}
	return ".", false
}

func GetFilesNames(args []string, storage *OutPutAnalyzerContainer) ([]string, bool, bool) {
	var file_name string
	var files_name []string = []string{}
	is_comp, is_fst_compil := checkCompillator(args, storage)
	if is_comp == true {
		for {
			if file_name = checkFile(args, storage, files_name); file_name != "" {
				files_name = append(files_name, file_name)
			} else {
				break
			}
		}
	}
	return files_name, is_comp, is_fst_compil
}

func makeStringAnalysis(line string, storage *OutPutAnalyzerContainer) {
	raw_string := configparser.SplitOwn(strings.Trim(line, " \n\t"))
	new_string := []string{}
	var skip int = 0
	var files_name []string
	for i := range raw_string {
		if skip == 1 {
			skip = 0
			continue
		}
		if len(raw_string[i]) > 0 {
			var utf8_value []rune
			utf8_value = []rune(raw_string[i])
			if utf8_value[0] == rune('-') && len(utf8_value) == 2 && chekExistsNext(raw_string, i+1) == true {
				new_string = append(new_string, raw_string[i]+raw_string[i+1])
				skip = 1
			} else {
				new_string = append(new_string, raw_string[i])
			}
		}
	}

	files_name, is_comp, is_fst_compil := GetFilesNames(new_string, storage)

	if dir_name_tmp, fnd_dir := checkForDirSwitching(new_string); fnd_dir == true {
		dir_name = dir_name_tmp
	}

	if len(files_name) > 0 {
		if def_values := getDefinitionList(new_string); len(def_values) > 0 {
			for _, f := range files_name {
				if _, found := storage.filesList[f]; found == false {
					storage.filesList[f] = &OutPutAnalyzerContainerItem{}
					storage.filesList[f].DefList = []string{}
					storage.filesList[f].IncludesList = []string{}
					storage.filesList[f].is_first_gcc = is_comp
					storage.filesList[f].number_of_params = len(new_string)
					storage.filesList[f].Raw = []string{}
				} else {
					if is_fst_compil == true {
						storage.filesList[f] = &OutPutAnalyzerContainerItem{}
						storage.filesList[f].DefList = []string{}
						storage.filesList[f].IncludesList = []string{}
						storage.filesList[f].is_first_gcc = is_comp
						storage.filesList[f].number_of_params = len(new_string)
						storage.filesList[f].Raw = []string{}
					}
				}
				storage.filesList[f].Raw = getRawWithoutC(new_string)
				storage.filesList[f].DefList = append(storage.filesList[f].DefList, getDefOnly(def_values)...)
				storage.filesList[f].IncludesList = append(storage.filesList[f].IncludesList, getIncOnly(def_values)...)
				if dir_name == "." {
					if path.Dir(f) == "." {
						storage.filesList[f].Dir = dir_name
					} else {
						storage.filesList[f].Dir = path.Dir(f)
					}
				} else {
					storage.filesList[f].Dir = dir_name
				}
				storage.filesList[f].File = path.Base(f)
			}
		} else {
			for _, f := range files_name {
				if _, found := storage.filesList[f]; found == false {
					storage.filesList[f] = &OutPutAnalyzerContainerItem{}
					storage.filesList[f].DefList = []string{}
					storage.filesList[f].IncludesList = []string{}
					storage.filesList[f].is_first_gcc = is_comp
					storage.filesList[f].number_of_params = len(new_string)
					storage.filesList[f].Raw = getRawWithoutC(new_string)
				} else {
					if is_fst_compil == true {
						storage.filesList[f] = &OutPutAnalyzerContainerItem{}
						storage.filesList[f].DefList = []string{}
						storage.filesList[f].IncludesList = []string{}
						storage.filesList[f].is_first_gcc = is_comp
						storage.filesList[f].number_of_params = len(new_string)
						storage.filesList[f].Raw = getRawWithoutC(new_string)
					}
				}
				if path.Dir(f) == "." {
					storage.filesList[f].Dir = dir_name
				} else {
					storage.filesList[f].Dir = path.Dir(f)
				}
				storage.filesList[f].File = path.Base(f)
			}
		}
	}

}

/*
* Конструктор объекта по умолчанию
 */
func CreateDefaultAnalyzer(config *configparser.ConfigparserContainer) *OutPutAnalyzerContainer {
	return &OutPutAnalyzerContainer{map[string]*OutPutAnalyzerContainerItem{}, config}
}

/*
* Процедура запуска компилятора, все что относится к команде компиллятора
* должно идти за параметром: cmd или -cmd или --cmd
 */
func (storage *OutPutAnalyzerContainer) ExcecuteComplilationProcess(args []string) error {
	var clean_cmd string
	var cmd *exec.Cmd
	start_cmd := args
	if len(start_cmd) > 0 {
		clean_cmd = strings.Join(start_cmd, " ")
		start_cmd = append(storage.config.GetBash(), strings.Join(start_cmd, " "))
		if len(start_cmd) == 1 {
			cmd = exec.Command(start_cmd[0])
		} else {
			cmd = exec.Command(start_cmd[0], start_cmd[1:]...)
		}
		var err error
		var stdout io.ReadCloser
		var stderr io.ReadCloser
		if stdout, err = cmd.StdoutPipe(); err != nil {
			return fmt.Errorf("Analyzer error: open stdout pipe error %s\n", err)
		}
		if stderr, err = cmd.StderrPipe(); err != nil {
			return fmt.Errorf("Analyzer error: open stdout pipe error %s\n", err)
		}
		scanner_out := bufio.NewScanner(stdout)
		scanner_err := bufio.NewScanner(stderr)

		if storage.config.GetStdError() == true {
			go func() {
				for scanner_err.Scan() {
					fmt.Println(scanner_err.Text())
				}
				if err = scanner_err.Err(); err != nil {
					fmt.Printf("Analyzer error: stderr scanner error %s\n", err)
				}
			}()
		}
		makeStringAnalysis(clean_cmd, storage)

		if err = cmd.Start(); err != nil {
			return fmt.Errorf("Analyzer error: start command error %s\n", err)
		}

		for scanner_out.Scan() {
			fmt.Println(scanner_out.Text())
			makeStringAnalysis(scanner_out.Text(), storage)
		}
		if storage.config.GetStdError() == true {
			if err = scanner_out.Err(); err != nil {
				fmt.Printf("Analyzer error: stdout scanner error %s\n", err)
			}
		}

		if err = cmd.Wait(); err != nil {
			return fmt.Errorf("Analyzer error: wait command error %s\n", err)
		}
		return nil
	}
	return fmt.Errorf("Analyzer error: empty compile command\n")
}

func (storage OutPutAnalyzerContainer) String() string {
	var tmp_result string = ""
	for k := range storage.filesList {
		tmp_result += "[" + k + "]=" + storage.filesList[k].String() + "\n"
	}
	return tmp_result
}

func (storage OutPutAnalyzerContainer) GetParsedFilesList() *map[string]*OutPutAnalyzerContainerItem {
	return &storage.filesList
}

func (storage *OutPutAnalyzerContainer) GetStringOfCmd(args []string) []string {
	var start_cmd []string
	args_len := len(args)
	for i, args_value := range args {
		if args_value == "cmd" || args_value == "-cmd" || args_value == "--cmd" {
			if args_len-1 != i {
				start_cmd = args[i+1:]
			} else {
				start_cmd = []string{}
			}
			break
		}
	}
	return start_cmd
}

func TransformFileName(fileName string, dir string) string {
	if path.IsAbs(fileName) == true {
		if _, err := os.Stat(fileName); err == nil {
			return fileName
		}
	}
	f_n := dir + "/" + fileName
	f_n = path.Clean(f_n)
	if _, err := os.Stat(f_n); os.IsNotExist(err) {
		return ""
	}
	return f_n
}
