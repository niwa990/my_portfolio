package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/niwa990/my_portfolio/app"
	"github.com/niwa990/my_portfolio/handlers"
)

func main() {
	// データベースへ接続
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/niwa_portfolio?parseTime=true")
	if err != nil {
		panic(err)
	}

	skills := app.NewSkills(db)
	if err := skills.CreateSkillsTable(); err != nil {
		log.Fatal(err)
	}

	users := app.NewUsers(db)
	if err := users.CreateUsersTable(); err != nil {
		log.Fatal(err)
	}

	// template css jsの読み込み
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))

	// ハンドラの登録
	r := mux.NewRouter()
	skillHandle := handlers.NewSkillHandlers(skills)
	r.HandleFunc("/", skillHandle.ListHandler)
	r.HandleFunc("/insert", skillHandle.InsertHandler)
	r.HandleFunc("/skill/{id:[0-9]+}/edit", skillHandle.EditHandler)
	r.HandleFunc("/skill/{id:[0-9]+}/update", skillHandle.UpdateHandler)
	r.HandleFunc("/skill/{id:[0-9]+}/delete", skillHandle.DeleteHandler)
	r.HandleFunc("/admin", skillHandle.AdminHandler)

	userHandle := handlers.NewUserHandlers(users)
	r.HandleFunc("/user", userHandle.UserHandler)
	r.HandleFunc("/singup", userHandle.SignUpHandler)
	r.HandleFunc("/login", userHandle.LoginHandler)
	r.HandleFunc("/logout", userHandle.LogoutHandler)

	http.Handle("/", r)

	fmt.Println("http://localhost:8090 で起動中...")
	// HTTPサーバを起動する
	log.Fatal(http.ListenAndServe(":8090", nil))
}
