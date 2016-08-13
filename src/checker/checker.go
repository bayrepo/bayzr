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

package checker

import (
	"bufio"
	"bytes"
	"configparser"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	WS
	IDENT
	BEGIN
	END
)

type preScanned struct {
	Tok Token
	Val string
}

var eof = rune(0)
var inner_tag bool = false

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	if ch == '"' {
		inner_tag = !inner_tag
	}
	return ((inner_tag == false && (ch != '{' && ch != '}')) || inner_tag == true) && ch != eof
}

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func (s *Scanner) Scan() (tok Token, lit string) {
	ch := s.read()

	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	}

	switch ch {
	case eof:
		return EOF, ""
	case '}':
		return END, string(ch)
	case '{':
		return BEGIN, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (s *Scanner) scanIdent() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return IDENT, buf.String()
}

type PluginInfoDataContainerVars struct {
	Def     string
	Def_cmd string
	Cmd     string
}

type PluginInfoDataContainerLang struct {
	Lang string
	Ext  []string
}

func (obj PluginInfoDataContainerLang) String() string {
	return fmt.Sprintf("[ lang: %s  ext: %s ]", obj.Lang, obj.Ext)
}

type PluginInfoDataContainer struct {
	name          string
	description   string
	id            string
	pType         string
	options       []string
	lang          []PluginInfoDataContainerLang
	lang_default  string
	result        []string
	result_delim  rune
	defs          []string
	includes      []string
	filename      string
	cmd           string
	arch          []PluginInfoDataContainerLang
	arch_default  string
	plat          []PluginInfoDataContainerLang
	plat_default  string
	variables     map[string]PluginInfoDataContainerVars
	stream        string
	subdescr      bool
	result_levels map[string][]string
	dontstop      bool
	autoinclude   string
	result_spaces int
	result_clean  string
	precommand    string
	fresh         bool
	compose       bool
	extraoptions  map[string][]string

	//tmp no printable
	cmd_for_fresh          []string
	accumulated_extra_opts []string
}

var listOfPlugins []*PluginInfoDataContainer

func pluginInfoDataContainer_make() *PluginInfoDataContainer {
	return &PluginInfoDataContainer{}
}

func check_prescan_tokens(preScan []preScanned) bool {
	beg := 0
	end := 0
	for _, item := range preScan {
		if item.Tok == BEGIN {
			beg += 1
		}
		if item.Tok == END {
			end += 1
		}
	}
	return (beg == end)
}

func findOneValue(preScan []preScanned, name string) (string, error) {
	var fnd_tag bool = false
	var fnd_open bool = false
	var result string = ""
	var fnd_val = false
	for _, item := range preScan {
		if fnd_tag == false {
			val := strings.Trim(item.Val, " \n\t")
			if item.Tok == IDENT && val == name {
				fnd_tag = true
			}
		} else {
			if fnd_open == false {
				if item.Tok == BEGIN {
					fnd_open = true
				}
			} else {
				if item.Tok == BEGIN {
					return "", fmt.Errorf("Literal { inside context %s section\n", name)
				} else if item.Tok == IDENT {
					if fnd_val == false {
						result = strings.Trim(item.Val, " \n\t")
						fnd_val = true
					} else {
						return "", fmt.Errorf("Find two bodies into %s section\n", name)
					}
				} else if item.Tok == END {
					fnd_open = false
					break
				}
			}
		}
	}
	if fnd_open == true {
		return "", fmt.Errorf("Not closed %s section\n", name)
	}
	if fnd_val == true {
		return strings.Replace(result, "\n", " ", -1), nil
	}
	return "", fmt.Errorf("Can't find %s tag\n", name)
}

func findSliceValue(preScan []preScanned, name string) ([]string, error) {
	var fnd_tag bool = false
	var fnd_open bool = false
	var result []string = []string{}
	var fnd_val = false
	for _, item := range preScan {
		if fnd_tag == false {
			val := strings.Trim(item.Val, " \n\t")
			if item.Tok == IDENT && val == name {
				fnd_tag = true
			}
		} else {
			if fnd_open == false {
				if item.Tok == BEGIN {
					fnd_open = true
				}
			} else {
				if item.Tok == BEGIN {
					return []string{}, fmt.Errorf("Literal { inside context %s section\n", name)
				} else if item.Tok == IDENT {
					if fnd_val == false {
						result = strings.Split(item.Val, "\n")
						for i := range result {
							result[i] = strings.Trim(result[i], " \t\n")
						}
						fnd_val = true
					} else {
						return []string{}, fmt.Errorf("Find two bodies into %s section\n", name)
					}
				} else if item.Tok == END {
					fnd_open = false
					break
				}
			}
		}
	}
	if fnd_open == true {
		return []string{}, fmt.Errorf("Not closed %s section\n", name)
	}
	if fnd_val == true {
		return result, nil
	}
	return []string{}, fmt.Errorf("Can't find %s tag\n", name)
}

