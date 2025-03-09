package author_domain

type Author struct {
	id   string
	name string
}

func NewAuthor(id, name string) *Author {
	return &Author{id: id, name: name}
}

func (a *Author) ID() string {
	return a.id
}

func (a *Author) Name() string {
	return a.name
}

func (a *Author) Update(name string) {
	a.name = name
}
