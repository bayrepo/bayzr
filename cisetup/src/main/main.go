package main

import (
	"bufio"
	"bytes"
	"data"
	"executor"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runner"
	"server"
	"sonarapi"
	"strconv"
	"strings"
	"time"
)

var db_passwd = "sonarPASS1234"

type ActivationSaverFunc func(this interface{}) error

type ActionSaver struct {
	name         string
	createAction ActivationSaverFunc
	removeAction ActivationSaverFunc
	removed      bool
	paramsList   map[string]interface{}
}

func MakeActionSaver(nm string, cA ActivationSaverFunc, rA ActivationSaverFunc) *ActionSaver {
	return &ActionSaver{nm, cA, rA, true, map[string]interface{}{}}
}

func (this *ActionSaver) Activate() error {
	fmt.Println("Begin execution stage: " + this.name)
	err_act := this.createAction(this)
	if err_act == nil {
		this.removed = false
		return nil
	} else {
		fmt.Println(err_act)
	}
	err_deact := this.removeAction(this)
	if err_deact != nil {
		fmt.Println("Wow I even can't remove previous actions for " + this.name + ". Looks like it is seriously trouble")
		os.Exit(1)
	}
	this.removed = true
	return err_act
}

func (this *ActionSaver) Deactivate() {
	if this.removed == false {
		fmt.Println("Begin removing stage: " + this.name)
		err_deact := this.removeAction(this)
		if err_deact != nil {
			fmt.Println("Wow I even can't remove previous actions for " + this.name + ". Looks like it is seriously trouble")
			os.Exit(1)
		}
		this.removed = true
	}
}

func (this *ActionSaver) SetParam(param interface{}, name string) {
	this.paramsList[name] = param
}

func (this *ActionSaver) GetParam(name string) interface{} {
	if _, ok := this.paramsList[name]; ok == true {
		return this.paramsList[name]
	}
	return nil
}

func welcomeMessage() int {
	var response string
	fmt.Println("Welcome to management script of continue integration system")
	fmt.Println("Chose what should to do next:")
	fmt.Println("  1) Install all needed components")
	fmt.Println("  2) Create backups of installed components")
	fmt.Println("  3) Restore from backup(bkp.zip)")
	fmt.Println("1/2/3 or q(quit) or e(exit) for next step")
	for {
		_, err := fmt.Scanf("%s", &response)
		if err != nil {
			fmt.Println("Input error")
			os.Exit(255)
		}
		formattedResponse := strings.ToLower(response)
		if formattedResponse == "q" || formattedResponse == "quit" || formattedResponse == "e" || formattedResponse == "exit" {
			os.Exit(0)
		}
		result, err_fmt := strconv.Atoi(formattedResponse)
		if err_fmt != nil || result < 1 || result > 3 {
			fmt.Println("Answer should be 1/2/3 or q(quit) or e(exit) for next step")
		} else {
			return result
		}
	}
}

func makeArgsFromStrings(cmd ...string) []string {
	return append([]string{}, cmd...)
}

func makeArgsFromString(cmd string) []string {
	return strings.Split(cmd, " ")
}