func findVariables(preScan []preScanned) (map[string]PluginInfoDataContainerVars, error) {
	var fnd_tag bool = false
	var fnd_open bool = false
	var fnd_val bool = false
	var result map[string]PluginInfoDataContainerVars = map[string]PluginInfoDataContainerVars{}
	var var_name string
	for _, item := range preScan {
		if fnd_tag == false {
			val := strings.Trim(item.Val, " \n\t")
			if item.Tok == IDENT && val[0] == '_' {
				fnd_tag = true
				var_name = val[1:]
			}
		} else {
			if fnd_open == false {
				if item.Tok == BEGIN {
					fnd_open = true
				}
			} else {
				if item.Tok == BEGIN {
					return map[string]PluginInfoDataContainerVars{}, fmt.Errorf("Literal { inside context %s section\n", var_name)
				} else if item.Tok == IDENT {
					if fnd_val == false {
						var chg bool = false
						var container PluginInfoDataContainerVars
						result_tmp := strings.Split(item.Val, "\n")
						for i := range result_tmp {
							result_tmp[i] = strings.Trim(result_tmp[i], " \t\n")
							if result_tmp[i] == "" {
								continue
							}
							result_tmp2 := strings.SplitN(result_tmp[i], "=", 2)
							if len(result_tmp2) > 1 {
								for i := range result_tmp2 {
									result_tmp2[i] = strings.Trim(result_tmp2[i], " \t\n")
								}
								if strings.ToLower(result_tmp2[0]) == "def" {
									container.Def_cmd = result_tmp2[1]
									chg = true
								}
								if strings.ToLower(result_tmp2[0]) == "cmd" {
									container.Cmd = result_tmp2[1]
									chg = true
								}
							} else {
								container.Def = result_tmp[i]
								chg = true
							}
						}
						fnd_val = true
						if chg == true {
							result[var_name] = container
						}
					} else {
						return map[string]PluginInfoDataContainerVars{}, fmt.Errorf("Find two bodies into %s section\n", var_name)
					}
				} else if item.Tok == END {
					fnd_open = false
					fnd_tag = false
					fnd_val = false
				}
			}
		}
	}
	return result, nil
}

func findLangLikeParam(preScan []preScanned, name string, default_val string) ([]PluginInfoDataContainerLang, string, error) {
	var result []string
	var err error
	var end_result []PluginInfoDataContainerLang
	def_lan := default_val
	if result, err = findSliceValue(preScan, name); err != nil {
		return []PluginInfoDataContainerLang{}, default_val, err
	}
	for _, value := range result {
		res := configparser.SplitOwnLongSep(value, []string{"="})
		if len(res) > 1 {
			slice_res := configparser.SplitOwn(res[1])
			if len(slice_res) > 0 {
				def := strings.Trim(strings.ToLower(res[0]), " \t\n")
				if def == "default" {
					def_lan = slice_res[0]
				} else {
					for index := range slice_res {
						if slice_res[index][0] != '.' {
							slice_res[index] = "." + slice_res[index]
						}
					}
					end_result = append(end_result, PluginInfoDataContainerLang{def, slice_res})
				}
			}
		}
	}
	return end_result, def_lan, nil
}

func findExtraOpt(preScan []preScanned) (map[string][]string, error) {
	var err error = fmt.Errorf("Extra Options not found")
	var result []string
	opts := map[string][]string{}
	if result, err = findSliceValue(preScan, "EXTRAOPTIONS"); err != nil {
		return opts, err
	}
	for _, value := range result {
		res := configparser.SplitOwnLongSep(value, []string{":"})
		if len(res) > 1 {
			opts[res[0]] = append(opts[res[0]], configparser.SplitOwn(strings.Trim(res[1], " \n\t"))...)
			opts[res[0]] = configparser.RemoveDuplicate(opts[res[0]])
		}
	}
	return opts, nil
}

