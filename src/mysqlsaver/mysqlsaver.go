package mysqlsaver

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"rullerlist"
	"strconv"
	"strings"
)

type MySQLSaver struct {
	db               *sql.DB
	ok               int
	current_build_id int64
	ruller           *rullerlist.RullerList
}

func (this *MySQLSaver) getStringLength(str string, length int) string {
	if len(str) > length {
		return str[:length]
	}
	return str
}

func (this *MySQLSaver) _checkTable(tbl_name string) (bool, error) {
	var table_name string = ""
	stmtOut, err := this.db.Prepare("SHOW TABLES LIKE '" + tbl_name + "'")
	if err != nil {
		return false, err
	}
	defer stmtOut.Close()
	err = stmtOut.QueryRow().Scan(&table_name)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if table_name == tbl_name {
		return true, nil
	}
	return false, nil
}

func (this *MySQLSaver) _executeSQLCpammnd(cmd string) error {
	stmtIns, err := this.db.Prepare(cmd)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (this *MySQLSaver) _checkAndCreateTables() error {
	fnd, err := this._checkTable("bayzr_last_check")
	if err != nil {
		return err
	}
	if fnd == false {
		if err := this._executeSQLCpammnd(
			`CREATE TABLE IF NOT EXISTS bayzr_last_check(
		        id INTEGER  NOT NULL AUTO_INCREMENT, 
		        checker VARCHAR(255), 
		        last_build_id INTEGER, 
		        PRIMARY KEY (id),
		        UNIQUE KEY(checker))`); err != nil {
			return err
		}
	}
	fnd, err = this._checkTable("bayzr_build_info")
	if err != nil {
		return err
	}
	if fnd == false {
		if err := this._executeSQLCpammnd(
			`CREATE TABLE IF NOT EXISTS bayzr_build_info(
		        id INTEGER  NOT NULL AUTO_INCREMENT, 
		        build_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
		        autor_of_build VARCHAR(255), 
		        name_of_build VARCHAR(255),
		        completed int,
		        PRIMARY KEY (id),
                INDEX build_date_I (build_date),
                INDEX autor_of_build_I (autor_of_build(15)))`); err != nil {
			return err
		}
	}
	fnd, err = this._checkTable("bayzr_err_extend")
	if err != nil {
		return err
	}
	if fnd == false {
		if err := this._executeSQLCpammnd(
			`CREATE TABLE IF NOT EXISTS bayzr_err_extend(
		        id INTEGER  NOT NULL AUTO_INCREMENT, 
		        file_name VARCHAR(255),
		        file_string TEXT,
		        file_pos INTEGER,
		        rec_type INTEGER,
		        plugin VARCHAR(50),
		        err_type INTEGER,
		        descript TEXT,
		        build_number INTEGER,
		        PRIMARY KEY (id),
                INDEX file_name_I (file_name(32)),
                INDEX err_type_I (err_type),
                INDEX plugin_I (plugin(12)),
                INDEX rec_type_I (rec_type),
                INDEX build_number_I (build_number))`); err != nil {
			return err
		}
	}
	fnd, err = this._checkTable("bayzr_err_list")
	if err != nil {
		return err
	}
	if fnd == false {
		if err := this._executeSQLCpammnd(
			`CREATE TABLE IF NOT EXISTS bayzr_err_list(
		        id INTEGER  NOT NULL AUTO_INCREMENT, 
		        plugin VARCHAR(50),
		        bayzr_err INTEGER,
		        severity VARCHAR(255),
		        file TEXT,
		        pos INTEGER,
		        descript TEXT,
		        build_number INTEGER,
		        PRIMARY KEY (id),
                INDEX plugin_I (plugin(12)),
                INDEX bayzr_err_I (bayzr_err),
                INDEX severity_I (severity(32)),
                INDEX pos_I (pos),
                INDEX build_number_I (build_number))`); err != nil {
			return err
		}
	}

	//Sever part
	should_add_admin := false
	fnd, err = this._checkTable("bayzr_USER")
	if err != nil {
		return err
	}
	if fnd == false {
		if err := this._executeSQLCpammnd(
			`CREATE TABLE IF NOT EXISTS bayzr_USER(
		        id INTEGER  NOT NULL AUTO_INCREMENT, 
		        login VARCHAR(255),
		        name VARCHAR(255),
		        password CHAR(35),
		        email VARCHAR(255),
		        group_id INTEGER,
		        PRIMARY KEY (id),
                INDEX USER_group_id_I (group_id),
                INDEX USER_login_I (login(20)),
                INDEX USER_name_I (name(20)),
                INDEX USER_email_I (email(32)))`); err != nil {
			return err
		}
		should_add_admin = true
	}
	fnd, err = this._checkTable("bayzr_GROUP")
	if err != nil {
		return err
	}
	if fnd == false {
		if err := this._executeSQLCpammnd(
			`CREATE TABLE IF NOT EXISTS bayzr_GROUP(
		        id INTEGER  NOT NULL AUTO_INCREMENT, 
		        name VARCHAR(255),
		        PRIMARY KEY (id),
                INDEX GROUP_name_I (name(20)))`); err != nil {
			return err
		}
	}
	fnd, err = this._checkTable("bayzr_RULE")
	if err != nil {
		return err
	}
	if fnd == false {
		if err := this._executeSQLCpammnd(
			`CREATE TABLE IF NOT EXISTS bayzr_RULE(
		        id INTEGER  NOT NULL AUTO_INCREMENT, 
		        name VARCHAR(10),
		        description VARCHAR(255),
		        PRIMARY KEY (id),
                INDEX RULE_name_I (name(10)))`); err != nil {
			return err
		}
	}
	fnd, err = this._checkTable("bayzr_GROUP_RULE")
	if err != nil {
		return err
	}
	if fnd == false {
		if err := this._executeSQLCpammnd(
			`CREATE TABLE IF NOT EXISTS bayzr_GROUP_RULE(
		        grp_id INTEGER,
		        rule_id INTEGER,
                INDEX GROUP_RULE_grp_I (grp_id),
                INDEX GROUP_RULE_rule_I (rule_id))`); err != nil {
			return err
		}
	}
	fnd, err = this._checkTable("bayzr_TASK")
	if err != nil {
		return err
	}
	if fnd == false {
		if err := this._executeSQLCpammnd(
			`CREATE TABLE IF NOT EXISTS bayzr_TASK(
		        id INTEGER  NOT NULL AUTO_INCREMENT, 
		        name VARCHAR(255),
		        task_type CHAR(1),
		        source VARCHAR(255),
		        pkgs_list TEXT,
		        build_cmds TEXT,
		        project_id VARCHAR(255),
		        auth_id VARCHAR(255),
		        period INTEGER,
		        start_time TIMESTAMP,
		        user_id INTEGER,
		        PRIMARY KEY (id),
                INDEX TASK_name_I (name(20)),
                INDEX TASK_task_type_I (task_type),
                INDEX TASK_source_I (source(32)),
                INDEX TASK_project_id_I (project_id(32)),
                INDEX TASK_user_id_I (user_id))`); err != nil {
			return err
		}
	}
	fnd, err = this._checkTable("bayzr_JOBS")
	if err != nil {
		return err
	}
	if fnd == false {
		if err := this._executeSQLCpammnd(
			`CREATE TABLE IF NOT EXISTS bayzr_JOBS(
		        id INTEGER  NOT NULL AUTO_INCREMENT, 
		        user_id INTEGER,
		        job_name VARCHAR(255),
		        commit VARCHAR(512),
		        build_date_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		        build_date_end TIMESTAMP,
		        build_id INTEGER,
		        task_id INTEGER,
		        PRIMARY KEY (id),
                INDEX JOBS_user_id_I (user_id),
                INDEX JOBS_task_id_id_I (task_id))`); err != nil {
			return err
		}
	}
	fnd, err = this._checkTable("bayzr_OUTPUT")
	if err != nil {
		return err
	}
	if fnd == false {
		if err := this._executeSQLCpammnd(
			`CREATE TABLE IF NOT EXISTS bayzr_OUTPUT(
		        job_id INTEGER, 
		        time_of_string TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		        line VARCHAR(255),
                INDEX OUTPUT_line_I (line(20)),
                INDEX OUTPUT_job_id_I (job_id))`); err != nil {
			return err
		}
	}
	if should_add_admin == true {
		res, err := this.db.Exec(`INSERT INTO bayzr_RULE(name, 
                                        description) VALUES('ru_admin', 'Admin permissions(Can do anything)')`)
		if err != nil {
			return err
		} else {
			id, err := res.LastInsertId()
			if err != nil {
				return err
			} else {
				res2, err2 := this.db.Exec(`INSERT INTO bayzr_GROUP(name) VALUES('Admin Group')`)
				if err2 != nil {
					return err2
				} else {
					id2, err2 := res2.LastInsertId()
					if err2 != nil {
						return err2
					} else {
						_, err3 := this.db.Exec(`INSERT INTO bayzr_GROUP_RULE(grp_id, rule_id) VALUES(?, ?)`, id2, id)
						if err3 != nil {
							return err3
						} else {
							_, err4 := this.db.Exec(`INSERT INTO bayzr_USER(
							   login,
		                       name,
		                       password,
		                       email,
		                       group_id
							) VALUES('su_admin', 'Admin User', '952f25e4337a61dee0cfa41e956e0124', 'admin@admin', ?)`, id2)
							if err4 != nil {
								return err4
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (this *MySQLSaver) Init(config string, ruller_param *rullerlist.RullerList) error {
	this.ok = 0
	this.ruller = ruller_param
	if config == "" {
		return nil
	}
	var err error
	this.db, err = sql.Open("mysql", config)
	if err != nil {
		return err
	}
	err = this.db.Ping()
	if err != nil {
		return err
	}
	err = this._checkAndCreateTables()
	if err != nil {
		return err
	}
	this.ok = 1
	return nil
}

func (this *MySQLSaver) EasyInit(config string) error {
	this.ok = 0
	var err error
	this.db, err = sql.Open("mysql", config)
	if err != nil {
		return err
	}
	err = this.db.Ping()
	if err != nil {
		return err
	}
	this.ok = 1
	return nil
}

func (this *MySQLSaver) Finalize() {
	if this.ok == 1 {
		this.db.Close()
		this.ok = 0
	}
}

func (this *MySQLSaver) CreateCurrentBuild(author string, name_of_build string) error {
	if this.ok == 1 {
		res, err := this.db.Exec(`INSERT INTO bayzr_build_info(autor_of_build, 
                                        name_of_build, completed) VALUES(?,?,0)`, author, name_of_build)
		if err != nil {
			return err
		} else {
			id, err := res.LastInsertId()
			if err != nil {
				return err
			} else {
				this.current_build_id = id
			}
		}
	}
	return nil
}

func (this *MySQLSaver) FinalizeCurrentBuild() error {
	if this.ok == 1 && this.current_build_id > 0 {
		_, err := this.db.Exec(`UPDATE bayzr_build_info SET completed = 1 WHERE id = ?`, this.current_build_id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *MySQLSaver) InsertInfo(plugin string, severity string,
	file_name string, position string, descr string, err_type int) error {
	if this.ok == 1 {
		pos, _ := strconv.ParseInt(strings.Trim(position, " \n\t"), 10, 64)
		sev := severity
		if this.ruller.IsInList(severity) == false {
			sev = ""
		}
		_, err := this.db.Exec(`INSERT INTO bayzr_err_list(plugin, bayzr_err, 
                                        severity, file, pos, descript, build_number) VALUES(?,?,?,?,?,?,?)`, this.getStringLength(plugin, 50), err_type, this.getStringLength(sev, 255),
			file_name, pos, descr, this.current_build_id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *MySQLSaver) InsertExtInfo(plugin string, file_name string, file_string string, rec_type int,
	err_type int, descr string, file_pos int64) error {
	if this.ok == 1 {

		_, err := this.db.Exec(`INSERT INTO bayzr_err_extend(plugin, file_name, file_string,
                                        rec_type, err_type, descript, build_number, file_pos) 
                                        VALUES(?,?,?,?,?,?,?,?)`, this.getStringLength(plugin, 50), this.getStringLength(file_name, 255),
			file_string, rec_type, err_type, descr, this.current_build_id, file_pos)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *MySQLSaver) IsDbConnect() bool {
	return this.ok == 1
}

func (this *MySQLSaver) CheckUser(user string, passwd string) (error, bool, []string, int) {
	stmtOut, err := this.db.Prepare(`SELECT t1.id FROM bayzr_USER as t1 
	WHERE t1.login=? AND t1.password=?`)
	if err != nil {
		return err, false, []string{}, 0
	}
	defer stmtOut.Close()
	var read_id int
	err = stmtOut.QueryRow(user, passwd).Scan(&read_id)
	if err != nil && err != sql.ErrNoRows {
		return err, false, []string{}, 0
	}
	if err == sql.ErrNoRows {
		return nil, false, []string{}, 0
	} else {
		short_list := []string{}
		err, list := this.GetUserPerms(read_id)
		if err != nil {
			return err, false, []string{}, 0
		}
		for _, val := range list {
			short_list = append(short_list, val[0])
		}

		return nil, true, short_list, read_id
	}
}

func (this *MySQLSaver) GetUserPerms(id int) (error, [][]string) {
	result := [][]string{}
	stmtOut, err := this.db.Prepare(`select t4.name, t4.description from bayzr_USER as t1 
	join bayzr_GROUP as t2 on t1.group_id = t2.id 
	join bayzr_GROUP_RULE as t3 on t2.id = t3.grp_id 
	join bayzr_RULE as t4 on t3.rule_id = t4.id 
	where t1.id = ?`)
	if err != nil {
		return err, result
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(id)
	if err != nil && err != sql.ErrNoRows {
		return err, result
	}
	for rows.Next() {
		var name string
		var descr string
		if err := rows.Scan(&name, &descr); err != nil {
			return err, [][]string{}
		}
		result = append(result, []string{name, descr})
	}
	return err, result
}

func (this *MySQLSaver) GetUserInfo(id int) (error, string, string, string, int, string) {
	stmtOut, err := this.db.Prepare(`SELECT t1.login, t1.name, t1.email, t2.id, t2.name 
	FROM bayzr_USER as t1 join bayzr_GROUP as t2 on t1.group_id = t2.id
	WHERE t1.id = ?`)
	if err != nil {
		return err, "", "", "", 0, ""
	}
	defer stmtOut.Close()
	var st_login string
	var st_name string
	var st_email string
	var st_grp_id int
	var st_grp string
	err = stmtOut.QueryRow(id).Scan(&st_login, &st_name, &st_email, &st_grp_id, &st_grp)
	if err != nil {
		return err, "", "", "", 0, ""
	}
	return nil, st_login, st_name, st_email, st_grp_id, st_grp
}