func executeCommand(command []string) (error, int, []string, []string) {
	var cmd *exec.Cmd
	var stderr_txt []string
	var stdout_txt []string
	fmt.Println("Ready to execute command: ", command)

	start_cmd := append([]string{"/usr/bin/bash", "-c"}, strings.Join(command, " "))
	if len(start_cmd) == 1 {
		cmd = exec.Command(start_cmd[0])
	} else {
		cmd = exec.Command(start_cmd[0], start_cmd[1:]...)
	}
	var err error
	var stdout io.ReadCloser
	var stderr io.ReadCloser
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return err, 255, stderr_txt, stdout_txt
	}
	if stderr, err = cmd.StderrPipe(); err != nil {
		return err, 255, stderr_txt, stdout_txt
	}
	scanner_out := bufio.NewScanner(stdout)
	scanner_err := bufio.NewScanner(stderr)

	go func() {
		for scanner_err.Scan() {
			stderr_txt = append(stderr_txt, scanner_err.Text())
			fmt.Println(scanner_err.Text())
		}
		if err = scanner_err.Err(); err != nil {
			fmt.Printf("stderr scanner error %s\n", err)
		}
	}()

	if err = cmd.Start(); err != nil {
		return err, 255, []string{}, []string{}
	}

	for scanner_out.Scan() {
		stdout_txt = append(stdout_txt, scanner_out.Text())
		fmt.Println(scanner_out.Text())
	}

	if err = scanner_out.Err(); err != nil {
		fmt.Printf("stdout scanner error %s\n", err)
	}

	if err = cmd.Wait(); err != nil {
		return err, 255, stderr_txt, stdout_txt
	}
	return nil, 0, stderr_txt, stdout_txt
}

/* BEGIN: wget install */
func CheckFirstPackageSet(parent interface{}) error {
	this := parent.(*ActionSaver)
	need_install := true
	need_install_u := true
	need_install_y := true
	err, _, _, stdout_l := executeCommand(makeArgsFromString("rpm -qa wget"))
	if err != nil {
		return err
	}
	for _, val := range stdout_l {
		need_install = (strings.Contains(val, "wget") == false)
		if need_install == false {
			break
		}
	}
	this.SetParam(need_install, "need_install")
	if need_install == true {
		err, _, _, stdout_l = executeCommand(makeArgsFromString("yum install wget -y"))
		if err != nil {
			return err
		}
	}
	err, _, _, stdout_l = executeCommand(makeArgsFromString("rpm -qa unzip"))
	if err != nil {
		return err
	}
	for _, val := range stdout_l {
		need_install_u = (strings.Contains(val, "unzip") == false)
		if need_install_u == false {
			break
		}
	}
	this.SetParam(need_install_u, "need_install_u")
	if need_install_u == true {
		err, _, _, stdout_l = executeCommand(makeArgsFromString("yum install unzip -y"))
		if err != nil {
			return err
		}
	}
	err, _, _, stdout_l = executeCommand(makeArgsFromString("rpm -qa yum-utils"))
	if err != nil {
		return err
	}
	for _, val := range stdout_l {
		need_install_y = (strings.Contains(val, "yum-utils") == false)
		if need_install_y == false {
			break
		}
	}
	this.SetParam(need_install_y, "need_install_y")
	if need_install_y == true {
		err, _, _, stdout_l = executeCommand(makeArgsFromString("yum install yum-utils -y"))
		if err != nil {
			return err
		}
	}
	return nil
}