func findResult(preScan []preScanned) ([]string, rune, string, map[string][]string, bool, int, string, error) {
	var result []string
	var err error
	var end_result []string
	def_delim := rune('|')
	def_stream := "stdout"
	def_keys := map[string][]string{"low": []string{}, "medium": []string{}, "high": []string{}}
	dntstop := false
	indent_spc := 0
	clean := ""
	if result, err = findSliceValue(preScan, "RESULT"); err != nil {
		return []string{}, def_delim, def_stream, def_keys, dntstop, indent_spc, clean, err
	}
	for _, value := range result {
		res := configparser.SplitOwnLongSep(value, []string{"="})
		if len(res) > 1 {
			if strings.ToLower(res[0]) == "delimit" {
				def_delim_tmp := []rune(strings.Trim(res[1], " \t\n"))
				if len(def_delim_tmp) > 0 {
					def_delim = def_delim_tmp[0]
				}
			}
			if strings.ToLower(res[0]) == "stream" {
				def_stream_tmp := strings.Trim(res[1], " \t\n")
				if len(def_stream_tmp) > 0 {
					def_stream = def_stream_tmp
				}
			}
			if strings.ToLower(res[0]) == "clean" {
				clean_tmp := strings.Trim(res[1], " \t\n")
				if len(clean_tmp) > 0 {
					clean = clean_tmp
				}
			}
			if strings.ToLower(res[0]) == "dontstop" {
				dntstop = true
			}
			if strings.ToLower(res[0]) == "spaces" {
				i, err := strconv.ParseInt(strings.Trim(res[1], " \n\t"), 10, 64)
				if err == nil {
					indent_spc = int(i)
				}
			}
			if strings.ToLower(res[0]) == "low" || strings.ToLower(res[0]) == "medium" || strings.ToLower(res[0]) == "high" {
				def_keys[strings.ToLower(res[0])] = configparser.SplitOwn(strings.Trim(res[1], " \t\n"))
				for key := range def_keys[strings.ToLower(res[0])] {
					def_keys[strings.ToLower(res[0])][key] = strings.ToLower(def_keys[strings.ToLower(res[0])][key])
				}
			}
		} else if len(res) > 0 {
			end_result = append(end_result, strings.Trim(res[0], " \t\n"))
		}
	}
	return end_result, def_delim, strings.ToLower(def_stream), def_keys, dntstop, indent_spc, clean, nil
}

func (obj *PluginInfoDataContainer) Parse(file_name string) error {
	var preScan []preScanned
	var tok Token
	var lit string
	if file_data, err := ioutil.ReadFile(file_name); err == nil {
		s := NewScanner(bytes.NewReader(file_data))
		for {
			if tok, lit = s.Scan(); tok == EOF {
				break
			}
			preScan = append(preScan, preScanned{tok, lit})
		}
		if check_prescan_tokens(preScan) == true {
			if nm, err := findOneValue(preScan, "NAME"); err != nil {
				return err
			} else {
				obj.name = nm
			}
			if nm, err := findOneValue(preScan, "DESCRIPTION"); err != nil {
				return err
			} else {
				obj.description = nm
			}
			if nm, err := findOneValue(preScan, "ID"); err != nil {
				return err
			} else {
				obj.id = nm
			}
			if nm, err := findOneValue(preScan, "TYPE"); err != nil {
				return err
			} else {
				obj.pType = nm
			}
			if nm, err := findOneValue(preScan, "FILENAME"); err != nil {
				return err
			} else {
				obj.filename = nm
			}
			if nm, err := findSliceValue(preScan, "OPTIONS"); err != nil {
				return err
			} else {
				obj.options = nm
			}
			if nm, err := findSliceValue(preScan, "DEFS"); err == nil {
				obj.defs = nm
			}
			if nm, err := findSliceValue(preScan, "INCLUDES"); err == nil {
				obj.includes = nm
			}
			if nm, err := findOneValue(preScan, "AUTOINCLUDE"); err != nil {
				return err
			} else {
				obj.autoinclude = nm
			}
			if nm, err := findOneValue(preScan, "CMD"); err != nil {
				return err
			} else {
				obj.cmd = nm
			}
			if nm, err := findOneValue(preScan, "BEFORECMD"); err != nil {
				obj.precommand = ""
			} else {
				obj.precommand = nm
			}
			if _, err := findOneValue(preScan, "FRESH"); err != nil {
				obj.fresh = false
			} else {
				obj.fresh = true
			}
			if _, err := findOneValue(preScan, "COMPOSE"); err != nil {
				obj.compose = false
			} else {
				obj.compose = true
			}
			if nm, def, err := findLangLikeParam(preScan, "LANG", "c"); err == nil {
				obj.lang = nm
				obj.lang_default = def
			}
			if nm, def, err := findLangLikeParam(preScan, "ARCH", "64"); err == nil {
				for _, val := range nm {
					for j := range val.Ext {
						if len(val.Ext[j]) > 1 {
							val.Ext[j] = val.Ext[j][1:]
						}
					}
				}
				obj.arch = nm
				obj.arch_default = def
			}
			if nm, def, err := findLangLikeParam(preScan, "PLAT", "unix64"); err == nil {
				for _, val := range nm {
					for j := range val.Ext {
						if len(val.Ext[j]) > 1 {
							val.Ext[j] = val.Ext[j][1:]
						}
					}
				}
				obj.plat = nm
				obj.plat_default = def
			}
			if nm, def, str, levels, dnt, spc, cln, err := findResult(preScan); err != nil {
				return err
			} else {
				obj.result = nm
				obj.result_delim = def
				obj.stream = str
				obj.result_levels = levels
				obj.dontstop = dnt
				obj.result_spaces = spc
				obj.result_clean = cln
			}
			if extraopt, err := findExtraOpt(preScan); err == nil {
				obj.extraoptions = extraopt
			}
			if nm, err := findVariables(preScan); err != nil {
				return err
			} else {
				obj.variables = nm
			}
			return nil
		} else {
			return fmt.Errorf("Incorrect syntaxis of file: ", file_name)
		}
	} else {
		return fmt.Errorf("Error reading file: ", file_name)
	}
	return fmt.Errorf("Unknown parsing error :(\n")
}

