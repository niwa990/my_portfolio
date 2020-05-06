package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v9"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Name1    string `validate:"required"`
	Name2    string `validate:"required"`
	Mail     string `validate:"required,email"`
	Admin    bool
	Password string `validate:"required,gte=5"`
}

type Users struct {
	db *sql.DB
}

// 新しいUsersを作成する
func NewUsers(db *sql.DB) *Users {
	return &Users{db: db}
}

// セッション
var store = sessions.NewCookieStore([]byte("secret-password"))

// テーブルがなかったら作成する
func (us *Users) CreateUsersTable() error {
	const sqlStr = `CREATE TABLE IF NOT EXISTS users(
		id         INTEGER PRIMARY KEY AUTO_INCREMENT,
		name1      TEXT NOT NULL,
		name2      TEXT NOT NULL,
		mail       TEXT NOT NULL,
		password   TEXT NOT NULL,
    admin      BOOL NOT NULL default false
	);`
	_, err := us.db.Exec(sqlStr)
	if err != nil {
		return err
	}

	return nil
}

// User取得
// func (us *Users) GetUsers(limit int) ([]*User, error) {
// }

// サインアップ
func (us *Users) SignUpUser(user *User) error {
	validate = validator.New()
	err := validate.Struct(user)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		// バリデーションエラーの場合の処理
		return
	}

	// パスワードのハッシュを生成
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	const sqlStr = `INSERT INTO users(name1, name2, mail, password, admin) VALUES (?,?,?,?,?);`
	_, err = us.db.Exec(sqlStr, user.Name1, user.Name2, user.Mail, user.Password, user.Admin)
	if err != nil {
		return err
	}

	user.Password = ""
	return nil
}

// ログイン
func (us *Users) LoginUser(user *User, w http.ResponseWriter, r *http.Request) error {
	validate = validator.New()
	err := validate.Struct(user)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		// バリデーションエラーの場合の処理
		return
	}

	// 入力値と登録値を比較する。
	inputPassword := user.Password
	var hashedPassword string
	const sqlStr = `SELECT password FROM users WHERE mail = ?;`
	err := us.db.QueryRow(sqlStr, user.Mail).Scan(&hashedPassword)
	if err != nil {
		return err
	}

	// 比較
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	if err != nil {
		return err
	}

	// セッション登録
	err = SessionCreate(w, r)
	if err != nil {
		return err
	}

	return nil
}

// session作成
func SessionCreate(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "session")
	if err != nil {
		return err
	}

	mail := r.FormValue("mail")
	if mail != "" {
		session.Values["mail"] = mail
	}

	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

// ログアウト
func (us *Users) LogoutUser(w http.ResponseWriter, r *http.Request) error {
	// セッション削除
	session, err := store.Get(r, "session")
	if err != nil {
		return err
	}

	// セッション情報のクリア
	session.Options = &sessions.Options{MaxAge: -1, Path: "/"}

	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}