func RemoveFirstpackageSet(parent interface{}) error {
	this := parent.(*ActionSaver)
	if var_status := this.GetParam("need_install"); var_status != nil {
		need_install := var_status.(bool)
		if need_install == true {
			err, _, _, _ := executeCommand(makeArgsFromString("rpm -e wget"))
			if err != nil {
				return err
			}
		}
	}
	if var_status := this.GetParam("need_install_u"); var_status != nil {
		need_install := var_status.(bool)
		if need_install == true {
			err, _, _, _ := executeCommand(makeArgsFromString("rpm -e unzip"))
			if err != nil {
				return err
			}
		}
	}
	if var_status := this.GetParam("need_install_y"); var_status != nil {
		need_install := var_status.(bool)
		if need_install == true {
			err, _, _, _ := executeCommand(makeArgsFromString("rpm -e yum-utils"))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/* END: wget install */

/* BEGIN: example install */
func ExampleActionInstall(parent interface{}) error {
	//this := parent.(*ActionSaver)
	return nil
}

func ExampleActionDelete(parent interface{}) error {
	//this := parent.(*ActionSaver)
	return nil
}

/* END: example install */

/* BEGIN: jenkins install */
func deactivationCommonCmd(this *ActionSaver, param_name string, command string) error {
	if var_status := this.GetParam(param_name); var_status != nil {
		jenkins_param := var_status.(string)
		if jenkins_param == "success" {
			err, _, _, _ := executeCommand(makeArgsFromString(command))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func JenkinsActionInstall(parent interface{}) error {
	this := parent.(*ActionSaver)
	this.SetParam("ok", "jenkins.repo")
	this.SetParam("ok", "jenkins")
	this.SetParam("ok", "jenkins_systemd")
	this.SetParam("ok", "jenkins_start")
	/*err, _, _, _ := executeCommand(makeArgsFromString("wget -O /etc/yum.repos.d/jenkins.repo http://pkg.jenkins-ci.org/redhat/jenkins.repo"))
	if err != nil {
		return err
	}
	this.SetParam("success", "jenkins.repo")
	err, _, _, _ = executeCommand(makeArgsFromString("rpm --import https://jenkins-ci.org/redhat/jenkins-ci.org.key"))
	if err != nil {
		return err
	}*/
	err, _, _, _ := executeCommand(makeArgsFromString("yum install java-1.8.0-openjdk java-1.8.0-openjdk-headless java-1.8.0-openjdk-devel -y"))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("yum install chkconfig -y"))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("alternatives --install /usr/bin/java java /opt/jdk1.8.0_101/bin/java 2"))
	if err != nil {
		return err
	}
	/*err, _, _, _ = executeCommand(makeArgsFromString("yum install jenkins -y"))
	if err != nil {
		return err
	}
	this.SetParam("success", "jenkins")
	err, _, _, _ = executeCommand(makeArgsFromString("systemctl enable jenkins"))
	if err != nil {
		return err
	}
	this.SetParam("success", "jenkins_systemd")
	err, _, _, _ = executeCommand(makeArgsFromString("service jenkins start"))
	if err != nil {
		return err
	}
	this.SetParam("success", "jenkins_start")
	time.Sleep(20000 * time.Millisecond)*/
	err, _, _, _ = executeCommand(makeArgsFromString("firewall-cmd --zone=public --add-port=8080/tcp --permanent"))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("firewall-cmd --zone=public --add-port=9000/tcp --permanent"))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("firewall-cmd --zone=public --add-port=11000/tcp --permanent"))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("firewall-cmd --zone=public --add-service=http --permanent"))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("firewall-cmd --reload"))
	if err != nil {
		return err
	}

	/*err1, _, _, stdout_l := executeCommand(makeArgsFromString("cat /var/lib/jenkins/secrets/initialAdminPassword"))
	if err1 != nil {
		return err1
	}
	if len(stdout_l) > 0 {
		fmt.Println("Jenkins initial password: ", stdout_l[0])
	} else {
		return fmt.Errorf("Can't catch Jenkins password")
	}*/
	return nil
}

func JenkinsActionDelete(parent interface{}) error {
	this := parent.(*ActionSaver)
	if err := deactivationCommonCmd(this, "jenkins_start", "service jenkins stop"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "jenkins_systemd", "systemctl disable jenkins"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "jenkins", "yum erase jenkins -y"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "jenkins.repo", "rm -f /etc/yum.repos.d/jenkins.repo"); err != nil {
		return err
	}
	return nil
}

/* END: jenkins install */

/* BEGIN: sonarqube install */
func replace_string(src []byte, what string, to string) []byte {
	var result = []byte{}
	what_b := []byte(what)
	to_b := []byte(to)
	i := bytes.Index(src, what_b)
	j := 0
	for i >= 0 {
		i = i + j
		result = append(result, src[j:i]...)
		result = append(result, to_b...)
		if i+len(what_b) < len(src) {
			j = i + len(what_b)
			i = bytes.Index(src[i+len(what_b):], what_b)

		} else {
			break
		}

	}
	result = append(result, src[j:]...)
	return result
}

