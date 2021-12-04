package person

import "testing"

func TestPersonGetDisplayname(t *testing.T) {
	tests := map[string]struct {
		Person Person
		want   string
	}{
		"Firstname only": {
			Person: Person{Firstname: "Firstname"},
			want:   "Firstname",
		},
		"Lastname only": {
			Person: Person{Lastname: "Lastname"},
			want:   "Lastname",
		},
		"Firstname Lastname": {
			Person: Person{
				Firstname: "Firstname",
				Lastname:  "Lastname"},
			want: "Firstname Lastname",
		},
		"Fullname only": {
			Person: Person{
				Firstname: "Firstname",
				Fullname:  "Full Name"},
			want: "Full Name",
		},
		"Fullname Firstname": {
			Person: Person{
				Firstname: "Firstname",
				Fullname:  "Full Name"},
			want: "Full Name",
		},
		"Fullname Lastname": {
			Person: Person{
				Lastname: "Lastname",
				Fullname: "Full Name"},
			want: "Full Name",
		},
		"Fullname Firstname Lastname": {
			Person: Person{
				Firstname: "Firstname",
				Lastname:  "Lastname",
				Fullname:  "Full Name"},
			want: "Full Name",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			displayname := test.Person.GetDisplayname()
			if displayname != test.want {
				t.Fail()
			}
		})
	}
}
