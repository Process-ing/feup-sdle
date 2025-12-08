package websocket

import (
	pb "sdle-server/proto"
)

type NodeInterface interface {
	HandleShoppingList(list *pb.ShoppingList) error
	GetShoppingList(id string) (*pb.ShoppingList, error)
}
