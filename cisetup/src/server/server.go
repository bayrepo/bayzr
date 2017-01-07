package server

import (
	"crypto/md5"
	"data"
	"fmt"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vaughan0/go-ini"
	"html/template"
	"io/ioutil"
	"mysqlsaver"
	"net/http"
	"strconv"
	"strings"
)

//Вспомогательные функции
type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{data.Asset, data.AssetDir, data.AssetInfo, root}
	return &binaryFileSystem{
		fs,
	}
}

//Конец вспомогательных функций

type CiServer struct {
	config string
	ses    sessions.Store
}

func (this *CiServer) readConfig(ini_file string) error {
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

func (this *CiServer) LoadTemplates() (multitemplate.Render, error) {
	templates := multitemplate.New()
	html_resource, err := data.Asset("cisetup/src/data/login.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("login", string(html_resource))
	html_resource, err = data.Asset("cisetup/src/data/profile.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("profile", string(html_resource))
	html_resource, err = data.Asset("cisetup/src/data/register.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("register", string(html_resource))
	html_resource, err = data.Asset("cisetup/src/data/users.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("users", string(html_resource))
	html_resource, err = data.Asset("cisetup/src/data/user.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("user", string(html_resource))
	html_resource, err = data.Asset("cisetup/src/data/task.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("task", string(html_resource))
	html_resource, err = data.Asset("cisetup/src/data/tasks.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("tasks", string(html_resource))
	html_resource, err = data.Asset("cisetup/src/data/taske.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("taske", string(html_resource))
	html_resource, err = data.Asset("cisetup/src/data/jobs.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("jobs", string(html_resource))
	html_resource, err = data.Asset("cisetup/src/data/job.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("job", string(html_resource))
	html_resource, err = data.Asset("cisetup/src/data/out.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("out", string(html_resource))
	templates.AddFromString("result", "{{.Cont}}")
	return templates, nil
}

func (this *CiServer) usePerms(s sessions.Session, who_can []string, hdr *gin.H) error {
	sess_perms := s.Get("perms")
	if sess_perms == nil {
		return fmt.Errorf("Please login to access this page")
	}
	perms := sess_perms.([]string)
	allow := false
	for _, val := range perms {
		for _, can := range who_can {
			if val == can {
				allow = true
				break
			}
		}
		if allow {
			break
		}
	}
	if allow == false {
		return fmt.Errorf("You can't access this page due the permissions. Ask Administartor to increase the rights")
	}
	for _, val := range perms {
		(*hdr)[val] = "y"
	}
	return nil
}

func (this *CiServer) PreRun(conf string) error {
	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		return fmt.Errorf("DataBase saving error %s\n", dbErr)
	}

	defer con.Finalize()
	return nil
}

func (this *CiServer) Run(port int, conf string) error {
	err := this.readConfig(conf)
	if err != nil {
		return err
	}
	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		return fmt.Errorf("DataBase saving error %s\n", dbErr)
	}
	err_tock, tockens := con.GetUsersTokenLists()
	if err_tock != nil {
		return fmt.Errorf("DataBase saving error %s\n", err_tock)
	}

	con.Finalize()

	router := gin.Default()

	if tpls, tpls_err := this.LoadTemplates(); tpls_err != nil {
		return tpls_err
	} else {
		router.HTMLRender = tpls
	}

	this.ses = sessions.NewCookieStore([]byte("bayzr-server"))
	router.Use(sessions.Sessions("bayzr-session", this.ses))
	router.Use(static.Serve("/css", BinaryFileSystem("cisetup/src/data/css/")))
	router.Use(static.Serve("/js", BinaryFileSystem("cisetup/src/data/js/")))
	router.Use(static.Serve("/js/i18n", BinaryFileSystem("cisetup/src/data/js/i18n")))
	router.Use(static.Serve("/i18n", BinaryFileSystem("cisetup/src/data/js/i18n")))
	router.Use(static.Serve("/fonts", BinaryFileSystem("cisetup/src/data/fonts/")))

	router.GET("/", this.root)
	router.POST("/", this.root)

	router.GET("/welcome", this.welcome)
	router.POST("/welcome", this.welcome)

	router.GET("/logout", this.logout)

	router.GET("/register", this.register)
	router.POST("/register", this.register_post)

	router.GET("/users", this.users)
	router.GET("/users/:uid", this.user)
	router.POST("/users/:uid", this.user)

	router.GET("/user/del/:uid", this.userdel)

	router.GET("/tasks/add", this.tasks)
	router.POST("/tasks/add", this.tasks_post)

	router.GET("/tasks", this.tasks_all)

	router.GET("/task/:tid", this.tasks_edit)
	router.POST("/task/:tid", this.tasks_edit_post)

	router.GET("/taskdel/:tid", this.taskdel)

	router.GET("/procs/*page", this.jobs)
	router.GET("/jobs/*page", this.jobs)

	router.GET("/procs/add", this.newjob)
	router.POST("/procs/add", this.newjob_post)
	router.GET("/jobs/add", this.newjob)
	router.POST("/jobs/add", this.newjob_post)
	router.GET("/job/add", this.newjob)
	router.POST("/job/add", this.newjob_post)

	router.GET("/jobdel/:jid", this.jobdel)

	router.GET("/output/:oid", this.out)
	router.GET("/result/:rid", this.result)

	pass := gin.Accounts{}
	for _, item := range tockens {
		pass[item[0]] = item[1]
	}
	if len(pass) > 0 {
		authorized := router.Group("/api", gin.BasicAuth(pass))

		authorized.GET("/ping", this.ping)
		authorized.POST("/jobjson", this.addtjson)

	}
	router.Run(fmt.Sprintf(":%d", port))

	return nil
}

func (this *CiServer) printSomethinWrong(c *gin.Context, code int, message string) {
	c.String(code, "Oops! Please retry.")
	c.Error(fmt.Errorf(message))
	c.Abort()
}

type LoginForm struct {
	User     string `form:"InputLogin"`
	Password string `form:"InputPasswd"`
	Send     string `form:"send"`
	Register string `form:"register"`
}

func (this *CiServer) root(c *gin.Context) {
	session := sessions.Default(c)
	sess_login := session.Get("login")
	if sess_login == nil {
		var con mysqlsaver.MySQLSaver
		dbErr := con.Init(this.config, nil)
		if dbErr != nil {
			this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
			return
		}
		defer con.Finalize()
		var form LoginForm
		if c.Bind(&form) == nil && form.User != "" {
			if err, fnd, perms, id := con.CheckUser(form.User, fmt.Sprintf("%x", md5.Sum([]byte(form.Password)))); err != nil {
				this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
				return
			} else {
				if fnd {
					session.Set("login", form.User)
					session.Set("id", id)
					session.Set("perms", perms)
					session.Save()
					c.Redirect(http.StatusSeeOther, "/welcome/")
					return
				} else {
					c.HTML(200, "login", gin.H{"FormUser": form.User,
						"ErrMSG": " Неверный логин или пароль "})
					return
				}
			}
		}
		c.HTML(200, "login", gin.H{"ErrMSG": ""})
		return
	}
	c.Redirect(http.StatusSeeOther, "/welcome/")
}

type ProfileForm struct {
	InputName      string `form:"InputName" binding:"required"`
	InputEmail1    string `form:"InputEmail1" binding:"required"`
	InputPassword1 string `form:"InputPassword1"`
	InputPassword2 string `form:"InputPassword2"`
}

func (this *CiServer) welcome(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin", "ru_task", "ru_result", "ru_norules"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}
	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()
	var form ProfileForm
	if c.Bind(&form) == nil {
		found_err := false
		if form.InputName == "" {
			hdr["InputName_err"] = "Name can't be empty"
			found_err = true
		}
		if form.InputEmail1 == "" {
			hdr["InputEmail1_err"] = "Email can't be empty"
			found_err = true
		}
		if form.InputPassword1 != "" {
			if form.InputPassword1 != form.InputPassword2 {
				hdr["InputPassword2_err"] = "Repeat of password is not same the main password"
				found_err = true
			}

		}

		if found_err == false {
			var pass string
			if form.InputPassword1 != "" {
				pass = ""
			} else {
				pass = fmt.Sprintf("%x", md5.Sum([]byte(form.InputPassword1)))
			}

			if err := con.UpdateUser(sess_id.(int), form.InputName, form.InputEmail1,
				pass); err != nil {
				this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
				return
			}
		}
	}

	err, st_login, st_name, st_email, _, st_grp, st_tocken := con.GetUserInfo(sess_id.(int))
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}
	hdr["User"] = st_login
	hdr["Name"] = st_name
	hdr["Email"] = st_email
	hdr["Group"] = st_grp
	hdr["Tocken"] = st_tocken
	hdr["Rules"] = session.Get("perms").([]string)

	err, cnt := con.GetUsersCount()
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	hdr["TaskCount"] = 0
	hdr["JobCount"] = 0
	hdr["UserCount"] = cnt
	c.HTML(200, "profile", hdr)
}

func (this *CiServer) logout(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	session.Delete("login")
	session.Delete("id")
	session.Delete("perms")
	session.Save()
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

type RegisterForm struct {
	InputLogin     string `form:"InputLogin""`
	InputName      string `form:"InputName""`
	InputEmail1    string `form:"InputEmail1"`
	InputPassword1 string `form:"InputPassword1"`
	InputPassword2 string `form:"InputPassword2"`
}

func (this *CiServer) register(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/logout")
		return
	}

	c.HTML(200, "register", gin.H{})
}

func (this *CiServer) register_post(c *gin.Context) {
	session := sessions.Default(c)
	sess_log := session.Get("login")
	if sess_log != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/logout")
		return
	}
	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()
	var form RegisterForm
	hdr := gin.H{}
	if c.Bind(&form) == nil {
		found_err := false
		if form.InputLogin == "" {
			hdr["InputLogin_err"] = "Loging can't be empty"
			found_err = true
		} else {
			if err, res := con.IsDupUser(form.InputLogin); err != nil {
				this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
				return
			} else {
				if res {
					hdr["InputLogin_err"] = "Loging already in use"
					found_err = true
				}
				hdr["InputLogin"] = form.InputLogin
			}
		}
		if form.InputName == "" {
			hdr["InputName_err"] = "Name can't be empty"
			found_err = true
		} else {
			hdr["InputName"] = form.InputName
		}
		if form.InputEmail1 == "" {
			hdr["InputEmail1_err"] = "Email can't be empty"
			found_err = true
		} else {
			hdr["InputEmail1"] = form.InputEmail1
		}
		if form.InputPassword1 == "" {
			hdr["InputPassword1_err"] = "Password can't be empty"
			found_err = true
		} else {
			hdr["InputPassword1"] = form.InputPassword1
			if form.InputPassword1 != form.InputPassword2 {
				hdr["InputPassword2_err"] = "Repeat of password is not same the main password"
				found_err = true
			}
			hdr["InputPassword2"] = form.InputPassword2
		}

		if found_err == false {
			if err, id := con.InsertUser(form.InputLogin, form.InputName, form.InputEmail1,
				fmt.Sprintf("%x", md5.Sum([]byte(form.InputPassword1)))); err != nil {
				this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
				return
			} else {
				if id == 0 {
					this.printSomethinWrong(c, 500, fmt.Sprintf("Wow no ID for new User"))
					return
				}
				if err, perms := con.GetUserPerms(id); err != nil {
					this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
					return
				} else {
					session.Set("login", form.InputLogin)
					session.Set("id", id)
					session.Set("perms", perms)
					session.Save()
					c.Redirect(http.StatusSeeOther, "/welcome/")
					return
				}
			}
		}
	}

	c.HTML(200, "register", hdr)
}

func (this *CiServer) users(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}
	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	if err, _, users_list := con.GetUsersList(); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	} else {
		hdr["Users"] = users_list
	}
	hdr["User"] = session.Get("login").(string)

	c.HTML(200, "users", hdr)
}

type ProfileFormUser struct {
	InputName      string `form:"InputName" binding:"required"`
	InputEmail1    string `form:"InputEmail1" binding:"required"`
	InputPassword1 string `form:"InputPassword1"`
	InputPassword2 string `form:"InputPassword2"`
	GrpId          string `form:"GroupInput" binding:"required"`
}

func (this *CiServer) user(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}
	uid := c.Param("uid")
	uid_number, err := strconv.Atoi(uid)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/users/")
		return
	}

	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	err, st_login, st_name, st_email, st_grp_id, st_grp, st_tocken := con.GetUserInfo(uid_number)
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	var form ProfileFormUser
	if c.Bind(&form) == nil {
		found_err := false
		if form.InputName == "" {
			hdr["InputName_err"] = "Name can't be empty"
			found_err = true
		}
		if form.InputEmail1 == "" {
			hdr["InputEmail1_err"] = "Email can't be empty"
			found_err = true
		}
		if form.InputPassword1 != "" {
			if form.InputPassword1 != form.InputPassword2 {
				hdr["InputPassword2_err"] = "Repeat of password is not same the main password"
				found_err = true
			}

		}

		if found_err == false {
			var pass string
			if form.InputPassword1 != "" {
				pass = ""
			} else {
				pass = fmt.Sprintf("%x", md5.Sum([]byte(form.InputPassword1)))
			}

			if err := con.UpdateUser(uid_number, form.InputName, form.InputEmail1,
				pass); err != nil {
				this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
				return
			}

			if grp_id_num, err := strconv.Atoi(form.GrpId); err == nil {
				if grp_id_num != st_grp_id {
					if err := con.UpdateUserGroup(uid_number, grp_id_num); err != nil {
						this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
						return
					}
				}
			}

		}
	}

	err, st_login, st_name, st_email, st_grp_id, st_grp, st_tocken = con.GetUserInfo(uid_number)
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	hdr["UID"] = uid_number
	hdr["User_m"] = st_login
	hdr["Name"] = st_name
	hdr["Email"] = st_email
	hdr["Group"] = st_grp
	hdr["Tocken"] = st_tocken
	hdr["Group_ID"] = fmt.Sprintf("%d", st_grp_id)

	if err, grps := con.GetGroups(); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	} else {
		hdr["Groups"] = grps
	}
	hdr["User"] = session.Get("login").(string)

	c.HTML(200, "user", hdr)
}

func (this *CiServer) userdel(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}
	uid := c.Param("uid")
	uid_number, err := strconv.Atoi(uid)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/users/")
		return
	}

	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	if err := con.DelUser(uid_number); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/users")
}

func (this *CiServer) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "alive", "status": http.StatusOK})
}

