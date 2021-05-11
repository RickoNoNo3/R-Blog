// Package data
// 此处存放所有与数据库交互的sql命令.
// sqlite里还写有多个触发器, 有些还是递归的
// 请务必与此处的sql配合阅读.
package data

import (
	"errors"

	"github.com/jmoiron/sqlx"

	"rickonono3/r-blog/mytype"
)

// GetEntityInfo
//
// 通过type和id获取一个Entity
//
// Entity和Dir/Article/File属于两个层次,
// Entity代表通用的基础信息
func GetEntityInfo(tx *sqlx.Tx, entityType, entityId int) (entity mytype.Entity, err error) {
	err = tx.Get(&entity, "SELECT id, type, title, createdT, modifiedT FROM layer_read WHERE type=? AND id=?", entityType, entityId)
	return
}

// GetDir 通过dirId获取一个dir
func GetDir(tx *sqlx.Tx, dirId int) (dir mytype.Dir, err error) {
	if dirId == 0 {
		return mytype.Dir{
			Entity: mytype.Entity{
				Id:    0,
				Type:  1,
				Title: "博客",
			},
		}, nil
	}
	dir.Entity, err = GetEntityInfo(tx, 0, dirId)
	return
}

// GetParentDir 获取某Entity所属的父dir, 注意不能为{0,0}(根目录)
func GetParentDir(tx *sqlx.Tx, entity mytype.Entity) (dirId int, err error) {
	if entity.Type == 0 && entity.Id == 0 {
		err = errors.New("get parent dir of {0,0}")
		return
	}
	err = tx.Get(&dirId, "SELECT dirId FROM layer WHERE id=? AND type=?", entity.Id, entity.Type)
	return
}

// GetArticle 通过articleId获取一个article
func GetArticle(tx *sqlx.Tx, articleId int) (article mytype.Article, err error) {
	article.Entity, err = GetEntityInfo(tx, 1, articleId)
	if err != nil {
		return
	}
	err = tx.Get(&article, "SELECT markdown, tags, voted, visited FROM article WHERE id=?", articleId)
	return
}

// GetFile 通过fileId获取一个file
func GetFile(tx *sqlx.Tx, fileId int) (file mytype.File, err error) {
	file.Entity, err = GetEntityInfo(tx, 2, fileId)
	return
}

// GetContents 通过dirId获取一层contents (使用视图layer_read)
func GetContents(tx *sqlx.Tx, dirId int) (contents []mytype.Entity, err error) {
	contents = make([]mytype.Entity, 0)
	err = tx.Select(&contents, "SELECT id, type, title, createdT, modifiedT FROM layer_read WHERE dirId=?", dirId)
	return
}

// GetWelcome 获取“建博客的一些心得”
func GetWelcome(tx *sqlx.Tx) (welcome string, err error) {
	err = tx.Get(&welcome, "SELECT markdown FROM article WHERE id=0")
	return
}

// CreateDir 创建目录
func CreateDir(tx *sqlx.Tx, title string, parentDirId int) (dirId int, err error) {
	res, err := tx.Exec("INSERT INTO dir (title) VALUES (?)", title)
	if err != nil {
		return
	}
	dirId64, err := res.LastInsertId()
	if err != nil {
		return
	}
	dirId = int(dirId64)
	err = CreateLayer(tx, mytype.Entity{
		Id:   dirId,
		Type: 0,
	}, parentDirId)
	return
}

// CreateArticle 创建文章
func CreateArticle(tx *sqlx.Tx, title, markdown string, dirId int) (articleId int, err error) {
	res, err := tx.Exec("INSERT INTO article (title, markdown) VALUES (?, ?)", title, markdown)
	if err != nil {
		return
	}
	articleId64, err := res.LastInsertId()
	if err != nil {
		return
	}
	articleId = int(articleId64)
	err = CreateLayer(tx, mytype.Entity{
		Id:   articleId,
		Type: 1,
	}, dirId)
	return
}

// CreateFile 创建文件
func CreateFile(tx *sqlx.Tx, filename string, dirId int) (fileId int, err error) {
	res, err := tx.Exec("INSERT INTO file (title) VALUES (?)", filename)
	if err != nil {
		return
	}
	fileId64, err := res.LastInsertId()
	if err != nil {
		return
	}
	fileId = int(fileId64)
	err = CreateLayer(tx, mytype.Entity{
		Id:   fileId,
		Type: 2,
	}, dirId)
	return
}

// CreateLayer 创建一层(并更新上层时间)
func CreateLayer(tx *sqlx.Tx, entity mytype.Entity, dirId int) (err error) {
	_, err = tx.Exec("INSERT INTO layer (id, type, dirId) VALUES (?,?,?)", entity.Id, entity.Type, dirId)
	if err != nil {
		return
	}
	err = UpdateLayer(tx, entity)
	return
}