func checkPluginName(plg *PluginInfoDataContainer) bool {
	for _, item := range listOfPlugins {
		if item.GetName() == plg.GetName() {
			return true
		}
	}
	return false
}

func MakePluginsList(dir_to_find string) {
	if info, err := os.Stat(dir_to_find); err == nil {
		if info.IsDir() == true {
			files, err := ioutil.ReadDir(dir_to_find)
			if err != nil {
				fmt.Printf("Error %s in reading of dir %s\n", err, dir_to_find)
				return
			}

			for _, file := range files {
				if path.Ext(file.Name()) == ".conf" {
					file_name := dir_to_find + "/" + file.Name()
					plg := pluginInfoDataContainer_make()
					if err := plg.Parse(file_name); err == nil {
						if checkPluginName(plg) == true {
							var listOfPlugins_tmp []*PluginInfoDataContainer
							for _, item := range listOfPlugins {
								if item.GetName() == plg.GetName() {
									continue
								} else {
									listOfPlugins_tmp = append(listOfPlugins_tmp, plg)
								}
							}
							listOfPlugins = listOfPlugins_tmp
						}
						listOfPlugins = append(listOfPlugins, plg)

					} else {
						fmt.Printf("Error %s in pasing of file %s\n", err, file_name)
					}
				}
			}
		}
	}
	return

}

func RemoveDuplicated() {
	listOfPlugins_tmp := []*PluginInfoDataContainer{}
	for _, val := range listOfPlugins {
		found := false
		for _, val_tmp := range listOfPlugins_tmp {
			if val_tmp.GetName() == val.GetName() {
				found = true
				break
			}
		}
		if found == false {
			listOfPlugins_tmp = append(listOfPlugins_tmp, val)
		}
	}
	listOfPlugins = listOfPlugins_tmp
}