//Example:
//curl -u su_admin:7c542b69b3f8af507e7808eed20b0d
// --data "user_token=7c542b69b3f8af507e7808eed20b0d&
//task_token=6fae706abecfbef0a8d1b3b7f775a7
//&commit=master&descr=my" http://127.0.0.1:11000/api/jobjson
func (this *CiServer) addtjson(c *gin.Context) {
	user_token := c.DefaultPostForm("user_token", "")
	task_token := c.DefaultPostForm("task_token", "")
	if user_token == "" || task_token == "" {
		c.JSON(http.StatusForbidden, gin.H{"status": "forbidden"})
		return
	}

	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	u_err, u_id, u_login := con.GetUserAuth(user_token)
	if u_err != nil {
		c.JSON(http.StatusOK, gin.H{"message": u_err, "status": http.StatusOK, "result": "error"})
		return
	}

	if u_id == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no user", "status": http.StatusOK, "result": "error"})
		return
	}

	t_err, t_id, t_ul := con.GetTaskAuth(task_token)
	if t_err != nil {
		c.JSON(http.StatusOK, gin.H{"message": u_err, "status": http.StatusOK, "result": "error"})
		return
	}

	if t_id == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no task", "status": http.StatusOK, "result": "error"})
		return
	}

	fnd_user := false
	for _, val := range t_ul {
		if val == u_login {
			fnd_user = true
			break
		}
	}

	if fnd_user == false {
		c.JSON(http.StatusForbidden, gin.H{"status": "forbidden"})
		return
	}

	job_name := "Job for " + u_login
	job_prior := "2"
	job_commit := c.DefaultPostForm("commit", "")
	job_task := t_id
	job_descr := c.DefaultPostForm("descr", "API BUILD")

	if err, j_id := con.InsertJob(u_id, job_name, job_commit, job_prior, job_task, job_descr); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err, "status": http.StatusOK, "result": "error"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "done", "status": http.StatusOK, "result": fmt.Sprintf("%d", j_id)})
	}
}

