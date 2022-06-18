package services

import (
	"testing"
	"time"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"

	"github.com/google/uuid"
)

func TestOrder_NewOrderService(t *testing.T) {

	ps, err := NewPersonService(WithMemoryPersonRepository())
	if err != nil {
		t.Error(err)
	}

	os, err := NewOrderService(ps, WithMemoryReserveRepository())

	if err != nil {
		t.Error(err)
	}

	p, err := person.NewPerson("Percy")
	if err != nil {
		t.Error(err)
	}

	_, err = os.PersonService.Add(p)
	if err != nil {
		t.Error(err)
	}

	duration, _ := time.ParseDuration("2h")
	_, err = os.CreateOrder(reserve.Reserve{
		Person: p, StartTime: time.Now(), EndTime: time.Now().Add(duration)},
		nil)

	if err != nil {
		t.Error(err)
	}
}

func CreateTestOrders() (os *OrderService, pl []person.Person, rl []reserve.Reserve, err error) {

	pl = []person.Person{
		{Fullname: "Percy", Id: uuid.New()},
		{Fullname: "Nelly", Id: uuid.New()},
	}

	rl = make([]reserve.Reserve, 4)

	rl[0], _ = reserve.NewReserve(pl[0],
		time.Date(2021, 12, 04, 12, 0, 0, 0, time.UTC),
		time.Date(2021, 12, 04, 14, 0, 0, 0, time.UTC),
	)

	rl[1], _ = reserve.NewReserve(pl[1],
		time.Date(2021, 12, 04, 16, 0, 0, 0, time.UTC),
		time.Date(2021, 12, 04, 20, 0, 0, 0, time.UTC),
	)

	rl[2], _ = reserve.NewReserve(pl[0],
		time.Date(2021, 12, 05, 10, 0, 0, 0, time.UTC),
		time.Date(2021, 12, 05, 11, 0, 0, 0, time.UTC),
	)

	rl[3], _ = reserve.NewReserve(pl[1],
		time.Date(2021, 12, 05, 20, 0, 0, 0, time.UTC),
		time.Date(2021, 12, 05, 23, 0, 0, 0, time.UTC),
	)

	os, err = NewOrderService(nil, WithMemoryReserveRepository())
	if err != nil {
		return
	}
	os.PersonService, err = NewPersonService(WithMemoryPersonRepository())

	if err != nil {
		return
	}

	for _, p := range pl {
		os.PersonService.Add(p)
	}

	for _, res := range rl {
		_, err = os.CreateOrder(res, nil)
		if err != nil {
			return
		}
	}
	return
}

func TestOrder_List(t *testing.T) {

	os, pl, rl, err := CreateTestOrders()

	if err != nil {
		t.Fail()
	}

	tests := map[string]struct {
		p     person.Person
		start time.Time
		end   time.Time
		want  map[uuid.UUID]reserve.Reserve
	}{
		"Date 05 list": {
			start: time.Date(2021, 12, 04, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2021, 12, 04, 23, 59, 0, 0, time.UTC),
			want:  map[uuid.UUID]reserve.Reserve{rl[0].Id: rl[0], rl[1].Id: rl[1]},
		},
		"Date 04 list": {
			start: time.Date(2021, 12, 05, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2021, 12, 05, 23, 59, 0, 0, time.UTC),
			want:  map[uuid.UUID]reserve.Reserve{rl[2].Id: rl[2], rl[3].Id: rl[3]},
		},
		"Date greater 04 list": {
			start: time.Date(2021, 12, 04, 0, 0, 0, 0, time.UTC),
			want:  map[uuid.UUID]reserve.Reserve{rl[0].Id: rl[0], rl[1].Id: rl[1], rl[2].Id: rl[2], rl[3].Id: rl[3]},
		},
		"Date less 04 list": {
			end:  time.Date(2021, 12, 04, 23, 59, 0, 0, time.UTC),
			want: map[uuid.UUID]reserve.Reserve{rl[0].Id: rl[0], rl[1].Id: rl[1]},
		},
		"Person 0 list": {
			p:    pl[0],
			want: map[uuid.UUID]reserve.Reserve{rl[0].Id: rl[0], rl[2].Id: rl[2]},
		},
		"Person 1 list": {
			p:    pl[1],
			want: map[uuid.UUID]reserve.Reserve{rl[1].Id: rl[1], rl[3].Id: rl[3]},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rlist, err := os.List(reserve.Reserve{Person: test.p, StartTime: test.start, EndTime: test.end}, false, nil)
			if err != nil {
				t.Fail()
			}

			for _, reserve := range rlist {
				checked := false
				for _, exp_reserve := range test.want {
					if reserve.Person.Id.ID() == exp_reserve.Person.Id.ID() {
						checked = true
					}
				}
				if !checked {
					t.Fail()
				}
			}
		})
	}

}
