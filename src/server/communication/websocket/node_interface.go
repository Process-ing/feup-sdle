package websocket

import (
	crdt "sdle-server/crdt/shopping"
	pb "sdle-server/proto"
)

type NodeInterface interface {
	HandleShoppingList(list *crdt.ShoppingList) error
	GetShoppingList(id string) (*pb.ShoppingList, error)
}
