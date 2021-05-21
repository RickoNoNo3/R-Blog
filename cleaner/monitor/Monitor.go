package monitor

import (
	"github.com/jmoiron/sqlx"
	"os"
	"regexp"
	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/datahelper"
	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/logger"
	"sync"
	"time"
)

/*********************************************
 文件监控器定时扫描数据库file表和article表内markdown字段内的资源链接，与resource文件夹的实际情况进行匹配.

   +-----+-------+----------+------------+----------------------+
   |IsTmp|InTable|InResource|描述        |做法                  |
   |-----|-------|----------|------------|----------------------|
   | 否  |  是   |    否    |无效file记录|什么也不做, 等手动删除|
   | 否  |  否   |    是    |多余的file  |删除资源文件          |
   | 是  |  是   |    否    |无效hash引用|什么也不做            |
   | 是  |  否   |    是    |多余的hash  |删除资源文件          |
   +-----+-------+----------+------------+----------------------+

*********************************************/

var (
	filePrompt    = regexp.MustCompile("^file_[0-9]+$")
	hashPrompt    = regexp.MustCompile("^hash_[0-9a-zA-z]+$")
	hashRefPrompt = regexp.MustCompile("hash_[0-9a-zA-z]+")
	resourcePath  = datahelper.GetResourcePathForServer()
	//ignoreNewFile = 30 * time.Second
	//interval      = 10 * time.Second
	ignoreNewFile = 24 * time.Hour
	interval      = 24 * time.Hour
)

func matchFiles(dbFileList []int, dbHashList []string) (matching map[string]bool) {
	matching = make(map[string]bool)
	// open folder
	path, err := os.Open(resourcePath)
	if err != nil {
		return
	}
	// read files from folder
	files, err := path.Readdir(-1)
	if err != nil {
		return
	}
	var (
		fileList    []string
		hashList    []string
		unknownList []string
	)
	for _, file := range files {
		fileName := file.Name()
		nowTime := time.Now()
		modTime := file.ModTime()
		if nowTime.Before(modTime) || nowTime.Sub(modTime) <= ignoreNewFile {
			logger.L.Debug("[CleanerMonitor]", fileName+" 因处于新文件忽略期内而不会被删除")
			matching[fileName] = false
		} else {
			matching[fileName] = true
		}
		if filePrompt.MatchString(fileName) {
			fileList = append(fileList, fileName)
		} else if hashPrompt.MatchString(fileName) {
			hashList = append(hashList, fileName)
		} else {
			unknownList = append(unknownList, fileName)
		}
	}
	// for fileList, if InResource and InTable, not delete it
	for _, fileId := range dbFileList {
		fileName := "file_" + typehelper.MustItoa(fileId)
		_, ok := matching[fileName]
		if ok {
			logger.L.Debug("[CleanerMonitor]", fileName+" 因正常关联于实体文件而不会被删除")
			matching[fileName] = false
		}
	}
	// for hashList, if InResource and InTable, not delete it
	for _, fileName := range dbHashList {
		_, ok := matching[fileName]
		if ok {
			logger.L.Debug("[CleanerMonitor]", fileName+" 因正常关联于文章内的引用而不会被删除")
			matching[fileName] = false
		}
	}
	// for unknownList, delete all
	return
}

func getDbFileList() (list []int) {
	data.DoTx("获取博客内所有文件实体", func(tx *sqlx.Tx) (err error) {
		list, err = data.GetAllIdByType(tx, 2)
		return
	})
	return
}

func getDbHashList() (list []string) {
	var markdownList []string
	data.DoTx("获取博客内所有文章实体的markdown", func(tx *sqlx.Tx) (err error) {
		markdownList, err = data.GetAllMarkdownInArticle(tx)
		return
	})
	list = make([]string, 0)
	for _, md := range markdownList {
		list = append(list, hashRefPrompt.FindAllString(md, -1)...)
	}
	return
}

var (
	cycle      = true
	cycleMutex sync.Mutex
	immChan    = make(chan byte, 100)
)

func Run() {
	acc := interval
	for {
		cycleMutex.Lock()
		if !cycle {
			cycleMutex.Unlock()
			break
		} else {
			cycleMutex.Unlock()
		}

		imm := false
		select {
		case <-immChan:
			imm = true
		default:
			time.Sleep(time.Second)
		}

		if imm || acc >= interval {
			go func() {
				logger.L.Info("[CleanerMonitor]", "开始执行文件统计...")
				matching := matchFiles(getDbFileList(), getDbHashList())
				logger.L.Info("[CleanerMonitor]", "文件统计结束")
				for fileName, doRemove := range matching {
					if doRemove {
						logger.L.Info("[CleanerMonitor]", "发现无效资源文件: "+fileName)
						datahelper.RemoveFileByName(fileName)
					}
				}
			}()
			acc = 0
		}
		acc += time.Second
	}
}

func Manually() {
	immChan <- 0
}

func Exit() {
	cycleMutex.Lock()
	cycle = false
	cycleMutex.Unlock()
	time.Sleep(time.Second)
}
