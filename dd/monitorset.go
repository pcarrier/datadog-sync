package dd

import (
	"github.com/meteor/datadog-sync/util"
)

type monitorSetIgnoringID map[string]Monitor

func (s monitorSetIgnoringID) add(x Monitor) {
	withoutID := x
	withoutID.ID = nil
	str, _ := util.Marshal(withoutID, util.JSON)
	s[str] = x
}

func (s monitorSetIgnoringID) contains(x Monitor) (Monitor, bool) {
	withoutID := x
	withoutID.ID = nil
	str, _ := util.Marshal(withoutID, util.JSON)
	r, ok := s[str]
	return r, ok
}

func (s monitorSetIgnoringID) remove(x Monitor) {
	withoutID := x
	withoutID.ID = nil
	str, _ := util.Marshal(withoutID, util.JSON)
	delete(s, str)
}

func (s monitorSetIgnoringID) entries(x Monitor) []Monitor {
	res := make([]Monitor, len(s))
	idx := 0
	for _, v := range s {
		res[idx] = v
		idx++
	}
	return res
}
