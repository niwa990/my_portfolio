package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/niwa990/my_portfolio/app"
)

// rootpath
func (skillHandle *Handlers) ListHandler(w http.ResponseWriter, r *http.Request) {
	skills, err := skillHandle.sk.GetSkills(10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isLogin, err := GetLoginState(r)
	if err != nil {
		return
	}

	var vd = &ViewData{Skill: skills, TempNum: 1, IsLogin: isLogin}

	// 取得したskillsをテンプレートに埋め込む
	template := loadTemplate()
	err = template.Execute(w, vd)
	if err != nil {
		return
	}
}

// Admin
func (skillHandle *Handlers) AdminHandler(w http.ResponseWriter, r *http.Request) {
	isLogin, err := GetLoginState(r)
	if err != nil {
		return
	}

	// ログインしていない場合root
	if !isLogin {
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
	template := loadTemplate()
	err = template.Execute(w, vd)
	if err != nil {
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

	isLogin, err := GetLoginState(r)
	if err != nil {
		return
	}

	skillname := r.FormValue("skillname")
	var period int
	if r.FormValue("period") != "" {
		period, err = strconv.Atoi(r.FormValue("period"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	skill := &app.Skill{
		SkillName: skillname,
		Period:    period,
	}

	var validateErr []string
	validateErr = app.SkillValidate(skill)
	if len(validateErr) == 0 {
		if err = skillHandle.sk.AddSkill(skill); err != nil {
			validateErr = append(validateErr, "")
			return
		}
	}

	var skills []*app.Skill
	skills = append(skills, skill)

	http.Redirect(w, r, "/", http.StatusFound)

	// 取得したskillをテンプレートに埋め込む
	var vd = &ViewData{Skill: skills, TempNum: 3, IsLogin: isLogin, Messages: validateErr}
	template := loadTemplate()
	err = template.Execute(w, vd)
	if err != nil {
		return
	}
}

func (skillHandle *Handlers) EditHandler(w http.ResponseWriter, r *http.Request) {
	// URLからIDを取得
	requestPath := strings.Split(r.URL.Path, "/")
	id := requestPath[2]

	skill, err := skillHandle.sk.GetSkill(id)
	if err != nil {
		return
	}

	isLogin, err := GetLoginState(r)
	if err != nil {
		return
	}

	var vd = &ViewData{Skill: skill, TempNum: 3, IsLogin: isLogin}

	// 取得したskillをテンプレートに埋め込む
	template := loadTemplate()
	err = template.Execute(w, vd)
	if err != nil {
		return
	}
}

// Skill Update
func (skillHandle *Handlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		return
	}

	// URLからIDを取得
	requestPath := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(requestPath[2]) // 固定値
	if err != nil {
		return
	}

	skillname := r.FormValue("skillname")
	var period int

	if r.FormValue("period") != "" {
		period, err = strconv.Atoi(r.FormValue("period"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	skill := &app.Skill{
		ID:        id,
		SkillName: skillname,
		Period:    period,
	}

	var validateErr []string
	validateErr = app.SkillValidate(skill)
	if len(validateErr) == 0 {
		if err := skillHandle.sk.UpdateSkill(skill); err != nil {
			validateErr = append(validateErr, "")
			return
		}
	}

	isLogin, err := GetLoginState(r)
	if err != nil {
		return
	}

	var skills []*app.Skill
	skills = append(skills, skill)

	http.Redirect(w, r, "/", http.StatusFound)

	// 取得したskillをテンプレートに埋め込む
	var vd = &ViewData{Skill: skills, TempNum: 3, IsLogin: isLogin, Messages: validateErr}
	template := loadTemplate()
	err = template.Execute(w, vd)
	if err != nil {
		return
	}
}

// Skill Delete
func (skillHandle *Handlers) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	requestPath := strings.Split(r.URL.Path, "/")
	id := requestPath[2]

	if err := skillHandle.sk.DeleteSkill(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
