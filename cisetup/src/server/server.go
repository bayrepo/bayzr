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
	html_resource, err := data.Asset("../cisetup/src/data/login.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("login", string(html_resource))
	html_resource, err = data.Asset("../cisetup/src/data/profile.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("profile", string(html_resource))
	html_resource, err = data.Asset("../cisetup/src/data/register.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("register", string(html_resource))
	html_resource, err = data.Asset("../cisetup/src/data/users.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("users", string(html_resource))
	html_resource, err = data.Asset("../cisetup/src/data/user.tpl")
	if err != nil {
		return templates, err
	}
	templates.AddFromString("user", string(html_resource))
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
	con.Finalize()

	router := gin.Default()

	if tpls, tpls_err := this.LoadTemplates(); tpls_err != nil {
		return tpls_err
	} else {
		router.HTMLRender = tpls
	}

	this.ses = sessions.NewCookieStore([]byte("bayzr-server"))
	router.Use(sessions.Sessions("bayzr-session", this.ses))
	router.Use(static.Serve("/css", BinaryFileSystem("../cisetup/src/data/css/")))
	router.Use(static.Serve("/js", BinaryFileSystem("../cisetup/src/data/js/")))
	router.Use(static.Serve("/fonts", BinaryFileSystem("../cisetup/src/data/fonts/")))

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

	err, st_login, st_name, st_email, _, st_grp := con.GetUserInfo(sess_id.(int))
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}
	hdr["User"] = st_login
	hdr["Name"] = st_name
	hdr["Email"] = st_email
	hdr["Group"] = st_grp
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

	err, st_login, st_name, st_email, st_grp_id, st_grp := con.GetUserInfo(uid_number)
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

	err, st_login, st_name, st_email, st_grp_id, st_grp = con.GetUserInfo(uid_number)
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
		return
	}

	hdr["UID"] = uid_number
	hdr["User_m"] = st_login
	hdr["Name"] = st_name
	hdr["Email"] = st_email
	hdr["Group"] = st_grp
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
