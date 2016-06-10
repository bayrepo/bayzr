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

package replacer

import (
	"bufio"
	"configparser"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"outputanalyzer"
	"strings"
)

type ReplacerContainer struct {
	program    string
	programcc  string
	programcxx string
}

func (this ReplacerContainer) String() string {
	return fmt.Sprintf("Program: %s\nCC: %s\nCXX: %s\n", this.program,
		this.programcc, this.programcxx)
}

func Make_ReplacerContainer() *ReplacerContainer {
	this := &ReplacerContainer{}
	if len(os.Args) > 0 {
		this.program = os.Args[0]
	} else {
		this.program = ""
	}
	this.programcc = "/usr/bin/gcc"
	this.programcxx = "/usr/bin/cc"
	return this
}

func (this *ReplacerContainer) SetGCCCompilers(command []string) []string {
	list := []string{}
	for _, val := range command {
		res := []rune(val)
		if len(res) > 1 && res[0] == rune('C') && res[1] == rune('C') {
			this.programcc = val[2:]
			continue
		}
		if len(res) > 2 && res[0] == rune('C') && res[1] == rune('X') && res[1] == rune('X') {
			this.programcxx = val[3:]
			continue
		}
		list = append(list, val)
	}
	os.Setenv("CC", this.program+" -cc cmd")
	os.Setenv("CXX", this.program+" -cxx cmd")
	list = append(list, "CC=\" --tag=CC "+this.program+" -cc cmd\"")
	list = append(list, "CXX=\" --tag=CXX "+this.program+" -cxx cmd\"")
	return list
}

var ccFlag *bool
var cxxFlag *bool

func init() {
	ccFlag = flag.Bool("cc", false, "Run as C compiler wrapper")
	cxxFlag = flag.Bool("cxx", false, "Run as C++ compiler wrapper")
}

func (this *ReplacerContainer) RunWrapper(storage *outputanalyzer.OutPutAnalyzerContainer, config *configparser.ConfigparserContainer) {
	compiler := this.programcc
	if *cxxFlag == true {
		compiler = this.programcxx
	}
	if len(os.Args) > 2 {
		currentPath := ""
		var err error
		if currentPath, err = os.Getwd(); err != nil {
			currentPath = ""
		}
		lst := []string{compiler}
		lst = append(lst, os.Args[3:]...)

		files_list, _, _ := outputanalyzer.GetFilesNames(lst, storage)
		new_args := []string{}

		for _, val := range os.Args[3:] {
			for _, fl := range files_list {
				if strings.Contains(val, fl) == true {
					if currentPath != "" {
						val = strings.Replace(val, fl, currentPath+"/"+fl, -1)
					}
				}
			}
			new_args = append(new_args, val)
		}
		cmd_str := compiler + " " + strings.Join(new_args, " ")
		cmd_str_orig := compiler + " " + strings.Join(os.Args[3:], " ")
		fmt.Println(cmd_str)
		var cmd *exec.Cmd
		start_cmd := append(config.GetBash(), cmd_str_orig)
		if len(start_cmd) == 1 {
			cmd = exec.Command(start_cmd[0])
		} else {
			cmd = exec.Command(start_cmd[0], start_cmd[1:]...)
		}
		var stdout io.ReadCloser
		var stderr io.ReadCloser
		if stdout, err = cmd.StdoutPipe(); err != nil {
			fmt.Fprintf(os.Stderr, "Replacer error: open stdout pipe error %s\n", err)
		}
		if stderr, err = cmd.StderrPipe(); err != nil {
			fmt.Fprintf(os.Stderr, "Replacer error: open stdout pipe error %s\n", err)
		}
		scanner_out := bufio.NewScanner(stdout)
		scanner_err := bufio.NewScanner(stderr)
		go func() {
			for scanner_err.Scan() {
				fmt.Fprintln(os.Stderr, scanner_err.Text())
			}

		}()
		if err = cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Replacer error: start command error %s\n", err)
		}
		for scanner_out.Scan() {
			fmt.Fprintln(os.Stderr, scanner_out.Text())
		}
		if err = scanner_out.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Replacer error: stdout scanner error %s\n", err)
		}

		if err = cmd.Wait(); err != nil {
			fmt.Fprintf(os.Stderr, "Analyzer error: wait command error %s\n", err)
		}
	}
}

func (this *ReplacerContainer) IsReplaced() bool {
	return *ccFlag == true || *cxxFlag == true
}