type TaskForm struct {
	TaskName       string   `form:"TaskName"`
	TaskType       string   `form:"TaskType"`
	TaskGit        string   `form:"TaskGit"`
	TaskPackGs     []string `form:"TaskPackGs"`
	TaskPackGsEarl []string `form:"TaskPackGsEarl"`
	TaskCmds       []string `form:"TaskCmds"`
	TaskPerType    string   `form:"TaskPerType"`
	TaskPeriod     string   `form:"TaskPeriod"`
	TaskUsers      []string `form:"TaskUsers"`
	TaskConfig     []string `form:"TaskConfig"`
	TaskBranch     string   `form:"TaskBranch"`
	TaskResult     string   `form:"TaskResult"`
	TaskBrn        string   `form:"TaskBrn"`
	TaskDiff       string   `form:"TaskDiff"`
	TaskPost       []string `form:"TaskPost"`
	TaskDir        string   `form:"TaskDir"`
}

func (this *CiServer) tasks(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin", "ru_task"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}

	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	err, pkg_lst := con.GetPackages()
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	p_pkg_list := [][]string{}
	for _, val := range pkg_lst {
		p_pkg_list = append(p_pkg_list, []string{val, ""})
	}

	err, _, u_list := con.GetUsersList()
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	p_u_list := [][]string{}
	for _, val := range u_list {
		p_u_list = append(p_u_list, []string{val[1], "", val[2]})
	}
	p_u_list = append(p_u_list, []string{"su_admin", "", "Admin"})

	example_config, err_reading := ioutil.ReadFile("/etc/bzr.conf")
	if err_reading != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("bzr.conf error reading %s\n", err_reading.Error()))
		return
	}

	hdr["TaskName"] = ""
	hdr["TaskType"] = "2"
	hdr["TaskGit"] = ""
	hdr["TaskBranch"] = "0"
	hdr["TaskPackGs"] = ""
	hdr["TaskPackGsEarl"] = p_pkg_list
	hdr["TaskCmds"] = ""
	hdr["TaskPerType"] = "5"
	hdr["TaskPeriod"] = ""
	hdr["TaskUsers"] = p_u_list
	hdr["TaskConfig"] = string(example_config)
	hdr["TaskBrn"] = ""
	hdr["TaskDiff"] = "n"
	hdr["TaskPost"] = ""
	hdr["TaskDir"] = ""

	hdr["User"] = session.Get("login").(string)

	c.HTML(200, "task", hdr)
}

