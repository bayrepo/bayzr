package runner

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vaughan0/go-ini"
	"log"
	"mysqlsaver"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type CiRunner struct {
	ci_threads      int64
	config          string
	ci_busy_threads int64
	ci_pids         []int
	nmb             string
	ci_time         int
	cmd             *exec.Cmd
}

func (this *CiRunner) readConfig(ini_file string) error {
	config_data, err := ini.LoadFile(ini_file)
	if err != nil {
		return err
	}
	config_tmp, ok := config_data.Get("mysql", "connect")
	if !ok {
		return fmt.Errorf("Can't read MySQL connect parameters")
	}
	this.config = config_tmp
	config_tmp, ok = config_data.Get("server", "workers")
	if !ok {
		return fmt.Errorf("Can't read workers number")
	}
	this.nmb = config_tmp
	config_tmp, ok = config_data.Get("server", "wait")
	if !ok {
		return fmt.Errorf("Can't read wait time")
	}
	nmb, err := strconv.Atoi(config_tmp)
	if err != nil {
		nmb = 30
		log.Printf("Use default 30 second wait period")
	}
	this.ci_time = nmb
	return nil
}

func (this *CiRunner) SetRunners(nmb int64) {
	this.ci_threads = nmb
}

func (this *CiRunner) MakeFakeOuptut(id int, message string) error {
	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		return fmt.Errorf("DataBase saving error %s\n", dbErr)
	}

	defer con.Finalize()

	err := con.InsertOutput(id, message)
	return err
}

func (this *CiRunner) MakeChild(id int) error {
	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		return fmt.Errorf("DataBase saving error %s\n", dbErr)
	}
	defer con.Finalize()

	err_db := con.TakeJob(id)
	if err_db != nil {
		return err_db
	}
	cmd := os.Args[0]
	binary, lookErr := exec.LookPath(cmd)
	if lookErr != nil {
		con.CompleteJob(id)
		return lookErr
	}

	fstdin := os.Stdin
	fstdout := os.Stdout
	fstderr := os.Stderr

	argv := []string{binary, fmt.Sprintf("-task=%d", id), "-task-run"}
	procAttr := syscall.ProcAttr{
		Dir:   ".",
		Files: []uintptr{fstdin.Fd(), fstdout.Fd(), fstderr.Fd()},
		Env:   []string{},
		Sys: &syscall.SysProcAttr{
			Foreground: false,
			Setpgid:    true,
			Pgid:       0,
		},
	}

	pid, err := syscall.ForkExec(binary, argv, &procAttr)
	err_msg := fmt.Sprintf("Started task %d for id %d", pid, id)
	if err != nil {
		con.InsertOutput(id, err_msg)
		con.CompleteJob(id)
	}
	this.ci_pids = append(this.ci_pids, pid)
	log.Println(err_msg)
	return err
}

func (this *CiRunner) GetNextId() (int, error) {
	id := 0
	var err error
	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		return 0, fmt.Errorf("DataBase saving error %s\n", dbErr)
	}

	defer con.Finalize()

	err, id = con.DetJobID()
	return id, err
}

func (this *CiRunner) RemovePID(pid int) {
	result := []int{}
	for _, val := range this.ci_pids {
		if val != pid {
			result = append(result, val)
		}
	}
	this.ci_pids = result
}

func (this *CiRunner) _scanStartedChilds(found_threads bool) (error, bool) {
	pid_to_remove := []int{}
	for _, pid := range this.ci_pids {

		var wstat syscall.WaitStatus
		pid_ret, err := syscall.Wait4(pid, &wstat, syscall.WNOHANG, nil)
		if err != nil {
			return err, found_threads
		}
		if pid_ret == pid && wstat.Exited() {
			this.ci_busy_threads -= 1
			pid_to_remove = append(pid_to_remove, pid)
			found_threads = true
		}

	}
	for _, pid := range pid_to_remove {
		this.RemovePID(pid)
	}
	return nil, found_threads
}

func (this *CiRunner) Run(conf string) error {
	this.ci_pids = []int{}
	this.ci_busy_threads = 0
	err := this.readConfig(conf)
	if err != nil {
		return err
	}
	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		return fmt.Errorf("DataBase saving error %s\n", dbErr)
	}

	con.Finalize()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("Got signal %s", sig.String())
			os.Exit(0)
		}
	}()

	found_threads := false
	for {
		if this.ci_busy_threads < this.ci_threads {
			id, err := this.GetNextId()
			if err != nil {
				if err2 := this.MakeFakeOuptut(id, fmt.Sprintf("%s", err.Error())); err2 != nil {
					log.Printf("Got error: %s", err2.Error())
					os.Exit(1)
				}
			}
			if id > 0 {
				if err = this.MakeChild(id); err != nil {
					if err2 := this.MakeFakeOuptut(id, fmt.Sprintf("%s", err.Error())); err2 != nil {
						log.Printf("Got error: %s", err2.Error())
						os.Exit(1)
					}
				}
				this.ci_busy_threads += 1
				found_threads = true
			} else {
				if err, found_threads = this._scanStartedChilds(found_threads); err != nil {
					return err
				}
			}
		} else {
			if err, found_threads = this._scanStartedChilds(found_threads); err != nil {
				return err
			}
		}
		if found_threads == false {
			time.Sleep(time.Duration(this.ci_time) * time.Second)
		}
		found_threads = false
	}

}

func (this *CiRunner) SelfRun(conf string) error {
	err := this.readConfig(conf)
	if err != nil {
		return err
	}
	this.cmd = exec.Command(os.Args[0], "-job-runner", this.nmb)
	go func() {
		for {
			log.Printf("Start job runner")

			err := this.cmd.Start()
			if err != nil {
				log.Fatal(err)
			}
			err = this.cmd.Wait()
			log.Printf("Command finished with error: %v", err)
			if err != nil {
				os.Exit(2)
			}
		}
	}()
	return nil
}

func (this *CiRunner) KillSelfRun() {
	if this.cmd.Process != nil {
		this.cmd.Process.Kill()
	}
}
