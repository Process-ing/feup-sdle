package crdt

import (
    "testing"
	"reflect"
)

// === Utility functions for tests ===

func contextsEqual(ctx1, ctx2 *DotContext) bool {
	if !reflect.DeepEqual(ctx1.compactContext, ctx2.compactContext) {
		return false
	}
	if !reflect.DeepEqual(ctx1.dots, ctx2.dots) {
		return false
	}
	return true
}

// === Tests ===

func TestDotContext_NewDotContext(t *testing.T) {
    ctx := NewDotContext()
    if ctx == nil {
        t.Fatal("NewDotContext() returned nil")
    }
    if len(ctx.compactContext) != 0 {
        t.Errorf("Expected compactContext to be empty, got %v", ctx.compactContext)
    }
    if ctx.dots.Size() != 0 {
        t.Errorf("Expected dots to be empty, got %v", ctx.dots)
    }
}

func TestDotContext_MakeDot(t *testing.T) {
    ctx := NewDotContext()
    dot := ctx.MakeDot("node1")

    if dot.id != "node1" {
        t.Errorf("Expected dot.id to be 'node1', got %s", dot.id)
    }
    if dot.seq != 1 {
        t.Errorf("Expected dot.seq to be 1, got %d", dot.seq)
    }

    dot2 := ctx.MakeDot("node1")
    if dot2.seq != 2 {
        t.Errorf("Expected dot2.seq to be 2, got %d", dot2.seq)
    }
}

func TestDotContext_Compact(t *testing.T) {
    ctx := NewDotContext()
    ctx.InsertDot(NewDot("node1", 1))
    ctx.InsertDot(NewDot("node1", 2))
    ctx.Compact()

    if seq, ok := ctx.compactContext["node1"]; !ok || seq != 2 {
        t.Errorf("Expected compactContext[node1] to be 2, got %d", seq)
    }
    if ctx.dots.Size() != 0 {
        t.Errorf("Expected dots to be empty after compact, got %v", ctx.dots)
    }
}

func TestDotContext_InsertDot(t *testing.T) {
    ctx := NewDotContext()
    dot := NewDot("node1", 1)
    ctx.InsertDot(dot)

    if _, ok := ctx.compactContext["node1"]; !ok {
        t.Errorf("Expected dots to contain %v, but it does not", dot)
    }
	if seq, ok := ctx.compactContext["node1"]; !ok || seq != 1 {
		t.Errorf("Expected compactContext[node1] to be 1, got %d", seq)
	}
}

func TestDotContext_Join(t *testing.T) {
    ctx1 := NewDotContext()
    ctx1.InsertDot(NewDot("node1", 1))
	ctx1.InsertDot(NewDot("node2", 1))
    ctx1.InsertDot(NewDot("node2", 2))

    ctx2 := NewDotContext()
    ctx2.InsertDot(NewDot("node1", 1))
    ctx2.InsertDot(NewDot("node1", 2))
    ctx2.InsertDot(NewDot("node1", 3))
    ctx2.InsertDot(NewDot("node3", 1))

    ctx1.Join(ctx2)

	if !ctx1.Knows(NewDot("node1", 3)) {
		t.Errorf("Expected ctx1 to know Dot(node1, 3)")
	}
	if !ctx1.Knows(NewDot("node2", 2)) {
		t.Errorf("Expected ctx1 to know Dot(node2, 2)")
	}
	if !ctx1.Knows(NewDot("node3", 1)) {
		t.Errorf("Expected ctx1 to know Dot(node3, 1)")
	}
	if ctx1.Knows(NewDot("node1", 4)) {
		t.Errorf("Did not expect ctx1 to know Dot(node1, 4)")
	}
	if ctx1.Knows(NewDot("node2", 3)) {
		t.Errorf("Did not expect ctx1 to know Dot(node2, 3)")
	}
	if ctx1.Knows(NewDot("node3", 2)) {
		t.Errorf("Did not expect ctx1 to know Dot(node3, 2)")
	}

	// Test commutativity

	ctx2.Join(ctx1)

	if !contextsEqual(ctx1, ctx2) {
		t.Errorf("Expected ctx1 and ctx2 to be equal after joining")
	}
}

func TestDotContext_Clone(t *testing.T) {
    ctx := NewDotContext()
    ctx.InsertDot(NewDot("node1", 1))
    ctx.InsertDot(NewDot("node2", 1))
	ctx.InsertDot(NewDot("node2", 2))

    clone := ctx.Clone()

	if !reflect.DeepEqual(ctx.compactContext, clone.compactContext) {
		t.Errorf("Expected clone.compactContext to equal to ctx.compactContext")
	}

	if !reflect.DeepEqual(ctx.dots, clone.dots) {
		t.Errorf("Expected clone.dots to equal to ctx.dots")
	}
}

func TestDotContext_Copy(t *testing.T) {
    ctx1 := NewDotContext()
    ctx1.InsertDot(NewDot("node1", 1))
    ctx1.InsertDot(NewDot("node2", 2))

    ctx2 := NewDotContext()
    ctx2.Copy(ctx1)

    if len(ctx2.compactContext) != len(ctx1.compactContext) {
        t.Errorf("Expected ctx2.compactContext to have length %d, got %d", len(ctx1.compactContext), len(ctx2.compactContext))
    }
    if ctx2.dots.Size() != ctx1.dots.Size() {
        t.Errorf("Expected ctx2.dots to have size %d, got %d", ctx1.dots.Size(), ctx2.dots.Size())
    }
}

func TestDotContext_String(t *testing.T) {
    ctx := NewDotContext()
    ctx.InsertDot(NewDot("node1", 1))
    ctx.InsertDot(NewDot("node2", 2))

    str := ctx.String()
    if str == "" {
        t.Errorf("Expected non-empty string representation, got empty string")
    }
}