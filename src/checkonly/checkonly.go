package checkonly

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CheckOnly struct {
	list_of_conditions []string
}

func (this CheckOnly) String() string {
	return fmt.Sprintf("list of check only cond: %s\n", strings.Join(this.list_of_conditions, ", "))
}

func Make_CheckOnly(lst []string) *CheckOnly {
	return &CheckOnly{lst}
}

func isDir(pth string) bool {
	fi, err := os.Stat(pth)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func ReadFirst100Byte(path string, pattern string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 100)
	_, err_read := f.Read(buf)
	if err_read != nil {
		return false
	}
	str := string(buf)
	return strings.ContainsAny(str, pattern)
}

func GetExtension(path string, pattern string) bool {
	ext := filepath.Ext(path)
	if ext == "" {
		return false
	}
	return strings.ContainsAny(ext, pattern[1:])
}

func (this *CheckOnly) Walk() []string {
	searchDir := "."

	fileList := []string{}
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if isDir(path) == false {
			is_fnd := false
			for _, val := range this.list_of_conditions {
				if val != "" {
					if val[0] == '.' {
						if GetExtension(path, val) == true {
							is_fnd = true
						}
					} else {
						if ReadFirst100Byte(path, val) == true {
							is_fnd = true
						}
					}
				}
			}
			if is_fnd == true {
				fileList = append(fileList, path)
			}
		}
		return nil
	})

	return fileList
}