func (this *CiServer) tasks_post(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin", "ru_task"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}

	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	err, pkg_lst := con.GetPackages()
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	p_pkg_list := [][]string{}
	for _, val := range pkg_lst {
		p_pkg_list = append(p_pkg_list, []string{val, ""})
	}

	err, _, u_list := con.GetUsersList()
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	p_u_list := [][]string{}
	for _, val := range u_list {
		p_u_list = append(p_u_list, []string{val[1], "", val[2]})
	}
	p_u_list = append(p_u_list, []string{"su_admin", "", "Admin"})

	hdr["TaskName"] = ""
	hdr["TaskType"] = "2"
	hdr["TaskGit"] = ""
	hdr["TaskBranch"] = "0"
	hdr["TaskPackGs"] = ""
	hdr["TaskPackGsEarl"] = p_pkg_list
	hdr["TaskCmds"] = ""
	hdr["TaskPerType"] = "5"
	hdr["TaskPeriod"] = ""
	hdr["TaskUsers"] = p_u_list
	hdr["TaskConfig"] = ""
	hdr["TaskResult"] = "result.html"
	hdr["TaskBrn"] = "result.html"
	hdr["TaskDiff"] = "n"
	hdr["TaskPost"] = ""
	hdr["TaskDir"] = ""

	var form TaskForm
	if c.Bind(&form) == nil {
		fnd_err := false
		if form.TaskName == "" {
			hdr["TaskName_err"] = "Task name can't be empty"
			fnd_err = true
		} else {
			hdr["TaskName"] = form.TaskName
		}
		hdr["TaskType"] = form.TaskType
		if form.TaskGit == "" {
			hdr["TaskGit_err"] = "Clonning command can't be empty"
			fnd_err = true
		} else {
			hdr["TaskGit"] = form.TaskGit
		}
		if len(form.TaskPackGs) > 0 {
			for _, val := range form.TaskPackGs {
				fnd := false
				for i, p_val := range p_pkg_list {
					if p_val[0] == val {
						p_pkg_list[i][1] = "selected"
						fnd = true
						break
					}
				}
				if fnd == false {
					p_pkg_list = append(p_pkg_list, []string{val, "selected"})
				}
			}
		}
		if len(form.TaskPackGsEarl) > 0 {
			for _, val := range form.TaskPackGsEarl {
				for i, p_val := range p_pkg_list {
					if p_val[0] == val {
						p_pkg_list[i][1] = "selected"
						break
					}
				}
			}
		}
		hdr["TaskPackGsEarl"] = p_pkg_list
		hdr["TaskCmds"] = strings.Join(form.TaskCmds, "\n")
		hdr["TaskPerType"] = form.TaskPerType
		if form.TaskPeriod == "" && form.TaskPerType != "5" {
			hdr["TaskPeriod_err"] = "Periid should by in format DD/MM/YYYY HH:MM:SS or type of period should be No period"
			fnd_err = true
		} else {
			hdr["TaskPeriod"] = form.TaskPeriod
		}
		if len(form.TaskUsers) > 0 {
			for _, val := range form.TaskUsers {
				for i, p_val := range p_u_list {
					if p_val[0] == val {
						p_u_list[i][1] = "selected"
						break
					}
				}
			}
		}
		hdr["TaskUsers"] = p_u_list

		if len(form.TaskConfig) == 0 {
			hdr["TaskConfig_err"] = "Config for analyzer can't be empty"
			fnd_err = true
		} else {
			hdr["TaskConfig"] = strings.Join(form.TaskConfig, "\n")
		}
		hdr["TaskBranch"] = form.TaskBranch

		if form.TaskResult == "" {
			hdr["TaskResult_err"] = "Task result file name can't be empty"
			fnd_err = true
		} else {
			hdr["TaskResult"] = form.TaskResult
		}
		hdr["TaskBrn"] = form.TaskBrn

		hdr["TaskDiff"] = form.TaskDiff
		hdr["TaskPost"] = strings.Join(form.TaskPost, "\n")
		hdr["TaskDir"] = form.TaskDir

		if fnd_err == false {
			i_taskBranch, _ := strconv.Atoi(hdr["TaskBranch"].(string))
			if err, _ := con.SaveTask(sess_id.(int), form.TaskName, form.TaskType, form.TaskGit,
				form.TaskPackGs, form.TaskPackGsEarl, form.TaskCmds, form.TaskPerType, form.TaskPeriod,
				form.TaskUsers, form.TaskConfig, i_taskBranch, form.TaskResult, form.TaskBrn, form.TaskDiff, form.TaskPost, form.TaskDir); err != nil {
				this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
				return
			} else {
				c.Redirect(http.StatusSeeOther, "/tasks/")
				return
			}
		}
	}

	hdr["User"] = session.Get("login").(string)

	c.HTML(200, "task", hdr)
}

