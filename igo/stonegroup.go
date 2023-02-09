package igo

type StoneGroup struct {
	Stone     Stone
	Liberties [][2]int
	Coords    [][2]int
}

func (sg *StoneGroup) Include(coord [2]int) bool {
	for _, sc := range sg.Coords {
		if sc == coord {
			return true
		}
	}

	return false
}

func (sg *StoneGroup) IncludeLiberty(coord [2]int) bool {
	for _, sc := range sg.Liberties {
		if sc == coord {
			return true
		}
	}

	return false
}
