package reserve

import (
	"reflect"
	"testing"
	"time"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
)

func CreateTestReserves() (repo *MemoryRepository, pl []person.Person, rl []Reserve, err error) {

	pl = []person.Person{
		{Fullname: "Percy", Id: uuid.New()},
		{Fullname: "Nelly", Id: uuid.New()},
	}

	rl = make([]Reserve, 4)

	rl[0], _ = NewReserve(pl[0],
		time.Date(2021, 12, 04, 12, 0, 0, 0, time.UTC),
		time.Date(2021, 12, 04, 14, 0, 0, 0, time.UTC),
	)

	rl[1], _ = NewReserve(pl[1],
		time.Date(2021, 12, 04, 16, 0, 0, 0, time.UTC),
		time.Date(2021, 12, 04, 20, 0, 0, 0, time.UTC),
	)

	rl[2], _ = NewReserve(pl[0],
		time.Date(2021, 12, 05, 10, 0, 0, 0, time.UTC),
		time.Date(2021, 12, 05, 11, 0, 0, 0, time.UTC),
	)

	rl[3], _ = NewReserve(pl[1],
		time.Date(2021, 12, 05, 20, 0, 0, 0, time.UTC),
		time.Date(2021, 12, 05, 23, 0, 0, 0, time.UTC),
	)

	repo = &MemoryRepository{}
	repo.reserves = map[uuid.UUID]Reserve{}
	for _, res := range rl {
		repo.reserves[res.Id] = res
	}
	return
}

func TestMemory_GetReserve(t *testing.T) {
	type testCase struct {
		name        string
		id          uuid.UUID
		expectedErr error
	}

	p, _ := person.NewPerson("Firstname")
	dur, _ := time.ParseDuration("2h")
	res, err := NewReserve(p, time.Now(), time.Now().Add(dur))

	if err != nil {
		t.Fatal(err)
	}
	id := res.Id

	repo := MemoryRepository{
		reserves: map[uuid.UUID]Reserve{res.Id: res},
	}

	testCases := []testCase{
		{
			name:        "No reserve By ID",
			id:          uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
			expectedErr: ErrReserveNotFound,
		}, {
			name:        "reserve By ID",
			id:          id,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			_, err := repo.Get(tc.id)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestMemory_AddReserve(t *testing.T) {
	type testCase struct {
		name        string
		duration    string
		expectedErr error
	}

	testCases := []testCase{
		{
			name:        "Add reserve",
			duration:    "2h",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := MemoryRepository{
				reserves: map[uuid.UUID]Reserve{},
			}

			duration, _ := time.ParseDuration(tc.duration)
			res, err := NewReserve(person.Person{Firstname: "Lily"}, time.Now(), time.Now().Add(duration))

			if err != nil {
				t.Fatal(err)
			}

			_, err = repo.Add(res)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}

			found, err := repo.Get(res.Id)
			if err != nil {
				t.Fatal(err)
			}
			if found.Id != res.Id {
				t.Errorf("Expected %v, got %v", res.Id, found.Id)
			}
		})
	}
}

func TestMemory_GetByFilter(t *testing.T) {

	repo, pl, rl, err := CreateTestReserves()

	if err != nil {
		t.Fail()
	}

	tests := map[string]struct {
		p     person.Person
		start time.Time
		end   time.Time
		want  map[uuid.UUID]Reserve
	}{
		"Date 05 list": {
			start: time.Date(2021, 12, 04, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2021, 12, 04, 23, 59, 0, 0, time.UTC),
			want:  map[uuid.UUID]Reserve{rl[0].Id: rl[0], rl[1].Id: rl[1]},
		},
		"Date 04 list": {
			start: time.Date(2021, 12, 05, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2021, 12, 05, 23, 59, 0, 0, time.UTC),
			want:  map[uuid.UUID]Reserve{rl[2].Id: rl[2], rl[3].Id: rl[3]},
		},
		"Date greater 04 list": {
			start: time.Date(2021, 12, 04, 0, 0, 0, 0, time.UTC),
			want:  map[uuid.UUID]Reserve{rl[0].Id: rl[0], rl[1].Id: rl[1], rl[2].Id: rl[2], rl[3].Id: rl[3]},
		},
		"Date less 04 list": {
			end:  time.Date(2021, 12, 04, 23, 59, 0, 0, time.UTC),
			want: map[uuid.UUID]Reserve{rl[0].Id: rl[0], rl[1].Id: rl[1]},
		},
		"Person 0 list": {
			p:    pl[0],
			want: map[uuid.UUID]Reserve{rl[0].Id: rl[0], rl[2].Id: rl[2]},
		},
		"Person 1 list": {
			p:    pl[1],
			want: map[uuid.UUID]Reserve{rl[1].Id: rl[1], rl[3].Id: rl[3]},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			frepo := NewMemoryRepository(&repo.reserves, Reserve{Person: test.p, StartTime: test.start, EndTime: test.end}, false)
			if err != nil {
				t.Fail()
			}

			if !reflect.DeepEqual(frepo.reserves, test.want) {
				t.Fail()
			}
		})
	}

}
