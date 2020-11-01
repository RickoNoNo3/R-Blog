package helper

import (
	"fmt"
	"os"

	"../data"
)

// 从输入流(stdin)获取每一个<type> <id>对, 组装成slice
func GetInputEntityList(inputFile *os.File) (list []data.Entity) {
	list = make([]data.Entity, 0)
	var entityType, id int
	for {
		if _, err := fmt.Fscan(inputFile, &entityType, &id); err != nil {
			break
		}
		list = append(list, data.Entity{
			Type: entityType,
			Id:   id,
		})
	}
	return
}

// 检查一个实体存在
func CheckEntity(entity data.Entity, allowRoot bool) bool {
	if entity.Type == 1 && entity.Id == 0 { // 根目录特判
		return allowRoot
	}
	var err error
	switch entity.Type {
	case 0:
		_, err = data.GetArticle(entity.Id)
	case 1:
		_, err = data.GetDir(entity.Id)
	default:
		return false
	}
	return err == nil
}
