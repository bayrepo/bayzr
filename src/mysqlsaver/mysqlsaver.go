package mysqlsaver

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

type MySQLSaver struct {
	db               *sql.DB
	ok               int
	current_build_id int64
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
	fnd, err := this._checkTable("bayzr_build_info")
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
                INDEX number_I (number),
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
	return nil
}

func (this *MySQLSaver) Init(config string) error {
	this.ok = 0
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

func (this *MySQLSaver) Finalize() {
	if this.ok == 1 {
		this.db.Close()
		this.ok = 0
	}
}

func (this *MySQLSaver) CreateCurrentBuild(author string, name_of_build string) error {
	if this.ok == 1 {
		res, err := this.db.Exec(`INSERT INTO bayzr_build_info(autor_of_build, 
                                        name_of_build) VALUES(?,?)`, author, name_of_build)
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

func (this *MySQLSaver) InsertInfo(plugin string, severity string,
	file_name string, position string, descr string) error {
	if this.ok == 1 {
		pos, _ := strconv.ParseInt(strings.Trim(position, " \n\t"), 10, 64)
		_, err := this.db.Exec(`INSERT INTO bayzr_err_list(plugin, bayzr_err, 
                                        severity, file, pos, descript, build_number) VALUES(?,0,?,?,?,?,?)`, this.getStringLength(plugin, 50), this.getStringLength(severity, 255),
			file_name, pos, descr, this.current_build_id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *MySQLSaver) InsertExtInfo(plugin string, file_name string, file_string string, rec_type int,
	err_type string, descr string, file_pos int64) error {
	if this.ok == 1 {
		_, err := this.db.Exec(`INSERT INTO bayzr_err_extend(plugin, file_name, file_string,
                                        rec_type, err_type, descript, build_number, file_pos) 
                                        VALUES(?,?,?,?,0,?,?,?)`, this.getStringLength(plugin, 50), this.getStringLength(file_name, 255),
			file_string, rec_type, descr, this.current_build_id, file_pos)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *MySQLSaver) IsDbConnect() bool {
	return this.ok == 1
}
