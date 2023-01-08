package volley

import (
	"time"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"

	"github.com/google/uuid"
)

type Activity int

const (
	Game       Activity = 0
	Training   Activity = 10
	Tournament Activity = 20
	Tennis     Activity = 30
)

func (a Activity) String() string {
	lnames := make(map[int]string)
	lnames[0] = "🏐 Игры"
	lnames[10] = "‼️ Тренировка"
	lnames[20] = "🌟 Турнир"
	lnames[30] = "🎾 Теннис"
	return lnames[int(a)]
}

func (a Activity) Emoji() string {
	lnames := make(map[int]string)
	lnames[0] = "🏐"
	lnames[10] = "‼️"
	lnames[20] = "🌟"
	lnames[30] = "🎾"
	return lnames[int(a)]
}

func NewVolley(p person.Person, start time.Time, end time.Time) Volley {
	return Volley{
		Reserve:    reserve.NewReserve(p, start, end),
		CourtCount: 1,
		MaxPlayers: 6,
		Members:    []Member{}}
}

type Volley struct {
	reserve.Reserve
	Activity   Activity `json:"activity"`
	MinLevel   int      `json:"min_level"`
	CourtCount int      `json:"court_count"`
	MaxPlayers int      `json:"max_players"`
	NetType    NetType  `json:"net_type"`
	Members    []Member `json:"members"`
}

func (res *Volley) Copy() (result Volley) {
	result = *res
	result.Id = uuid.New()
	return
}

func (v *Volley) HasPlayerByTelegramId(id int) bool {
	for _, mb := range v.Members {
		if mb.Person.TelegramId == id {
			return mb.Count > 0
		}
	}
	return false
}

func (v *Volley) PlayerCount(pid uuid.UUID) (count int) {
	for i, pl := range v.Members {
		if v.Members[i].Id != pid {
			count += pl.Count
		}
	}
	return
}

func (v *Volley) GetMember(pid uuid.UUID) (mb Member) {
	for _, mb := range v.Members {
		if mb.Id == pid {
			return mb
		}
	}
	return
}

func (v *Volley) GetMemberByTelegramId(tid int) (pl Member) {
	for _, pl := range v.Members {
		if pl.TelegramId == tid {
			return pl
		}
	}
	return
}

func (v *Volley) RemovePlayerByTelegramId(tid int) {
	newplist := []Member{}
	for _, pl := range v.Members {
		if pl.TelegramId != tid {
			newplist = append(newplist, pl)
		}
	}
	v.Members = newplist
}

func (v *Volley) PlayerInReserve(pid uuid.UUID) bool {
	count := 0
	for _, mb := range v.Members {
		if mb.Id == pid {
			return count >= v.MaxPlayers
		}
		count += mb.Count
	}
	return count >= v.MaxPlayers
}

func (v *Volley) JoinPlayer(mb Member) {
	for i, m := range v.Members {
		if m.Id == mb.Id {
			v.Members[i] = mb
			return
		}
	}
	v.Members = append(v.Members, mb)
}
