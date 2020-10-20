package app

import (
	"database/sql"
	"fmt"
	"log"
)

type Skills struct {
	db *sql.DB
}

type Skill struct {
	ID        int
	SkillName string `validate:"required"`
	Period    int    `validate:"required"`
}

func NewSkills(db *sql.DB) *Skills {
	return &Skills{db: db}
}

// テーブルがなかったら作成する
func (sk *Skills) CreateSkillsTable() error {
	const sqlStr = `CREATE TABLE IF NOT EXISTS skills(
		id         INTEGER PRIMARY KEY AUTO_INCREMENT,
		skillname  TEXT NOT NULL,
		period     INTEGER NOT NULL
	);`

	_, err := sk.db.Exec(sqlStr)
	if err != nil {
		return err
	}

	return nil
}

func (sk *Skills) GetSkill(id string) ([]*Skill, error) {
	// LIMITで件数を最大の取得する件数を絞る
	const sqlStr = `SELECT * FROM skills WHERE id = ?`
	rows, err := sk.db.Query(sqlStr, id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close() // 関数終了時にCloseが呼び出される

	var skills []*Skill
	for rows.Next() {
		var skill Skill
		err := rows.Scan(&skill.ID, &skill.SkillName, &skill.Period)
		if err != nil {
			return nil, err
		}
		skills = append(skills, &skill)
	}

	return skills, nil
}

func (sk *Skills) GetSkills(limit int) ([]*Skill, error) {
	// LIMITで件数を最大の取得する件数を絞る
	const sqlStr = `SELECT * FROM skills ORDER BY id DESC LIMIT ?`
	rows, err := sk.db.Query(sqlStr, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 関数終了時にCloseが呼び出される

	var skills []*Skill
	for rows.Next() {
		var skill Skill
		err := rows.Scan(&skill.ID, &skill.SkillName, &skill.Period)
		if err != nil {
			return nil, err
		}
		skills = append(skills, &skill)
	}

	return skills, nil
}

// データベースに新しいItemを追加する
func (sk *Skills) AddSkill(skill *Skill) error {
	const sqlStr = `INSERT INTO skills(skillname, period) VALUES (?,?);`
	_, err := sk.db.Exec(sqlStr, skill.SkillName, skill.Period)
	if err != nil {
		return err
	}
	return nil
}

// Skillを更新する
func (sk *Skills) UpdateSkill(skill *Skill) error {
	// 排他処理 楽観的排他制御, idを取得
	tx, err := sk.db.Begin()
	if err != nil {
		return err
	}

	const sqlStr = `UPDATE skills SET skillname = ?, period = ? WHERE id = ?;`
	_, err = tx.Exec(sqlStr, skill.SkillName, skill.Period, skill.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Skillを削除する
func (sk *Skills) DeleteSkill(id string) error {
	tx, err := sk.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("delete")

	// 削除区分更新
	const sqlStr = `DELETE FROM skills WHERE id = ?;`
	_, err = tx.Exec(sqlStr, id)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
