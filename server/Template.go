package server

import (
	"html/template"
	"io"
	"io/ioutil"

	"github.com/rickonono3/m2obj"

	"rickonono3/r-blog/helper/userhelper"
	"rickonono3/r-blog/objects"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

// Render 执行渲染, 在执行前在data中加入Global字段, 填充全局可用的各种属性供模板使用
//
// 注意dataI 必须为 *m2obj.Object 类型
func (t *Template) Render(w io.Writer, name string, dataI interface{}, c echo.Context) error {
	data := dataI.(*m2obj.Object)
	isMobile, isSpecial := userhelper.CheckUA(c.Request().UserAgent())
	isAdmin := userhelper.CheckAdmin(c)
	data.Set("Global", m2obj.Group{
		"IsMobile":  isMobile,
		"IsSpecial": !isMobile && isSpecial,
		"IsAdmin":   isAdmin,
		"Config":    objects.Config,
	})
	return t.templates.ExecuteTemplate(w, name, data.Staticize())
}

// listTemplateFiles 是为了解决template包下的ParseBlob方法只能识别仅仅一层目录下的模板文件的问题而编写的, 逻辑是将./view/下的所有文件递归列出, 最后生成一个路径列表, 这样目录套目录的结构也可以识别成模板.
func listTemplateFiles(dirName string) (list []string) {
	rootDirName := objects.CWD + "view/"
	list = make([]string, 0)
	if files, err := ioutil.ReadDir(rootDirName + dirName); err == nil {
		for _, file := range files {
			if file.IsDir() {
				list = append(list, listTemplateFiles(dirName+file.Name()+"/")...)
			} else {
				list = append(list, rootDirName+dirName+file.Name())
			}
		}
	} else {
		panic("Can not load views!")
	}
	return
}

// updateRender 初始化Echo框架模板接口
func updateRender() {
	E.Renderer = &Template{
		templates: template.Must(template.New("views").Funcs(sprig.FuncMap()).ParseFiles(listTemplateFiles("")...)),
	}
}
