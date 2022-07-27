package volley

import (
	"time"
	"volleybot/pkg/domain/person"
	"volleybot/pkg/domain/reserve"

	"github.com/google/uuid"
)

func NewVolley(p person.Person, start time.Time, end time.Time) Volley {
	return Volley{
		Reserve:    reserve.NewReserve(p, start, end),
		CourtCount: 1,
		MaxPlayers: 6,
		Players:    []person.Player{}}
}

type Volley struct {
	reserve.Reserve
	Activity   reserve.Activity `json:"activity"`
	MinLevel   int              `json:"min_level"`
	CourtCount int              `json:"court_count"`
	MaxPlayers int              `json:"max_players"`
	Players    []person.Player  `json:"players"`
}

func (res *Volley) Copy() (result Volley) {
	result = *res
	result.Id = uuid.New()
	return
}

func (v *Volley) HasPlayerByTelegramId(id int) bool {
	for _, pl := range v.Players {
		if pl.Person.TelegramId == id {
			return pl.Count > 0
		}
	}
	return false
}

func (v *Volley) PlayerCount(pid uuid.UUID) (count int) {
	for i, pl := range v.Players {
		if v.Players[i].Id != pid {
			count += pl.Count
		}
	}
	return
}

func (v *Volley) GetPlayer(pid uuid.UUID) (pl person.Player) {
	for _, pl := range v.Players {
		if pl.Id == pid {
			return pl
		}
	}
	return
}

func (v *Volley) GetPlayerByTelegramId(tid int) (pl person.Player) {
	for _, pl := range v.Players {
		if pl.TelegramId == tid {
			return pl
		}
	}
	return
}

func (v *Volley) RemovePlayerByTelegramId(tid int) {
	newplist := []person.Player{}
	for _, pl := range v.Players {
		if pl.TelegramId != tid {
			newplist = append(newplist, pl)
		}
	}
	v.Players = newplist
}

func (v *Volley) PlayerInReserve(pid uuid.UUID) bool {
	count := 0
	for _, pl := range v.Players {
		if pl.Id == pid {
			return count >= v.MaxPlayers
		}
		count += pl.Count
	}
	return count >= v.MaxPlayers
}

func (v *Volley) JoinPlayer(pl person.Player) {
	for i, p := range v.Players {
		if p.Id == pl.Id {
			v.Players[i] = pl
			return
		}
	}
	v.Players = append(v.Players, pl)
}
