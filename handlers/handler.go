package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/niwa990/my_portfolio/app"
)

type Handlers struct {
	sk *app.Skills
	us *app.Users
}

type ViewData struct {
	Skill    []*app.Skill
	User     []*app.User
	TempNum  int
	IsLogin  bool
	Messages []string
}

// Handlersを作成
func NewSkillHandlers(sk *app.Skills) *Handlers {
	return &Handlers{sk: sk}
}

func NewUserHandlers(us *app.Users) *Handlers {
	return &Handlers{us: us}
}

// テンプレート
var templates = make(map[string]*template.Template)

func loadTemplate() *template.Template {
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
		return nil
	}
	return t
}

// セッション
var store = sessions.NewCookieStore([]byte("secret-password"))

func GetLoginState(r *http.Request) (bool, error) {
	session, err := store.Get(r, "session")
	if err != nil {
		return false, err
	}

	return session.Values["mail"] != nil, nil
}
