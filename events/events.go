// The events package provides event handling facilities.
package events

import (
	"log";
	"container/vector";
	"sort";
)

// Event handlers are called in order of priority.
const (
	PrioFirst	= 0;
	PrioEarly	= 25;
	PrioNormal	= 50;
	PrioLate	= 75;
	PrioLast	= 100;
)
type Priority int

// An event handler function also has a name so that it can be removed again
// later.
type handler struct {
	name	string;
	prio	Priority;
	fn	interface{};
}

// Comparison operator to support sorting handlers by priority.
func (a *handler) Less(elem interface{}) bool {
	// Order according to priority, then by name.  The name comparison is
	// not necessary but ensures there is always a stable order independent
	// of the order in which handlers were added.
	b := elem.(*handler);
	if a.prio < b.prio {
		return true
	} else if a.prio == b.prio && a.name < b.name {
		return true
	}
	return false;
}

// Event handlers may stop propagation if they have consumed the event and do
// not wish further handlers to be called.
const (
	EventContinue	= true;
	EventStop	= false;
)
type EventAction bool

// Event handlers are called indirectly through an invoke function.  The invoke
// function has the opportunity to cast types so that handler functions do not
// need to cast.
type InvokeFunc func(fn interface{}, args ...) EventAction

// Event chains are identified by name.
type Chain struct {
	name		string;
	invoke		InvokeFunc;
	handlers	*vector.Vector;
}

// NewChain creates an empty event chain.
func NewChain(name string, invoke InvokeFunc) *Chain {
	return &Chain{name, invoke, vector.New(0)}
}

// AddHandler inserts a handler into the chain, positioned according to its priority.
func (chain *Chain) AddHandler(name string, prio Priority, fn interface{}) {
	handler := &handler{name, prio, fn};
	chain.handlers.Push(handler);

	// Restore priority ordering
	sort.Sort(chain.handlers);
}

// RemoveHandler removes a handler from a chain given its name.
func (chain *Chain) RemoveHandler(name string) {
	for i := 0; i < chain.handlers.Len(); i++ {
		if chain.handlers.At(i).(*handler).name == name {
			chain.handlers.Delete(i);
			return;
		}
	}
}

// Notify calls each event handler until the event is consumed or the end of
// the chain is reached.
func (chain *Chain) Notify(args ...) {
	for elem := range chain.handlers.Iter() {
		action := chain.invoke(elem.(*handler).fn, args);
		if action == EventStop {
			return
		}
	}
}

// Chains can be registered globally here.
var chains = make(map[string] *Chain)

func AddChain(name string, invoke InvokeFunc) {
	log.Stderrf("AddChain chain=%s\n", name);
	chains[name] = NewChain(name, invoke)
}

func RemoveChain(name string) {
	log.Stderrf("RemoveChain chain=%s\n", name);
	chains[name] = nil, false
}

func AddHandler(chainName string, name string, prio Priority, fn interface{}) {
	chain, ok := chains[chainName];
	if ok {
		log.Stderrf("AddHandler chain=%s name=%s\n", chainName, name);
		chain.AddHandler(name, prio, fn);
	}
}

func RemoveHandler(chainName string, name string) {
	chain, ok := chains[chainName];
	if ok {
		log.Stderrf("RemoveHandler chain=%s name=%s\n", chainName, name);
		chain.RemoveHandler(name);
	}
}

func Notify(name string, args ...) {
	chain, ok := chains[name];
	if ok {
		log.Stderrf("Notify %s\n", name);
		chain.Notify(args);
	}
}
