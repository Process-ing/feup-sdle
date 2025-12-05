package crdt

import "fmt"

type Dot[T comparable] struct {
	replicaID T
	seq       int
}

func NewDot[T comparable](replicaID T, seq int) Dot[T] {
	return Dot[T]{replicaID: replicaID, seq: seq}
}

func (d Dot[T]) String() string {
	return fmt.Sprintf("Dot{replicaID %v, seq %d}", d.replicaID, d.seq)
}