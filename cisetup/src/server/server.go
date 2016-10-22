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
	return templates, nil
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
		if c.Bind(&form) == nil {
			if err, fnd, perms, id := con.CheckUser(form.User, fmt.Sprintf("%x", md5.Sum([]byte(form.Password)))); err != nil {
				this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", err.Error()))
				return
			} else {
				if fnd {
					session.Set("login", form.User)
					session.Set("id", id)
					session.Set("perms", perms)
					session.Save()
					c.Redirect(http.StatusMovedPermanently, "/welcome/")
				} else {
					c.HTML(200, "login", gin.H{"FormUser": form.User,
						"ErrMSG": " Неверный логин или пароль "})
					return
				}
			}
		}
		c.HTML(200, "login", gin.H{})
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/welcome/")
}

type ProfileForm struct {
	InputName      string `form:"InputName" binding:"required"`
	InputEmail1    string `form:"InputEmail1" binding:"required"`
	InputPassword1 string `form:"InputPassword1" binding:"required"`
	InputPassword2 string `form:"InputPassword2" binding:"required"`
}

func (this *CiServer) welcome(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusMovedPermanently, "/")
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
	}
	err, st_login, st_name, st_email, _, st_grp := con.GetUserInfo(sess_id.(int))
	if err != nil {
		this.printSomethinWrong(c, 500, fmt.Sprintf("DataBase error %s\n", dbErr.Error()))
		return
	}
	c.HTML(200, "profile", gin.H{
		"User":      st_login,
		"Name":      st_name,
		"Email":     st_email,
		"Group":     st_grp,
		"Rules":     session.Get("perms").([]string),
		"TaskCount": 0,
		"JobCount":  0,
		"UserCount": 0})
}

func (this *CiServer) logout(c *gin.Context) {
	session := sessions.Default(c)
	sess_id := session.Get("id")
	if sess_id == nil {
		c.Redirect(http.StatusMovedPermanently, "/")
	}
	session.Delete("login")
	session.Delete("id")
	session.Delete("perms")
	session.Save()
	c.Redirect(http.StatusMovedPermanently, "/")
}
