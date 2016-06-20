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
	"flag"
	"fmt"
	"gccincludes"
	"os"
	"os/user"
	"outputanalyzer"
	"replacer"
	"reporter"
	"resultanalyzer"
	"templater"
)

const APP_VERSION = "0.1"

var versionFlag *bool
var configFile string
var listOfPlugins *bool
var useOnlyLocalFiles *bool
var listOfOutputFiles string

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
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("    bay [options] cmd ...\n")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {

	var err error

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
	
	analyzer := outputanalyzer.CreateDefaultAnalyzer(config)

	if replace.IsReplaced() == true {
		replace.RunWrapper(analyzer, config)
		os.Exit(0)
	}

	if *listOfPlugins {
		fmt.Println(checker.GetPluginContainerList())
		os.Exit(0)
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
				for file_name, value := range *analyzer.GetParsedFilesList() {
					if config.IsFileIgnored(file_name) == false {
						c_f, cpp_f := gcc.MakeHeaders(obj_item.GetAutoIncludes())
						new_incc_list := append(gcc.GetC(), c_f)
						new_inccpp_list := append(gcc.GetCPP(), cpp_f)
						file_name = outputanalyzer.TransformFileName(file_name, value.Dir)
						if file_name == "" {
						    continue
						}
						if cmd, err := obj_item.GetPluginCMD(file_name, value.IncludesList, value.DefList, config, new_incc_list, new_inccpp_list, value.Raw); err != nil {
							fmt.Println("Got error when cmd parsed ", err)
							os.Exit(1)
						} else {
							result_analyzer := resultanalyzer.Make_ResultAnalyzerConatiner(file_name, obj_item, current_analyzer_path)
							list_of_result = append(list_of_result, result_analyzer)
							if err := result_analyzer.ParseResultOfCommand(cmd, config); err != nil {
								fmt.Println("Got error when checker result parsed ", err)
								os.Exit(1)
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
				result_analyzer := resultanalyzer.Make_ResultAnalyzerConatiner(obj_item.GetName(), obj_item, current_analyzer_path)
				list_of_result = append(list_of_result, result_analyzer)
				result_analyzer.MakePreCommand(config)
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
	
	config.SetFilesList(listOfOutputFiles)

	report := reporter.Make_ReporterContainer(config, &list_of_result)
	if path, fnd := report.CreateReport(); fnd == true {
		tpl := templater.MakeTemplater()
		tpl.PropogateData(report, path, config)
	}
}
