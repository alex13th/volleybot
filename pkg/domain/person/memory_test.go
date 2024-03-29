package person

import (
	"testing"

	"github.com/google/uuid"
)

func TestMemory_GetPerson(t *testing.T) {
	type testCase struct {
		name        string
		id          uuid.UUID
		expectedErr error
	}

	p := NewPerson("Firstname")
	id := p.Id

	repo := MemoryRepository{
		persons: map[uuid.UUID]Person{p.Id: p},
	}

	testCases := []testCase{
		{
			name:        "No person By ID",
			id:          uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
			expectedErr: ErrPersonNotFound,
		}, {
			name:        "person By ID",
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

func TestMemory_AddPerson(t *testing.T) {
	type testCase struct {
		name        string
		firstname   string
		expectedErr error
	}

	testCases := []testCase{
		{
			name:        "Add person",
			firstname:   "Percy",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := MemoryRepository{
				persons: map[uuid.UUID]Person{},
			}

			p := NewPerson(tc.firstname)
			_, err := repo.Add(p)
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}

			found, err := repo.Get(p.Id)
			if err != nil {
				t.Fatal(err)
			}
			if found.Id != p.Id {
				t.Errorf("Expected %v, got %v", p.Id, found.Id)
			}
		})
	}
}
