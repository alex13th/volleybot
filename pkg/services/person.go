package services

import (
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PersonConfiguration func(serv *PersonService) error

func WithPersonRepository(pr person.PersonRepository) PersonConfiguration {
	return func(serv *PersonService) error {
		serv.persons = pr
		return nil
	}
}

func WithMemoryPersonRepository() PersonConfiguration {
	pr := person.NewMemoryRepository()
	return WithPersonRepository(pr)
}

func WithPgPersonRepository(dbpool *pgxpool.Pool) PersonConfiguration {
	pr, _ := person.NewPgRepository(dbpool)
	pr.UpdateDB()
	return WithPersonRepository(&pr)
}

func NewPersonService(cfgs ...PersonConfiguration) (serv *PersonService, err error) {
	serv = &PersonService{}
	for _, cfg := range cfgs {
		err = cfg(serv)
		if err != nil {
			return
		}
	}
	return
}

type PersonService struct {
	persons person.PersonRepository
}

func (serv PersonService) Add(p person.Person) (person.Person, error) {
	return serv.persons.Add(p)
}

func (serv PersonService) Get(pid uuid.UUID) (person.Person, error) {
	return serv.persons.Get(pid)
}

func (serv PersonService) GetByTelegramId(tid int) (person.Person, error) {
	return serv.persons.GetByTelegramId(tid)
}

func (serv PersonService) Update(p person.Person) error {
	return serv.persons.Update(p)
}
