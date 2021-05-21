package datahelper

import (
	"github.com/jmoiron/sqlx"

	"rickonono3/r-blog/data"
	"rickonono3/r-blog/mytype"
)

type DirInfo struct {
	Id    int
	Title string
	Layer int
}

func GetDirInfoList() ([]DirInfo, error) {
	type stackNode struct {
		DirInfo  DirInfo
		Contents []mytype.Entity
		Ite      int
	}
	var (
		stk  = make([]stackNode, 0)
		list = make([]DirInfo, 0)
	)
	err := data.DoTx("递归搜索博客内的目录结构", func(tx *sqlx.Tx) (err error) {
		var rootContents []mytype.Entity
		if rootContents, err = data.GetContents(tx, 0); err == nil {
			rootNode := stackNode{
				DirInfo: DirInfo{
					Id:    0,
					Title: "博客",
					Layer: 0,
				},
				Contents: rootContents,
			}
			stk = append(stk, rootNode)
			list = append(list, rootNode.DirInfo)
			for len(stk) > 0 {
				top := stk[len(stk)-1]
				if top.Ite >= len(top.Contents) {
					stk = stk[:len(stk)-1]
					continue
				}
				entity := top.Contents[top.Ite]
				stk[len(stk)-1].Ite++
				if entity.Type != 0 {
					continue
				}
				if entityContents, err := data.GetContents(tx, entity.Id); err == nil {
					node := stackNode{
						DirInfo: DirInfo{
							Id:    entity.Id,
							Title: entity.Title,
							Layer: top.DirInfo.Layer + 1,
						},
						Contents: entityContents,
						Ite:      0,
					}
					stk = append(stk, node)
					list = append(list, node.DirInfo)
				} else {
					return err
				}
			}
		}
		return
	})
	return list, err
}
