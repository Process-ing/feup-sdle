package crdt

import "fmt"

type Dot struct {
	id  string
	seq uint32
}

func NewDot(id string, seq uint32) Dot {
	return Dot{id: id, seq: seq}
}

func (d Dot) String() string {
	return fmt.Sprintf("Dot{id %v, seq %d}", d.id, d.seq)
}
