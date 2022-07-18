package res

type DateTimeResources struct {
	DayCount    int
	DateMessage string
	DateButton  string
	TimeMessage string
	TimeButton  string
}

type CourtResources struct {
	Message    string
	Button     string
	Min        int
	Max        int
	MaxPlayers int
}

type ActivityResources struct {
	Message string
	Button  string
	Min     int
	Max     int
}

type SetResources struct {
	Message string
	Button  string
	Min     int
	Max     int
}

type MaxPlayerResources struct {
	Message          string
	CountError       string
	GroupChatWarning string
	Button           string
	Min              int
	Max              int
}

type DescriptionResources struct {
	Message     string
	DoneMessage string
	Button      string
}

type JoinPlayerResources struct {
	Message         string
	Button          string
	ArriveButton    string
	LeaveButton     string
	MultiButton     string
	MultiButtonText string
}

type PriceResources struct {
	Message string
	Button  string
	Min     int
	Max     int
	Step    int
}

type CancelResources struct {
	Message string
	Button  string
	Confirm string
	Abort   string
}
