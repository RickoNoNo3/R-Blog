package data

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"rickonono3/r-blog/mytype"
)

func InitDBTest(clear bool) {
	OpenDB("blog_test.db")
	if clear {
		DoTx(func(tx *sqlx.Tx) (err error) {
			tableList := make([]string, 0)
			if err = tx.Select(&tableList, "select name from sqlite_master where type=?", "table"); err != nil {
				panic(err)
			}
			for _, table := range tableList {
				if _, err = tx.Exec("delete from " + table); err != nil {
					panic(err)
				}
			}
			return
		})
	}
}

func TestCreateDir(t *testing.T) {
	InitDBTest(true)
	defer CloseDB()
	type TestData struct {
		Title       string
		ParentDirId int
		IdExpect    int
		IdActual    int
		Err         error
	}
	testData := []TestData{
		{
			Title:       "Dir1",
			ParentDirId: 0,
			IdExpect:    1,
		}, {
			Title:       "Dir2",
			ParentDirId: 0,
			IdExpect:    2,
		}, {
			Title:       "Dir3",
			ParentDirId: 2,
			IdExpect:    3,
		}, {
			Title:       "Dir4",
			ParentDirId: 1,
			IdExpect:    4,
		}, {
			Title:       "Dir5",
			ParentDirId: 0,
			IdExpect:    5,
		},
	}
	DoTx(func(tx *sqlx.Tx) error {
		for i, data := range testData {
			testData[i].IdActual, testData[i].Err = CreateDir(tx, data.Title, data.ParentDirId)
		}
		return nil
	})
	for _, data := range testData {
		assert.NoError(t, data.Err)
		assert.Equal(t, data.IdExpect, data.IdActual)
	}
	type TableInfo struct {
		Id    int
		Title string
	}
	tableInfos := make([]TableInfo, 0)
	for _, data := range testData {
		tableInfos = append(tableInfos, TableInfo{
			Id:    data.IdExpect,
			Title: data.Title,
		})
	}
	assert.Equal(
		t,
		tableInfos,
		func() (list []TableInfo) {
			list = make([]TableInfo, 0)
			err := DoTx(func(tx *sqlx.Tx) (err error) {
				err = tx.Select(&list, "select id, title from dir")
				return
			})
			if err != nil {
				panic(err)
			}
			return
		}(),
	)
}

func TestCreateArticle(t *testing.T) {
	InitDBTest(true)
	defer CloseDB()
	type TestData struct {
		Title          string
		Markdown       string
		Id             int
		Err            error
		TitleActual    string
		MarkdownActual string
	}
	testData := []TestData{
		{
			Title:    "Title1",
			Markdown: "# Title1\nThis is a demo.",
		},
		{
			Title:    "Title2",
			Markdown: "# Title2\nThis is a demo too.",
		},
		{
			Title:    "Title3",
			Markdown: "# Title3\n包含各种字符😘",
		},
	}
	DoTx(func(tx *sqlx.Tx) error {
		for i, data := range testData {
			testData[i].Id, testData[i].Err = CreateArticle(
				tx,
				data.Title,
				data.Markdown,
				0,
			)
		}
		return nil
	})
	for i, data := range testData {
		err := DoTx(func(tx *sqlx.Tx) (err error) {
			return tx.QueryRowx("select title, markdown from article where id=?", data.Id).Scan(&testData[i].TitleActual, &testData[i].MarkdownActual)
		})
		assert.NoError(t, err)
	}
	for i, data := range testData {
		assert.NoError(t, data.Err)
		assert.Equal(t, i+1, data.Id)
		assert.Equal(t, data.Title, data.TitleActual)
		assert.Equal(t, data.Markdown, data.MarkdownActual)
	}
}

func TestCreateFile(t *testing.T) {
	InitDBTest(true)
	defer CloseDB()
	var (
		id  [2]int
		err [2]error
	)
	DoTx(func(tx *sqlx.Tx) error {
		id[0], err[0] = CreateFile(tx, "File", 0)
		id[1], err[1] = CreateFile(tx, "Name", 0)
		return nil
	})
	assert.NoError(t, err[0])
	assert.NoError(t, err[1])
	assert.Equal(t, 1, id[0])
	assert.Equal(t, 2, id[1])
	assert.Equal(
		t,
		[]string{
			"File",
			"Name",
		},
		func() (titleList []string) {
			titleList = make([]string, 0)
			if err := DoTx(func(tx *sqlx.Tx) (err error) {
				return tx.Select(&titleList, "select title from file")
			}); err != nil {
				return []string{}
			}
			return
		}(),
	)
}

