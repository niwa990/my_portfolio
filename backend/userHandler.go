package userHandler

import (
	"html/template"
	"log"
	"net/http"
)

// HTTPハンドラを集めた型
type Handlers struct {
	us *Users
}

type ViewData struct {
	Skill   []*Skill
	TempNum int
	IsLogin bool
}

// Handlersを作成
func NewUserHandlers(us *Users) *Handlers {
	return &Handlers{us: us}
}

// テンプレート
var templates = make(map[string]*template.Template)

// テンプレートを形成
func loadTemplate() *template.Template {
	// template must を使う？
	t, err := template.ParseFiles(
		"templates/application.html",
		"templates/_header.html",
		"templates/_footer.html",
		// 可変となるテンプレート
		"templates/root.html",
		"templates/user.html",
		"templates/admin.html",
		"templates/login.html",
	)

	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	return t
}

// rootpath
func (skillHandle *Handlers) ListHandler(w http.ResponseWriter, r *http.Request) {
	skills, err := skillHandle.sk.GetSkills(10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := store.Get(r, "session")
	if err != nil {
		return
	}

	isLogin := false
	if session.Values["mail"] != nil {
		isLogin = true
	}

	var vd = &ViewData{Skill: skills, TempNum: 1, IsLogin: isLogin}

	// 取得したskillsをテンプレートに埋め込む
	templates["root"] = loadTemplate()
	err = templates["root"].Execute(w, vd)
	if err != nil {
		log.Fatal("Cannot Get View ", err)
		return
	}
}

// Admin
func (skillHandle *Handlers) AdminHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		return
	}

	// ログインしていない場合root
	if session.Values["mail"] == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// 最新の10件を取得し、skillsに入れる
	skills, err := skillHandle.sk.GetSkills(10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var vd = &ViewData{Skill: skills, TempNum: 2, IsLogin: true}

	// 取得したskillsをテンプレートに埋め込む
	templates["admin"] = loadTemplate()
	err = templates["admin"].Execute(w, vd)
	if err != nil {
		log.Fatal("Cannot Get View ", err)
		return
	}
}

// User
func (userHandle *Handlers) UserHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		return
	}

	// ログインしていない場合root
	if session.Values["mail"] == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	var vd = &ViewData{TempNum: 3, IsLogin: true}

	// 取得したskillsをテンプレートに埋め込む
	templates["root"] = loadTemplate()
	err = templates["root"].Execute(w, vd)
	if err != nil {
		log.Fatal("Cannot Get View ", err)
		return
	}
}

// User Create
func (userHandle *Handlers) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if r.Method != http.MethodPost {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		return
	}

	name1 := r.FormValue("name1")
	name2 := r.FormValue("name2")
	mail := r.FormValue("mail")
	// admin := r.FormValue("admin")
	password := r.FormValue("password")

	user := &User{
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

	http.Redirect(w, r, "/admin", http.StatusFound)
}

// Login
func (userHandle *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session, err := store.Get(r, "session")
		if err != nil {
			return
		}

		isLogin := false
		if session.Values["mail"] != nil {
			isLogin = true
		}

		var vd = &ViewData{TempNum: 4, IsLogin: isLogin}
		templates["login"] = loadTemplate()
		err = templates["login"].Execute(w, vd)
		return

	} else if r.Method != http.MethodPost {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		return
	}

	mail := r.FormValue("mail")
	password := r.FormValue("password")

	user := &User{
		Mail: mail,
		// Admin:    admin,
		Password: password,
	}

	if err := userHandle.us.LoginUser(user, w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// Logout
func (userHandle *Handlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	if err := userHandle.us.LogoutUser(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
