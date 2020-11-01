// 此处存放所有与数据库交互的sql命令.
// sqlite里还写有多个触发器, 有些还是递归的
// 请务必与此处的sql配合阅读.
package data

import (
	"github.com/jmoiron/sqlx"
)

// 通过markdown添加一条article
func CreateArticle(md, title string, dirId int) (id int, err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		res, err := tx.Exec("insert into article(markdown, title) values (?, ?)", md, title)
		if err != nil {
			return
		}
		id = getSqlResultId(&res)
		_, err = tx.Exec("update layer set dirId = ? where id = ? and type = 0", dirId, id)
		return
	})
	return
}

// 通过name添加一条dir
func CreateDir(name string, dirId int) (id int, err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		res, err := tx.Exec("insert into dir(name) values (?)", name)
		if err != nil {
			return
		}
		id = getSqlResultId(&res)
		_, err = tx.Exec("update layer set dirId = ? where id = ? and type = 1", dirId, id)
		return
	})
	return
}

// 通过dirId获取一条dir
func GetDir(dirId int) (dir Dir, err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		err = tx.Get(&dir, "select * from dir where id=?", dirId)
		return
	})
	return
}

// 通过articleId获取一条article
func GetArticle(articleId int) (article Article, err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		err = tx.Get(&article, "select id, title, markdown, html from article where id=?", articleId)
		return
	})
	return
}

// 通过dirId获取一层layer (使用视图layer_read)
func GetLayer(dirId int) (contents []LayerContent, err error) {
	contents = make([]LayerContent, 0)
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		err = tx.Select(&contents, "select id, type, text, createdT, modifiedT from layer_read where dirId=?", dirId)
		return
	})
	return
}

// 获取某Entity所属的dir. 注意不能为{1,0}(根目录)
func GetParentDir(entity Entity) (dirId int, err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		err = tx.Get(&dirId, "select dirId from layer where id=? and type=?", entity.Id, entity.Type)
		return
	})
	return
}

// 获得某条article的drewT和createdT和modifiedT
func GetArticleTime(articleId int) (drewT, createdT, modifiedT string, err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		err = tx.QueryRow("select drewT, createdT, modifiedT from article A, layer L where A.id=L.id and A.id=? and L.type=0", articleId).Scan(&drewT, &createdT, &modifiedT)
		return
	})
	return
}

// 获得某条dir的createdT和modifiedT
func GetDirTime(dirId int) (createdT, modifiedT string, err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		err = tx.QueryRow("select createdT, modifiedT from layer where id=? and type=1", dirId).Scan(&createdT, &modifiedT)
		return
	})
	return
}

// 修改某条article的markdown
func SetArticleMarkdownAndTitle(articleId int, md, title string) (err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		_, err = tx.Exec("update article set markdown=?, title=? where id=?", md, title, articleId)
		return
	})
	return
}

// 修改某条article的title和html
func SetArticleHtml(articleId int, htmlres string) (err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		_, err = tx.Exec("update article set html=? where id=?", htmlres, articleId)
		return
	})
	return
}

// 修改某条dir的name
func SetDirName(dirId int, name string) (err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		_, err = tx.Exec("update dir set name=? where id=?", name, dirId)
		return
	})
	return
}

// 修改某entity的layer dirId
func SetEntityLayer(entity Entity, dirId int) (err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		_, err = tx.Exec("update layer set dirId=? where id=? and type=?", dirId, entity.Id, entity.Type)
		return
	})
	return
}

// 删除某条dir
func DeleteDir(dirId int) (err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		_, err = tx.Exec("delete from dir where id=?", dirId)
		return
	})
	return
}

// 删除某条article
func DeleteArticle(articleId int) (err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		_, err = tx.Exec("delete from article where id=?", articleId)
		return
	})
	return
}