func SonarActionInstall(parent interface{}) error {
	this := parent.(*ActionSaver)
	sonar_version := "5.6.1"
	this.SetParam("ok", "sonar_mv")
	this.SetParam("ok", "sonarsc_mv")
	this.SetParam("ok", "sonar_ln")
	this.SetParam("ok", "sonarsc_ln")
	this.SetParam("ok", "sonar_rpm")
	this.SetParam("ok", "sonar_mysqL-install")
	this.SetParam("ok", "sonar_mysqL-install-en")
	this.SetParam("ok", "sonar_mysqL-install-st")
	this.SetParam("ok", "sonar_mysqL-dbcreate")
	this.SetParam("ok", "sonar_service")
	this.SetParam("ok", "mysql_db")
	fmt.Println("Chose version of SonarQube will be installed:")
	fmt.Println("   1. SonarQube 6.0")
	fmt.Println("   2. SonarQube 5.6")
	fmt.Println("Press 1 or 2")
	sonar_ch := "1"
	fmt.Scan(&sonar_ch)
	if strings.Trim(sonar_ch, " \n\t") == "1" {
		sonar_version = "6.0"
	}
	fmt.Println("Set password for sonarqube database:")
	sonar_ch = "1"
	fmt.Scan(&db_passwd)
	db_passwd = strings.Trim(db_passwd, " \n\t")
	this.SetParam(sonar_version, "sonar_version")

	err, _, _, _ := executeCommand(makeArgsFromString("wget https://sonarsource.bintray.com/Distribution/sonar-scanner-cli/sonar-scanner-2.6.1.zip"))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("unzip sonar-scanner-2.6.1.zip"))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("wget https://sonarsource.bintray.com/Distribution/sonarqube/sonarqube-" + sonar_version + ".zip"))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("unzip sonarqube-" + sonar_version + ".zip"))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("mv -n sonarqube-" + sonar_version + " /usr/local"))
	if err != nil {
		return err
	}
	this.SetParam("success", "sonar_mv")
	err, _, _, _ = executeCommand(makeArgsFromString("ln -s /usr/local/sonarqube-" + sonar_version + "/ /usr/local/sonar"))
	if err != nil {
		return err
	}
	this.SetParam("success", "sonar_ln")
	err, _, _, _ = executeCommand(makeArgsFromString("mv -n sonar-scanner-2.6.1 /usr/local"))
	if err != nil {
		return err
	}
	this.SetParam("success", "sonarsc_mv")
	err, _, _, _ = executeCommand(makeArgsFromString("ln -s /usr/local/sonar-scanner-2.6.1 /usr/local/sonar-scanner"))
	if err != nil {
		return err
	}
	this.SetParam("success", "sonarsc_ln")
	err, _, _, _ = executeCommand(makeArgsFromString("rpm -ihv http://dev.mysql.com/get/mysql57-community-release-el7-8.noarch.rpm"))
	if err != nil {
		return err
	}
	this.SetParam("success", "sonar_rpm")
	err, _, _, _ = executeCommand(makeArgsFromString("yum-config-manager --enable mysql56-community"))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("yum-config-manager --disable mysql57-community"))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("yum -y install mysql-community-common mysql-community-bench mysql-community-server mysql-connector-java mysql-connector-odbc mysql-community-client mysql-connector-python mysql-community-libs mysql-community-devel"))
	if err != nil {
		return err
	}
	this.SetParam("success", "sonar_mysqL-install")
	err, _, _, _ = executeCommand(makeArgsFromString("systemctl enable mysqld"))
	if err != nil {
		return err
	}
	this.SetParam("success", "sonar_mysqL-install-en")
	err, _, _, _ = executeCommand(makeArgsFromString("systemctl start mysqld"))
	if err != nil {
		return err
	}
	this.SetParam("success", "sonar_mysqL-install-st")
	this.SetParam("success", "mysql_db")
	err, _, _, _ = executeCommand(makeArgsFromString("mysql -e \"CREATE DATABASE sonar CHARACTER SET utf8 COLLATE utf8_general_ci;\""))
	if err != nil {
		return err
	}
	this.SetParam("success", "sonar_mysqL-dbcreate")
	err, _, _, _ = executeCommand(makeArgsFromString("mysql -e \"CREATE USER 'sonar' IDENTIFIED BY '" + db_passwd + "';\""))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("mysql -e \"GRANT ALL ON sonar.* TO 'sonar'@'%' IDENTIFIED BY '" + db_passwd + "';\""))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("mysql -e \"GRANT ALL ON sonar.* TO 'sonar'@'localhost' IDENTIFIED BY '" + db_passwd + "';\""))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("mysql -e \"FLUSH PRIVILEGES;\""))
	if err != nil {
		return err
	}

	data1, err := data.Asset("../cisetup/src/data/sonar")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/init.d/sonar", data1, 0755)
	if err != nil {
		return err
	}

	data2_tmp, err := data.Asset("../cisetup/src/data/sonar.properties")
	if err != nil {
		return err
	}
	data2 := replace_string(data2_tmp, "sonarPASSSWD", db_passwd)
	err = ioutil.WriteFile("/usr/local/sonar/conf/sonar.properties", data2, 0644)
	if err != nil {
		return err
	}
	
	data2_2_tmp, err := data.Asset("../cisetup/src/data/sonar.properties")
	if err != nil {
		return err
	}
	data2_2 := replace_string(data2_2_tmp, "sonarPASSSWD", db_passwd)
	err = ioutil.WriteFile("/usr/local/sonar-scanner/conf/sonar-scanner.properties", data2_2, 0644)
	if err != nil {
		return err
	}

	err, _, _, _ = executeCommand(makeArgsFromString("chkconfig --add sonar"))
	if err != nil {
		return err
	}
	this.SetParam("success", "sonar_service")
	
	err, _, _, _ = executeCommand(makeArgsFromString("/sbin/chkconfig sonar on"))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/usr/local/sonar/conf/.pass", []byte(db_passwd), 0400)
	if err != nil {
		return err
	}

	return nil
}

