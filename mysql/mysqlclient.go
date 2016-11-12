/*
 *       Filename:  mysqlclient.php
 *    Description:
 *         Author:  liuyoubin@gumpcome.com
 *        Created:  2016-08-27 17:37:33
 */
package mysqlclient

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

type MysqlClient struct{
    db *sql.DB
    host string
    port int
    username string
    password string
    database string
}

type MysqlError struct{
    Errno int
    Errmsg string
    Err error
}

func (e *MysqlError) Error() string{
	errstr := "Errno:" + string(e.Errno)+",Errmsg:" + e.Errmsg;
    if e.Err != nil {
        errstr = errstr + ",Error:" + e.Err.Error()
    }
    return errstr
}

func CheckError(err error){
    if err != nil {
		fmt.Println(err)
        panic(err)
    }
}

func NewMysqlClient(host string, port int,
	username string, password string, database string) (client *MysqlClient){
    client = &MysqlClient{}
    client.database = database
    client.host = host
    client.port = port
    client.username = username
    client.password = password
    client.Connect()
    return client
}

func (client *MysqlClient) Connect() {
    
    if client.host == "" {
        CheckError(&MysqlError{Errno:1,Errmsg:"host is not setted"})
    }
    var err error
    host_intro := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&collation=utf8_general_ci", client.username,
        client.password, client.host, client.port, client.database)
	client.db, err = sql.Open("mysql", host_intro)
    CheckError(err)

    err = client.db.Ping()
    CheckError(err)
}

func (client MysqlClient) Close(){
    client.db.Close()
}

func (client MysqlClient) CheckConnect(){
    err := client.db.Ping();
    if err != nil {
        client.Connect()
    }
}

func (client MysqlClient) Query(sql string, args ...interface{}) (result []map[string]string){

    client.CheckConnect()

    var err error

    stmt, err := client.db.Prepare(sql)
    CheckError(err)

    rows, err := stmt.Query(args...)
    CheckError(err)
	
	defer func(){
		stmt.Close()
	}()
    columns, _ := rows.Columns()

    scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
    result = make([]map[string]string, 0)

	for i := range values {
		scanArgs[i] = &values[i]
	}

    i := 0
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
            if col == nil {
                record[columns[i]] = "";
				continue
            }
            v, ok := col.(int64)
            if ok {
                record[columns[i]] = fmt.Sprintf("%d", v)
                continue
            }
			record[columns[i]] = string(col.([]byte))
		}
		result = append(result, record)
        i++
	}
    
    return result
}


func (client MysqlClient) Operate(sql string, args ...interface{}) (affectedRows int64){
	
    client.CheckConnect()

    var err error

    stmt, err := client.db.Prepare(sql)
    CheckError(err)

	defer func(){
		stmt.Close()
	}()

    res, err := stmt.Exec(args...)
    CheckError(err)

	affectedRows, err = res.RowsAffected()
	CheckError(err)

	return affectedRows
}