func (obj PluginInfoDataContainer) String() string {
	buf := fmt.Sprintf("Name: %s\nDesc: %s\nID: %s\nType: %s\nOptions: %s\nResult: %s\nDelim: %s\nDefs: %s\nIncs: %s\nFName: %s\nCmd: %s\nLangdefault: %s\nLang:\n",
		obj.name, obj.description, obj.id, obj.pType, strings.Join(obj.options, ", "),
		strings.Join(obj.result, ", "), string(obj.result_delim),
		strings.Join(obj.defs, ", "), strings.Join(obj.includes, ", "), obj.filename, obj.cmd,
		obj.lang_default)

	for _, key := range obj.lang {
		buf = buf + fmt.Sprintf("%s\n", key)
	}
	buf = buf + fmt.Sprintf("ArchDefault: %s\nArch:\n", obj.arch_default)
	for _, key := range obj.arch {
		buf = buf + fmt.Sprintf("%s\n", key)
	}
	buf = buf + fmt.Sprintf("PlatDefault: %s\nPlat:\n", obj.plat_default)
	for _, key := range obj.plat {
		buf = buf + fmt.Sprintf("%s\n", key)
	}
	buf = buf + "Variables:\n"
	for key := range obj.variables {
		buf = buf + fmt.Sprintf("   %s [def: %s  def_cmd: %s  cmd: %s]\n", key,
			obj.variables[key].Def, obj.variables[key].Def_cmd, obj.variables[key].Cmd)
	}
	buf = buf + fmt.Sprintf("Result stream: %s\n", obj.stream)
	for key := range obj.result_levels {
		buf += fmt.Sprintf(" Level %s: %s\n", key, obj.result_levels[key])
	}
	buf = buf + fmt.Sprintf("Result dontstop: %s\n", obj.dontstop)
	buf = buf + fmt.Sprintf("Cmd autoinclude: %s\n", obj.autoinclude)
	buf = buf + fmt.Sprintf("Result spaces: %d\n", obj.result_spaces)
	buf = buf + fmt.Sprintf("Result clean: %d\n", obj.result_clean)
	buf = buf + fmt.Sprintf("Precommand: %d\n", obj.precommand)
	buf = buf + fmt.Sprintf("Fresh: %s\n", obj.fresh)
	return buf
}

func (obj *PluginInfoDataContainer) GetPluginCMD(file_name string, incs []string, defs []string, conf *configparser.ConfigparserContainer, cincs []string, cppincs []string, raw []string) (string, error) {
	return parseCMD(obj.cmd, file_name, incs, defs, obj, conf, cincs, cppincs, raw)
}

