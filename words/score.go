package words

import (
	"fmt"
	"strconv"
	"strings"
)

type Score struct {
	Word       string
	ByPosition []int
	Total      int
}

func (s *Score) String() string {
	ss := make([]string, 0, len(s.ByPosition))
	for _, s := range s.ByPosition {
		ss = append(ss, strconv.Itoa(s))
	}
	return fmt.Sprintf("%s %d [%s]", s.Word, s.Total, strings.Join(ss, ", "))
}
