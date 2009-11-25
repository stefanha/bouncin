package events

import "testing"

type prioCase struct {
	name	string;
	prio	Priority;
}

type prioTest struct {
	cases	[]prioCase;
	order	[]string;
}

// Test cases for handler priority ordering.
var prioTests = []prioTest {
	prioTest {
		[]prioCase {
			prioCase{"bouncer", PrioLast},
			prioCase{"replay", PrioNormal},
			prioCase{"log", PrioNormal},
		},
		[]string {"log", "replay", "bouncer"},
	},
	prioTest {
		[]prioCase {
			prioCase{"a1", PrioLate},
			prioCase{"a2", PrioEarly},
			prioCase{"a3", PrioLast},
			prioCase{"a4", PrioFirst},
			prioCase{"a5", PrioNormal},
		},
		[]string {"a4", "a2", "a5", "a1", "a3"},
	},
	prioTest {
		[]prioCase {
			prioCase{"a1", PrioNormal},
			prioCase{"a2", PrioNormal},
			prioCase{"a3", PrioNormal},
			prioCase{"a4", PrioNormal},
			prioCase{"a5", PrioNormal},
		},
		[]string {"a1", "a2", "a3", "a4", "a5"},
	},
}

func (chain *Chain) checkOrder(t *testing.T, tnum int, order []string) bool {
	i := 0;
	for _, name := range order {
		handler := chain.handlers.At(i).(*handler);
		if handler.name != name {
			t.Errorf("test %d expected order: %#v at %d got \"%s\"\n",
				tnum, order, i, handler.name);
			return false;
		}
		i++;
	}
	return true;
}

func TestPriority(t *testing.T) {
	for tnum, prioTest := range prioTests {
		chain := NewChain("test", nil);
		for _, prioCase := range prioTest.cases {
			chain.AddHandler(prioCase.name, prioCase.prio, nil)
		}
		if !chain.checkOrder(t, tnum, prioTest.order) {
			return
		}
	}
}
