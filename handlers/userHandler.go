package handlers

import (
	"fmt"
	"net/http"

	"github.com/niwa990/my_portfolio/app"
)

// User
func (userHandle *Handlers) UserHandler(w http.ResponseWriter, r *http.Request) {
	isLogin, err := GetLoginState(r)

	// ログインしていない場合root
	if isLogin == false {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	vd := &ViewData{TempNum: 3, IsLogin: true}
	template := loadTemplate()
	err = template.Execute(w, vd)
	if err != nil {
		return
	}
}

// User Create
func (userHandle *Handlers) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// GET --------------------------------------------
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if r.Method == http.MethodPost {
		// POST --------------------------------------------
		name1 := r.FormValue("name1")
		name2 := r.FormValue("name2")
		mail := r.FormValue("mail")
		// admin := r.FormValue("admin")
		password := r.FormValue("password")

		user := &app.User{
			Name1: name1,
			Name2: name2,
			Mail:  mail,
			// Admin:    admin,
			Password: password,
		}

		if err := userHandle.us.SignUpUser(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var validateErr []string
		validateErr = app.UserInsertValidete(user)

		if len(validateErr) == 0 {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		isLogin, err := GetLoginState(r)
		if err != nil {
			return
		}

		vd := &ViewData{TempNum: 2, IsLogin: isLogin, Messages: validateErr}
		template := loadTemplate()
		err = template.Execute(w, vd)
		if err != nil {
			fmt.Println("err")
			return
		}

	} else {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		return
	}
}

// Login
func (userHandle *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// GET --------------------------------------------
		isLogin, err := GetLoginState(r)
		if err != nil {
			return
		}

		var vd = &ViewData{TempNum: 4, IsLogin: isLogin}
		templates := loadTemplate()
		err = templates.Execute(w, vd)
		return
	} else if r.Method == http.MethodPost {
		// POST --------------------------------------------
		mail := r.FormValue("mail")
		password := r.FormValue("password")

		user := &app.User{
			Mail:     mail,
			Password: password,
		}

		var validateErr []string
		validateErr = app.UserLoginValidate(user)

		if len(validateErr) == 0 {
			if err := userHandle.us.LoginUser(user, w, r); err != nil {
				validateErr = append(validateErr, "メールアドレスとパスワードが一致しません。")
			}
		}

		isLogin, err := GetLoginState(r)
		if err != nil {
			return
		}

		// http headerに情報を追加する
		// http status codeを指定する
		// response bodyにデータを書き込む
		if isLogin {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		vd := &ViewData{TempNum: 4, IsLogin: isLogin, Messages: validateErr}
		template := loadTemplate()
		err = template.Execute(w, vd)
		if err != nil {
			fmt.Println("err")
			return
		}
	} else {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		return
	}
}

// Logout
func (userHandle *Handlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	if err := userHandle.us.LogoutUser(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
