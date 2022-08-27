package tactics

import "strings"

const (
	Black  = "black"
	Red    = "red"
	Orange = "orange"
	Yellow = "yellow"
	Green  = "green"
	Blue   = "blue"
	Purple = "purple"
	Brown  = "brown"
	White  = "white"
)

var discordMap map[string]string = map[string]string{
	Black:  ":black_circle:",
	Red:    ":red_circle:",
	Orange: ":orange_circle:",
	Yellow: ":yellow_circle:",
	Green:  ":green_circle:",
	Blue:   ":blue_circle:",
	Purple: ":purple_circle:",
	Brown:  ":brown_circle:",
	White:  ":white_circle:",
}

type Player struct {
	Tokens int
	Range  int
}

type Tile struct {
	Value string
}

type Map struct {
	Grid [][]Tile
}

func NewMap(width int, height int) Map {
	columns := make([][]Tile, width)

	for i := range columns {
		columns[i] = make([]Tile, height)
		for j := range columns[i] {
			columns[i][j] = Tile{Black}
		}
	}

	return Map{columns}
}

func (m Map) ToDiscordString() string {
	sb := strings.Builder{}
	for _, col := range m.Grid {
		for _, val := range col {
			str, ok := discordMap[val.Value]
			if ok {
				sb.WriteString(str)
			} else {
				sb.WriteString(discordMap[Black])
			}
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}
