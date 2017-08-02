package executor

import (
	"bufio"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vaughan0/go-ini"
	"io"
	"io/ioutil"
	"log"
	"logger"
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
	ci_id           int
	config          string
	con             mysqlsaver.MySQLSaver
	build_id        string
	sIp             string
	days_for_delete int
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
		logger.LogString(err.Error())
		return err
	}
	config_tmp, ok := config_data.Get("mysql", "connect")
	if !ok {
		logger.LogString("Can't read MySQL connect parameters")
		return fmt.Errorf("Can't read MySQL connect parameters")
	}
	this.config = config_tmp

	config_tmp, ok = config_data.Get("server", "ip")
	if !ok {
		this.sIp = "xx.xx.xx.xx:yyyy"
	} else {
		this.sIp = config_tmp
	}

	this.days_for_delete = 30

	days_tmp, ok := config_data.Get("mysql", "clean")
	if ok {
		days_tmp_i, err_i := strconv.Atoi(days_tmp)
		if err_i == nil {
			this.days_for_delete = days_tmp_i
		}
	}

	return nil
}

func (this *CiExec) GetTaskInfo() (map[string]string, error) {
	err, result := this.con.GetTaskFullInfo(this.ci_id)
	if err != nil {
		logger.LogString(err.Error())
		return result, err
	}
	return result, nil
}

