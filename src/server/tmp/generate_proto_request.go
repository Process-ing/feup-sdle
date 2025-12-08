package main

import (
	"fmt"
	"log"
	"os"
	pb "sdle-server/proto"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Connect to WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8000/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	// Create request
	req := &pb.ClientRequest{
		RequestType: &pb.ClientRequest_GetShoppingList_{
			GetShoppingList_: &pb.GetShoppingListRequest{
				Id: "my-shopping-list-id",
			},
		},
	}
	out, err := proto.Marshal(req)
	if err != nil {
		fmt.Println("Failed to marshal request:", err)
		os.Exit(1)
	}

	// Send message
	err = conn.WriteMessage(websocket.BinaryMessage, out)
	if err != nil {
		log.Println("write:", err)
		return
	}
	fmt.Println("Sent GetShoppingListRequest for ID: my-shopping-list-id")

	// Read response
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}

	// Unmarshal and print response
	var shoppingList pb.ShoppingList
	if err := proto.Unmarshal(message, &shoppingList); err != nil {
		log.Fatalln("Failed to unmarshal shopping list:", err)
	}

	fmt.Printf("Received Shopping List:\n%+v\n", &shoppingList)
}