func TestGetDir(t *testing.T) {
	InitDBTest(true)
	defer CloseDB()
	dirId := -1
	err := DoTx(func(tx *sqlx.Tx) (err error) {
		_, _ = CreateDir(tx, "Dir", 0)
		dirId, err = CreateDir(tx, "Dir", 1)
		return
	})
	assert.NoError(t, err)
	assert.Equal(
		t,
		mytype.Dir{
			Entity: mytype.Entity{
				Id:        dirId,
				Type:      0,
				Title:     "Dir",
				CreatedT:  "",
				ModifiedT: "",
			},
		},
		func() (dir mytype.Dir) {
			DoTx(func(tx *sqlx.Tx) (err error) {
				dir, err = GetDir(tx, dirId)
				return
			})
			// These attributes doesn't be tested
			dir.Entity.CreatedT = ""
			dir.Entity.ModifiedT = ""
			return
		}(),
	)
}

func TestGetDir2(t *testing.T) {
	InitDBTest(true)
	defer CloseDB()
	dir := mytype.Dir{}
	DoTx(func(tx *sqlx.Tx) (err error) {
		dir, err = GetDir(tx, 0)
		assert.NoError(t, err)
		return
	})
	assert.Equal(
		t,
		dir,
		mytype.Dir{
			Entity: mytype.Entity{
				Id:    0,
				Type:  1,
				Title: "博客",
			},
		},
	)
}

func TestGetArticle(t *testing.T) {
	InitDBTest(true)
	defer CloseDB()
	title := "# Title"
	markdown := "# Title\n" +
		"## SubTitle\n" +
		"This is a markdown text file!\n" +
		"\n" +
		"See [Help Document](/help.html) to learn more.\n"
	articleId := -1
	err := DoTx(func(tx *sqlx.Tx) (err error) {
		articleId, err = CreateArticle(
			tx,
			title,
			markdown,
			0,
		)
		return
	})
	assert.NoError(t, err)
	assert.Equal(
		t,
		mytype.Article{
			Entity: mytype.Entity{
				Id:        articleId,
				Type:      1,
				Title:     title,
				CreatedT:  "",
				ModifiedT: "",
			},
			Markdown: markdown,
			Tags:     "",
			Voted:    0,
			Visited:  0,
		},
		func() (article mytype.Article) {
			DoTx(func(tx *sqlx.Tx) (err error) {
				article, err = GetArticle(tx, articleId)
				return
			})
			// These attributes doesn't be tested
			article.Entity.CreatedT = ""
			article.Entity.ModifiedT = ""
			article.Tags = ""
			article.Voted = 0
			article.Visited = 0
			return
		}(),
	)
}

func TestGetFile(t *testing.T) {
	InitDBTest(true)
	defer CloseDB()
	fileId := -1
	err := DoTx(func(tx *sqlx.Tx) (err error) {
		fileId, err = CreateFile(tx, "File", 0)
		return
	})
	assert.NoError(t, err)
	assert.Equal(
		t,
		mytype.File{
			Entity: mytype.Entity{
				Id:        fileId,
				Type:      2,
				Title:     "File",
				CreatedT:  "",
				ModifiedT: "",
			},
		},
		func() (file mytype.File) {
			DoTx(func(tx *sqlx.Tx) (err error) {
				file, err = GetFile(tx, fileId)
				return
			})
			// These attributes doesn't be tested
			file.Entity.CreatedT = ""
			file.Entity.ModifiedT = ""
			return
		}(),
	)
}

