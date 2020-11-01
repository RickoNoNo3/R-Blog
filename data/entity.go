package data

type Entity struct {
	Type int
	Id   int
}

type Article struct {
	Id       int
	Title    string
	Markdown string
	Html     string
}

type Dir struct {
	Id   int
	Name string
}

type LayerContent struct {
	Id        int    `json:"id"`
	Type      int    `json:"type"`
	Text      string `json:"text"`
	CreatedT  string `db:"createdT" json:"createdT"`
	ModifiedT string `db:"modifiedT" json:"modifiedT"`
}
