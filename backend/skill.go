package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-playground/validator/v9"
)

type Skill struct {
	ID        int
	SkillName string `validate:"required"`
	Period    int    `validate:"required"`
}

type Skills struct {
	db *sql.DB
}

func NewSkills(db *sql.DB) *Skills {
	return &Skills{db: db}
}

// テーブルがなかったら作成する
func (ab *Skills) CreateSkillsTable() error {
	const sqlStr = `CREATE TABLE IF NOT EXISTS skills(
		id         INTEGER PRIMARY KEY AUTO_INCREMENT,
		skillname  TEXT NOT NULL,
		period     INTEGER NOT NULL
	);`

	_, err := ab.db.Exec(sqlStr)
	if err != nil {
		return err
	}

	return nil
}

// 最近追加したものを最大limit件だけItemを取得する
// エラーが発生したら第2戻り値で返す
func (ab *Skills) GetSkills(limit int) ([]*Skill, error) {
	// ORDER BY id DESCでidの降順（大きい順）=最近追加したものが先にくる
	// LIMITで件数を最大の取得する件数を絞る
	const sqlStr = `SELECT * FROM skills ORDER BY id DESC LIMIT ?`
	rows, err := ab.db.Query(sqlStr, limit)
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return skills, nil
}

// データベースに新しいItemを追加する
func (ab *Skills) AddSkill(skill *Skill) error {
	validate = validator.New()
	err := validate.Struct(skill)
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
		fmt.Println(err)
	}

	const sqlStr = `INSERT INTO skills(skillname, period) VALUES (?,?);`
	_, err := ab.db.Exec(sqlStr, skill.SkillName, skill.Period)
	if err != nil {
		return err
	}
	return nil
}

// Skillを更新する
func (ab *Skills) UpdateSkill(skill *Skill) error {
	// 排他処理 楽観的排他制御, idを取得
	tx, err := ab.db.Begin()
	if err != nil {
		log.Fatal(err)
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
func (ab *Skills) DeleteSkill(id int) error {
	tx, err := ab.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(id)

	// 削除区分更新
	const sqlStr = `DELETE FROM skills WHERE id = ?;`
	_, err = tx.Exec(sqlStr, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