func TestGetEntityInfo(t *testing.T) {
	InitDBTest(true)
	defer CloseDB()
	type TestData struct {
		Title string
		DirId int
		Id    int
		Type  int
	}
	testData := []TestData{
		{
			Title: "Dir",
			DirId: 0,
			Type:  0,
		}, {
			Title: "Article",
			DirId: 1,
			Type:  1,
		}, {
			Title: "File",
			DirId: 1,
			Type:  2,
		},
	}
	DoTx(func(tx *sqlx.Tx) (err error) {
		testData[0].Id, err = CreateDir(tx, testData[0].Title, testData[0].DirId)
		assert.NoError(t, err)
		testData[1].Id, err = CreateArticle(tx, testData[1].Title, "# "+testData[1].Title, testData[1].DirId)
		assert.NoError(t, err)
		testData[2].Id, err = CreateFile(tx, testData[2].Title, testData[2].DirId)
		assert.NoError(t, err)
		return
	})
	for _, data := range testData {
		entity := mytype.Entity{}
		DoTx(func(tx *sqlx.Tx) (err error) {
			entity, err = GetEntityInfo(tx, data.Type, data.Id)
			assert.NoError(t, err)
			return
		})
		assert.Equal(
			t,
			mytype.Entity{
				Id:    data.Id,
				Type:  data.Type,
				Title: data.Title,
			},
			mytype.Entity{
				Id:    entity.Id,
				Type:  entity.Type,
				Title: entity.Title,
			},
		)
	}
}

func TestGetParentDir(t *testing.T) {
	InitDBTest(true)
	defer CloseDB()
	type TestData struct {
		Title string
		DirId int
		Id    int
		Type  int
	}
	testData := []TestData{
		{
			Title: "Dir1",
			DirId: 0,
			Type:  0,
		}, {
			Title: "Dir2",
			DirId: 1,
			Type:  0,
		}, {
			Title: "Article",
			DirId: 1,
			Type:  1,
		}, {
			Title: "File",
			DirId: 2,
			Type:  2,
		},
	}
	DoTx(func(tx *sqlx.Tx) (err error) {
		testData[0].Id, err = CreateDir(tx, testData[0].Title, testData[0].DirId)
		assert.NoError(t, err)
		testData[1].Id, err = CreateDir(tx, testData[1].Title, testData[1].DirId)
		assert.NoError(t, err)
		testData[2].Id, err = CreateArticle(tx, testData[2].Title, "# "+testData[2].Title, testData[2].DirId)
		assert.NoError(t, err)
		testData[3].Id, err = CreateFile(tx, testData[3].Title, testData[3].DirId)
		assert.NoError(t, err)
		return
	})
	for _, data := range testData {
		dirId := -1
		DoTx(func(tx *sqlx.Tx) (err error) {
			dirId, err = GetParentDir(tx, mytype.Entity{
				Id:    data.Id,
				Type:  data.Type,
				Title: data.Title,
			})
			assert.NoError(t, err)
			return
		})
		assert.Equal(
			t,
			data.DirId,
			dirId,
		)
	}
}

func TestGetParentDir2(t *testing.T) {
	InitDBTest(true)
	defer CloseDB()
	DoTx(func(tx *sqlx.Tx) (err error) {
		_, err = GetParentDir(tx, mytype.Entity{
			Id:   0,
			Type: 0,
		})
		assert.Error(t, err)
		return
	})
}

func TestGetContents(t *testing.T) {
	InitDBTest(true)
	defer CloseDB()
	type TestData struct {
		Title string
		DirId int
		Id    int
		Type  int
	}
	testData := []TestData{
		{
			Title: "Dir1",
			DirId: 0,
			Type:  0,
		}, {
			Title: "Dir2",
			DirId: 1,
			Type:  0,
		}, {
			Title: "Dir3",
			DirId: 1,
			Type:  0,
		}, {
			Title: "Article",
			DirId: 1,
			Type:  1,
		}, {
			Title: "File",
			DirId: 2,
			Type:  2,
		},
	}
	wants := [][]mytype.Entity{
		{
			{
				Id:    1,
				Type:  0,
				Title: "Dir1",
			},
		}, {
			{
				Id:    2,
				Type:  0,
				Title: "Dir2",
			},
			{
				Id:    3,
				Type:  0,
				Title: "Dir3",
			},
			{
				Id:    1,
				Type:  1,
				Title: "Article",
			},
		}, {
			{
				Id:    1,
				Type:  2,
				Title: "File",
			},
		}, {},
	}
	DoTx(func(tx *sqlx.Tx) (err error) {
		testData[0].Id, err = CreateDir(tx, testData[0].Title, testData[0].DirId)
		assert.NoError(t, err)
		testData[1].Id, err = CreateDir(tx, testData[1].Title, testData[1].DirId)
		assert.NoError(t, err)
		testData[2].Id, err = CreateDir(tx, testData[2].Title, testData[2].DirId)
		assert.NoError(t, err)
		testData[3].Id, err = CreateArticle(tx, testData[3].Title, "# "+testData[3].Title, testData[3].DirId)
		assert.NoError(t, err)
		testData[4].Id, err = CreateFile(tx, testData[4].Title, testData[4].DirId)
		assert.NoError(t, err)
		return
	})
	for i, want := range wants {
		contents := make([]mytype.Entity, 0)
		DoTx(func(tx *sqlx.Tx) (err error) {
			contents, err = GetContents(tx, i)
			assert.NoError(t, err)
			return
		})
		for j := range contents {
			contents[j].ModifiedT = ""
			contents[j].CreatedT = ""
		}
		assert.Equal(
			t,
			want,
			contents,
		)
	}
}

