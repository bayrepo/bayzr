package executor

import (
	"bufio"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vaughan0/go-ini"
	"io"
	"io/ioutil"
	"mysqlsaver"
	"os/exec"
	"strings"
	"log"
)

const (
	chroot_path = "/mnt/chroot/"
)

type CiExec struct {
	ci_id  int
	config string
	con    mysqlsaver.MySQLSaver
}

func (this *CiExec) readConfig(ini_file string) error {
	config_data, err := ini.LoadFile(ini_file)
	if err != nil {
		return err
	}
	config_tmp, ok := config_data.Get("mysql", "connect")
	if !ok {
		return fmt.Errorf("Can't read MySQL connect parameters")
	}
	this.config = config_tmp

	return nil
}

func (this *CiExec) GetTaskInfo() (map[string]string, error) {
	err, result := this.con.GetTaskFullInfo(this.ci_id)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (this *CiExec) MakeFakeOuptut(message string) error {
	err := this.con.InsertOutput(this.ci_id, message)
	return err
}

func (this *CiExec) Exc(args []string) error {
	var clean_cmd string
	var cmd *exec.Cmd
	start_cmd := args
	if len(start_cmd) > 0 {
		clean_cmd = strings.Join(start_cmd, " ")
		start_cmd = append([]string{"/usr/bin/bash", "-c"}, strings.Join(start_cmd, " "))
		if len(start_cmd) == 1 {
			cmd = exec.Command(start_cmd[0])
		} else {
			cmd = exec.Command(start_cmd[0], start_cmd[1:]...)
		}
		var err error
		this.MakeFakeOuptut("+++:" + clean_cmd)
		var stdout io.ReadCloser
		var stderr io.ReadCloser
		if stdout, err = cmd.StdoutPipe(); err != nil {
			return fmt.Errorf("open stdout pipe error %s\n", err)
		}
		if stderr, err = cmd.StderrPipe(); err != nil {
			return fmt.Errorf("open stdout pipe error %s\n", err)
		}
		scanner_out := bufio.NewScanner(stdout)
		scanner_err := bufio.NewScanner(stderr)

		go func() {
			for scanner_err.Scan() {
				this.MakeFakeOuptut("Error:" + scanner_err.Text())
			}
		}()

		if err = cmd.Start(); err != nil {
			return fmt.Errorf("start command error %s\n", err)
		}

		for scanner_out.Scan() {
			this.MakeFakeOuptut(scanner_out.Text())
		}

		if err = cmd.Wait(); err != nil {
			return fmt.Errorf("Analyzer error: wait command error %s\n", err)
		}
		return nil
	}
	return fmt.Errorf("empty compile command\n")
}

func (this *CiExec) Run(id int, conf string) error {
	err := this.readConfig(conf)
	if err != nil {
		return err
	}
	dbErr := this.con.Init(this.config, nil)
	if dbErr != nil {
		return fmt.Errorf("DataBase saving error %s\n", dbErr)
	}
	defer this.con.Finalize()
	defer this.con.CompleteJob(this.ci_id)

	this.ci_id = id

	taskInfo, err := this.GetTaskInfo()
	if err != nil {
		log.Printf("Task runner got error %s", err.Error())
		return err
	}

	err = this.Exc([]string{"/usr/bin/sudo", "/usr/sbin/yumbootstrap", "--verbose", "centos-7-mod", fmt.Sprintf("%scentos-7-mod.%d", chroot_path, id)})
	if err != nil {
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	defer func() {
		err = this.Exc([]string{"/usr/bin/sudo", "/usr/sbin/yumbootstrap", "--uninstall", "centos-7-mod", fmt.Sprintf("%scentos-7-mod.%d", chroot_path, id)})
		if err != nil {
			this.MakeFakeOuptut("Error: " + err.Error())
		}

	}()

	err = this.Exc([]string{"/usr/sbin/chroot", fmt.Sprintf("%scentos-7-mod.%d", chroot_path, id), "/bin/env", "-i", "HOME=/home/checker", "TERM=\"$TERM\"", "PS1='\\u:\\w\\$ '", "PATH=/bin:/usr/bin:/sbin:/usr/sbin", "/bin/bash", "--login", "+h"})
	if err != nil {
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	defer func() {
		err = this.Exc([]string{"exit"})
		if err != nil {
			this.MakeFakeOuptut("Error: " + err.Error())
		}

	}()

	packages_list_str := strings.Trim(taskInfo["pkgs"], " \n,")
	if len(packages_list_str) > 0 {
		packages_list_raw := strings.Split(packages_list_str, ",")
		packages_list := []string{}
		for _, val := range packages_list_raw {
			package_ := strings.Trim(val, " \n,")
			if package_ != "" {
				packages_list = append(packages_list, package_)
			}
		}
		if len(packages_list) > 0 {
			cmds := append([]string{"/usr/bin/yum", "-y", "--nogpgcheck"}, packages_list...)
			err = this.Exc(cmds)
			if err != nil {
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
		}
	}

	err = this.Exc([]string{"/usr/bin/mkdir", "chkdir"})
	if err != nil {
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	err = this.Exc([]string{"/usr/bin/cd", "chkdir"})
	if err != nil {
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	src_raw_str := strings.Trim(taskInfo["source"], " \n,")
	if len(src_raw_str) > 0 {
		src_raw_str_raw := strings.Split(src_raw_str, " ")
		git_list := []string{}
		for _, val := range src_raw_str_raw {
			item_ := strings.Trim(val, " \n,")
			if item_ != "" {
				git_list = append(git_list, item_)
			}
		}
		if len(git_list) > 0 {
			err = this.Exc(git_list)
			if err != nil {
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
			err = this.Exc([]string{"/usr/bin/cd", "*"})
			if err != nil {
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
			err = this.Exc([]string{"/usr/bin/git", "checkout", "-b", "checkcommit", strings.Trim(taskInfo["commit"], " \n")})
			if err != nil {
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
		}
	}

	task_keys_str := strings.Trim(taskInfo["task_name"], " \n,")
	task_type := strings.Trim(taskInfo["task_type"], " \n,")
	task_keys_raw := strings.Split(task_keys_str, ":")
	task_keys := []string{}
	if len(task_keys_raw) >= 3 {
		task_keys = []string{task_keys_raw[0], task_keys_raw[1], task_keys_raw[2]}
	} else if len(task_keys_raw) == 2 {
		task_keys = []string{task_keys_raw[0], task_keys_raw[1], task_keys_raw[1]}
	} else if len(task_keys_raw) == 1 {
		task_keys = []string{task_keys_raw[0], task_keys_raw[0], task_keys_raw[0]}
	} else {
		this.MakeFakeOuptut("Error: project has no name")
		return err
	}

	sona_config := []byte(fmt.Sprintf("sonar.projectKey=%s:%s\nsonar.projectName=%s\nsonar.projectVersion=%s\nsonar.sources=.\n",
		task_keys[0], task_keys[1], task_keys[0], task_keys[2]))
	err = ioutil.WriteFile("sonar-project.properties", sona_config, 0644)

	if err != nil {
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	err = this.Exc([]string{"/usr/bin/cat", "sonar-project.properties"})
	if err != nil {
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	sona_config = []byte(taskInfo["config"] + fmt.Sprintf(`
	[database]
     connecturl=%s
	`, this.config))
	err = ioutil.WriteFile("bzr.conf", sona_config, 0644)

	err = this.Exc([]string{"/usr/bin/cat", "bzr.conf"})
	if err != nil {
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	cmds_raw := strings.Split(taskInfo["cmds"], "\n")

	for _, val := range cmds_raw {
		cmd_macros := strings.Trim(val, " \n\t")
		cmd := strings.Replace(cmd_macros, "{{CHECK}}",
			fmt.Sprintf("/usr/bin/bayzr -build-author %s -build-name \"%s.%s\" cmd ", taskInfo["login"], taskInfo["task_name"], taskInfo["id"]),
			1)
		cmds := strings.Split(cmd, " ")
		cmds_no_empty := []string{}
		for _, item := range cmds {
			itm := strings.Trim(item, " \n\t")
			if item != "" {
				cmds_no_empty = append(cmds_no_empty, itm)
			}
		}
		if len(cmds_no_empty) > 0 {
			err = this.Exc(cmds_no_empty)
			if err != nil {
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
		}
	}

	sonar_tp := task_type
	if sonar_tp == "1" {
		err = this.Exc([]string{"/usr/local/sonar-scanner/bin/sonar-runner"})
		if err != nil {
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}
	}

	return nil
}
