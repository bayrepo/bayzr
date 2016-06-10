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

package gccincludes

import (
	"bufio"
	"configparser"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	cFile   = "cHeaderFileGeneratedByBzrPrg.h"
	cppFile = "cHeaderFileGeneratedByBzrPrg.hpp"
)

type GCCContainer struct {
	c     []string
	cpp   []string
	c_h   []string
	cpp_h []string
}

func (this GCCContainer) String() string {
	return fmt.Sprintf("C: %s\nC++: %s\n", strings.Join(this.c, ", "), strings.Join(this.cpp, ", "))
}

func Make_GCCContainer() *GCCContainer {
	return &GCCContainer{[]string{}, []string{}, []string{}, []string{}}
}

func (this *GCCContainer) makeStringAnalysis(line string, list []string) []string {
	line = strings.Trim(line, " \n\t")
	res := []rune(line)
	if len(res) > 0 && res[0] == rune('/') {
		if _, err := os.Stat(line); err == nil {
			return append(list, line)
		}
	}
	return list
}

func (this *GCCContainer) makeStringAnalysisDef(line string, list []string) []string {
	line = strings.Trim(line, " \n\t")
	return append(list, line)
}

func (this *GCCContainer) propogateCmd(cmd_in string, config *configparser.ConfigparserContainer, tp bool) ([]string, error) {
	var start_cmd []string
	var cmd *exec.Cmd
	list := []string{}
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
		if tp == true {
			if stdout, err = cmd.StdoutPipe(); err != nil {
				return []string{}, fmt.Errorf("GCCIncluder error: open stdout pipe error %s\n", err)
			}
		} else {
			if stdout, err = cmd.StderrPipe(); err != nil {
				return []string{}, fmt.Errorf("GCCIncluder error: open stdout pipe error %s\n", err)
			}
		}

		var scanner_out *bufio.Scanner
		scanner_out = bufio.NewScanner(stdout)

		if err = cmd.Start(); err != nil {
			return []string{}, fmt.Errorf("Result analyzer error: start command error %s\n", err)
		}

		for scanner_out.Scan() {
			if tp == true {
				list = this.makeStringAnalysisDef(scanner_out.Text(), list)
			} else {
				list = this.makeStringAnalysis(scanner_out.Text(), list)
			}
		}
		if err = scanner_out.Err(); err != nil {
			fmt.Printf("Result analyzer error: stdout scanner error %s\n", err)
		}

		if err = cmd.Wait(); err != nil {
			return []string{}, fmt.Errorf("Result analyzer error: wait command error %s\n", err)
		}
		return list, nil
	}
	return []string{}, fmt.Errorf("Result analyzer error: empty check command\n")
}

func (this *GCCContainer) GetGCCIncludes(config *configparser.ConfigparserContainer) {
	if cpp, err := this.propogateCmd("echo | /usr/bin/cpp -x c++ -Wp,-v -", config, false); err != nil {
		fmt.Println("Error in cpp default includes ", err)
		os.Exit(1)
	} else {
		this.cpp = cpp
	}
	if c, err := this.propogateCmd("echo | /usr/bin/gcc -v -x c -E -", config, false); err != nil {
		fmt.Println("Error in c default includes ", err)
		os.Exit(1)
	} else {
		this.c = c
	}
	if cpp, err := this.propogateCmd("/usr/bin/cpp -x c++ -dM -E - < /dev/null", config, false); err != nil {
		fmt.Println("Error in cpp default includes ", err)
		os.Exit(1)
	} else {
		this.cpp_h = cpp
	}
	if c, err := this.propogateCmd("/usr/bin/gcc -dM -E - < /dev/null", config, true); err != nil {
		fmt.Println("Error in c default includes ", err)
		os.Exit(1)
	} else {
		this.c_h = c
	}
}

func (this *GCCContainer) GetC() []string {
	return this.c
}

func (this *GCCContainer) GetCPP() []string {
	return this.cpp
}

func (this *GCCContainer) MakeHeaders(autoinclude string) (string, string) {
	if _, err := os.Stat(cFile); err != nil {
		file, err := os.Create(cFile)
		if err != nil {
			fmt.Printf("Can't create file %s: %s", cFile, err)
			os.Exit(1)
		}
		defer file.Close()
		for _, value := range this.c_h {
			file.WriteString(value + "\n")
		}
	}
	if _, err := os.Stat(cppFile); err != nil {
		file, err := os.Create(cppFile)
		if err != nil {
			fmt.Printf("Can't create file %s: %s", cppFile, err)
			os.Exit(1)
		}
		defer file.Close()
		for _, value := range this.cpp_h {
			file.WriteString(value + "\n")
		}
	}
	return autoinclude + cFile, autoinclude + cppFile
}