func SonarActionDelete(parent interface{}) error {
	this := parent.(*ActionSaver)
	sonar_version := "5.6.1"
	if var_status := this.GetParam("need_install"); var_status != nil {
		sonar_version = var_status.(string)
	}

	if err := deactivationCommonCmd(this, "sonar_service", "chkconfig --del sonar"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "sonar_mysqL-dbcreate", "mysql -e \"DROP DATABASE sonar;\""); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "sonar_mysqL-install-st", "systemctl stop mysqld"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "sonar_mysqL-install-en", "systemctl disable mysqld"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "sonar_mysqL-install", "yum erase -y mysql-community-common mysql-community-bench mysql-community-server mysql-connector-java mysql-connector-odbc mysql-community-client mysql-connector-python mysql-community-libs mysql-community-devel"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "sonar_rpm", "yum erase -y mysql57-community-release"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "mysql_db", "rm -rf /var/lib/mysql"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "sonarsc_ln", "rm -rf /usr/local/sonar-scanner"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "sonarsc_mv", "rm -rf /usr/local/sonar-scanner-2.6.1"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "sonar_ln", "rm -rf /usr/local/sonar"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "sonar_mv", "rm -rf /usr/local/sonarqube-"+sonar_version); err != nil {
		return err
	}
	return nil
}

/* END: sonarqube install */

