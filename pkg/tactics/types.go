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
	Black:  ":black_large_square:",
	Red:    ":red_square:",
	Orange: ":orange_square:",
	Yellow: ":yellow_square:",
	Green:  ":green_square:",
	Blue:   ":blue_square:",
	Purple: ":purple_square:",
	Brown:  ":brown_square:",
	White:  ":white_large_square:",
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
	sb.WriteString("Map status:\n")
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