// UpdateLayer 更新上层时间
func UpdateLayer(tx *sqlx.Tx, entity mytype.Entity) (err error) {
	if entity.Type == 0 && entity.Id == 0 {
		err = errors.New("get parent dir of {0,0}")
		return
	}
	tEntity := entity
	for {
		_, err = tx.Exec("UPDATE layer SET modifiedT=CURRENT_TIMESTAMP WHERE id=? AND type=?", tEntity.Id, tEntity.Type)
		if err != nil {
			return
		}
		var dirId int
		dirId, err = GetParentDir(tx, tEntity)
		if err != nil {
			return
		}
		tEntity = mytype.Entity{
			Id:   dirId,
			Type: 0,
		}
		if tEntity.Type == 0 && tEntity.Id == 0 {
			break
		}
	}
	return
}

// MoveLayer 移动层到其他dirId, 并把dirId的时间更新
func MoveLayer(tx *sqlx.Tx, entity mytype.Entity, dirId int) (err error) {
	// 处理移动目录时的一些问题
	if entity.Type == 0 {
		// 待移动目录为目标目录本身, 无效操作
		if entity.Id == dirId {
			return
		}
		// 目标目录为待移动目录的直接父目录, 无效操作
		if entity.Id != 0 {
			tDirId := -1
			tDirId, err = GetParentDir(tx, entity)
			if err != nil {
				return
			}
			if tDirId == dirId {
				return
			}
		}
		// 待移动目录为目标目录的祖先(待移动目录为0时恒为真), 报错
		// 不断获取目标目录 dirId 的父目录 tDirId, 直到根目录
		// 目标目录的任何父目录 tDirId 不能是待移动目录 entity.Id
		{
			tDirId := dirId
			for tDirId != 0 {
				// 获取父目录
				tDirId, err = GetParentDir(tx, mytype.Entity{
					Id:   tDirId,
					Type: 0,
				})
				if err != nil {
					return
				}
				// 父目录Id与待移动目录Id比较, 相同则返回错误
				if tDirId == entity.Id {
					err = errors.New("move into it's child")
					return
				}
			}
		}
	}
	_ = UpdateLayer(tx, mytype.Entity{
		Id:   dirId,
		Type: 0,
	})
	_, err = tx.Exec("UPDATE layer SET dirId=? WHERE id=? AND type=?", dirId, entity.Id, entity.Type)
	return
}

// RemoveLayer 删除层, 并把时间更新
func RemoveLayer(tx *sqlx.Tx, entity mytype.Entity) (err error) {
	if entity.Type == 0 && entity.Id == 0 {
		err = errors.New("remove of {0,0}")
		return
	}
	_ = UpdateLayer(tx, entity)
	// BFS获取所有待删除项, 存入 removedList
	removedList := make([]mytype.Entity, 0)
	{
		tEntity := entity
		que := mytype.EntityQueue{}
		que.Enqueue(tEntity)
		for !que.IsEmpty() {
			if front, ok := que.Dequeue(); ok {
				removedList = append(removedList, front)
				if front.Type == 0 {
					if contents, err := GetContents(tx, front.Id); err == nil {
						for _, v := range contents {
							que.Enqueue(v)
						}
					} else {
						break
					}
				}
			} else {
				break
			}
		}
	}
	// 删除 removedList 中的每一项的本体数据和层数据
	for _, v := range removedList {
		if _, err = tx.Exec("DELETE FROM layer WHERE id=? AND type=?", v.Id, v.Type); err != nil {
			return
		}
		switch v.Type {
		case 0:
			_, err = tx.Exec("DELETE FROM dir WHERE id=?", v.Id)
			if err != nil {
				return
			}
		case 1:
			_, err = tx.Exec("DELETE FROM article WHERE id=?", v.Id)
			if err != nil {
				return
			}
		case 2:
			_, err = tx.Exec("DELETE FROM file WHERE id=?", v.Id)
			if err != nil {
				return
			}
		}
	}
	return
}

func EditDir(tx *sqlx.Tx, id int, title string) (err error) {
	_, err = tx.Exec("UPDATE dir SET title=? WHERE id=?", title, id)
	if err != nil {
		return
	}
	err = UpdateLayer(tx, mytype.Entity{
		Id:    id,
		Type:  0,
		Title: title,
	})
	return
}

func EditArticle(tx *sqlx.Tx, id int, title, markdown string) (err error) {
	_, err = tx.Exec("UPDATE article SET title=?, markdown=? WHERE id=?", title, markdown, id)
	if err != nil {
		return
	}
	err = UpdateLayer(tx, mytype.Entity{
		Id:    id,
		Type:  1,
		Title: title,
	})
	return
}

func EditFile(tx *sqlx.Tx, id int, title string) (err error) {
	_, err = tx.Exec("UPDATE file SET title=? WHERE id=?", title, id)
	if err != nil {
		return
	}
	err = UpdateLayer(tx, mytype.Entity{
		Id:    id,
		Type:  2,
		Title: title,
	})
	return
}

func GetAllIdByType(tx *sqlx.Tx, entityType int) (list []int, err error) {
	list = make([]int, 0)
	err = tx.Select(&list, "SELECT id FROM layer WHERE type=? ORDER BY id", entityType)
	return
}

func GetAllMarkdownInArticle(tx *sqlx.Tx) (list []string, err error) {
	list = make([]string, 0)
	err = tx.Select(&list, "SELECT markdown FROM article")
	return
}