type LayerTestData struct {
	Entity   mytype.Entity
	DirId    int
	Duration time.Duration
}

/**
 *             [ d2--[ a2
 *             [
 * (d0)--[ d1--[ d3--[ d4
 *       [     [
 *       [ f1  [ a1
 *
 * [0:d1, 1:d2, 2:d3, 3:d4, 4:a1, 5:a2, 6:f1]
 */
var layerTestData = []LayerTestData{
	{
		Entity: mytype.Entity{
			Id:    1,
			Type:  0,
			Title: "Dir1",
		},
		DirId:    0,
		Duration: 0,
	},
	{
		Entity: mytype.Entity{
			Id:    2,
			Type:  0,
			Title: "Dir2",
		},
		DirId:    1,
		Duration: time.Second,
	},
	{
		Entity: mytype.Entity{
			Id:    3,
			Type:  0,
			Title: "Dir3",
		},
		DirId:    1,
		Duration: time.Second,
	},
	{
		Entity: mytype.Entity{
			Id:    4,
			Type:  0,
			Title: "Dir4",
		},
		DirId:    3,
		Duration: time.Second,
	},
	{
		Entity: mytype.Entity{
			Id:    1,
			Type:  1,
			Title: "Article1",
		},
		DirId:    1,
		Duration: 0,
	},
	{
		Entity: mytype.Entity{
			Id:    2,
			Type:  1,
			Title: "Article2",
		},
		DirId:    2,
		Duration: 2 * time.Second,
	},
	{
		Entity: mytype.Entity{
			Id:    1,
			Type:  2,
			Title: "File1",
		},
		DirId:    0,
		Duration: time.Second,
	},
}

func GetLayerTime(t *testing.T, tx *sqlx.Tx, i int) (modifiedT, createdT string, err error) {
	row := tx.QueryRowx("select modifiedT, createdT from layer where id=? and type=?", layerTestData[i].Entity.Id, layerTestData[i].Entity.Type)
	assert.NoError(t, row.Err())
	err = row.Scan(&modifiedT, &createdT)
	assert.NoError(t, err)
	assert.NotEmpty(t, modifiedT)
	assert.NotEmpty(t, createdT)
	return
}

func ImmediatelyCreateTestLayer() {
	InitDBTest(true)
	defer CloseDB()
	DoTx(func(tx *sqlx.Tx) (err error) {
		for _, data := range layerTestData {
			switch data.Entity.Type {
			case 0:
				_, err = CreateDir(tx, data.Entity.Title, data.DirId)
			case 1:
				_, err = CreateArticle(tx, data.Entity.Title, "# "+data.Entity.Title, data.DirId)
			case 2:
				_, err = CreateFile(tx, data.Entity.Title, data.DirId)
			}
		}
		return
	})
}

func TestCreateLayer(t *testing.T) {
	InitDBTest(true)
	defer CloseDB()
	DoTx(func(tx *sqlx.Tx) (err error) {
		for i, data := range layerTestData {
			// sleep
			time.Sleep(data.Duration)
			// create
			switch data.Entity.Type {
			case 0:
				_, err = CreateDir(tx, data.Entity.Title, data.DirId)
			case 1:
				_, err = CreateArticle(tx, data.Entity.Title, "# "+data.Entity.Title, data.DirId)
			case 2:
				_, err = CreateFile(tx, data.Entity.Title, data.DirId)
			}
			assert.NoError(t, err)
			// test modifiedT and createdT when created ONE
			var modifiedT, createdT string
			modifiedT, createdT, err = GetLayerTime(t, tx, i)
			assert.NoError(t, err)
			assert.Equal(t, modifiedT, createdT)
		}
		return
	})
	// test modifiedT and createdT when created ALL
	modifiedTs := make([]string, len(layerTestData))
	createdTs := make([]string, len(layerTestData))
	for i := range layerTestData {
		DoTx(func(tx *sqlx.Tx) (err error) {
			modifiedTs[i], createdTs[i], err = GetLayerTime(t, tx, i)
			assert.NoError(t, err)
			return
		})
	}
	// !!! d1m = d2m = a2m = a2c
	assert.Equal(t, modifiedTs[0], modifiedTs[1])
	assert.Equal(t, modifiedTs[0], modifiedTs[5])
	assert.Equal(t, modifiedTs[0], createdTs[5])
	// !!! d3m = d4m = d4c
	assert.Equal(t, modifiedTs[2], modifiedTs[3])
	assert.Equal(t, modifiedTs[2], createdTs[3])
	// !!! a1m = a1c
	assert.Equal(t, modifiedTs[4], createdTs[4])
	// !!! f1m = f1c
	assert.Equal(t, modifiedTs[6], createdTs[6])
}

