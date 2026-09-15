package ensigncycle

import (
	"reflect"
	"sort"
	"testing"
)

func TestLiveScheduleOrder(t *testing.T) {
	jobs := []liveScheduledTest{{name: "short", hint: 10}, {name: "z", hint: 30}, {name: "long", hint: 40}, {name: "a", hint: 30}}
	orderLiveTests(jobs)
	var names []string
	for _, job := range jobs {
		names = append(names, job.name)
	}
	if want := []string{"long", "a", "z", "short"}; !reflect.DeepEqual(names, want) {
		t.Fatalf("order = %v, want %v", names, want)
	}
}

type liveScheduledTest struct {
	name string
	hint int
	run  func(*testing.T)
}

func orderLiveTests(jobs []liveScheduledTest) {
	sort.Slice(jobs, func(i, j int) bool {
		if jobs[i].hint == jobs[j].hint {
			return jobs[i].name < jobs[j].name
		}
		return jobs[i].hint > jobs[j].hint
	})
}
