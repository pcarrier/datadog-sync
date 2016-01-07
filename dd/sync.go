package dd

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/Sirupsen/logrus"
	"github.com/meteor/datadog-sync/util"
)

type update struct {
	from, to Monitor
}

func (m *Monitor) shortDescription() string {
	if m.ID != nil {
		return fmt.Sprintf("[#%d] %s", *m.ID, m.Name)
	}
	return m.Name
}

// SyncMonitors creates Datadog metrics from local that aren't in Datadog,
// then update metrics in Datadog to metrics from local whose ID matches,
// then deletes metrics in Datadog that weren't accounted for in local
func SyncMonitors(local, remote []Monitor, client *http.Client, dryRun, verbose bool) error {
	remoteSetIgnoringID := make(monitorSetIgnoringID)
	remotesByID := make(map[id]Monitor)
	localSet := make(monitorSetIgnoringID)
	var toCreate []Monitor
	var toUpdate []update

	for _, r := range remote {
		remoteSetIgnoringID.add(r)
		remotesByID[*r.ID] = r
	}

	for _, l := range local {
		if l.ID == nil {
			if r, ok := remoteSetIgnoringID.contains(l); ok {
				delete(remotesByID, *r.ID)
			} else {
				toCreate = append(toCreate, l)
			}
		} else { // l has an ID
			if r, ok := remotesByID[*l.ID]; ok {
				delete(remotesByID, *r.ID)
				if !reflect.DeepEqual(l, r) {
					toUpdate = append(toUpdate, update{from: l, to: r})
				}
			} else {
				return fmt.Errorf("no remote alert #%d", *l.ID)
			}
		}
		localSet.add(l)
	}

	creations := len(toCreate)
	updates := len(toUpdate)
	deletions := len(remotesByID)
	total := creations + updates + deletions

	logrus.Infof("%d creations, %d updates, %d deletions", creations, updates, deletions)

	for i, m := range toCreate {
		logrus.Infof("CREATE %d/%d/%d: %s", i, creations, total, m.shortDescription())
		if !dryRun {
			if err := m.create(client); err != nil {
				return err
			}
		}
		if verbose {
			repr, _ := util.Marshal(m, util.YAML)
			logrus.Debug(repr)
		}
	}

	for i, u := range toUpdate {
		logrus.Infof("UPDATE %d/%d/%d: %s", i, updates, total, u.from.shortDescription())
		if !dryRun {
			if err := u.from.update(client, &u.to); err != nil {
				return err
			}
		}
		if verbose {
			f, _ := util.Marshal(u.from, util.YAML)
			t, _ := util.Marshal(u.to, util.YAML)
			logrus.Debugf("%s\n=>\n%s", f, t)
		}
	}

	idx := 0
	for _, m := range remotesByID {
		logrus.Infof("DELETE %d/%d/%d: %s", idx, deletions, total, m.shortDescription())
		idx++
		if !dryRun {
			if err := m.delete(client); err != nil {
				return err
			}
		}
		if verbose {
			repr, _ := util.Marshal(m, util.YAML)
			logrus.Debug(repr)
		}
	}

	return nil
}
