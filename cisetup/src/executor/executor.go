package executor

import (
	"bufio"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vaughan0/go-ini"
	"io"
	"io/ioutil"
	"log"
	"mysqlsaver"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

const (
	chroot_path = "/mnt/chroot/"
)

type CiExec struct {
	ci_id    int
	config   string
	con      mysqlsaver.MySQLSaver
	build_id string
	sIp      string
}

func (this *CiExec) IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err == nil {
		return fileInfo.IsDir(), err
	} else {
		return false, err
	}
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

	config_tmp, ok = config_data.Get("server", "ip")
	if !ok {
		this.sIp = "xx.xx.xx.xx:yyyy"
	} else {
		this.sIp = config_tmp
	}

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
		env := os.Environ()
		env = append(env, fmt.Sprintf("INJAIL=yes"))
		env = append(env, fmt.Sprintf("BUILDID=%s", this.build_id))
		env = append(env, fmt.Sprintf("PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/root/bin"))
		cmd.Env = env

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
			return fmt.Errorf("wait command error %s\n", err)
		}
		return nil
	}
	return fmt.Errorf("empty compile command\n")
}

func (this *CiExec) Run(id int, conf string) error {
	this.ci_id = id
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

	taskInfo, err := this.GetTaskInfo()
	if err != nil {
		log.Printf("Task runner got error %s", err.Error())
		return err
	}
	this.build_id = taskInfo["task_name"] + "." + taskInfo["id"]

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

	/*err = this.Exc([]string{"/usr/sbin/chroot", fmt.Sprintf("%scentos-7-mod.%d", chroot_path, id), "/bin/env", "-i", "HOME=/home/checker", "TERM=\"$TERM\"", "PS1='\\u:\\w\\$ '", "PATH=/bin:/usr/bin:/sbin:/usr/sbin", "/bin/bash", "--login", "+h"})
	if err != nil {
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	defer func() {
		err = this.Exc([]string{"exit"})
		if err != nil {
			this.MakeFakeOuptut("Error: " + err.Error())
		}

	}()*/

	d, err1 := syscall.Open("/", syscall.O_RDONLY, 0)
	if err1 != nil {
		this.MakeFakeOuptut("Error: " + err1.Error())
		return err1
	}
	defer syscall.Close(d)

	dir := fmt.Sprintf("%scentos-7-mod.%d", chroot_path, id)
	err = syscall.Chroot(dir)
	this.MakeFakeOuptut(fmt.Sprintf("+++: chroot %s", dir))
	if err != nil {
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	defer func() {
		this.MakeFakeOuptut("+++: Back to real root")
		err2 := syscall.Fchdir(d)
		if err2 != nil {
			this.MakeFakeOuptut("Error: " + err2.Error())
			return
		}
		err2 = syscall.Chroot(".")
		if err2 != nil {
			this.MakeFakeOuptut("Error: " + err2.Error())
			return
		}
	}()

	err = os.Chdir("/home/checker")
	this.MakeFakeOuptut(fmt.Sprintf("+++: chdir to /home/checker"))
	if err != nil {
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

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
			cmds := append([]string{"/bin/sudo", "/usr/bin/yum", "-y", "--nogpgcheck", "install"}, packages_list...)
			err = this.Exc(cmds)
			if err != nil {
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
		}
	}

	err = os.Mkdir("chkdir", 0755)
	this.MakeFakeOuptut(fmt.Sprintf("+++: mkdir /home/checker/chkdir"))
	if err != nil {
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	err = os.Chdir("/home/checker/chkdir")
	this.MakeFakeOuptut(fmt.Sprintf("+++: chdir to /home/checker"))
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
			files, _ := ioutil.ReadDir("/home/checker/chkdir")
			for _, f := range files {
				f_name := "/home/checker/chkdir/" + f.Name()
				is_d, d_err := this.IsDirectory(f_name)
				if d_err != nil {
					this.MakeFakeOuptut("Error: " + d_err.Error())
					return err
				}
				if f_name != "." && f_name != ".." && is_d == true {
					err = os.Chdir(f_name)
					this.MakeFakeOuptut(fmt.Sprintf("+++: chdir to %s", f_name))
					if err != nil {
						this.MakeFakeOuptut("Error: " + err.Error())
						return err
					}
				}
			}

			if taskInfo["use_branch"] == "0" {
				err = this.Exc([]string{"/usr/bin/git", "checkout", "-b", "checkcommit", strings.Trim(taskInfo["commit"], " \n")})
				if err != nil {
					this.MakeFakeOuptut("Error: " + err.Error())
					return err
				}
			} else {
				err = this.Exc([]string{"/usr/bin/git", "checkout", "-b", "checkbranch", "remotes/origin/" + strings.Trim(taskInfo["commit"], " \n")})
				if err != nil {
					this.MakeFakeOuptut("Error: " + err.Error())
					return err
				}
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

	sona_config = []byte(strings.Replace(taskInfo["config"], "\r", "", -1) + fmt.Sprintf(`
[database]
connecturl=%s
	`, this.config))
	err = ioutil.WriteFile("bzr.conf", sona_config, 0644)

	err = this.Exc([]string{"/usr/bin/cat", "bzr.conf"})
	if err != nil {
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	need_diff := ""
	if taskInfo["diff"] == "y" {
		need_diff = "-diff patch_f.patch"
		err = this.Exc([]string{"/usr/bin/git", "format-patch", "-1", strings.Trim(taskInfo["commit"], ">patch_f.patch")})
		if err != nil {
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}
		err = this.Exc([]string{"/usr/bin/cat", "patch_f.patch"})
		if err != nil {
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}
	}

	cmds_raw := strings.Split(strings.Replace(taskInfo["cmds"], "\r", "", -1), "\n")

	for _, val := range cmds_raw {
		cmd_macros := strings.Trim(val, " \n\t")
		cmd := strings.Replace(cmd_macros, "{{CHECK}}",
			fmt.Sprintf("/usr/bin/bayzr -build-author %s -build-name \"%s.%s\" %s cmd ", taskInfo["login"], taskInfo["task_name"], taskInfo["id"], need_diff),
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

	this.MakeFakeOuptut("+++: Save result to " + taskInfo["result_file"])
	if _, err := os.Stat(taskInfo["result_file"]); err == nil {
		if err := this.con.InsertExtInfoFromResult(taskInfo["result_file"], taskInfo["task_name"]+"."+taskInfo["id"]); err != nil {
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}
	} else {
		this.MakeFakeOuptut("Error: " + taskInfo["result_file"] + " not found")
	}

	sonar_tp := task_type
	if sonar_tp == "1" {
		err = this.Exc([]string{"/usr/local/sonar-scanner/bin/sonar-scanner"})
		if err != nil {
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}
	}

	post_cmds_raw := strings.Split(strings.Replace(taskInfo["post"], "\r", "", -1), "\n")

	if len(post_cmds_raw) > 0 {
		err, result := this.con.GetJob(this.ci_id)
		if err != nil {
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}
		if len(result) > 0 {

			build_id, err_i := strconv.Atoi(result[0][5])
			if err_i != nil {
				this.MakeFakeOuptut("Error: " + err_i.Error())
				return err
			}

			err, nmb_errors := this.con.GetBuildErrors(build_id)
			if err != nil {
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}

			url_report := fmt.Sprintf("%s/result/%d", this.sIp, build_id)
			url_output := fmt.Sprintf("%s/output/%d", this.sIp, this.ci_id)

			post_script := "#!/bin/bash\n\n"

			for _, val := range post_cmds_raw {
				cmd_macros := strings.Trim(val, " \n\t")
				post_script = post_script + cmd_macros + "\n"
			}
			err = ioutil.WriteFile("post_execute", []byte(post_script), 0755)

			err = this.Exc([]string{"/usr/bin/cat", "post_execute"})
			if err != nil {
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}

			err = this.Exc([]string{"./post_execute", nmb_errors, url_report, url_output})
			if err != nil {
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
		} else {
			this.MakeFakeOuptut("Error: no build info found")
		}

	}

	return nil
}
