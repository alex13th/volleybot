package person

import (
	"testing"
)

func TestTelegramView(t *testing.T) {

	tests := map[string]struct {
		p    Person
		want string
	}{
		"Fullname": {
			p:    Person{Fullname: "Full Name"},
			want: "Full Name",
		},
		"Firstname": {
			p:    Person{Firstname: "Firstname"},
			want: "Firstname",
		},
		"Firstname_Lastname": {
			p: Person{
				Firstname: "Firstname",
				Lastname:  "Lastname"},
			want: "Firstname Lastname",
		},
		"Fullname_Firstname_Lastname": {
			p: Person{
				Fullname:  "Full Name",
				Firstname: "Firstname",
				Lastname:  "Lastname"},
			want: "Full Name",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tgv := TelegramView{Person: test.p}
			text := tgv.GetText()
			if text != test.want {
				t.Fail()
			}
		})
	}
}
