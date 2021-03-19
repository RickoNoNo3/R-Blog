package objects

import (
	"html/template"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	"rickonono3/r-blog/mytype"
)

func init() {
	rand.Seed(time.Now().Unix())
	Cache = mytype.NewGroup(mytype.ObjectList{
		"AdminHash": mytype.NewValue(strconv.Itoa(rand.Int())),
	})
	Config = mytype.NewGroup(mytype.ObjectList{
		"Cwd": mytype.NewValue(makeCwd()),
		"Blog": mytype.NewGroup(mytype.ObjectList{
			// "CDN":   mytype.NewValue("http://81.70.142.172/"),
			"CDN":   mytype.NewValue("http://localhost:13808/"),
			"BGImg": mytype.NewValue("img/bg.jpg"),
			"Icons": mytype.NewGroup(mytype.ObjectList{
				"Dir":     mytype.NewValue(template.HTML("&#xe806;")),
				"Article": mytype.NewValue(template.HTML("&#xe809;")),
				"File":    mytype.NewValue(template.HTML("&#xe809;")),
			}),
		}),
		"AdminPSWD": mytype.NewValue("tydhc199963+"),
		"IsInDebug": mytype.NewValue(true),
	})
	if blogRoot := os.Getenv("BlogRoot"); blogRoot != "" {
		blogRoot = regexp.MustCompile("/?$").ReplaceAllString(blogRoot, "/")
		Config.Set("Cwd", mytype.NewValue(blogRoot))
	}
}
