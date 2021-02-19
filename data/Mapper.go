// 此处存放所有与数据库交互的sql命令.
// sqlite里还写有多个触发器, 有些还是递归的
// 请务必与此处的sql配合阅读.
package data

import (
	"github.com/jmoiron/sqlx"

	"rickonono3/r-blog/mytype"
)

// 通过type和id获取一个Entity
// Entity和Dir/Article/File属于两个层次的东西
// 一个代表通用的基础信息, 另一个代表特定的详细信息
func GetEntityInfo(entityType, entityId int) (entity mytype.Entity, err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		err = tx.Get(&entity, "select id, type, title, createdT, modifiedT from layer_read where type=? and id=?", entityType, entityId)
		return
	})
	return
}

// 通过dirId获取一个dir
func GetDir(dirId int) (dir mytype.Dir, err error) {
	if dirId == 0 {
		return mytype.Dir{
			Entity: mytype.Entity{
				Id:    0,
				Type:  1,
				Title: "博客",
			},
		}, nil
	}
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		dir.Entity, err = GetEntityInfo(0, dirId)
		// dir暂无详情信息
		// err = tx.Get(&dir, "select * from dir where id=?", dirId)
		return
	})
	return
}

// 获取某Entity所属的父dir, 注意不能为{1,0}(根目录)
func GetParentDir(entity mytype.Entity) (dirId int, err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		err = tx.Get(&dirId, "select dirId from layer where id=? and type=?", entity.Id, entity.Type)
		return
	})
	return
}

// 通过articleId获取一个article
func GetArticle(articleId int) (article mytype.Article, err error) {
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		article.Entity, err = GetEntityInfo(1, articleId)
		if err != nil {
			return
		}
		err = tx.Get(&article, "select markdown, tags, voted, visited from article where id=?", articleId)
		return
	})
	return
}

// 通过dirId获取一层contents (使用视图layer_read)
func GetContents(dirId int) (contents []mytype.Entity, err error) {
	contents = make([]mytype.Entity, 0)
	err = DoTx(func(tx *sqlx.Tx) (err error) {
		err = tx.Select(&contents, "select id, type, title, createdT, modifiedT from layer_read where dirId=?", dirId)
		return
	})
	return
}
