package server

import (
	"html/template"
	"io"
	"io/ioutil"

	"rickonono3/r-blog/config"
	"rickonono3/r-blog/helper/userhelper"
	"rickonono3/r-blog/mytype"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

// Render 执行渲染, 在执行前在data中加入Global字段, 填充全局可用的各种属性供模板使用
func (t *Template) Render(w io.Writer, name string, dataI interface{}, c echo.Context) error {
	data := dataI.(*mytype.Object)
	isMobile, isSpecial := userhelper.CheckUA(c.Request().UserAgent())
	// TODO: IsAdmin
	data.Set("Global", mytype.NewGroup(mytype.ObjectList{
		"IsMobile":  mytype.NewValue(isMobile),
		"IsSpecial": mytype.NewValue(isSpecial),
		"IsAdmin":   mytype.NewValue(false),
		"Config":    config.Config,
	}))
	return t.templates.ExecuteTemplate(w, name, data.Staticize())
}

// listTemplateFiles 将rootDir下的所有文件递归写入一个模板文件路径列表
func listTemplateFiles(dirName string) (list []string) {
	rootDirName := config.Get("Cwd").Val.(string) + "view/"
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