/* BEGIN: bayzr install */
func BayZRActionInstall(parent interface{}) error {
	this := parent.(*ActionSaver)
	this.SetParam("ok", "bayzr.repo")
	this.SetParam("ok", "pkginstall")
	this.SetParam("ok", "sonarstart")
	this.SetParam("ok", "bzrc")
	err, _, _, _ := executeCommand(makeArgsFromString("wget -O /etc/yum.repos.d/home:bayrepo.repo http://download.opensuse.org/repositories/home:/bayrepo/CentOS_7/home:bayrepo.repo"))
	if err != nil {
		return err
	}
	this.SetParam("success", "bayzr.repo")
	err, _, _, _ = executeCommand(makeArgsFromString("yum install -y auto-buildrequires yumbootstrap bayzr"))
	if err != nil {
		return err
	}
	this.SetParam("success", "pkginstall")
	err, _, _, _ = executeCommand(makeArgsFromString("cp /usr/share/bzr.java/bayzr-plugin-0.0.1-rel1.jar /usr/local/sonar/extensions/plugins/"))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("service sonar start"))
	if err != nil {
		return err
	}

	err, _, _, _ = executeCommand(makeArgsFromString("mysql -e \"CREATE DATABASE bayzr CHARACTER SET utf8 COLLATE utf8_general_ci;\""))
	if err != nil {
		return err
	}
	this.SetParam("success", "bzrc")
	err, _, _, _ = executeCommand(makeArgsFromString("mysql -e \"CREATE USER 'bayzr' IDENTIFIED BY '" + db_passwd + "';\""))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("mysql -e \"GRANT ALL ON bayzr.* TO 'bayzr'@'%' IDENTIFIED BY '" + db_passwd + "';\""))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("mysql -e \"GRANT ALL ON bayzr.* TO 'bayzr'@'localhost' IDENTIFIED BY '" + db_passwd + "';\""))
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("mysql -e \"FLUSH PRIVILEGES;\""))
	if err != nil {
		return err
	}

	this.SetParam("success", "sonarstart")
	return nil
}

func BayZRActionDelete(parent interface{}) error {
	this := parent.(*ActionSaver)
	if err := deactivationCommonCmd(this, "sonarstart", "service sonar stop"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "pkginstall", "yum erase auto-buildrequires yumbootstrap bayzr -y"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "bayzr.repo", "rm -rf /etc/yum.repos.d/home:bayrepo.repo"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "bzrc", "mysql -e \"DROP DATABASE bayzr;\""); err != nil {
		return err
	}
	return nil
}

/* END: bayzr install */

/* BEGIN: plugins install */
func PluginActionInstall(parent interface{}) error {
	this := parent.(*ActionSaver)
	this.SetParam("ok", "pylint")
	fmt.Println("Waiting for sonarqube start")
	time.Sleep(180000 * time.Millisecond)
	sa := sonarapi.MakeSonarApi("http://127.0.0.1:9000", "admin", "admin")
	err := sa.Connect()
	if err != nil {
		return err
	}
	aPlugins, err := sa.GetPluginList(true)
	if err != nil {
		return err
	}
	iPlugins, err := sa.GetPluginList(false)
	if err != nil {
		return err
	}

	if err := sa.InstallPlugin("genericcoverage", aPlugins, iPlugins, true); err != nil {
		return err
	}
	if err := sa.InstallPlugin("python", aPlugins, iPlugins, true); err != nil {
		return err
	}
	if err := sa.InstallPlugin("l10nru", aPlugins, iPlugins, true); err != nil {
		return err
	}
	if err := sa.InstallPlugin("widgetlab", aPlugins, iPlugins, true); err != nil {
		return err
	}
	if err := sa.InstallPlugin("scmgit", aPlugins, iPlugins, true); err != nil {
		return err
	}

	if err := sa.InstallPlugin("csharp", aPlugins, iPlugins, false); err != nil {
		return err
	}
	if err := sa.InstallPlugin("java", aPlugins, iPlugins, false); err != nil {
		return err
	}
	if err := sa.InstallPlugin("javascript", aPlugins, iPlugins, false); err != nil {
		return err
	}

	err, _, _, _ = executeCommand(makeArgsFromString("yum -y install pylint"))
	if err != nil {
		return err
	}
	this.SetParam("success", "pylint")

	err, _, _, _ = executeCommand(makeArgsFromString("service sonar stop"))
	if err != nil {
		return err
	}
	fmt.Println("Waiting for sonarqube stop")
	time.Sleep(20000 * time.Millisecond)
	err, _, _, _ = executeCommand(makeArgsFromString("service sonar start"))
	if err != nil {
		return err
	}
	fmt.Println("Waiting for sonarqube start")
	time.Sleep(180000 * time.Millisecond)

	if err := sa.SetSonarOption("sonar.python.pylint", "/usr/bin/pylint"); err != nil {
		return err
	}
	if err := sa.SetSonarOption("sonar.python.pylint.reportPath", "pylint_report.txt"); err != nil {
		return err
	}
	if err := sa.SetSonarOption("sonar.bayzr.pass", db_passwd); err != nil {
		return err
	}

	return nil
}

