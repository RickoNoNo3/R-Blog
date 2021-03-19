package admin

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/bloghelper"
	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/mytype"
)

// 参数:
//   Title
//   Option (="edit"|"about")
// dir专属:
//   Dir
//   DirParentId
//   Contents

func Edit(c echo.Context) (err error) {
	entityTypeStr := c.QueryParam("type")
	entityIdStr := c.QueryParam("id")
	entityParentIdStr := c.QueryParam("parentId")
	entityType := typehelper.MustAtoi(entityTypeStr)
	entityId := typehelper.MustAtoi(entityIdStr)
	entityParentId := typehelper.MustAtoi(entityParentIdStr)
	err = data.DoTx(func(tx *sqlx.Tx) (err error) {
		switch entityType {
		case 1:
			if entityId >= 0 { // 是编辑文章
				var article mytype.Article
				if article, err = data.GetArticle(tx, entityId); err == nil {
					var parentId int
					if parentId, err = data.GetParentDir(tx, article.Entity); err == nil {
						return c.Render(http.StatusOK, "admin_edit_article", mytype.NewGroup(mytype.ObjectList{
							"Title":    mytype.NewValue(bloghelper.MakeTitle("编辑文章")),
							"Option":   mytype.NewValue("edit"),
							"IsNew":    mytype.NewValue(false),
							"Article":  mytype.NewValue(article),
							"ParentId": mytype.NewValue(parentId),
						}))
					}
				}
			} else { // 是新建文章
				return c.Render(http.StatusOK, "admin_edit_article", mytype.NewGroup(mytype.ObjectList{
					"Title":    mytype.NewValue(bloghelper.MakeTitle("新建文章")),
					"Option":   mytype.NewValue("edit"),
					"IsNew":    mytype.NewValue(true),
					"ParentId": mytype.NewValue(entityParentId),
				}))
			}
		default:
			var dir mytype.Dir
			if dir, err = data.GetDir(tx, entityId); err == nil {
				var contents []mytype.Entity
				if contents, err = data.GetContents(tx, entityId); err == nil {
					parentId := -1
					if entityId != 0 {
						parentId, _ = data.GetParentDir(tx, dir.Entity)
					}
					return c.Render(http.StatusOK, "admin_edit_dir", mytype.NewGroup(mytype.ObjectList{
						"Title":    mytype.NewValue(bloghelper.MakeTitle(dir.Entity.Title)),
						"Option":   mytype.NewValue("edit"),
						"Dir":      mytype.NewValue(dir),
						"ParentId": mytype.NewValue(parentId),
						"Contents": mytype.NewValue(contents),
					}))
				}
			}
		}
		return
	})
	return err
}
