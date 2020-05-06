package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// logを集積する
	// テンプレート化
	// ログイン状態で画面切り替え
	// 画面遷移管理

	// データベースへ接続
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/niwa_portfolio?parseTime=true")
	if err != nil {
		panic(err)
	}

	skills := NewSkills(db)
	if err := skills.CreateSkillsTable(); err != nil {
		log.Fatal(err)
	}

	users := NewUsers(db)
	if err := users.CreateUsersTable(); err != nil {
		log.Fatal(err)
	}

	// template css jsの読み込み
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates/"))))

	// ハンドラの登録
	skillHandle := NewSkillHandlers(skills)
	http.HandleFunc("/", skillHandle.ListHandler)
	http.HandleFunc("/insert", skillHandle.InsertHandler)
	http.HandleFunc("/update", skillHandle.UpdateHandler)
	http.HandleFunc("/delete", skillHandle.DeleteHandler)
	http.HandleFunc("/admin", skillHandle.AdminHandler)

	userHandle := NewUserHandlers(users)
	http.HandleFunc("/user", userHandle.UserHandler)
	http.HandleFunc("/singup", userHandle.SignUpHandler)
	http.HandleFunc("/login", userHandle.LoginHandler)
	http.HandleFunc("/logout", userHandle.LogoutHandler)

	fmt.Println("http://localhost:8090 で起動中...")
	// HTTPサーバを起動する
	log.Fatal(http.ListenAndServe(":8090", nil))
}
