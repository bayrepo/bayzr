package mysqlsaver

import (
	"bufio"
	"crypto/rand"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"rullerlist"
	"sort"
	"strconv"
	"strings"
)

type MySQLSaver struct {
	db               *sql.DB
	ok               int
	current_build_id int64
	ruller           *rullerlist.RullerList
}

func randToken() string {
	b := make([]byte, 15)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
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
	fnd, err = this._checkTable("bayzr_err_extend_file")
	if err != nil {
		return err
	}
	if fnd == false {
		if err := this._executeSQLCpammnd(
			`CREATE TABLE IF NOT EXISTS bayzr_err_extend_file(
		        id INTEGER  NOT NULL AUTO_INCREMENT, 
		        file_string TEXT,
		        build_number INTEGER,
		        PRIMARY KEY (id),
                INDEX build_number_I_file (build_number))`); err != nil {
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
		        login VARCHAR(255) DEFAULT "",
		        name VARCHAR(255) DEFAULT "",
		        password CHAR(35) DEFAULT "",
		        email VARCHAR(255) DEFAULT "",
		        group_id INTEGER,
		        api_auth_tocken VARCHAR(30) DEFAULT "",
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
		        name VARCHAR(255) DEFAULT "",
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
		        name VARCHAR(10) DEFAULT "",
		        description VARCHAR(255) DEFAULT "",
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
		        name VARCHAR(255) DEFAULT "",
		        task_type CHAR(1) DEFAULT "",
		        source VARCHAR(255) DEFAULT "",
		        pkgs_list TEXT,
		        build_cmds TEXT,
		        period CHAR(1) DEFAULT "",
		        start_time VARCHAR(20) DEFAULT "",
		        user_id INTEGER,
		        check_config TEXT,
		        users_list TEXT,
		        auth_tocken VARCHAR(30) DEFAULT "",
		        use_branch INT DEFAULT 0,
		        result_file VARCHAR(128) DEFAULT "result.html",
		        branch VARCHAR(255) DEFAULT "",
		        PRIMARY KEY (id),
                INDEX TASK_name_I (name(20)),
                INDEX TASK_task_type_I (task_type),
                INDEX TASK_source_I (source(32)),
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
		        job_name VARCHAR(255) DEFAULT "",
		        commit VARCHAR(512) DEFAULT "",
		        create_date_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		        build_date_start TIMESTAMP,
		        build_date_end TIMESTAMP,
		        build_id INTEGER,
		        priority CHAR DEFAULT "",
		        task_id INTEGER,
		        descr TEXT,
		        PRIMARY KEY (id),
                INDEX JOBS_user_id_I (user_id),
                INDEX JOBS_priority_I (priority),
                INDEX JOBS_create_date_start_I (create_date_start),
                INDEX JOBS_task_id_id_I (task_id),
                INDEX JOBS_descr_I (descr(50)),
                INDEX JOBS_build_date_start_I (build_date_start))`); err != nil {
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
		                       group_id,
		                       api_auth_tocken
							) VALUES('su_admin', 'Admin User', '952f25e4337a61dee0cfa41e956e0124', 'admin@admin', ?, ?)`, id2, randToken())
							if err4 != nil {
								return err4
							}
							_, err4 = this.db.Exec(`INSERT INTO bayzr_USER(
							   login,
		                       name,
		                       password,
		                       email,
		                       group_id,
		                       api_auth_tocken
							) VALUES('su_checker', 'Auto Checker', 'nopasswd', 'admin@admin', ?, ?)`, id2, randToken())
							if err4 != nil {
								return err4
							}
						}
					}
				}
			}
		}
		res, err = this.db.Exec(`INSERT INTO bayzr_RULE(name, 
                                        description) VALUES('ru_task', 'Can create new project and watch the results')`)
		if err != nil {
			return err
		}
		ru_task, err := res.LastInsertId()
		if err != nil {
			return err
		}
		res, err = this.db.Exec(`INSERT INTO bayzr_RULE(name, 
                                        description) VALUES('ru_result', 'Can watch the results')`)
		if err != nil {
			return err
		}
		ru_result, err := res.LastInsertId()
		if err != nil {
			return err
		}
		res, err = this.db.Exec(`INSERT INTO bayzr_RULE(name, 
                                        description) VALUES('ru_norules', 'Only new user')`)
		if err != nil {
			return err
		}
		ru_norules, err := res.LastInsertId()
		if err != nil {
			return err
		}

		res2, err2 := this.db.Exec(`INSERT INTO bayzr_GROUP(name) VALUES('User Group')`)
		if err2 != nil {
			return err2
		} else {
			id2, err2 := res2.LastInsertId()
			if err2 != nil {
				return err2
			} else {
				_, err3 := this.db.Exec(`INSERT INTO bayzr_GROUP_RULE(grp_id, rule_id) VALUES(?, ?)`, id2, ru_norules)
				if err3 != nil {
					return err3
				}
			}
		}

		res2, err2 = this.db.Exec(`INSERT INTO bayzr_GROUP(name) VALUES('Viewer Group')`)
		if err2 != nil {
			return err2
		} else {
			id2, err2 := res2.LastInsertId()
			if err2 != nil {
				return err2
			} else {
				_, err3 := this.db.Exec(`INSERT INTO bayzr_GROUP_RULE(grp_id, rule_id) VALUES(?, ?)`, id2, ru_result)
				if err3 != nil {
					return err3
				}
			}
		}

		res2, err2 = this.db.Exec(`INSERT INTO bayzr_GROUP(name) VALUES('Developer Group')`)
		if err2 != nil {
			return err2
		} else {
			id2, err2 := res2.LastInsertId()
			if err2 != nil {
				return err2
			} else {
				_, err3 := this.db.Exec(`INSERT INTO bayzr_GROUP_RULE(grp_id, rule_id) VALUES(?, ?)`, id2, ru_task)
				if err3 != nil {
					return err3
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
	/*if this.ok == 1 {

			_, err := this.db.Exec(`INSERT INTO bayzr_err_extend(plugin, file_name, file_string,
	                                        rec_type, err_type, descript, build_number, file_pos)
	                                        VALUES(?,?,?,?,?,?,?,?)`, this.getStringLength(plugin, 50), this.getStringLength(file_name, 255),
				file_string, rec_type, err_type, descr, this.current_build_id, file_pos)
			if err != nil {
				return err
			}
		}*/
	return nil
}

func (this *MySQLSaver) InsertExtInfoFromResult(file_path string, build_name string) error {
	if this.ok == 1 {
		berr, id := this.GetBuildId(build_name)
		if berr != nil {
			return berr
		}
		if id > 0 {

			file, errf := os.Open(file_path)
			if errf != nil {
				return errf
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				_, err := this.db.Exec(`INSERT INTO bayzr_err_extend_file(file_string, build_number)
	                                        VALUES(?,?)`, scanner.Text(), id)
				if err != nil {
					return err
				}
			}

			if erre := scanner.Err(); erre != nil {
				return erre
			}

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

func (this *MySQLSaver) GetUserInfo(id int) (error, string, string, string, int, string, string) {
	stmtOut, err := this.db.Prepare(`SELECT t1.login, t1.name, t1.email, t2.id, t2.name, t1.api_auth_tocken 
	FROM bayzr_USER as t1 join bayzr_GROUP as t2 on t1.group_id = t2.id
	WHERE t1.id = ?`)
	if err != nil {
		return err, "", "", "", 0, "", ""
	}
	defer stmtOut.Close()
	var st_login string
	var st_name string
	var st_email string
	var st_grp_id int
	var st_grp string
	var st_tocken string
	err = stmtOut.QueryRow(id).Scan(&st_login, &st_name, &st_email, &st_grp_id, &st_grp, &st_tocken)
	if err != nil {
		return err, "", "", "", 0, "", ""
	}
	return nil, st_login, st_name, st_email, st_grp_id, st_grp, st_tocken
}

func (this *MySQLSaver) IsDupUser(user string) (error, bool) {
	stmtOut, err := this.db.Prepare(`SELECT t1.id FROM bayzr_USER as t1 
	WHERE t1.login=?`)
	if err != nil {
		return err, true
	}
	defer stmtOut.Close()
	var read_id int
	err = stmtOut.QueryRow(user).Scan(&read_id)
	if err != nil && err != sql.ErrNoRows {
		return err, true
	}
	if err == sql.ErrNoRows {
		return nil, false
	} else {
		return nil, true
	}
}

func (this *MySQLSaver) InsertUser(login string, name string, email string, password string) (error, int) {
	var id int64 = 0
	if this.ok == 1 {
		if err, id_tmp := this.GetGrpIdByName("User Group"); err != nil {
			return err, 0
		} else {
			res, err2 := this.db.Exec(`INSERT INTO bayzr_USER(login, name, password, email, group_id, api_auth_tocken) 
                                        VALUES(?,?,?,?,?,?)`, login, name, password, email, id_tmp, randToken())
			if err2 != nil {
				return err2, 0
			}

			id, err = res.LastInsertId()
			if err != nil {
				return err, 0
			}
		}
	}

	return nil, int(id)
}

func (this *MySQLSaver) UpdateUser(id int, name string, email string, password string) error {
	if this.ok == 1 {

		_, err2 := this.db.Exec(`UPDATE bayzr_USER SET name = ?, email = ? WHERE id = ?`,
			name, email, id)
		if err2 != nil {
			return err2
		}

		if password != "" {
			_, err2 := this.db.Exec(`UPDATE bayzr_USER SET password = ? WHERE id = ?`,
				password, id)
			if err2 != nil {
				return err2
			}
		}

	}

	return nil
}

func (this *MySQLSaver) GetGrpIdByName(name string) (error, int) {
	stmtOut, err := this.db.Prepare(`SELECT t1.id FROM bayzr_GROUP as t1 
	WHERE t1.name=?`)
	if err != nil {
		return err, 0
	}
	defer stmtOut.Close()
	var read_id int
	err = stmtOut.QueryRow(name).Scan(&read_id)
	if err != nil && err != sql.ErrNoRows {
		return err, 0
	}
	if err == sql.ErrNoRows {
		return nil, 0
	} else {
		return nil, read_id
	}
}

func (this *MySQLSaver) GetUsersList() (error, int, [][]string) {
	result := [][]string{}
	counter := 0
	stmtOut, err := this.db.Prepare(`SELECT t1.id, t1.login, t1.name, t1.email, 
	t1.group_id, t2.name, t1.api_auth_tocken FROM bayzr_USER as t1 JOIN bayzr_GROUP as t2 on t1.group_id = t2.id 
	WHERE t1.login <> 'su_admin' ORDER BY t1.login, t1.name`)
	if err != nil {
		return err, 0, result
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query()
	if err != nil && err != sql.ErrNoRows {
		return err, 0, result
	}
	for rows.Next() {
		var id string
		var login string
		var name string
		var email string
		var group_id string
		var group_name string
		var st_tocken string
		if err := rows.Scan(&id, &login, &name, &email, &group_id, &group_name, &st_tocken); err != nil {
			return err, 0, [][]string{}
		}
		result = append(result, []string{id, login, name, email, group_id, group_name, st_tocken})
		counter += 1
	}
	return err, counter, result
}

func (this *MySQLSaver) GetUsersTokenLists() (error, [][]string) {
	result := [][]string{}
	stmtOut, err := this.db.Prepare(`SELECT t1.login, t1.api_auth_tocken 
	FROM bayzr_USER as t1 ORDER BY t1.login`)
	if err != nil {
		return err, result
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query()
	if err != nil && err != sql.ErrNoRows {
		return err, result
	}
	for rows.Next() {
		var login string
		var tocken string
		if err := rows.Scan(&login, &tocken); err != nil {
			return err, [][]string{}
		}
		result = append(result, []string{login, tocken})
	}
	return err, result
}

func (this *MySQLSaver) GetUsersCount() (error, int) {
	counter := 0
	stmtOut, err := this.db.Prepare(`SELECT count(*) FROM bayzr_USER`)
	if err != nil {
		return err, 0
	}
	defer stmtOut.Close()
	err = stmtOut.QueryRow().Scan(&counter)
	if err != nil && err != sql.ErrNoRows {
		return err, 0
	}
	if err == sql.ErrNoRows {
		return nil, 0
	} else {
		return nil, counter
	}
}

func (this *MySQLSaver) GetGroups() (error, [][]string) {
	result := [][]string{}
	stmtOut, err := this.db.Prepare(`select t1.id, t1.name from bayzr_GROUP as t1 ORDER BY t1.name`)
	if err != nil {
		return err, result
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query()
	if err != nil && err != sql.ErrNoRows {
		return err, result
	}
	for rows.Next() {
		var id string
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err, [][]string{}
		}
		result = append(result, []string{id, name})
	}
	return err, result
}

func (this *MySQLSaver) UpdateUserGroup(id int, grp_id int) error {
	if this.ok == 1 {

		_, err2 := this.db.Exec(`UPDATE bayzr_USER SET group_id = ? WHERE id = ?`,
			grp_id, id)
		if err2 != nil {
			return err2
		}

	}

	return nil
}

func (this *MySQLSaver) DelUser(id int) error {
	if this.ok == 1 {

		_, err2 := this.db.Exec(`DELETE FROM bayzr_USER WHERE id = ?`, id)
		if err2 != nil {
			return err2
		}

	}

	return nil
}

func (this *MySQLSaver) GetPackages() (error, []string) {
	result := []string{}
	stmtOut, err := this.db.Prepare(`select pkgs_list from bayzr_TASK`)
	if err != nil {
		return err, result
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query()
	if err != nil && err != sql.ErrNoRows {
		return err, result
	}
	for rows.Next() {
		var pkgs_list string
		if err := rows.Scan(&pkgs_list); err != nil {
			return err, []string{}
		}
		result = append(result, strings.Split(pkgs_list, ",")...)
	}
	result_prc := []string{}
	for _, val := range result {
		if strings.Trim(val, " ") == "" {
			continue
		}
		fnd := false
		for _, val_2 := range result_prc {
			if strings.Trim(val, " ") == val_2 {
				fnd = true
				break
			}
		}
		if fnd == false {
			result_prc = append(result_prc, strings.Trim(val, " "))
		}
	}
	sort.Strings(result_prc)
	return err, result_prc
}

func (this *MySQLSaver) SaveTask(owner_id int, name string, Ttype string, git string, new_pkgs []string,
	old_pkgs []string, cmds []string, Ptype string, period string, users []string, cfg []string, use_branch int, repname string, brn string) (error, int) {
	var id int64 = 0
	if this.ok == 1 {
		task_type := Ttype
		pkgs_list := ""
		if len(old_pkgs) > 0 {
			pkgs_list += strings.Join(old_pkgs, ",")
		}
		if len(new_pkgs) > 0 {
			c_list := []string{}
			for _, val := range new_pkgs {
				fnd := false
				for _, o_val := range old_pkgs {
					if o_val == val {
						fnd = true
						break
					}
				}
				if fnd == false {
					c_list = append(c_list, val)
				}
			}
			if len(c_list) > 0 {
				if pkgs_list != "" {
					pkgs_list = pkgs_list + "," + strings.Join(c_list, ",")
				} else {
					pkgs_list = strings.Join(c_list, ",")
				}
			}
		}
		per_type := Ptype
		res, err2 := this.db.Exec(`INSERT INTO bayzr_TASK(name, task_type, source, pkgs_list,
										build_cmds, period, start_time, user_id, check_config,
										users_list, auth_tocken, use_branch, result_file, branch) 
                                        VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, name, string(task_type), git, pkgs_list,
			strings.Join(cmds, "\n"), string(per_type), string(period), owner_id,
			strings.Join(cfg, "\n"), strings.Join(users, ","), randToken(), use_branch, repname, brn)
		if err2 != nil {
			return err2, 0
		}

		var err error
		id, err = res.LastInsertId()
		if err != nil {
			return err, 0
		}

	}

	return nil, int(id)

}

func (this *MySQLSaver) UpdateTask(id int, owner_id int, name string, Ttype string, git string, new_pkgs []string,
	old_pkgs []string, cmds []string, Ptype string, period string, users []string, cfg []string, use_branch int, repname string, brn string) error {
	if this.ok == 1 {
		task_type := Ttype[0]
		pkgs_list := ""
		if len(old_pkgs) > 0 {
			pkgs_list += strings.Join(old_pkgs, ",")
		}
		if len(new_pkgs) > 0 {
			c_list := []string{}
			for _, val := range new_pkgs {
				fnd := false
				for _, o_val := range old_pkgs {
					if o_val == val {
						fnd = true
						break
					}
				}
				if fnd == false {
					c_list = append(c_list, val)
				}
			}
			if len(c_list) > 0 {
				if pkgs_list != "" {
					pkgs_list = pkgs_list + "," + strings.Join(c_list, ",")
				} else {
					pkgs_list = strings.Join(c_list, ",")
				}
			}
		}
		per_type := Ptype[0]
		_, err2 := this.db.Exec(`UPDATE bayzr_TASK SET name = ?, task_type = ?, source = ?, pkgs_list = ?,
										build_cmds = ?, period = ?, start_time = ?, check_config = ?,
										users_list = ?, use_branch = ?, result_file = ?, branch = ? WHERE id = ?`, name, string(task_type), git, pkgs_list,
			strings.Join(cmds, "\n"), string(per_type), string(period),
			strings.Join(cfg, "\n"), strings.Join(users, ","), use_branch, repname, brn, id)
		if err2 != nil {
			return err2
		}

	}

	return nil

}

func (this *MySQLSaver) GetTasks() (error, [][]string) {
	result := [][]string{}
	stmtOut, err := this.db.Prepare(`select t1.id, t1.name, t1.task_type, t1.source, t1.pkgs_list,
	t1.build_cmds, t1.period, t1.start_time, t2.name, t1.check_config, t1.users_list,
	t1.auth_tocken, t1.branch from bayzr_TASK as t1 
	join bayzr_USER as t2 on t1.user_id = t2.id order by t1.id`)
	if err != nil {
		return err, result
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query()
	if err != nil && err != sql.ErrNoRows {
		return err, result
	}
	for rows.Next() {
		var (
			t1_id           int
			t1_task_type    string
			t1_name         string
			t1_source       string
			t1_pkgs_list    string
			t1_build_cmds   string
			t1_period       string
			t1_start_time   string
			t2_name         string
			t1_check_config string
			t1_users_list   string
			t1_auth_tocken  string
			ti_branch       string
		)
		if err := rows.Scan(&t1_id, &t1_name, &t1_task_type, &t1_source, &t1_pkgs_list,
			&t1_build_cmds, &t1_period, &t1_start_time, &t2_name, &t1_check_config, &t1_users_list,
			&t1_auth_tocken, &ti_branch); err != nil {
			return err, [][]string{}
		}
		result = append(result, []string{fmt.Sprintf("%d", t1_id), t1_name, t1_task_type, t1_source, t1_pkgs_list,
			t1_build_cmds, t1_period, t1_start_time, t2_name, t1_check_config, t1_users_list,
			t1_auth_tocken, ti_branch})
	}
	return err, result
}

func (this *MySQLSaver) GetTask(id int) (error, []string) {

	stmtOut, err := this.db.Prepare(`select t1.id, t1.name, t1.task_type, t1.source, t1.pkgs_list,
	t1.build_cmds, t1.period, t1.start_time, t2.name, t1.check_config, t1.users_list,
	t1.auth_tocken, t1.use_branch, t1.result_file, t1.branch from bayzr_TASK as t1 
	join bayzr_USER as t2 on t1.user_id = t2.id where t1.id = ?`)
	if err != nil {
		return err, []string{}
	}
	defer stmtOut.Close()

	var (
		t1_id           int
		t1_name         string
		t1_task_type    string
		t1_source       string
		t1_pkgs_list    string
		t1_build_cmds   string
		t1_period       string
		t1_start_time   string
		t2_name         string
		t1_check_config string
		t1_users_list   string
		t1_auth_tocken  string
		t1_use_branch   string
		t1_result_file  string
		t1_branch       string
	)

	err = stmtOut.QueryRow(id).Scan(&t1_id, &t1_name, &t1_task_type, &t1_source, &t1_pkgs_list,
		&t1_build_cmds, &t1_period, &t1_start_time, &t2_name, &t1_check_config, &t1_users_list,
		&t1_auth_tocken, &t1_use_branch, &t1_result_file, &t1_branch)
	if err != nil && err != sql.ErrNoRows {
		return err, []string{}
	}
	if err == sql.ErrNoRows {
		return nil, []string{}
	} else {
		return nil, []string{fmt.Sprintf("%d", t1_id), t1_name, t1_task_type, t1_source, t1_pkgs_list,
			t1_build_cmds, t1_period, t1_start_time, t2_name, t1_check_config, t1_users_list,
			t1_auth_tocken, t1_use_branch, t1_result_file, t1_branch}
	}

}

func (this *MySQLSaver) DelTask(id int) error {
	if this.ok == 1 {

		_, err2 := this.db.Exec(`DELETE FROM bayzr_TASK WHERE id = ?`, id)
		if err2 != nil {
			return err2
		}

	}

	return nil
}

func (this *MySQLSaver) InsertJob(owner_id int, name string, commit string, priority string, task_id int, descr string) (error, int) {
	var id int64 = 0
	if this.ok == 1 {

		res, err2 := this.db.Exec(`INSERT INTO bayzr_JOBS(user_id, job_name, commit, priority, task_id, descr) 
                                        VALUES(?,?,?,?,?,?)`, owner_id, name, commit, priority, task_id, descr)
		if err2 != nil {
			return err2, 0
		}
		var err error

		id, err = res.LastInsertId()
		if err != nil {
			return err, 0
		}

	}

	return nil, int(id)
}

func (this *MySQLSaver) GetJobs() (error, [][]string) {
	result := [][]string{}
	stmtOut, err := this.db.Prepare(`select t1.id, t1.job_name, t1.commit, ifnull(t1.build_date_start,""), ifnull(t1.build_date_end,""),
	ifnull(t1.build_id,0), t1.priority, t2.name, t3.name, t1.descr from bayzr_JOBS as t1 
	join bayzr_USER as t2 on t1.user_id = t2.id join bayzr_TASK as t3 on t3.id = t1.task_id order by t1.id desc`)
	if err != nil {
		return err, result
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query()
	if err != nil && err != sql.ErrNoRows {
		return err, result
	}
	for rows.Next() {
		var (
			t1_id               int
			t1_job_name         string
			t1_commit           string
			t1_build_date_start string
			t1_build_date_end   string
			t1_build_id         string
			t1_priority         string
			t2_user_name        string
			t3_task_name        string
			t1_descr            string
		)
		if err := rows.Scan(&t1_id, &t1_job_name, &t1_commit, &t1_build_date_start, &t1_build_date_end,
			&t1_build_id, &t1_priority, &t2_user_name, &t3_task_name, &t1_descr); err != nil {
			return err, [][]string{}
		}
		result = append(result, []string{fmt.Sprintf("%d", t1_id), t1_job_name, t1_commit, t1_build_date_start, t1_build_date_end,
			t1_build_id, t1_priority, t2_user_name, t3_task_name, t1_descr})
	}
	return err, result
}

func (this *MySQLSaver) DelJob(id int) error {
	if this.ok == 1 {

		_, err2 := this.db.Exec(`DELETE FROM bayzr_JOBS WHERE id = ?`, id)
		if err2 != nil {
			return err2
		}

	}

	return nil
}

func (this *MySQLSaver) DetJobID() (error, int) {
	stmtOut, err := this.db.Prepare(`select id from bayzr_JOBS where build_date_start=0 order by priority desc, create_date_start asc limit 1`)
	if err != nil {
		return err, 0
	}
	defer stmtOut.Close()
	var id int
	err = stmtOut.QueryRow().Scan(&id)
	if err == sql.ErrNoRows {
		return nil, 0
	}
	if err != nil {
		return err, 0
	}
	return nil, id
}

func (this *MySQLSaver) InsertOutput(id int, message string) error {
	if this.ok == 1 {

		_, err2 := this.db.Exec(`INSERT INTO bayzr_OUTPUT(job_id, line) 
                                        VALUES(?,?)`, id, message)
		if err2 != nil {
			return err2
		}

	}

	return nil
}

func (this *MySQLSaver) GetTaskFullInfo(id int) (error, map[string]string) {
	if this.ok == 1 {

		stmtOut, err := this.db.Prepare(`select t1.id, t1.user_id, t2.login, t2.name, t1.job_name, 
			t1.commit, t3.name, t3.task_type, t3.source, t3.pkgs_list, t3.build_cmds, t3.check_config, 
			t3.result_file, t3.use_branch from bayzr_JOBS as t1 join bayzr_USER as t2 on t2.id = t1.user_id 
			join bayzr_TASK as t3 on t3.id = t1.task_id where t1.id = ?`)
		if err != nil {
			return err, map[string]string{}
		}
		defer stmtOut.Close()

		var (
			t1_id           string
			t1_user_id      string
			t2_login        string
			t2_name         string
			t1_job_name     string
			t1_commit       string
			t3_name         string
			t3_task_type    string
			t3_source       string
			t3_pkgs_list    string
			t3_build_cmds   string
			t3_check_config string
			t3_result_file  string
			t3_use_branch   string
		)

		err = stmtOut.QueryRow(id).Scan(&t1_id, &t1_user_id, &t2_login, &t2_name, &t1_job_name,
			&t1_commit, &t3_name, &t3_task_type, &t3_source, &t3_pkgs_list, &t3_build_cmds, &t3_check_config,
			&t3_result_file, &t3_use_branch)
		if err != nil && err != sql.ErrNoRows {
			return err, map[string]string{}
		}
		if err == sql.ErrNoRows {
			return fmt.Errorf("Empty job"), map[string]string{}
		} else {
			res := map[string]string{
				"id":          t1_id,
				"user_id":     t1_user_id,
				"login":       t2_login,
				"name":        t2_name,
				"job_name":    t1_job_name,
				"commit":      t1_commit,
				"task_name":   t3_name,
				"task_type":   t3_task_type,
				"source":      t3_source,
				"pkgs":        t3_pkgs_list,
				"cmds":        t3_build_cmds,
				"config":      t3_check_config,
				"result_file": t3_result_file,
				"use_branch":  t3_use_branch,
			}
			return nil, res
		}

	}

	return fmt.Errorf("No database connetion"), map[string]string{}
}

func (this *MySQLSaver) TakeJob(id int) error {
	if this.ok == 1 {

		_, err2 := this.db.Exec(`update bayzr_JOBS set build_date_start = CURRENT_TIMESTAMP() where id = ?`, id)
		if err2 != nil {
			return err2
		}

	}

	return nil
}

func (this *MySQLSaver) CompleteJob(id int) error {
	if this.ok == 1 {

		_, err2 := this.db.Exec(`update bayzr_JOBS set build_date_end = CURRENT_TIMESTAMP() where id = ?`, id)
		if err2 != nil {
			return err2
		}

		_, err3 := this.db.Exec(`update bayzr_JOBS as b1 
		join 
		(select tt1.id as idb, tt2.id as idj from bayzr_build_info as tt1 
		join 
		(select concat(t2.name, ".", t1.id) as nm, t1.id from bayzr_JOBS as t1 join bayzr_TASK as t2 
		on t1.task_id = t2.id) as tt2 on tt1.name_of_build = tt2.nm and tt1.completed = 1) as b2 
		on b1.id = b2.idj set b1.build_id = b2.idb`)
		if err3 != nil {
			return err3
		}

	}

	return nil
}

func (this *MySQLSaver) GetOut(id int) (error, [][]string) {
	result := [][]string{}
	stmtOut, err := this.db.Prepare(`select t1.time_of_string, t1.line from bayzr_OUTPUT as t1 where t1.job_id = ? order by t1.time_of_string`)
	if err != nil {
		return err, [][]string{}
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(id)
	if err != nil && err != sql.ErrNoRows {
		return err, result
	}
	for rows.Next() {
		var (
			t1_time string
			t1_line string
		)
		if err := rows.Scan(&t1_time, &t1_line); err != nil {
			return err, [][]string{}
		}
		t1_res := ""
		if strings.Contains(t1_line, "+++:") {
			t1_res = "success"
		}
		if strings.Contains(t1_line, "Error:") {
			t1_res = "danger"
		}
		result = append(result, []string{t1_time, t1_line, t1_res})
	}

	return err, result
}

func (this *MySQLSaver) GetBuildId(name string) (error, int) {
	stmtOut, err := this.db.Prepare(`select id from bayzr_build_info where name_of_build = ?`)
	if err != nil {
		return err, 0
	}
	defer stmtOut.Close()
	var read_id int
	err = stmtOut.QueryRow(name).Scan(&read_id)
	if err != nil && err != sql.ErrNoRows {
		return err, 0
	}
	if err == sql.ErrNoRows {
		return nil, 0
	} else {
		return nil, read_id
	}
}

func (this *MySQLSaver) GetResult(id int) (error, []string) {
	result := []string{}
	stmtOut, err := this.db.Prepare(`select file_string from bayzr_err_extend_file where build_number = ? order by id`)
	if err != nil {
		return err, []string{}
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(id)
	if err != nil && err != sql.ErrNoRows {
		return err, result
	}
	for rows.Next() {
		var (
			t1_line string
		)
		if err := rows.Scan(&t1_line); err != nil {
			return err, []string{}
		}
		result = append(result, t1_line)
	}

	return err, result
}

func (this *MySQLSaver) GetUserAuth(auth string) (error, int, string) {

	stmtOut, err := this.db.Prepare(`select t1.id, t1.login from bayzr_USER as t1 where t1.api_auth_tocken = ?`)
	if err != nil {
		return err, 0, ""
	}
	defer stmtOut.Close()

	var (
		t1_id    int
		t1_login string
	)

	err = stmtOut.QueryRow(auth).Scan(&t1_id, &t1_login)
	if err != nil && err != sql.ErrNoRows {
		return err, 0, ""
	}
	if err == sql.ErrNoRows {
		return nil, 0, ""
	} else {
		return nil, t1_id, t1_login
	}

}

func (this *MySQLSaver) GetTaskAuth(auth string) (error, int, []string) {

	stmtOut, err := this.db.Prepare(`select t1.id, t1.users_list from bayzr_TASK as t1 where t1.auth_tocken = ?`)
	if err != nil {
		return err, 0, []string{}
	}
	defer stmtOut.Close()

	var (
		t1_id         int
		t1_users_list string
	)

	err = stmtOut.QueryRow(auth).Scan(&t1_id, &t1_users_list)
	if err != nil && err != sql.ErrNoRows {
		return err, 0, []string{}
	}
	if err == sql.ErrNoRows {
		return nil, 0, []string{}
	} else {
		u_list_raw := strings.Split(t1_users_list, ",")
		u_list := []string{}
		for _, val := range u_list_raw {
			res := strings.Trim(val, ", \n\t\r")
			if res != "" {
				u_list = append(u_list, res)
			}
		}
		return nil, t1_id, u_list
	}

}

func (this *MySQLSaver) GetListOfTimedId() (error, [][]string) {
	result := [][]string{}
	stmtOut, err := this.db.Prepare(`select id, period, start_time, branch, name from bayzr_TASK where period <> '5' and use_branch = 1 and branch <> ''`)
	if err != nil {
		return err, result
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query()
	if err != nil && err != sql.ErrNoRows {
		return err, result
	}
	for rows.Next() {
		var (
			t1_id int
			t1_period string
			t1_start_time string
			t1_branch string
			t1_name string
		)
		if err := rows.Scan(&t1_id, &t1_period, &t1_start_time, &t1_branch, &t1_name); err != nil {
			return err, [][]string{}
		}
		result = append(result, []string{fmt.Sprintf("%d", t1_id), t1_period, t1_start_time, t1_branch, t1_name})
	}

	return err, result
}