func (obj *PluginInfoDataContainer) getVariableValue(name string, file string, incs []string, defs []string, conf *configparser.ConfigparserContainer, cincs []string, cppincs []string, raw []string) (string, bool, error) {
	if name == "RAW" {
		if len(raw) > 0 {
			return strings.Join(raw, " ") + " -c " + file, true, nil
		} else {
			return "-c " + file, true, nil
		}
	} else if name == "FRESH" {
		return strings.Join(obj.cmd_for_fresh, " "), true, nil
	} else if name == "CUSTOMINCLUDES" {
		result := ""
		lang_res := obj.lang_default
		for _, val := range obj.lang {
			for _, ext := range val.Ext {
				if len(file)-len(ext) > 0 && len(file)-len(ext) < len(file) && file[len(file)-len(ext):] == ext {
					lang_res = val.Lang
					break
				}
			}
		}
		res := strings.Join(obj.includes, " ")
		res_new := strings.Split(res, " ")
		for i := range res_new {
			res_new[i] = strings.Trim(res_new[i], " \n\t")
		}
		res_slice := []string{}
		i := 0
		compil_result := ""
		if lang_res == "c" {
			for _, val := range cincs {
				if strings.HasPrefix(val, obj.autoinclude) == false {
					res_slice = append(res_slice, strings.Replace(res_new[i%len(res_new)], "$:", val, 1))
					i++
				} else {
					compil_result += " " + val
				}
			}
		}
		i = 0
		if lang_res == "c++" {
			for _, val := range cppincs {
				if strings.HasPrefix(val, obj.autoinclude) == false {
					res_slice = append(res_slice, strings.Replace(res_new[i%len(res_new)], "$:", val, 1))
					i++
				} else {
					compil_result += " " + val
				}
			}
		}
		if len(compil_result) > 0 || len(res_slice) > 0 {
			result += compil_result + " " + strings.Join(res_slice, " ")
		}
		return result, true, nil
	} else if name == "CUSTOMDEFS" {
		result := conf.GetFileAddInfo(file)
		return result, true, nil
	} else if name == "LANG" {
		lang_res := obj.lang_default
		for _, val := range obj.lang {
			for _, ext := range val.Ext {
				if len(file)-len(ext) > 0 && len(file)-len(ext) < len(file) && file[len(file)-len(ext):] == ext {
					lang_res = val.Lang
					break
				}
			}
		}
		return lang_res, true, nil
	} else if name == "DEFS" {
		if len(defs) > 0 {
			res := strings.Join(obj.defs, " ")
			res_new := strings.Split(res, " ")
			for i := range res_new {
				res_new[i] = strings.Trim(res_new[i], " \n\t")
			}
			res_slice := []string{}
			for i, val := range defs {
				res_slice = append(res_slice, strings.Replace(res_new[i%len(res_new)], "$:", val, 1))
			}
			return strings.Join(res_slice, " "), true, nil
		} else {
			return "", true, nil
		}
	} else if name == "INCLUDES" {
		if len(incs) > 0 {
			res := strings.Join(obj.includes, " ")
			res_new := strings.Split(res, " ")
			for i := range res_new {
				res_new[i] = strings.Trim(res_new[i], " \n\t")
			}
			res_slice := []string{}
			for i, val := range incs {
				res_slice = append(res_slice, strings.Replace(res_new[i%len(res_new)], "$:", val, 1))
			}
			return strings.Join(res_slice, " "), true, nil
		} else {
			return "", true, nil
		}
	} else if name == "FILENAME" {
		res := strings.Replace(obj.filename, "$FILE", file, -1)
		return res, true, nil
	} else if name == "PLAT" {
		res := obj.plat_default
		if runtime.GOARCH == "amd64" {
			for _, val := range obj.plat {
				if val.Lang == "64" && len(val.Ext) > 0 {
					res = val.Ext[0]
					break
				}
			}
		} else {
			for _, val := range obj.plat {
				if val.Lang == "32" && len(val.Ext) > 0 {
					res = val.Ext[0]
					break
				}
			}
		}
		return res, true, nil
	} else if name == "ARCH" {
		res := obj.arch_default
		if runtime.GOARCH == "amd64" {
			for _, val := range obj.arch {
				if val.Lang == "64" && len(val.Ext) > 0 {
					res = val.Ext[0]
					break
				}
			}
		} else {
			for _, val := range obj.arch {
				if val.Lang == "32" && len(val.Ext) > 0 {
					res = val.Ext[0]
					break
				}
			}
		}
		return res, true, nil
	} else if name == "FILE" {
		return file, true, nil
	} else if name == "OPTIONS" {
		res := strings.Join(obj.options, " ")
		if len(obj.accumulated_extra_opts) > 0 {
			res = res + strings.Join(obj.accumulated_extra_opts, " ")
		}
		return res, true, nil
	} else {
		res := ""
		var err error
		val, found := obj.variables[name]
		if found == true {
			if val.Cmd != "" {
				if res, err = executeCommandLine(val.Cmd, conf); err == nil {
					return res, true, nil
				}
				fmt.Println("Error for command ", val.Cmd, " with error ", err)
			}
			if val.Def_cmd != "" {
				if res, err = executeCommandLine(val.Def_cmd, conf); err == nil {
					return res, true, nil
				}
				fmt.Println("Error for command ", val.Def_cmd, " with error ", err)
			}
			return val.Def, true, err
		}
	}
	return "", false, nil
}

