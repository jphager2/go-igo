package igo

type StoneGroup struct {
	Stone     Stone
	Liberties []Coord
	Coords    []Coord
}

func (sg *StoneGroup) Include(coord Coord) bool {
	for _, sc := range sg.Coords {
		if sc == coord {
			return true
		}
	}

	return false
}

func (sg *StoneGroup) IncludeLiberty(coord Coord) bool {
	for _, sc := range sg.Liberties {
		if sc == coord {
			return true
		}
	}

	return false
}
