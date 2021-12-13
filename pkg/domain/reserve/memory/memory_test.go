package memory

import (
	"testing"
	"time"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"

	"github.com/google/uuid"
)

func TestMemory_GetReserve(t *testing.T) {
	type testCase struct {
		name        string
		id          uuid.UUID
		expectedErr error
	}

	p, _ := person.NewPerson("Firstname")
	dur, _ := time.ParseDuration("2h")
	res, err := reserve.NewReserve(p, time.Now(), time.Now().Add(dur))

	if err != nil {
		t.Fatal(err)
	}
	id := res.Id

	repo := MemoryRepository{
		reserves: map[uuid.UUID]reserve.Reserve{res.Id: res},
	}

	testCases := []testCase{
		{
			name:        "No reserve By ID",
			id:          uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
			expectedErr: reserve.ErrReserveNotFound,
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
				reserves: map[uuid.UUID]reserve.Reserve{},
			}

			duration, _ := time.ParseDuration(tc.duration)
			res, err := reserve.NewReserve(person.Person{Firstname: "Lily"}, time.Now(), time.Now().Add(duration))

			if err != nil {
				t.Fatal(err)
			}

			err = repo.Add(res)
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
