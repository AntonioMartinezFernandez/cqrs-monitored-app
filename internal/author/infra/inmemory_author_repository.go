package author_infra

import (
	"errors"
	"slices"
	"sync"

	author_domain "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/author/domain"
)

type InmemoryAuthorRepository struct {
	mut     *sync.Mutex
	authors []author_domain.Author
}

func NewInmemoryAuthorRepository() *InmemoryAuthorRepository {
	return &InmemoryAuthorRepository{
		mut:     &sync.Mutex{},
		authors: []author_domain.Author{},
	}
}

func (iar *InmemoryAuthorRepository) FindAll() ([]author_domain.Author, error) {
	iar.mut.Lock()
	defer iar.mut.Unlock()

	return iar.authors, nil
}

func (iar *InmemoryAuthorRepository) FindByID(id string) (*author_domain.Author, error) {
	iar.mut.Lock()
	defer iar.mut.Unlock()

	for _, author := range iar.authors {
		if author.ID() == id {
			return &author, nil
		}
	}

	return nil, errors.New("author not found")
}

func (iar *InmemoryAuthorRepository) Save(newAuthor author_domain.Author) error {
	iar.mut.Lock()
	defer iar.mut.Unlock()

	for _, author := range iar.authors {
		if author.ID() == newAuthor.ID() {
			return errors.New("author already exists")
		}
	}

	iar.authors = append(iar.authors, newAuthor)
	return nil
}

func (iar *InmemoryAuthorRepository) Update(author author_domain.Author) error {
	iar.mut.Lock()
	defer iar.mut.Unlock()

	for i, a := range iar.authors {
		if a.ID() == author.ID() {
			iar.authors[i] = author
			return nil
		}
	}

	return errors.New("author not found")
}

func (iar *InmemoryAuthorRepository) Delete(id string) error {
	iar.mut.Lock()
	defer iar.mut.Unlock()

	for i, author := range iar.authors {
		if author.ID() == id {
			iar.authors = slices.Delete(iar.authors, i, i+1)
			return nil
		}
	}

	return errors.New("author not found")
}