func parseCMD(cmd string, file_name string, incs []string, defs []string, obj *PluginInfoDataContainer, conf *configparser.ConfigparserContainer, cincs []string, cppincs []string, raw []string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(cmd))
	scanner.Split(bufio.ScanWords)
	var variables []string = []string{}
	var cmd_result string = cmd
	for scanner.Scan() {
		vr := strings.Trim(scanner.Text(), " \n\t")
		if vr != "" {
			if len(vr) > 1 {
				var vr_nm string = ""
				var fnd_var = false
				for _, ch := range vr {
					if ch == '$' {
						fnd_var = true
						if vr_nm != "" {
							variables = append(variables, vr_nm)
							vr_nm = ""
						}
					}
					if fnd_var == true && ch != '$' {
						vr_nm = vr_nm + string(ch)
					}
				}
				if vr_nm != "" {
					variables = append(variables, vr_nm)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if len(variables) > 0 {
		variables = configparser.RemoveDuplicate(variables)
		for _, cVr := range variables {
			str, found, _ := obj.getVariableValue(cVr, file_name, incs, defs, conf, cincs, cppincs, raw)
			if found == true {
				str, _ = parseCMD(str, file_name, incs, defs, obj, conf, cincs, cppincs, raw)
			}
			cmd_result = strings.Replace(cmd_result, "$"+cVr, str, -1)
		}
	}
	return cmd_result, nil
}

func executeCommandLine(cmd string, conf *configparser.ConfigparserContainer) (string, error) {
	var err error
	var out []byte
	if len(cmd) > 0 {
		start_cmd := append(conf.GetBash(), cmd)
		if len(start_cmd) == 1 {
			out, err = exec.Command(start_cmd[0]).CombinedOutput()
		} else {
			out, err = exec.Command(start_cmd[0], start_cmd[1:]...).CombinedOutput()
		}
		if err != nil {
			return "", fmt.Errorf("Get error %s on command starting %s\n", err, cmd)
		}

		out_res := strings.Trim(strings.Join(configparser.SplitOwnLongSepNoTrimSep(string(out), []string{" ", "\t", "\n"}), " "), " \n\t")
		return out_res, nil
	}
	return "", fmt.Errorf("Empty command in plugin config\n")
}

func GetPluginContainerById(id string) *PluginInfoDataContainer {
	for _, obj := range listOfPlugins {
		if obj.id == id {
			return obj
		}
	}
	return nil
}

func GetPluginContainerList() string {
	result := fmt.Sprintf("%15s|%40s|%10s\n", "Name", "Description", "Id")
	result = result + fmt.Sprintf("---------------+----------------------------------------+----------\n")
	for _, obj := range listOfPlugins {
		var_name := obj.name
		if len(var_name) > 15 {
			var_name = var_name[:14] + "+"
		}
		var_id := obj.id
		var_buf := strings.NewReader(obj.description)
		fst := true
	loop:
		var var_result []rune = []rune{}
		for i := 0; i < 40; i++ {
			ch, size, err := var_buf.ReadRune()
			if size > 0 && err == nil {
				var_result = append(var_result, ch)
			} else {
				break
			}
		}
		if fst == true {
			result += fmt.Sprintf("%-15s|%-40s|%-10s\n", var_name, string(var_result), var_id)
		} else {
			result += fmt.Sprintf("%-15s|%-40s|%-10s\n", " ", string(var_result), " ")
		}
		fst = false
		if var_buf.Len() > 0 {
			goto loop
		}
		result = result + fmt.Sprintf("---------------+----------------------------------------+----------\n")

	}
	return result
}

func (obj *PluginInfoDataContainer) GetNameId() string {
	return obj.id
}

func (obj *PluginInfoDataContainer) GetName() string {
	return obj.name
}

func (obj *PluginInfoDataContainer) GetStream() string {
	return obj.stream
}

func (obj *PluginInfoDataContainer) GetDelim() rune {
	return obj.result_delim
}

func (obj *PluginInfoDataContainer) GetResult() []string {
	return obj.result
}

func (obj *PluginInfoDataContainer) GetType() string {
	return obj.pType
}

func (obj *PluginInfoDataContainer) GetResultLevels(ident string) string {
	ident = strings.ToLower(strings.Trim(ident, " \n\t"))
	for key := range obj.result_levels {
		for _, val := range obj.result_levels[key] {
			if val == ident {
				return strings.ToUpper(key)
			}
		}
	}
	return "NORMAL"
}

func (obj *PluginInfoDataContainer) GetResultStop() bool {
	return obj.dontstop
}

func (obj *PluginInfoDataContainer) GetAutoIncludes() string {
	return obj.autoinclude
}

func (obj *PluginInfoDataContainer) GetSpaces() int {
	return obj.result_spaces
}

func (obj *PluginInfoDataContainer) GetFresh() bool {
	return obj.fresh
}

func (obj *PluginInfoDataContainer) GetPreCommand() string {
	return obj.precommand
}

func (obj *PluginInfoDataContainer) SetCmdForFresh(lst []string) {
	obj.cmd_for_fresh = lst
}

func (obj *PluginInfoDataContainer) GetClean() string {
	return obj.result_clean
}

func GetFullPluginList() []*PluginInfoDataContainer {
	return listOfPlugins
}

func (obj *PluginInfoDataContainer) GetCompose() bool {
	return obj.compose
}

func (obj *PluginInfoDataContainer) GetExtraOptions(file_name string) ([]string, bool) {
	for f_n, value := range obj.extraoptions {
		if strings.Contains(file_name, f_n) == true {
			return value, true
		}
	}
	return []string{}, false
}

func (obj *PluginInfoDataContainer) SetAccumulatedExtraOptions(opts []string) {
	obj.accumulated_extra_opts = append(obj.accumulated_extra_opts, opts...)
	obj.accumulated_extra_opts = configparser.RemoveDuplicate(obj.accumulated_extra_opts)
}