func TestRemoveLayer(t *testing.T) {
	ImmediatelyCreateTestLayer()
	InitDBTest(false)
	defer CloseDB()
	// 测试步骤:
	//   1. 删除d0, 报错
	//   2. 删除d2, 结构变为
	//      d0--[d1--[d3--[d4
	//          [f1  [a1
	//   3. 删除d1, 结构变为
	//      d0--[f1
	//   4. 删除f1, 结构变为
	//      d0
	countQuery := "select count(*) from layer where type=? and id=?"
	countQueryAll := "select count(*) from layer"
	t.Run("Remove d0", func(t *testing.T) {
		DoTx(func(tx *sqlx.Tx) (err error) {
			err = RemoveLayer(tx, mytype.Entity{
				Id:   0,
				Type: 0,
			})
			assert.Error(t, err)
			return
		})
	})
	t.Run("Remove d2", func(t *testing.T) {
		DoTx(func(tx *sqlx.Tx) (err error) {
			err = RemoveLayer(tx, mytype.Entity{
				Id:   2,
				Type: 0,
			})
			assert.NoError(t, err)
			var count int
			err = tx.QueryRowx(countQuery, 0, 2).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
			err = tx.QueryRowx(countQuery, 1, 2).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
			err = tx.QueryRowx(countQueryAll).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 5, count)
			return
		})
	})
	t.Run("Remove d1", func(t *testing.T) {
		DoTx(func(tx *sqlx.Tx) (err error) {
			err = RemoveLayer(tx, mytype.Entity{
				Id:   1,
				Type: 0,
			})
			assert.NoError(t, err)
			var count int
			err = tx.QueryRowx(countQuery, 0, 1).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
			err = tx.QueryRowx(countQuery, 0, 3).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
			err = tx.QueryRowx(countQuery, 0, 4).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
			err = tx.QueryRowx(countQuery, 1, 1).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
			err = tx.QueryRowx(countQuery, 2, 1).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 1, count)
			err = tx.QueryRowx(countQueryAll).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 1, count)
			return
		})
	})
	t.Run("Remove f1", func(t *testing.T) {
		DoTx(func(tx *sqlx.Tx) (err error) {
			err = RemoveLayer(tx, mytype.Entity{
				Id:   1,
				Type: 2,
			})
			assert.NoError(t, err)
			var count int
			err = tx.QueryRowx(countQuery, 2, 1).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
			err = tx.QueryRowx(countQueryAll).Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
			return
		})
	})
}

