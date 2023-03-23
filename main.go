package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Job struct {
	Log        string //log类型
	JobID      int
	User       string
	JobName    string
	CTime      time.Time
	QTime      time.Time
	ETime      time.Time
	StartTime  time.Time //开始运行时刻
	EndTime    time.Time //结束时刻
	ExitStatus string    //程序退出码
	Platform   string
	Exec_host  string
}

func main() {
	// 解析命令行参数
	logFilePath := "./test_data/testdata"
	dbConfigPath := "./test_db_info.json"
	// ...

	// 读取数据库配置
	db, err := readDBConfig(dbConfigPath)
	if err != nil {
		log.Fatalf("Failed to read database configuration: %v", err)
	}
	defer db.Close()

	// 打开日志文件
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	// 创建数据库表格
	if err := createDBTable(db); err != nil {
		log.Fatalf("Failed to create database table: %v", err)
	}

	// 逐行读取日志文件并分析
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		job, err := parseJob(line)
		if err != nil {
			log.Printf("Failed to parse job from line: %s", line)
			continue
		}

		// 将作业信息写入数据库
		if err := writeJobToDB(db, job); err != nil {
			log.Printf("Failed to write job to database: %v", err)
		}
	}
}

func readDBConfig(configPath string) (*sql.DB, error) {
	// 读取配置文件
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config struct {
		MySQLIP   string `json:"mysql_ip"`
		MySQLPort int    `json:"mysql_port"`
		MySQLUser string `json:"mysql_user"`
		MySQLPass string `json:"mysql_pass"`
		DefaultDB string `json:"default_db"`
	}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	// 连接数据库
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.MySQLUser, config.MySQLPass, config.MySQLIP, config.MySQLPort, config.DefaultDB)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func createDBTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS hpc_job (
			job_id VARCHAR(32),
			user VARCHAR(64),
			job_name VARCHAR(255),
			log VARCHAR(4),
			c_time DATETIME,
			q_time DATETIME,
			e_time DATETIME,
			start_time DATETIME,
			end_time DATETIME,
			exit_status VARCHAR(4),
			platform VARCHAR(32),
		    exec_host VARCHAR(4096),
			PRIMARY KEY (job_id)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create job table: %w", err)
	}
	return nil
}

func parseJob(line string) (*Job, error) {

	match := strings.FieldsFunc(line, func(r rune) bool {
		return r == ';' || r == ' '
	})
	if len(match) < 15 {
		return nil, nil //这行不用解析 返回nil 然后在之后的writeToJob中不作处理
	}

	//处理log类型
	if match[2] == "D" || match[2] == "L" {
		return nil, nil //如果是D或L就返回nil不对这行作处理
	}
	//如果不是D或L，就继续执行 然后在后面会把match[2]填到这个类型里

	//处理jobid
	str := (strings.Split(match[3], "."))[0]
	jobid, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println(err)
	}

	//处理log主体里到内容
	user_string := ""
	jobname_string := ""
	ctime_string := time.Time{}
	qtime_string := time.Time{}
	etime_string := time.Time{}
	start_string := time.Time{}
	end_string := time.Time{}
	exit_status := ""
	platform_str := ""
	exechost_str := ""
	for _, content := range match {
		//
		if !strings.Contains(content, "=") {
			continue
		}
		if strings.Contains(content, "Resource_List") {
			//Resource_List.select=3:ncpus=36:mpiprocs=36:pas_applications_enabled=Dyna:platform=G6254-36c-384G
			strList := strings.Split(content, ":")
			for _, str := range strList {
				if strings.Contains(str, "platform") {
					platform_str = (strings.Split(str, "="))[1]
				}
			}
		}
		key := (strings.Split(content, "="))[0]
		value := (strings.Split(content, "="))[1]
		switch key {
		case "user":
			user_string = value
		case "jobname":
			jobname_string = value
		case "ctime":
			ctime_string, err = parseTimeField(value)
		case "qtime":
			qtime_string, err = parseTimeField(value)
		case "etime":
			etime_string, err = parseTimeField(value)
		case "start":
			start_string, err = parseTimeField(value)
		case "end":
			end_string, err = parseTimeField(value)
		case "Exit_status":
			exit_status = value
		case "exec_host":
			exechost_str = value
		default:
			// 默认代码块
		}

		//
	}
	job := &Job{

		JobID:      jobid,
		User:       user_string,    //user=bilizhi
		JobName:    jobname_string, //jobname=000_Main_C009_h3d
		Log:        match[2],
		CTime:      ctime_string, //ctime=1678037519
		QTime:      qtime_string, //qtime=1678037519
		ETime:      etime_string, //etime=1678037541
		StartTime:  start_string, //start=1678037544
		EndTime:    end_string,   //etime=1678037541
		ExitStatus: exit_status,  //Exit_status=0
		Platform:   platform_str,
		Exec_host:  exechost_str,
	}

	return job, nil
}

func writeJobToDB(db *sql.DB, job *Job) error {
	if job == nil {
		return nil
	}
	stmt, err := db.Prepare(`
    INSERT INTO hpc_job (job_id, user, job_name, log, c_time, q_time, e_time, start_time, end_time, exit_status, platform, exec_host) 
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    ON DUPLICATE KEY UPDATE
        user = VALUES(user),
        job_name = VALUES(job_name),
        log = VALUES(log),
        c_time = VALUES(c_time),
        q_time = VALUES(q_time),
        e_time = VALUES(e_time),
        start_time = VALUES(start_time),
        end_time = VALUES(end_time),
        exit_status = VALUES(exit_status),
        platform = VALUES(platform),
		exec_host = VALUES(exec_host);
`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		job.JobID,
		job.User,
		job.JobName,
		job.Log,
		job.CTime,
		job.QTime,
		job.ETime,
		job.StartTime,
		job.EndTime,
		job.ExitStatus,
		job.Platform,
		job.Exec_host,
	)
	if err != nil {
		return fmt.Errorf("failed to write job to database: %w", err)
	}
	return nil
}

func parseTimeField(timeStr string) (time.Time, error) {
	// 解析日志文件中的时间字段
	// 将时间戳字符串转换为int64类型的时间戳
	timestamp, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	// 使用time.Unix()方法将时间戳转换为time.Time类型的时间对象
	return time.Unix(timestamp, 0), nil
}
