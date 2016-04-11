package main

import (
	"crypto/md5"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"openvpn-mgr/gencfg"
)

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func OpenDB(dbFile string) (db *sql.DB, err error) {
	return sql.Open("sqlite3", dbFile)
}

type Account struct {
	dbFile string
}

func (a Account) Auth(name, password string) bool {
	db, err := OpenDB(a.dbFile)
	defer db.Close()
	if err != nil {
		return false
	}
	log.SetFlags(log.LstdFlags)
	row := db.QueryRow(fmt.Sprintf(
		`select name from dial_account
         where
            name='%s'
            and password='%s'`,
		name, password))
	var falt string
	row.Scan(&falt)
	if falt == name {
		log.Printf("%s用户认证成功", name)
		return true
	} else {
		log.Printf("%s用户认证失败", name)
		return false
	}
}

func (a Account) UserAdd(name, password string) bool {
	db, err := OpenDB(a.dbFile)
	defer db.Close()
	if err != nil {
		return false
	}

	_, err = db.Exec(fmt.Sprintf(
		`insert into dial_account
         (name,password,created_at)
         values('%s','%s', '%s')`, name, password, time.Now()))
	if err != nil {
		log.Printf("用户%s添加失败\n", name)
		return false
	} else {
		log.Printf("用户%s添加成功\n", name)
		return true
	}
}

func (a Account) UserDel(name string) bool {
	db, err := OpenDB(a.dbFile)
	defer db.Close()
	if err != nil {
		return false
	}
	_, err = db.Exec(fmt.Sprintf(
		`delete from dial_account
         where name='%s'`, name))
	if err != nil {
		log.Printf("用户%s删除失败\n", name)
		return false
	} else {
		log.Printf("用户%s删除成功\n", name)
		return true
	}
}

func (a Account) UserModify(name, password string) bool {
	db, err := OpenDB(a.dbFile)
	defer db.Close()
	if err != nil {
		return false
	}

	_, err = db.Exec(fmt.Sprintf(
		`update dial_account
         set password='%s'
         where name='%s'`, password, name))
	if err != nil {
		log.Printf("用户%s修改失败\n", name)
		return false
	} else {
		log.Printf("用户%s修改成功\n", name)
		return true
	}
}

func (a Account) Search(name string) bool {
	db, err := OpenDB(a.dbFile)
	defer db.Close()
	if err != nil {
		return false
	}

	row := db.QueryRow(fmt.Sprintf(
		`select name from dial_account
         where
            name='%s'`,
		name))
	var falt string
	row.Scan(&falt)
	if falt != "" {
		log.Printf("查找到用户%s\n", name)
		return true
	} else {
		log.Printf("未查找到用户%s\n", name)
		return false
	}
}

func (a Account) ListAllUser() {
	db, err := OpenDB(a.dbFile)
	defer db.Close()
	if err != nil {
		log.Println("显示所有用户")
	}

	row, err := db.Query(`select name from dial_account`)
	if err != nil {
		log.Println("显示所有用户")
	}
	fmt.Println("UserInfo:")
	i := 0
	for row.Next() {
		var username string
		row.Scan(&username)
		i++
		fmt.Printf("%d. %s\n",i,username)
	}
	fmt.Printf("总共 %d 个用户.\n",i)
}

func checkUser(username string) {
	if username == "" {
		log.Println("缺少用户名参数")
		os.Exit(1)
	}
}

func checkPassword(password string) {
	if password == "" {
		log.Println("缺少密码参数")
		os.Exit(1)
	}
}

func main() {
	err := os.Chdir(path.Dir(os.Args[0]))
	Fatal(err)
	dbFile := "instance/website.db"
	fileinfo, err := os.Stat(dbFile)
	Fatal(err)
	if !fileinfo.Mode().IsRegular() {
		Fatal(fmt.Errorf("数据库文件异常."))
	}

	var (
		cmd      string
		username string
		password string
	)
	flag.StringVar(&cmd, "c", "auth", "指定执行命令(auth list search add del modify initcfg")
	flag.StringVar(&username, "u", "", "用户")
	flag.StringVar(&password, "p", "", "密码")
	flag.Parse()

	auth := new(Account)
	auth.dbFile = dbFile

	if password != "" {
		password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	}
	switch cmd {
	case "initcfg":
		gencfg.GenCfg()
	case "auth":
		password = fmt.Sprintf("%x", md5.Sum([]byte(os.Getenv("password"))))
		username = os.Getenv("username")
		checkUser(username)
		checkPassword(password)
		ok := auth.Auth(username, password)
		if ok {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	case "list":
		auth.ListAllUser()
		os.Exit(0)
	case "search":
		checkUser(username)
		if auth.Search(username) {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	case "del":
		checkUser(username)
		if !auth.Search(username) {
			os.Exit(1)
		}
		if auth.UserDel(username) {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	case "modify":
		checkUser(username)
		checkPassword(password)
		if !auth.Search(username) {
			os.Exit(1)
		}
		if auth.UserModify(username, password) {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	case "add":
		checkUser(username)
		checkPassword(password)
		if auth.Search(username) {
			os.Exit(1)
		}
		if auth.UserAdd(username, password) {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	default:
		flag.Usage()
	}

}
