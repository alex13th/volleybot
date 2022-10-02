package volley

import (
	"time"
	"volleybot/pkg/domain/location"
	"volleybot/pkg/domain/person"
)

func NewPlayer(prsn person.Person) Player {
	return Player{
		Person: prsn,
	}
}

type Player struct {
	person.Person
	Level PlayerLevel `json:"level"`
}

func (pl Player) String() string {
	return pl.Person.String()
}

func (pl Player) CheckLocationRole(l location.Location, role string) bool {
	for _, r := range pl.LocationRoles[l.Id] {
		if r == role {
			return true
		}
	}
	return false
}

type Member struct {
	Player
	MemberId   int
	Count      int
	ArriveTime time.Time
	Paid       bool
}

type PlayerLevel int

const (
	Nothing      PlayerLevel = 0
	Novice       PlayerLevel = 10
	Begginer     PlayerLevel = 20
	BeginnerPlus PlayerLevel = 30
	MiddleMinus  PlayerLevel = 40
	Middle       PlayerLevel = 50
	MiddlePlus   PlayerLevel = 60
	Advanced     PlayerLevel = 70
	Proffesional PlayerLevel = 80
)

func (pl PlayerLevel) String() string {
	lnames := make(map[int]string)
	lnames[0] = "Не определен"
	lnames[10] = "Новичок"
	lnames[20] = "Начальный"
	lnames[30] = "Начальный➕"
	lnames[40] = "Средний➖"
	lnames[50] = "Средний"
	lnames[60] = "Средний➕"
	lnames[70] = "Уверенный"
	lnames[80] = "Профи"
	return lnames[int(pl)]
}

func (pl PlayerLevel) Emoji() string {
	lnames := make(map[int]string)
	lnames[0] = ""
	lnames[10] = "🙌"
	lnames[20] = "👏"
	lnames[30] = "🤝"
	lnames[40] = "👌"
	lnames[50] = "👍"
	lnames[60] = "💪"
	lnames[70] = "⭐️"
	lnames[80] = "👑"
	return lnames[int(pl)]
}

type NetType int

const (
	Undefined NetType = 0
	Female    NetType = 10
	Male      NetType = 20
)

func (nt NetType) String() string {
	names := make(map[int]string)
	names[0] = "Не определен"
	names[10] = "♂️ Мужская"
	names[20] = "♀️ Женская"
	return names[int(nt)]
}
