package crdt

import (
	"fmt"
	"sdle-server/utils"
)

type DotContext struct {
	compactContext map[string]int
	dots           utils.Set[Dot]
}

func NewDotContext() *DotContext {
	return &DotContext{
		compactContext: make(map[string]int),
		dots:           utils.NewSet[Dot](),
	}
}

func (dc *DotContext) Knows(dot Dot) bool {
	if localSeq, ok := dc.compactContext[dot.id]; ok {
		return dot.seq <= localSeq
	}

	return false
}

func (dc *DotContext) MakeDot(replicaID string) Dot {
	if localSeq, ok := dc.compactContext[replicaID]; ok {
		localSeq++
		dc.compactContext[replicaID] = localSeq
		return NewDot(replicaID, localSeq)
	} else {
		dc.compactContext[replicaID] = 1
		return NewDot(replicaID, 1)
	}
}

func (dc *DotContext) InsertDot(dot Dot) {
	dc.InsertDotCompact(dot, true)
}

func (dc *DotContext) InsertDotCompact(dot Dot, compact bool) {
	dc.dots.Add(dot)
	if compact {
		dc.Compact()
	}
}

func (dc *DotContext) Compact() {
	dotAdded := true
	for dotAdded {
		dotAdded = false

		for dot := range dc.dots {
			if localSeq, ok := dc.compactContext[dot.id]; ok { // Has entry in compact context
				if dot.seq == localSeq + 1 { // Dot is sequentially after, can compact
					dc.compactContext[dot.id] = dot.seq
					dc.dots.Remove(dot)
					dotAdded = true

				} else if dot.seq <= localSeq { // Dot is already included, ignore
					dc.dots.Remove(dot)
				}

			} else { // No entry in compact context exists
				if dot.seq == 1 { // Can compact
					delete(dc.compactContext, dot.id)
					dc.dots.Remove(dot)
					dotAdded = true
				}
			}
		}
	}
}

func (dc *DotContext) Join(other *DotContext) {
	if dc == other {
		return
	}

	for replicaID, otherSeq := range other.compactContext {
		if localSeq, ok := dc.compactContext[replicaID]; ok {
			dc.compactContext[replicaID] = max(localSeq, otherSeq)
		} else {
			dc.compactContext[replicaID] = otherSeq
		}
	}

	for otherDot := range other.dots {
		dc.dots.Add(otherDot)
	}

	dc.Compact()
}

func (dc *DotContext) Clone() *DotContext {
	clone := NewDotContext()

	for replicaID, seq := range dc.compactContext {
		clone.compactContext[replicaID] = seq
	}

	for dot := range dc.dots {
		clone.dots.Add(dot)
	}

	return clone
}

func (dc *DotContext) String() string {
	return fmt.Sprintf("DotContext{compactContext: %v, dots: %v}", dc.compactContext, dc.dots)
}