func (this *CiServer) tasks_all(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin", "ru_task"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}

	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	if err, lst := con.GetTasks(); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	} else {
		hdr["Tasks"] = lst
	}

	hdr["User"] = session.Get("login").(string)

	c.HTML(200, "tasks", hdr)
}

func (this *CiServer) tasks_edit(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin", "ru_task"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}

	tid := c.Param("tid")
	tid_number, err := strconv.Atoi(tid)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/tasks/")
		return
	}

	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	err, pkg_lst := con.GetPackages()
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	p_pkg_list := [][]string{}
	for _, val := range pkg_lst {
		p_pkg_list = append(p_pkg_list, []string{val, ""})
	}

	err, _, u_list := con.GetUsersList()
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	p_u_list := [][]string{}
	for _, val := range u_list {
		p_u_list = append(p_u_list, []string{val[1], "", val[2]})
	}
	p_u_list = append(p_u_list, []string{"su_admin", "", "Admin"})

	err, lst := con.GetTask(tid_number)
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}

	if len(lst[10]) > 0 {
		for _, val := range strings.Split(lst[10], ",") {
			for i, p_val := range p_u_list {
				if p_val[0] == strings.Trim(val, " ") {
					p_u_list[i][1] = "selected"
					break
				}
			}
		}
	}

	if len(lst[4]) > 0 {
		for _, val := range strings.Split(lst[4], ",") {
			for i, p_val := range p_pkg_list {
				if p_val[0] == strings.Trim(val, " ") {
					p_pkg_list[i][1] = "selected"
					break
				}
			}
		}
	}

	hdr["TaskId"] = lst[0]
	hdr["TaskName"] = lst[1]
	hdr["TaskBranch"] = lst[12]
	hdr["TaskType"] = lst[2]
	hdr["TaskGit"] = lst[3]
	hdr["TaskPackGs"] = ""
	hdr["TaskPackGsEarl"] = p_pkg_list
	hdr["TaskCmds"] = lst[5]
	hdr["TaskPerType"] = lst[6]
	hdr["TaskPeriod"] = lst[7]
	hdr["TaskUsers"] = p_u_list
	hdr["TaskConfig"] = lst[9]
	hdr["TaskToken"] = lst[11]
	hdr["TaskResult"] = lst[13]
	hdr["TaskBrn"] = lst[14]
	hdr["TaskDiff"] = lst[15]
	hdr["TaskPost"] = lst[16]
	hdr["TaskDir"] = lst[17]

	hdr["User"] = session.Get("login").(string)

	c.HTML(200, "taske", hdr)
}

