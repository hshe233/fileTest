package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"io/ioutil"
	"strings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var path string = "D:/temp11111"

	//获取目标路径下的所有文件名
	dir_list := fileList(path)

	//准备好数据库连接
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/mysql?charset=utf8")
	checkErr(err)
	stmt, err := db.Prepare("INSERT loginfo SET file=?,time=?,duration=?")
	checkErr(err)

	//循环解析文件并将结果入库
	for _, file := range dir_list {

		inputFile, inputErr := os.Open(path + "/" + file)
		if inputErr != nil {
			fmt.Println(inputErr)
			continue
		}

		inputReader := bufio.NewReader(inputFile)
		for {
			inputString, readerError := inputReader.ReadString('\n')
			time, duration := parser(inputString)
			if time != "" {
				_, err := stmt.Exec(file, time, duration)
				checkErr(err)
			}
			if readerError == io.EOF {
				break
			}
		}
		inputFile.Close()
	}
	db.Close()
}

func fileList(path string) []string {
	var list []string = nil
	dir_list, e := ioutil.ReadDir(path)

	if e != nil {
		fmt.Println("Reading directory error", ":", e)
		return nil
	}

	for _, v := range dir_list {
		list = append(list, v.Name())
	}

	if list == nil {
		panic("Given directory(" + path + ") is empty!")
	}

	return list
}

func parser(text string) (string, string) {
	duration := Between(text, "FINISH [", "ms] TAG")
	if duration == "" {
		return "", ""
	}
	time := text[0:8]

	return time, duration
}

func Between(str, starting, ending string) string {
	s := strings.Index(str, starting)
	if s < 0 {
		return ""
	}
	s += len(starting)
	e := strings.Index(str[s:], ending)
	if e < 0 {
		return ""
	}
	return str[s: s+e]
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
