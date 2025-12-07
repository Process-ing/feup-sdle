package crdt

import "sdle-server/crdt/generic"

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

func (sl *ShoppingList) PutItem(itemID string, name string, quantityDiff int64, acquiredDiff int64) *ShoppingList {
	delta := NewShoppingList(sl.crdtID, "", name)

	itemsDelta := sl.items.Apply(itemID, func(item *ShoppingItem) *ShoppingItem {
		if item.Name() == "" {
			item.SetName(name)
		}

		itemDelta := item.IncQuantity(quantityDiff)
		itemDelta.Join(item.IncAcquired(acquiredDiff))

		return itemDelta
	})

	delta.SetContext(itemsDelta.Context())

	return delta
}