func (this *CiServer) tasks_edit_post(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin", "ru_task"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}

	tid := c.Param("tid")
	tid_number, err := strconv.Atoi(tid)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/tasks/")
		return
	}

	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	err, pkg_lst := con.GetPackages()
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	p_pkg_list := [][]string{}
	for _, val := range pkg_lst {
		p_pkg_list = append(p_pkg_list, []string{val, ""})
	}

	err, _, u_list := con.GetUsersList()
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	p_u_list := [][]string{}
	for _, val := range u_list {
		p_u_list = append(p_u_list, []string{val[1], "", val[2]})
	}
	p_u_list = append(p_u_list, []string{"su_admin", "", "Admin"})

	err, lst := con.GetTask(tid_number)
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}

	if len(lst[10]) > 0 {
		for _, val := range strings.Split(lst[10], ",") {
			for i, p_val := range p_u_list {
				if p_val[0] == strings.Trim(val, " ") {
					p_u_list[i][1] = "selected"
					break
				}
			}
		}
	}

	if len(lst[4]) > 0 {
		for _, val := range strings.Split(lst[4], ",") {
			for i, p_val := range p_pkg_list {
				if p_val[0] == strings.Trim(val, " ") {
					p_pkg_list[i][1] = "selected"
					break
				}
			}
		}
	}

	hdr["TaskId"] = lst[0]
	hdr["TaskName"] = lst[1]
	hdr["TaskBranch"] = lst[12]
	hdr["TaskType"] = lst[2]
	hdr["TaskGit"] = lst[3]
	hdr["TaskPackGs"] = ""
	hdr["TaskPackGsEarl"] = p_pkg_list
	hdr["TaskCmds"] = lst[5]
	hdr["TaskPerType"] = lst[6]
	hdr["TaskPeriod"] = lst[7]
	hdr["TaskUsers"] = p_u_list
	hdr["TaskConfig"] = lst[9]
	hdr["TaskToken"] = lst[11]
	hdr["TaskResult"] = lst[13]
	hdr["TaskBrn"] = lst[14]
	hdr["TaskDiff"] = lst[15]
	hdr["TaskPost"] = lst[16]
	hdr["TaskDir"] = lst[17]

	var form TaskForm
	if c.Bind(&form) == nil {
		fnd_err := false
		if form.TaskName == "" {
			hdr["TaskName_err"] = "Task name can't be empty"
			fnd_err = true
		} else {
			hdr["TaskName"] = form.TaskName
		}
		hdr["TaskType"] = form.TaskType
		if form.TaskGit == "" {
			hdr["TaskGit_err"] = "Clonning command can't be empty"
			fnd_err = true
		} else {
			hdr["TaskGit"] = form.TaskGit
		}
		if len(form.TaskPackGs) > 0 {
			for _, val := range form.TaskPackGs {
				fnd := false
				for i, p_val := range p_pkg_list {
					if p_val[0] == val {
						p_pkg_list[i][1] = "selected"
						fnd = true
						break
					}
				}
				if fnd == false {
					p_pkg_list = append(p_pkg_list, []string{val, "selected"})
				}
			}
		}
		if len(form.TaskPackGsEarl) > 0 {
			for _, val := range form.TaskPackGsEarl {
				for i, p_val := range p_pkg_list {
					if p_val[0] == val {
						p_pkg_list[i][1] = "selected"
						break
					}
				}
			}
		}
		hdr["TaskPackGsEarl"] = p_pkg_list
		hdr["TaskCmds"] = strings.Join(form.TaskCmds, "\n")
		hdr["TaskPerType"] = form.TaskPerType
		if form.TaskPeriod == "" && form.TaskPerType != "5" {
			hdr["TaskPeriod_err"] = "Periid should by in format DD/MM/YYYY HH:MM:SS or type of period should be No period"
			fnd_err = true
		} else {
			hdr["TaskPeriod"] = form.TaskPeriod
		}
		if len(form.TaskUsers) > 0 {
			for _, val := range form.TaskUsers {
				for i, p_val := range p_u_list {
					if p_val[0] == val {
						p_u_list[i][1] = "selected"
						break
					}
				}
			}
		}
		hdr["TaskUsers"] = p_u_list

		if len(form.TaskConfig) == 0 {
			hdr["TaskConfig_err"] = "Config for analyzer can't be empty"
			fnd_err = true
		} else {
			hdr["TaskConfig"] = strings.Join(form.TaskConfig, "\n")
		}

		hdr["TaskBranch"] = form.TaskBranch

		if form.TaskResult == "" {
			hdr["TaskResult_err"] = "Task result file name can't be empty"
			fnd_err = true
		} else {
			hdr["TaskResult"] = form.TaskResult
		}
		hdr["TaskBrn"] = form.TaskBrn
		hdr["TaskDiff"] = form.TaskDiff
		hdr["TaskPost"] = strings.Join(form.TaskPost, "\n")
		hdr["TaskDir"] = form.TaskDir

		if fnd_err == false {
			i_taskBranch, _ := strconv.Atoi(hdr["TaskBranch"].(string))
			if err := con.UpdateTask(tid_number, sess_id.(int), form.TaskName, form.TaskType, form.TaskGit,
				form.TaskPackGs, form.TaskPackGsEarl, form.TaskCmds, form.TaskPerType, form.TaskPeriod,
				form.TaskUsers, form.TaskConfig, i_taskBranch, form.TaskResult, form.TaskBrn, form.TaskDiff, form.TaskPost, form.TaskDir); err != nil {
				this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
				return
			} else {
				c.Redirect(http.StatusSeeOther, "/tasks/")
				return
			}
		}
	}

	hdr["User"] = session.Get("login").(string)

	c.HTML(200, "taske", hdr)
}