func PluginActionDelete(parent interface{}) error {
	this := parent.(*ActionSaver)
	if err := deactivationCommonCmd(this, "pylint", "yum erase pylint -y"); err != nil {
		return err
	}
	return nil
}

/* END: plugins install */

/* BEGIN: squid and yumbootstrap install */
func SquidActionInstall(parent interface{}) error {
	this := parent.(*ActionSaver)
	this.SetParam("ok", "suqid")
	this.SetParam("ok", "squid_en")
	this.SetParam("ok", "squid_st")
	err, _, _, _ := executeCommand(makeArgsFromString("yum install -y squid"))
	if err != nil {
		return err
	}
	this.SetParam("success", "squid")
	err, _, _, _ = executeCommand(makeArgsFromString("systemctl enable squid"))
	if err != nil {
		return err
	}
	this.SetParam("success", "squid_en")

	data0, err := data.Asset("../cisetup/src/data/squid.conf")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/squid/squid.conf", data0, 0640)
	if err != nil {
		return err
	}
	err, _, _, _ = executeCommand(makeArgsFromString("systemctl start squid"))
	if err != nil {
		return err
	}
	this.SetParam("success", "squid_st")

	data2, err := data.Asset("../cisetup/src/data/centos-7-mod.suite")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/yumbootstrap/suites/centos-7-mod.suite", data2, 0644)
	if err != nil {
		return err
	}

	data3, err := data.Asset("../cisetup/src/data/centos-7-mod.list")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/yumbootstrap/suites/packages/centos-7-mod.list", data3, 0644)
	if err != nil {
		return err
	}

	data4, err := data.Asset("../cisetup/src/data/repomd.xml.key")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/yumbootstrap/suites/gpg/repomd.xml.key", data4, 0644)
	if err != nil {
		return err
	}

	data5, err := data.Asset("../cisetup/src/data/addbayzr.py")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/yumbootstrap/suites/scripts/addbayzr.py", data5, 0755)
	if err != nil {
		return err
	}

	return nil
}

func SquidActionDelete(parent interface{}) error {
	this := parent.(*ActionSaver)
	if err := deactivationCommonCmd(this, "squid_st", "systemctl stop squid"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "squid_en", "systemctl disable squid"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "squid", "yum erase squid -y"); err != nil {
		return err
	}
	return nil
}

/* END: squid and yumbootstrap install */

