package crdt

import (
	"fmt"
	"sdle-server/crdt/generic"
)

type ShoppingList struct {
	crdtID     string
	listID     string
	name       string
	dotContext *crdt.DotContext
	items      *crdt.ORMap[string, *ShoppingItem]
}

func NewShoppingList(crdtID string, listID string, name string) *ShoppingList {
	dotContext := crdt.NewDotContext()

	items := crdt.NewORMap[string, *ShoppingItem](crdtID)
	items.SetContext(dotContext)

	return &ShoppingList{
		crdtID:     crdtID,
		listID:     listID,
		name:       name,
		dotContext: dotContext,
		items:      items,
	}
}

func (sl *ShoppingList) CRDTID() string {
	return sl.crdtID
}

func (sl *ShoppingList) Name() string {
	return sl.name
}

func (sl *ShoppingList) Context() *crdt.DotContext {
	return sl.dotContext
}

func (sl *ShoppingList) SetContext(ctx *crdt.DotContext) {
	sl.dotContext = ctx
	sl.items.SetContext(ctx)
}


func (sl *ShoppingList) Items() []ShoppingItem {
	result := []ShoppingItem{}

	ids := sl.items.Keys()
	for _, id := range ids {
		item := sl.items.Get(id)

		// Due to resets, some items may be null
		if !item.IsNull() {
			result = append(result, *item.Clone())
		}
	}

	return result
}

func (sl *ShoppingList) GetItem(itemID string) *ShoppingItem {
	item := sl.items.Get(itemID)
	if item.IsNull() {
		return nil
	}
	return item.Clone()
}

func (sl *ShoppingList) PutItem(itemID string, name string, quantityDiff int64, acquiredDiff int64) *ShoppingList {
	delta := NewShoppingList(sl.crdtID, sl.listID, "")

	itemsDelta := sl.items.Apply(itemID, func(item *ShoppingItem) *ShoppingItem {
		item.SetItemID(itemID)
		item.SetName(name)

		itemDelta := item.IncQuantity(quantityDiff)
		itemDelta.Join(item.IncAcquired(acquiredDiff))

		return itemDelta
	})

	delta.items = itemsDelta
	delta.SetContext(itemsDelta.Context())

	return delta
}

func (sl *ShoppingList) RemoveItem(itemID string) *ShoppingList {
	delta := NewShoppingList(sl.crdtID, sl.listID, "")

	itemsDelta := sl.items.Remove(itemID)
	delta.items = itemsDelta
	delta.SetContext(itemsDelta.Context())

	return delta
}

func (sl *ShoppingList) Join(other *ShoppingList) {
	sl.items.Join(other.items)
	sl.dotContext.Join(other.dotContext)
}

func (sl *ShoppingList) Clone() *ShoppingList {
	clone := NewShoppingList(sl.crdtID, sl.listID, sl.name)
	clone.dotContext = sl.dotContext.Clone()
	clone.items = sl.items.Clone()
	return clone
}

func (sl *ShoppingList) String() string {
	return fmt.Sprintf("ShoppingList{crdtID: %s, listID: %s, name: %s, dotContext: %v, items: %v}",
		sl.crdtID, sl.listID, sl.name, sl.dotContext, sl.items)
}