func TestMoveLayer(t *testing.T) {
	ImmediatelyCreateTestLayer()
	InitDBTest(false)
	defer CloseDB()
	// 测试步骤:
	//   1. 移动d0到d0, 结构不变
	//               [d2--[a2
	//      d0--[d1--[d3--[d4
	//          [f1  [a1
	//   2. 移动d2到d3, 结构变为
	//      d0--[d1--[d3--[d4
	//          [f1  [a1  [d2--[a2
	//   3. 移动d1到d0, 结构不变
	//      d0--[d1--[d3--[d4
	//          [f1  [a1  [d2--[a2
	//   4. 移动d0到d3, 报错
	//      d0--[d1--[d3--[d4
	//          [f1  [a1  [d2--[a2
	//   5. 移动a1到d0, 结构变为
	//      d0--[d1--[d3--[d4
	//          [f1       [d2--[a2
	//          [a1
	countQueryDir := "select count(*) from layer where dirId=?"
	dirQuery := "select dirId from layer where type=? and id=?"
	t.Run("Move d0 to d0", func(t *testing.T) {
		DoTx(func(tx *sqlx.Tx) (err error) {
			err = MoveLayer(tx, mytype.Entity{
				Id:   0,
				Type: 0,
			}, 0)
			assert.NoError(t, err)
			var num int
			// d0: count(children)=2
			err = tx.QueryRowx(countQueryDir, 0).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 2, num)
			// d1: count(children)=3
			err = tx.QueryRowx(countQueryDir, 1).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 3, num)
			// d2: count(children)=1
			err = tx.QueryRowx(countQueryDir, 2).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 1, num)
			// d3: count(children)=1
			err = tx.QueryRowx(countQueryDir, 3).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 1, num)
			// d4: count(children)=0
			err = tx.QueryRowx(countQueryDir, 4).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 0, num)
			return
		})
	})
	t.Run("Move d2 to d3", func(t *testing.T) {
		DoTx(func(tx *sqlx.Tx) (err error) {
			err = MoveLayer(tx, mytype.Entity{
				Id:   2,
				Type: 0,
			}, 3)
			assert.NoError(t, err)
			var num int
			// d0: count(children)=2
			err = tx.QueryRowx(countQueryDir, 0).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 2, num)
			// d1: count(children)=2
			err = tx.QueryRowx(countQueryDir, 1).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 2, num)
			// d2: count(children)=1
			err = tx.QueryRowx(countQueryDir, 2).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 1, num)
			// d3: count(children)=2
			err = tx.QueryRowx(countQueryDir, 3).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 2, num)
			// d4: count(children)=0
			err = tx.QueryRowx(countQueryDir, 4).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 0, num)
			// d2: dirId=3
			err = tx.QueryRowx(dirQuery, 0, 2).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 3, num)
			// a2: dirId=2
			err = tx.QueryRowx(dirQuery, 1, 2).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 2, num)
			return
		})
	})
	t.Run("Move d1 to d0", func(t *testing.T) {
		DoTx(func(tx *sqlx.Tx) (err error) {
			err = MoveLayer(tx, mytype.Entity{
				Id:   1,
				Type: 0,
			}, 0)
			assert.NoError(t, err)
			var num int
			// d0: count(children)=2
			err = tx.QueryRowx(countQueryDir, 0).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 2, num)
			// d1: count(children)=2
			err = tx.QueryRowx(countQueryDir, 1).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 2, num)
			// d2: count(children)=1
			err = tx.QueryRowx(countQueryDir, 2).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 1, num)
			// d3: count(children)=2
			err = tx.QueryRowx(countQueryDir, 3).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 2, num)
			// d4: count(children)=0
			err = tx.QueryRowx(countQueryDir, 4).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 0, num)
			// d1: dirId=0
			err = tx.QueryRowx(dirQuery, 0, 1).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 0, num)
			return
		})
	})
	t.Run("Move d0 to d3", func(t *testing.T) {
		DoTx(func(tx *sqlx.Tx) (err error) {
			err = MoveLayer(tx, mytype.Entity{
				Id:   0,
				Type: 0,
			}, 3)
			assert.EqualError(t, err, "move into it's child")
			return
		})
	})
	t.Run("Move a1 to d0", func(t *testing.T) {
		DoTx(func(tx *sqlx.Tx) (err error) {
			err = MoveLayer(tx, mytype.Entity{
				Id:   1,
				Type: 1,
			}, 0)
			assert.NoError(t, err)
			var num int
			// d0: count(children)=3
			err = tx.QueryRowx(countQueryDir, 0).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 3, num)
			// d1: count(children)=1
			err = tx.QueryRowx(countQueryDir, 1).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 1, num)
			// d2: count(children)=1
			err = tx.QueryRowx(countQueryDir, 2).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 1, num)
			// d3: count(children)=2
			err = tx.QueryRowx(countQueryDir, 3).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 2, num)
			// d4: count(children)=0
			err = tx.QueryRowx(countQueryDir, 4).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 0, num)
			// d1: dirId=0
			err = tx.QueryRowx(dirQuery, 0, 1).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 0, num)
			// a1: dirId=0
			err = tx.QueryRowx(dirQuery, 1, 1).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 0, num)
			// f1: dirId=0
			err = tx.QueryRowx(dirQuery, 2, 1).Scan(&num)
			assert.NoError(t, err)
			assert.Equal(t, 0, num)
			return
		})
	})
}
