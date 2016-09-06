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

package main

import (
	"checker"
	"configparser"
	"crypto/md5"
	"diffanalyzer"
	"flag"
	"fmt"
	"gccincludes"
	"io"
	"mysqlsaver"
	"os"
	"os/user"
	"outputanalyzer"
	"path/filepath"
	"replacer"
	"reporter"
	"resultanalyzer"
	"strings"
	"templater"
	"visualmenu"
)

const APP_VERSION = "0.1"

var versionFlag *bool
var configFile string
var listOfPlugins *bool
var useOnlyLocalFiles *bool
var listOfOutputFiles string
var listOfDiffFiles string
var printAnalizerCommands *bool
var printVisualMenu *bool
var dryRun *bool
var author string
var build_info string

const (
	config_file_name = "bzr.conf"
	plugin_dir       = "bzr.d"
)

func is_value_in_array(val string, lst []string) bool {
	for _, item := range lst {
		if val == item {
			return true
		}
	}
	return false
}

func init() {
	const (
		defaultFile      = "/etc/" + config_file_name
		defaultFileusage = "path to configuration file"
	)
	versionFlag = flag.Bool("version", false, "Print the version number.")
	flag.StringVar(&configFile, "config", defaultFile, defaultFileusage)
	listOfPlugins = flag.Bool("list", false, "Show list of available plugins")
	useOnlyLocalFiles = flag.Bool("not-only-local", false, "Show in result errors not only for project files")
	flag.StringVar(&listOfOutputFiles, "files", "*", "List of files should be inserted to report or * for all files(by dafault)")
	flag.StringVar(&listOfDiffFiles, "diff", "", "List of patch file for get list of patched files")
	printAnalizerCommands = flag.Bool("debug-commands", false, "Show list of generated static analizers options and commands")
	printVisualMenu = flag.Bool("menu", false, "Show console menu for project options configuring")
	dryRun = flag.Bool("dry-run", false, "Show list of generated static analizers options and commands without analitic tool starting")
	flag.StringVar(&author, "build-author", "build-master", "Set of build's author")
	flag.StringVar(&build_info, "build-name", "autobuild", "Set of build's name")
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("    bayzr [options] cmd ...\n")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {

	var err error
	var list_of_analyzer_commands map[string][]string = map[string][]string{}

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
		os.Exit(0)
	}

	replace := replacer.Make_ReplacerContainer()

	current_analyzer_path := "/"

	config := configparser.CreateDefaultConfig()
	//ищем стандартный конфигурационный файл
	config.ReadConfig(configFile)
	checker.MakePluginsList("/etc/" + plugin_dir)
	config.AddFndPath("/etc/" + plugin_dir + "/")
	//ищем конфигурационный файл в домашней директории
	var user_info *user.User
	if user_info, err = user.Current(); err == nil {
		config.ReadConfig(user_info.HomeDir + "/" + config_file_name)
		checker.MakePluginsList(user_info.HomeDir + "/" + plugin_dir)
		config.AddFndPath(user_info.HomeDir + "/" + plugin_dir + "/")
	}
	//ищем конфигурационный файл в текущей диретории
	if currentPath, err := os.Getwd(); err == nil {
		if *useOnlyLocalFiles == false {
			current_analyzer_path = currentPath + "/"
		}
		config.ReadConfig(currentPath + "/" + config_file_name)
		checker.MakePluginsList(currentPath + "/" + plugin_dir)
		config.AddFndPath(currentPath + "/" + plugin_dir + "/")
	}
	//дозаполняем дефолтными значениями если после чтения конфигурации ничего не нашлось
	config.DefaultPropogate()

	checker.RemoveDuplicated()

	if *printVisualMenu == true {
		menu := visualmenu.Make_VisualMenu(config, checker.GetFullPluginList())
		menu.Show()
		os.Exit(0)
	}

	analyzer := outputanalyzer.CreateDefaultAnalyzer(config)

	if replace.IsReplaced() == true {
		replace.RunWrapper(analyzer, config)
		os.Exit(0)
	}

	var DBase mysqlsaver.MySQLSaver
	dbErr := DBase.Init(config.Connector())
	if dbErr != nil {
		fmt.Printf("DataBase saving error %s\n", dbErr)
		os.Exit(1)
	}
	defer DBase.Finalize()

	if *listOfPlugins {
		fmt.Println(checker.GetPluginContainerList())
		os.Exit(0)
	}

	if len(listOfDiffFiles) > 0 {
		patch_list := diffanalyzer.Make_DiffAnalyzerContainer()
		patch_list.ParseFilesList(listOfDiffFiles)
		files_from_patch, linenembers_form_patch := patch_list.GetFoundList()
		config.SetFilesList_list(files_from_patch, linenembers_form_patch)
	} else {
		config.SetFilesList(listOfOutputFiles)
	}

	gcc := gccincludes.Make_GCCContainer()
	gcc.GetGCCIncludes(config)

	list_of_fresh := []string{}

	for _, plg_fst := range config.GetListOfPlugins() {
		obj_item := checker.GetPluginContainerById(plg_fst)
		if obj_item == nil {
			continue
		}
		if obj_item.GetFresh() == true {
			list_of_fresh = append(list_of_fresh, plg_fst)
		}
	}

	fmt.Println("--------------------Process of gathering source information is begun-----------------")

	cmd_for_fresh := analyzer.GetStringOfCmd(os.Args)
	cmd_mod_ := cmd_for_fresh

	if config.Replacer() == true {
		cmd_mod_ = replace.SetGCCCompilers(cmd_mod_)
	}

	if err = analyzer.ExcecuteComplilationProcess(cmd_mod_); err != nil {
		fmt.Println(err)
	}

	fmt.Println("--------------------Process of source analyzing is begun-----------------------------")

	list_of_result := []*resultanalyzer.ResultAnalyzerConatiner{}
	if len(config.GetListOfPlugins()) > len(list_of_fresh) {
		for _, plg := range config.GetListOfPlugins() {
			if is_value_in_array(plg, list_of_fresh) == true {
				continue
			}
			obj_item := checker.GetPluginContainerById(plg)
			if obj_item == nil {
				continue
			}
			fmt.Printf("--------------------Process of source analyzing is begun by plugin %s----------------\n", obj_item.GetName())
			if obj_item != nil {
				if obj_item.GetCompose() == true {
					files_list := map[string]map[string]map[string][]string{}
					c_f, cpp_f := gcc.MakeHeaders(obj_item.GetAutoIncludes())
					new_incc_list := append(gcc.GetC(), c_f)
					new_inccpp_list := append(gcc.GetCPP(), cpp_f)
					val_new_inc := map[string][]string{}
					val_new_defs := map[string][]string{}
					extaopts_for_files := map[string][]string{}
					for file_name, value := range *analyzer.GetParsedFilesList() {
						if config.IsFileIgnored(file_name) == false {
							file_name = outputanalyzer.TransformFileName(file_name, value.Dir)
							if file_name == "" {
								continue
							}
							hash := md5.New()
							io.WriteString(hash, strings.Join(value.IncludesList, " ")+strings.Join(value.DefList, " "))
							hash_str := fmt.Sprintf("%x", hash.Sum(nil))
							if _, ok := val_new_inc[hash_str]; ok == false {
								val_new_inc[hash_str] = []string{}
							}
							if _, ok := val_new_defs[hash_str]; ok == false {
								val_new_defs[hash_str] = []string{}
							}
							if _, ok := files_list[hash_str]; ok == false {
								files_list[hash_str] = map[string]map[string][]string{}
							}
							val_new_inc[hash_str] = append(val_new_inc[hash_str], value.IncludesList...)
							val_new_defs[hash_str] = append(val_new_defs[hash_str], value.DefList...)
							val_new_inc[hash_str] = configparser.RemoveDuplicate(val_new_inc[hash_str])
							val_new_defs[hash_str] = configparser.RemoveDuplicate(val_new_defs[hash_str])
							dir := filepath.Dir(file_name)
							ext := filepath.Ext(file_name)
							base := filepath.Base(file_name)
							if extra, fnd := obj_item.GetExtraOptions(file_name); fnd == true {
								extaopts_for_files[hash_str] = append(extaopts_for_files[hash_str], extra...)
							}
							if _, ok := files_list[hash_str][dir]; ok == false {
								files_list[hash_str][dir] = map[string][]string{}
							}
							if ext == "" {
								files_list[hash_str][dir]["noextensionsinfilename"] = append(files_list[hash_str][dir]["noextensionsinfilename"], file_name)
							} else {
								files_list[hash_str][dir][ext] = append(files_list[hash_str][dir][ext], base)
							}
						}
					}
					files_list_last := map[string][]string{}
					for hash, value_tmp_h := range files_list {
						for dir_path, value_tmp := range value_tmp_h {
							for ext_name, file_value := range value_tmp {

								fnd := false
								for hash_tmp, value_tmp_tmp_h := range files_list {
									if hash_tmp != hash {
										for dir_path_tmp, _ := range value_tmp_tmp_h {
											if dir_path_tmp == dir_path {
												fnd = true
												break
											}
										}
									}
									if fnd == true {
										break
									}
								}
								if fnd == true {
									ext_name = "noextensionsinfilename"
								}

								if ext_name == "noextensionsinfilename" {
									for _, f_name := range file_value {
										files_list_last[hash] = append(files_list_last[hash], f_name)
									}
								} else {
									if len(file_value) > 1 {
										files_list_last[hash] = append(files_list_last[hash], dir_path+"/*"+ext_name)
									} else {
										for _, f_name := range file_value {
											files_list_last[hash] = append(files_list_last[hash], f_name)
										}
									}
								}
							}
						}
					}
					if len(files_list_last) == 0 {
						continue
					}
					for hash, f_list := range files_list_last {
						obj_item.SetAccumulatedExtraOptions(extaopts_for_files[hash])
						if cmd, err := obj_item.GetPluginCMD(strings.Join(f_list, " "), val_new_inc[hash], val_new_defs[hash], config, new_incc_list, new_inccpp_list, []string{}); err != nil {
							fmt.Println("Got error when cmd parsed ", err)
							os.Exit(1)
						} else {
							result_analyzer := resultanalyzer.Make_ResultAnalyzerConatiner(obj_item.GetName(), obj_item, current_analyzer_path, *config)
							list_of_result = append(list_of_result, result_analyzer)
							if *printAnalizerCommands == true || *dryRun == true {
								fmt.Println(cmd)
								if _, ok := list_of_analyzer_commands[obj_item.GetName()]; ok == false {
									list_of_analyzer_commands[obj_item.GetName()] = []string{}
								}
								list_of_analyzer_commands[obj_item.GetName()] = append(list_of_analyzer_commands[obj_item.GetName()], cmd)
								if *dryRun {
									continue
								}
							}
							if err := result_analyzer.ParseResultOfCommand(cmd, config); err != nil {
								fmt.Println("Got error when checker result parsed ", err)
								os.Exit(1)
							}
						}
					}
				} else {
					for file_name, value := range *analyzer.GetParsedFilesList() {
						if config.IsFileIgnored(file_name) == false {
							c_f, cpp_f := gcc.MakeHeaders(obj_item.GetAutoIncludes())
							new_incc_list := append(gcc.GetC(), c_f)
							new_inccpp_list := append(gcc.GetCPP(), cpp_f)
							file_name = outputanalyzer.TransformFileName(file_name, value.Dir)
							if file_name == "" {
								continue
							}
							extaopts_for_files := []string{}
							if extra, fnd := obj_item.GetExtraOptions(file_name); fnd == true {
								extaopts_for_files = extra
							}
							obj_item.SetAccumulatedExtraOptions(extaopts_for_files)
							if cmd, err := obj_item.GetPluginCMD(file_name, value.IncludesList, value.DefList, config, new_incc_list, new_inccpp_list, value.Raw); err != nil {
								fmt.Println("Got error when cmd parsed ", err)
								os.Exit(1)
							} else {
								result_analyzer := resultanalyzer.Make_ResultAnalyzerConatiner(file_name, obj_item, current_analyzer_path, *config)
								list_of_result = append(list_of_result, result_analyzer)
								if *printAnalizerCommands == true || *dryRun == true {
									fmt.Println(cmd)
									if _, ok := list_of_analyzer_commands[obj_item.GetName()]; ok == false {
										list_of_analyzer_commands[obj_item.GetName()] = []string{}
									}
									list_of_analyzer_commands[obj_item.GetName()] = append(list_of_analyzer_commands[obj_item.GetName()], cmd)
									if *dryRun {
										continue
									}
								}
								if err := result_analyzer.ParseResultOfCommand(cmd, config); err != nil {
									fmt.Println("Got error when checker result parsed ", err)
									os.Exit(1)
								}
							}
						}

					}
				}
			} else {
				fmt.Printf("Can't find plugin %s\n", plg)
				os.Exit(1)
			}
		}
	}

	for _, plg := range list_of_fresh {
		obj_item := checker.GetPluginContainerById(plg)
		if obj_item == nil {
			continue
		}
		fmt.Printf("--------------------Process of source analyzing is begun by plugin %s----------------\n", obj_item.GetName())
		if obj_item != nil {
			obj_item.SetCmdForFresh(cmd_for_fresh)
			if cmd, err := obj_item.GetPluginCMD(obj_item.GetName(), []string{}, []string{}, config, []string{}, []string{}, []string{}); err != nil {
				fmt.Println("Got error when cmd parsed ", err)
				os.Exit(1)
			} else {
				result_analyzer := resultanalyzer.Make_ResultAnalyzerConatiner(obj_item.GetName(), obj_item, current_analyzer_path, *config)
				list_of_result = append(list_of_result, result_analyzer)
				result_analyzer.MakePreCommand(config)
				if *printAnalizerCommands == true || *dryRun == true {
					fmt.Println(cmd)
					if _, ok := list_of_analyzer_commands[obj_item.GetName()]; ok == false {
						list_of_analyzer_commands[obj_item.GetName()] = []string{}
					}
					list_of_analyzer_commands[obj_item.GetName()] = append(list_of_analyzer_commands[obj_item.GetName()], cmd)
					if *dryRun {
						continue
					}
				}
				if err := result_analyzer.ParseResultOfCommand(cmd, config); err != nil {
					fmt.Println("Got error when checker result parsed ", err)
					os.Exit(1)
				}
				result_analyzer.RemoveGarbage(config)
			}
		} else {
			fmt.Printf("Can't find plugin %s\n", plg)
			os.Exit(1)
		}
	}

	dbErr = DBase.CreateCurrentBuild(author, build_info)
	if dbErr != nil {
		fmt.Printf("DataBase saving error %s\n", dbErr)
		os.Exit(1)
	}

	report := reporter.Make_ReporterContainer(config, &list_of_result, list_of_analyzer_commands, &DBase)
	if path, fnd := report.CreateReport(); fnd == true {
		tpl := templater.MakeTemplater()
		tpl.PropogateData(report, path, config)
	}
}