func (this *CiServer) taskdel(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin", "ru_task"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}
	tid := c.Param("tid")
	tid_number, err := strconv.Atoi(tid)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/tasks/")
		return
	}

	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	if err := con.DelTask(tid_number); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/tasks")
}

func (this *CiServer) jobs(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin", "ru_task", "ru_result"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}

	page := c.Param("page")
	page_number, err := strconv.Atoi(page)
	if err != nil {
		page_number = 0
	}

	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	hdr["Jobs"] = []string{}

	if err, lst, pg := con.GetJobs(page_number); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	} else {
		hdr["Jobs"] = lst
		hdr["Page"] = page_number
		page_numbers := []int{}
		for i := 1; i <= pg; i++ {
			page_numbers = append(page_numbers, i)
		}
		hdr["PageNmbrs"] = page_numbers
	}

	hdr["User"] = session.Get("login").(string)

	c.HTML(200, "jobs", hdr)
}

func (this *CiServer) newjob(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin", "ru_task", "ru_result"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}

	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	tasks := [][]string{}
	if err, lst := con.GetTasks(); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	} else {
		for _, val := range lst {
			tasks = append(tasks, []string{val[0], "", val[1]})
		}
	}

	hdr["JobName"] = ""
	hdr["JobPrior"] = "2"
	hdr["JobCommit"] = ""
	hdr["JobTask"] = tasks
	hdr["JobDescr"] = ""

	hdr["User"] = session.Get("login").(string)

	c.HTML(200, "job", hdr)
}

