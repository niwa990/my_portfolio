package skillHandler

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// HTTPハンドラを集めた型
type Handlers struct {
	sk *Skills
}

type ViewData struct {
	Skill   []*Skill
	TempNum int
	IsLogin bool
}

// Handlersを作成
func NewSkillHandlers(sk *Skills) *Handlers {
	return &Handlers{sk: sk}
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

// Skill Create
func (skillHandle *Handlers) InsertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		return
	}

	skillname := r.FormValue("skillname")
	period, err := strconv.Atoi(r.FormValue("period"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	skill := &Skill{
		SkillName: skillname,
		Period:    period,
	}

	if err := skillHandle.sk.AddSkill(skill); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// Skill Update
func (skillHandle *Handlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))

	skillname := r.FormValue("skillname")
	period, err := strconv.Atoi(r.FormValue("period"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	skill := &Skill{
		ID:        id,
		SkillName: skillname,
		Period:    period,
	}

	if err := skillHandle.sk.UpdateSkill(skill); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// Skill Delete
func (skillHandle *Handlers) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))

	if err := skillHandle.sk.DeleteSkill(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
