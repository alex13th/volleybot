package person

import (
	"testing"
)

func TestTelegramView(t *testing.T) {

	tests := map[string]struct {
		p    Person
		str  string
		text string
	}{
		"Fullname": {
			p:    Person{Fullname: "Full Name"},
			str:  "👤 Full Name",
			text: "*Имя*: \n*Фамилия*: \n*Полное имя*: Full Name\n*Пол*: Не определен",
		},
		"Firstname": {
			p:    Person{Firstname: "Firstname"},
			str:  "👤 Firstname",
			text: "*Имя*: Firstname\n*Фамилия*: \n*Полное имя*: Firstname\n*Пол*: Не определен",
		},
		"Firstname_Lastname": {
			p: Person{
				Firstname: "Firstname",
				Lastname:  "Lastname"},
			str:  "👤 Firstname Lastname",
			text: "*Имя*: Firstname\n*Фамилия*: Lastname\n*Полное имя*: Firstname Lastname\n*Пол*: Не определен",
		},
		"Fullname_Firstname_Lastname": {
			p: Person{
				Fullname:  "Full Name",
				Firstname: "Firstname",
				Lastname:  "Lastname",
				Sex:       1},
			str:  "👦🏻 Full Name",
			text: "*Имя*: Firstname\n*Фамилия*: Lastname\n*Полное имя*: Full Name\n*Пол*: 👦🏻 мальчик",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tgv := TelegramView{Person: test.p}
			str := tgv.String()
			text := tgv.GetText()
			if str != test.str {
				t.Fail()
			}
			if text != test.text {
				t.Fail()
			}
		})
	}
}
