package mytype

type Entity struct {
	Id        int
	Type      int
	Title     string
	CreatedT  string `db:"createdT"`
	ModifiedT string `db:"modifiedT"`
}
