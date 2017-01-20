package diffanalyzer

import (
	"bufio"
	"configparser"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type DiffAnalyzerContainer struct {
	patchFilesList            []string
	foundFilesList            []string
	foundFilesList_LineNumber map[string][]configparser.DiffAnalyzerContainer_x_y
}

func Make_DiffAnalyzerContainer() *DiffAnalyzerContainer {
	return &DiffAnalyzerContainer{[]string{}, []string{}, map[string][]configparser.DiffAnalyzerContainer_x_y{}}
}

func (this *DiffAnalyzerContainer) ParseFilesList(in string) {
	this.patchFilesList = configparser.SplitOwnLongSepNoTrimSep(in, []string{","})

	for _, fileName := range this.patchFilesList {
		file_source, err := os.Open(fileName)
		if err != nil {
			continue
		}
		defer file_source.Close()
		reader := bufio.NewReader(file_source)
		file_found_name := ""

		for {
			line, err := reader.ReadString('\n')
			if line != "" && len(line) > 2 {
				if (line[:3] == "+++" || line[:3] == "---") && len(line) > 3 {
					new_line := line[4:]
					result := configparser.SplitOwnLongSepNoTrimSep(new_line, []string{"\t", " "})
					if len(result) > 0 {
						path_file := strings.Trim(result[0], " \n\t")
						if path_file != "." && path_file != "/" && path_file != "\\" && path_file != "" {
							is_fnd := false
							for _, value := range this.foundFilesList {
								if value == path_file {
									is_fnd = true
								}
							}
							if is_fnd == false {
								result := strings.Split(path_file, "/")
								final_path := ""
								if len(result) > 0 {
									if len(result) == 1 {
										final_path = path_file
									} else {
										len_r := len(result)
										for i := len_r; i > 0; i-- {
											final_path = strings.Join(result, "/")
											if _, err := os.Stat(final_path); err == nil {
												break
											}
											final_path = ""
											result2 := result[1:]
											result = result2
										}
									}
									if final_path != "" {
										file_found_name = final_path
										this.foundFilesList = append(this.foundFilesList, file_found_name)
									}
								}
							}
						}
					}
				}
				if line[:3] == "@@ " {
					new_line := line[4:]
					result := configparser.SplitOwnLongSepNoTrimSep(new_line, []string{"\t", " "})
					if len(result) > 2 {
						result2 := configparser.SplitOwnLongSepNoTrimSep(result[1], []string{","})
						if len(result2) > 1 {
							beg, err1 := strconv.ParseInt(strings.Trim(result2[0], " \n\t"), 10, 64)
							if err1 != nil {
								continue
							}
							end, err2 := strconv.ParseInt(strings.Trim(result2[1], " \n\t"), 10, 64)
							if err2 != nil {
								continue
							}
							if file_found_name != "" {
								this.foundFilesList_LineNumber[file_found_name] = append(this.foundFilesList_LineNumber[file_found_name], configparser.DiffAnalyzerContainer_x_y{beg, beg + end})
							}
						} else if len(result2) == 1 {
						    beg = 0
							end, err2 := strconv.ParseInt(strings.Trim(result2[0], " \n\t"), 10, 64)
							if err2 != nil {
								continue
							}
							if file_found_name != "" {
								this.foundFilesList_LineNumber[file_found_name] = append(this.foundFilesList_LineNumber[file_found_name], configparser.DiffAnalyzerContainer_x_y{beg, beg + end})
							}
						}
					}
				}
			}
			if err != nil {
				break
			}
		}
	}
	this.foundFilesList = configparser.RemoveDuplicate(this.foundFilesList)
}

func (this DiffAnalyzerContainer) String() string {
	result := fmt.Sprintf("patchFilesList: %s\nfoundFilesList: %s\n", strings.Join(this.patchFilesList, ", "), strings.Join(this.foundFilesList, ", "))
	for fname, value := range this.foundFilesList_LineNumber {
		result += fmt.Sprintf("[%s]\n", fname)
		for _, value_inner := range value {
			result += fmt.Sprintf("     %s\n", value_inner.String())
		}
	}
	return result
}

func (this DiffAnalyzerContainer) GetFoundList() ([]string, map[string][]configparser.DiffAnalyzerContainer_x_y) {
	return this.foundFilesList, this.foundFilesList_LineNumber
}