/* BEGIN: ciserver install */
func CiActionInstall(parent interface{}) error {
	this := parent.(*ActionSaver)
	this.SetParam("ok", "ciserver")
	err, _, _, _ := executeCommand(makeArgsFromString("yum install -y ciserver"))
	if err != nil {
		return err
	}
	this.SetParam("success", "ciserver")
	err, _, _, _ = executeCommand(makeArgsFromString("systemctl enable ciserver"))
	if err != nil {
		return err
	}
	this.SetParam("success", "ciserver_en")

	data0, err := data.Asset("../cisetup/src/data/sudoers")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/sudoers", data0, 0440)
	if err != nil {
		return err
	}

	err, _, _, _ = executeCommand(makeArgsFromString("setcap cap_sys_chroot+ep /usr/sbin/chroot"))
	if err != nil {
		return err
	}
	
	err, _, _, _ = executeCommand(makeArgsFromString("setcap cap_sys_chroot+ep /usr/sbin/citool"))
	if err != nil {
		return err
	}


	err, _, _, _ = executeCommand(makeArgsFromString("systemctl start ciserver"))
	if err != nil {
		return err
	}
	this.SetParam("success", "ciserver_st")

	return nil
}

func CiActionDelete(parent interface{}) error {
	this := parent.(*ActionSaver)
	if err := deactivationCommonCmd(this, "ciserver_st", "systemctl stop ciserver"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "ciserver_en", "systemctl disable ciserver"); err != nil {
		return err
	}
	if err := deactivationCommonCmd(this, "ciserver", "yum erase ciserver -y"); err != nil {
		return err
	}
	return nil
}

/* END: ciserver install */

var serverRun *bool
var jobRunner uint64
var taskRunner uint64
var taskRun *bool

func init() {
	serverRun = flag.Bool("server-run", false, "Start program as server")
	flag.Uint64Var(&jobRunner, "job-runner", 0, "Number of simultaneous tasks")
	flag.Uint64Var(&taskRunner, "task", 0, "Number if task id to execute")
	taskRun = flag.Bool("task-run", false, "Task flag(internal use)")
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("    ciutil [options] cmd ...\n")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {

	if *serverRun == true {
		var rnr runner.CiRunner
		var srv server.CiServer
		if err := srv.PreRun("/etc/citool.ini"); err != nil {
			fmt.Println(err)
			os.Exit(255)
		}
		err := rnr.SelfRun("/etc/citool.ini")
		defer rnr.KillSelfRun()
		if err == nil {
			if err := srv.Run(11000, "/etc/citool.ini"); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
		os.Exit(0)
	}

	if jobRunner > 0 {
		var rnr runner.CiRunner
		rnr.SetRunners(int64(jobRunner))
		rnr.Run("/etc/citool.ini")
		os.Exit(0)
	}

	if taskRunner > 0 {
		log.Printf("Running task %d", taskRunner)
		var ex executor.CiExec
		ex.Run(int(taskRunner), "/etc/citool.ini")
		os.Exit(0)
	}
	
	if *taskRun == true {
	    os.Exit(0)
	}

	actionsList := []*ActionSaver{}
	result := welcomeMessage()
	switch result {
	case 1:
		actionsList = append(actionsList, MakeActionSaver("Initial packages list installation", CheckFirstPackageSet, RemoveFirstpackageSet))
		actionsList = append(actionsList, MakeActionSaver("Install Java 8 and Jenkins", JenkinsActionInstall, JenkinsActionDelete))
		actionsList = append(actionsList, MakeActionSaver("Install SonarQube", SonarActionInstall, SonarActionDelete))
		actionsList = append(actionsList, MakeActionSaver("Install BayZR", BayZRActionInstall, BayZRActionDelete))
		actionsList = append(actionsList, MakeActionSaver("Install SonarQube plugins", PluginActionInstall, PluginActionDelete))
		actionsList = append(actionsList, MakeActionSaver("Install Squid and YumBootsTrap plugins", SquidActionInstall, SquidActionDelete))
		actionsList = append(actionsList, MakeActionSaver("Install Ci Server", CiActionInstall, CiActionDelete))
		break
	case 2:
		break
	case 3:
		break
	}
	for _, action := range actionsList {
		if err := action.Activate(); err != nil {
			fmt.Println("==================================================It is looks like error happened==========================================================")
			for i := len(actionsList) - 1; i >= 0; i-- {
				actionsList[i].Deactivate()
			}
			break
		}
	}
	os.Exit(0)
}