func (this *CiExec) MakeFakeOuptut(message string) error {
	err := this.con.InsertOutput(this.ci_id, message)
	if err != nil {
		logger.LogString(err.Error())
	}
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
		env = append(env, fmt.Sprintf("TASKID=%d", this.ci_id))
		env = append(env, fmt.Sprintf("PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/root/bin"))
		cmd.Env = env

		var err error
		this.MakeFakeOuptut("+++:" + clean_cmd)
		var stdout io.ReadCloser
		var stderr io.ReadCloser
		if stdout, err = cmd.StdoutPipe(); err != nil {
			logger.Log("open stdout pipe error %s\n", err.Error())
			return fmt.Errorf("open stdout pipe error %s\n", err)
		}
		if stderr, err = cmd.StderrPipe(); err != nil {
			logger.Log("open stdout pipe error %s\n", err.Error())
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
			logger.Log("start command error %s\n", err.Error())
			return fmt.Errorf("start command error %s\n", err)
		}

		for scanner_out.Scan() {
			this.MakeFakeOuptut(scanner_out.Text())
		}

		if err = cmd.Wait(); err != nil {
			logger.Log("wait command error %s\n", err.Error())
			return fmt.Errorf("wait command error %s\n", err)
		}
		return nil
	}
	logger.LogString("empty compile command\n")
	return fmt.Errorf("empty compile command\n")
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func (this *CiExec) Run(id int, conf string) error {
	this.ci_id = id
	err := this.readConfig(conf)
	if err != nil {
		logger.LogString(err.Error())
		return err
	}
	dbErr := this.con.Init(this.config, nil)
	if dbErr != nil {
		logger.Log("DataBase saving error %s\n", dbErr.Error())
		return fmt.Errorf("DataBase saving error %s\n", dbErr)
	}
	defer this.con.Finalize()
	defer this.con.CompleteJob(this.ci_id)

	taskInfo, err := this.GetTaskInfo()
	if err != nil {
		logger.Log("Task runner got error %s", err.Error())
		log.Printf("Task runner got error %s", err.Error())
		return err
	}

	if taskInfo["task_id"] == "-1" {
		dEr, oV, rV, eV, bV, jV, eR, eE, eB, eO := this.con.CleanMySQL(this.days_for_delete)
		if dEr != nil {
			this.con.UpdateJobState(this.ci_id, 1)
			this.MakeFakeOuptut("Error: " + dEr.Error())
			return dEr
		}
		this.MakeFakeOuptut(fmt.Sprintf("+++: Number of OUTPUTS items deleted %d", oV))
		this.MakeFakeOuptut(fmt.Sprintf("+++: Number of REPORTS items deleted %d", rV))
		this.MakeFakeOuptut(fmt.Sprintf("+++: Number of ERRORS items deleted %d", eV))
		this.MakeFakeOuptut(fmt.Sprintf("+++: Number of BUILDS items deleted %d", bV))
		this.MakeFakeOuptut(fmt.Sprintf("+++: Number of JOBS items deleted %d", jV))
		this.MakeFakeOuptut(fmt.Sprintf("+++: Number of empty REPORTS items deleted %d", eR))
		this.MakeFakeOuptut(fmt.Sprintf("+++: Number of empty ERRORS items deleted %d", eE))
		this.MakeFakeOuptut(fmt.Sprintf("+++: Number of empty BUILDS items deleted %d", eB))
		this.MakeFakeOuptut(fmt.Sprintf("+++: Number of empty OUTPUTS items deleted %d", eO))
		return nil
	}

	this.build_id = taskInfo["task_name"] + "." + taskInfo["id"]

	err = this.Exc([]string{"/usr/bin/sudo", "/usr/sbin/yumbootstrap", "--verbose", "centos-7-mod", fmt.Sprintf("%scentos-7-mod.%d", chroot_path, id)})
	if err != nil {
		this.con.UpdateJobState(this.ci_id, 1)
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	defer func() {
		err = this.Exc([]string{"/usr/bin/sudo", "/usr/sbin/yumbootstrap", "--uninstall", "centos-7-mod", fmt.Sprintf("%scentos-7-mod.%d", chroot_path, id)})
		if err != nil {
			this.MakeFakeOuptut("Error: " + err.Error())
		}

	}()

	d, err1 := syscall.Open("/", syscall.O_RDONLY, 0)
	if err1 != nil {
		this.con.UpdateJobState(this.ci_id, 1)
		this.MakeFakeOuptut("Error: " + err1.Error())
		return err1
	}
	defer syscall.Close(d)

	dir := fmt.Sprintf("%scentos-7-mod.%d", chroot_path, id)
	err = syscall.Chroot(dir)
	this.MakeFakeOuptut(fmt.Sprintf("+++: chroot %s", dir))
	if err != nil {
		this.con.UpdateJobState(this.ci_id, 1)
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}
	ci_id_unchange := this.ci_id
	chroot_dir_unchange := dir

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
		arch_path := fmt.Sprintf("%s/home/checker/%d.tar.gz", chroot_dir_unchange, ci_id_unchange)
		dst_path := fmt.Sprintf("/usr/share/citool/%d.tar.gz", ci_id_unchange)
		if _, err_p := os.Stat(arch_path); !os.IsNotExist(err_p) {
			copyFileContents(arch_path, dst_path)
		}
	}()

	err = os.Chdir("/home/checker")
	this.MakeFakeOuptut(fmt.Sprintf("+++: chdir to /home/checker"))
	if err != nil {
		this.con.UpdateJobState(this.ci_id, 1)
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
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
		}
	}

	pre_build_cmd_list_str := strings.Trim(taskInfo["pre_build_cmd"], " \n,")
	if len(pre_build_cmd_list_str) > 0 {
		pre_build_cmd_list_raw := strings.Split(pre_build_cmd_list_str, ",")
		pre_script := "#!/bin/bash\n\n"
		pre_script = pre_script + "export LANG=en_US.UTF-8\n"
		for _, val := range pre_build_cmd_list_raw {
			cmd_macros := strings.Trim(val, " \n\t")
			pre_script = pre_script + cmd_macros + "\n"
		}
		err = ioutil.WriteFile("/home/checker/pre_execute", []byte(pre_script), 0755)

		err = this.Exc([]string{"/usr/bin/cat", "/home/checker/pre_execute"})
		if err != nil {
			this.con.UpdateJobState(this.ci_id, 1)
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}

		err = this.Exc([]string{"/bin/sudo", "/home/checker/pre_execute"})
		if err != nil {
			this.con.UpdateJobState(this.ci_id, 1)
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}
	}

	err = os.Mkdir("chkdir", 0755)
	this.MakeFakeOuptut(fmt.Sprintf("+++: mkdir /home/checker/chkdir"))
	if err != nil {
		this.con.UpdateJobState(this.ci_id, 1)
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	err = os.Chdir("/home/checker/chkdir")
	this.MakeFakeOuptut(fmt.Sprintf("+++: chdir to /home/checker/chkdir"))
	if err != nil {
		this.con.UpdateJobState(this.ci_id, 1)
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	src_raw_str := strings.Trim(taskInfo["source"], " \n,")

	commits_list := strings.Split(strings.Trim(taskInfo["commit"], " \n"), ",")
	commit_last := ""
	commit_first := ""
	if len(commits_list) == 1 {
		commit_last = commits_list[0]
	} else if len(commits_list) > 1 {
		commit_last = commits_list[1]
		commit_first = commits_list[0]
	}

	copy_catalog := ""
	origin_catalog := ""

	if len(src_raw_str) > 0 {
		git_project := src_raw_str
		err = this.Exc([]string{"/usr/bin/git", "clone", git_project})
		if err != nil {
			this.con.UpdateJobState(this.ci_id, 1)
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}
		files, _ := ioutil.ReadDir("/home/checker/chkdir")
		for _, f := range files {
			f_name := "/home/checker/chkdir/" + f.Name()
			is_d, d_err := this.IsDirectory(f_name)
			if d_err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + d_err.Error())
				return err
			}
			if f_name != "." && f_name != ".." && is_d == true {
				copy_catalog = f_name + ".copy"
				origin_catalog = f_name
				err = os.Chdir(f_name)
				this.MakeFakeOuptut(fmt.Sprintf("+++: chdir to %s", f_name))
				if err != nil {
					this.con.UpdateJobState(this.ci_id, 1)
					this.MakeFakeOuptut("Error: " + err.Error())
					return err
				}
			}
		}

		if taskInfo["use_branch"] == "0" {
			err = this.Exc([]string{"/usr/bin/git", "checkout", "-b", "checkcommit", strings.Trim(commit_last, " \n")})
			if err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
		} else if taskInfo["use_branch"] == "2" && git_project != "" {
			err = this.Exc([]string{"/usr/bin/git", "fetch", git_project, strings.Trim(commit_last, " \n")})
			if err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
			err = this.Exc([]string{"/usr/bin/git", "checkout", "FETCH_HEAD"})
			if err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
		} else {
			err = this.Exc([]string{"/usr/bin/git", "checkout", "-b", "checkbranch", "remotes/origin/" + strings.Trim(commit_last, " \n")})
			if err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
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
		this.con.UpdateJobState(this.ci_id, 1)
		this.MakeFakeOuptut("Error: project has no name")
		return err
	}

	sona_config := []byte(strings.Replace(taskInfo["config"], "\r", "", -1) + fmt.Sprintf(`
[database]
connecturl=%s
	`, this.config))
	err = ioutil.WriteFile("/home/checker/bzr.conf", sona_config, 0644)

	err = this.Exc([]string{"/usr/bin/cat", "/home/checker/bzr.conf"})
	if err != nil {
		this.con.UpdateJobState(this.ci_id, 1)
		this.MakeFakeOuptut("Error: " + err.Error())
		return err
	}

	if origin_catalog != "" && copy_catalog != "" {
		err = this.Exc([]string{"/usr/bin/cp", "-Rv", origin_catalog, copy_catalog})
		if err != nil {
			this.con.UpdateJobState(this.ci_id, 1)
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}
	}

	need_diff := ""
	if taskInfo["diff"] == "1" {
		need_diff = "-diff /home/checker/patch_f.patch"
		if commit_first == "" {
			if taskInfo["use_branch"] == "2" {
				err = this.Exc([]string{"/usr/bin/git", "diff", "HEAD^", "HEAD", ">/home/checker/patch_f.patch"})
				if err != nil {
					this.con.UpdateJobState(this.ci_id, 1)
					this.MakeFakeOuptut("Error: " + err.Error())
					return err
				}
			} else {
				err = this.Exc([]string{"/usr/bin/git", "format-patch", "-1", strings.Trim(commit_last, " \n"), ">/home/checker/patch_f.patch"})
				if err != nil {
					this.con.UpdateJobState(this.ci_id, 1)
					this.MakeFakeOuptut("Error: " + err.Error())
					return err
				}
			}
		} else {
			err = this.Exc([]string{"/usr/bin/git", "diff", strings.Trim(commit_first, " \n"), strings.Trim(commit_last, " \n"), ">/home/checker/patch_f.patch"})
			if err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
		}
		err = this.Exc([]string{"/usr/bin/cat", "/home/checker/patch_f.patch"})
		if err != nil {
			this.con.UpdateJobState(this.ci_id, 1)
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}
	}

	bld_err, cov_id := this.con.CreateCoverityBuild(this.ci_id, origin_catalog)
	if bld_err != nil {
		this.con.UpdateJobState(this.ci_id, 1)
		this.MakeFakeOuptut("Error: " + bld_err.Error())
		return err
	}
	if cov_id > 0 {
		err = this.Exc([]string{"/usr/bin/tar", "zcf", fmt.Sprintf("/home/checker/%d.tar.gz", this.ci_id), origin_catalog})
		if err != nil {
			this.con.UpdateJobState(this.ci_id, 1)
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}
	}

	cmds_raw := strings.Split(strings.Replace(taskInfo["cmds"], "\r", "", -1), "\n")

	if len(cmds_raw) > 0 {
		cmd_script := "#!/bin/bash\n\n"
		cmd_script = cmd_script + "export LANG=en_US.UTF-8\n"
		for _, val := range cmds_raw {
			cmd_macros := strings.Trim(val, " \n\t")
			check_fnd := false
			if strings.Contains(cmd_macros, "{{CHECK}}") {
				check_fnd = true
			}
			cmd := strings.Replace(cmd_macros, "{{CHECK}}",
				fmt.Sprintf("/usr/bin/bayzr -debug-commands -build-author %s -build-name \"%s.%s\" %s cmd ", taskInfo["login"], taskInfo["task_name"], taskInfo["id"], need_diff),
				1)
			if strings.Contains(cmd_macros, "{{CHECK_CLEAR}}") {
				check_fnd = true
			}
			cmd = strings.Replace(cmd_macros, "{{CHECK_CLEAR}}",
				fmt.Sprintf("/usr/bin/bayzr -no-build -debug-commands -build-author %s -build-name \"%s.%s\" %s cmd ", taskInfo["login"], taskInfo["task_name"], taskInfo["id"], need_diff),
				1)

			if cmd != "" {
				cmd_script = cmd_script + cmd + "\n"
			}

			if check_fnd {
				cmd_script = cmd_script + "if [ $? -ne 0 ]; then\n"
				cmd_script = cmd_script + "exit 255\n"
				cmd_script = cmd_script + "fi\n"
			}
		}
		err = ioutil.WriteFile("cmd_execute", []byte(cmd_script), 0755)

		err = this.Exc([]string{"/usr/bin/cat", "cmd_execute"})
		if err != nil {
			this.con.UpdateJobState(this.ci_id, 1)
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}

		err = this.Exc([]string{"./cmd_execute"})
		if err != nil {
			this.con.UpdateJobState(this.ci_id, 1)
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}
	}

	if cov_id == 0 {
		this.MakeFakeOuptut("+++: Save result to " + taskInfo["result_file"])
		if _, err := os.Stat(taskInfo["result_file"]); err == nil {
			if err := this.con.InsertExtInfoFromResult(taskInfo["result_file"], taskInfo["task_name"]+"."+taskInfo["id"]); err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
		} else {
			this.MakeFakeOuptut("Error: " + taskInfo["result_file"] + " not found")
		}
	}

	sonar_tp := task_type
	if sonar_tp == "1" {

		if origin_catalog != "" && copy_catalog != "" {
			err = this.Exc([]string{"/usr/bin/mv", origin_catalog, copy_catalog + ".garbage"})
			if err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
			err = this.Exc([]string{"/usr/bin/cp", "-Rv", copy_catalog, origin_catalog})
			if err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
			err = os.Chdir(origin_catalog)
			this.MakeFakeOuptut(fmt.Sprintf("+++: chdir to " + origin_catalog))
			if err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
		}

		if taskInfo["dir_to_execute"] != "" {
			err = os.Chdir(taskInfo["dir_to_execute"])
			this.MakeFakeOuptut(fmt.Sprintf("+++: chdir to " + taskInfo["dir_to_execute"]))
			if err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
		}

		sona_config_s := fmt.Sprintf("sonar.projectKey=%s:%s\nsonar.projectName=%s\nsonar.projectVersion=%s\nsonar.sources=.\nsonar.sourceEncoding=UTF-8\nsonar.import_unknown_files=true\n",
			task_keys[0], task_keys[1], task_keys[0], task_keys[2])

		err = ioutil.WriteFile("/home/checker/sonar-project.properties", []byte(sona_config_s), 0644)

		if err != nil {
			this.con.UpdateJobState(this.ci_id, 1)
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}

		err = this.Exc([]string{"/usr/bin/cat", "/home/checker/sonar-project.properties"})
		if err != nil {
			this.con.UpdateJobState(this.ci_id, 1)
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}

		err = this.Exc([]string{"/usr/local/sonar-scanner/bin/sonar-scanner", "-Dproject.settings=/home/checker/sonar-project.properties"})
		if err != nil {
			this.con.UpdateJobState(this.ci_id, 1)
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}
	} else {
		if taskInfo["dir_to_execute"] != "" {
			err = os.Chdir(taskInfo["dir_to_execute"])
			this.MakeFakeOuptut(fmt.Sprintf("+++: chdir to " + taskInfo["dir_to_execute"]))
			if err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
		}
	}

	post_cmds_raw := strings.Split(strings.Replace(taskInfo["post"], "\r", "", -1), "\n")

	if len(post_cmds_raw) > 0 {
		this.con.AddBuildInfoBeforePost(this.ci_id)
		err, result := this.con.GetJob(this.ci_id)
		if err != nil {
			this.con.UpdateJobState(this.ci_id, 1)
			this.MakeFakeOuptut("Error: " + err.Error())
			return err
		}
		if len(result) > 0 {

			build_id, err_i := strconv.Atoi(result[0][5])
			if err_i != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err_i.Error())
				return err_i
			}

			err, nmb_errors := this.con.GetBuildErrors(build_id)
			if err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}

			url_report := fmt.Sprintf("%s/result/%d", this.sIp, this.ci_id)
			url_output := fmt.Sprintf("%s/output/%d", this.sIp, this.ci_id)

			post_script := "#!/bin/bash\n\n"
			post_script = post_script + "export LANG=en_US.UTF-8\n"

			for _, val := range post_cmds_raw {
				cmd_macros := strings.Trim(val, " \n\t")
				post_script = post_script + cmd_macros + "\n"
			}
			err = ioutil.WriteFile("post_execute", []byte(post_script), 0755)

			err = this.Exc([]string{"/usr/bin/cat", "post_execute"})
			if err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}

			err = this.Exc([]string{"./post_execute", nmb_errors, url_report, url_output, fmt.Sprintf("%d", this.ci_id)})
			if err != nil {
				this.con.UpdateJobState(this.ci_id, 1)
				this.MakeFakeOuptut("Error: " + err.Error())
				return err
			}
		} else {
			this.MakeFakeOuptut("Error: no build info found")
		}

	}

	return nil
}