type JobForm struct {
	JobName   string `form:"JobName"`
	JobPrior  string `form:"JobPrior"`
	JobCommit string `form:"JobCommit"`
	JobTask   int    `form:"JobTask"`
	JobDescr  string `form:"JobDescr"`
}

func (this *CiServer) newjob_post(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin", "ru_task", "ru_result"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}

	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	tasks := [][]string{}
	if err, lst := con.GetTasks(); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	} else {
		for _, val := range lst {
			tasks = append(tasks, []string{val[0], "", val[1]})
		}
	}

	hdr["JobName"] = ""
	hdr["JobPrior"] = "2"
	hdr["JobCommit"] = ""
	hdr["JobTask"] = tasks
	hdr["JobDescr"] = ""

	var form JobForm
	if c.Bind(&form) == nil {
		fnd_err := false
		if form.JobName == "" {
			hdr["JobName_err"] = "Job name can't be empty"
			fnd_err = true
		} else {
			hdr["JobName"] = form.JobName
		}
		if form.JobCommit == "" {
			hdr["JobCommit_err"] = "Job commit name can't be empty"
			fnd_err = true
		} else {
			hdr["JobCommit"] = form.JobCommit
		}

		task_id_str := fmt.Sprintf("%d", form.JobTask)
		for i, val := range tasks {
			if val[0] == task_id_str {
				tasks[i][1] = "selected"
				break
			}
		}
		hdr["JobTask"] = tasks

		hdr["JobPrior"] = form.JobPrior

		hdr["JobDescr"] = form.JobDescr

		if fnd_err == false {
			if err, _ := con.InsertJob(sess_id.(int), form.JobName, form.JobCommit, form.JobPrior, form.JobTask, form.JobDescr); err != nil {
				this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
				return
			} else {
				c.Redirect(http.StatusSeeOther, "/jobs/")
				return
			}
		}
	}

	hdr["User"] = session.Get("login").(string)

	c.HTML(200, "job", hdr)
}

func (this *CiServer) jobdel(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin", "ru_task"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}
	jid := c.Param("jid")
	jid_number, err := strconv.Atoi(jid)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/procs/")
		return
	}

	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	if err := con.DelJob(jid_number); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/procs")
}

func (this *CiServer) out(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin", "ru_task", "ru_result"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}
	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	hdr["out"] = []string{}
	oid := c.Param("oid")
	oid_number, err := strconv.Atoi(oid)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/jobs/")
		return
	}

	if err, lst := con.GetOut(oid_number); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	} else {
		hdr["out"] = lst
	}

	hdr["User"] = session.Get("login").(string)

	c.HTML(200, "out", hdr)
}

func (this *CiServer) result(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	hdr := gin.H{}
	if err := this.usePerms(session, []string{"ru_admin", "ru_task", "ru_result"}, &hdr); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	}
	var con mysqlsaver.MySQLSaver
	dbErr := con.Init(this.config, nil)
	if dbErr != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	defer con.Finalize()

	oid := c.Param("rid")
	oid_number, err := strconv.Atoi(oid)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/jobs/")
		return
	}

	if err, lst := con.GetResult(oid_number); err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("%s\n", err.Error()))
		return
	} else {
		hdr := gin.H{}
		hdr["Cont"] = template.HTML(strings.Join(lst, "\n"))
		c.HTML(200, "result", hdr)
	}
}
