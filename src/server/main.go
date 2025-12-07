package main

import (
	"fmt"
	"os"
	"os/signal"
	"sdle-server/node"
	"syscall"
	"time"
)

func create_node(id string) *node.Node {
	dataDir := "./data/" + id

	node, err := node.New(id, dataDir)

	if err != nil {
		panic(err)
	}

	return node
}

func main() {
	// wsHandler := websocket.NewWebSocketHandler()
	//
	// // Register handlers
	// http.Handle("/ws", wsHandler)
	//
	// log.Println("Starting server on :8080")
	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	log.Fatal("Failed to start server on port 8080: ", err)
	// }

	nodes := []*node.Node{
		create_node("localhost:5000"),
		create_node("localhost:5001"),
		create_node("localhost:5002"),
		create_node("localhost:5003"),
	}

	errCh := make(chan error, 2)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Start nodes concurrently so Start() doesn't block the rest of main
	for i, nd := range nodes {
		go func(i int, n *node.Node) {
			if err := n.StartReceiving(); err != nil {
				errCh <- err
			}
		}(i, nd)

		time.Sleep(300 * time.Millisecond)

		nd.JoinToRing(nodes[0].GetAddress())

	}

	// Let the ring stabilize
	time.Sleep(2 * time.Second)

	fmt.Println("\n--- Ring stabilized, starting test ---")

	// Pick any node to be our "client" coordinator, e.g., nodes[0]
	coordinatorNode := nodes[0]
	testKey := "my-awesome-key"
	testValue := []byte("my-awesome-value")

	// 1. HAS (should be false)
	fmt.Printf("\n[1] Client sending HAS for key '%s' to coordinator %s\n", testKey, coordinatorNode.GetID())
	has, err := coordinatorNode.Has(testKey)
	if err != nil {
		// Expecting a "key not found" error, which the store might return.
		// A more robust implementation would check for specific error types.
		fmt.Printf("HAS check completed (key does not exist).\n")
	} else if !has {
		fmt.Printf("HAS check successful! Key does not exist.\n")
	} else {
		fmt.Printf("HAS check failed! Key should not exist yet.\n")
	}
	time.Sleep(500 * time.Millisecond)

	// 2. PUT
	fmt.Printf("\n[2] Client sending PUT for key '%s' to coordinator %s\n", testKey, coordinatorNode.GetID())
	if err := coordinatorNode.Put(testKey, testValue); err != nil {
		fmt.Println("PUT failed:", err)
	} else {
		fmt.Println("PUT successful!")
	}
	time.Sleep(500 * time.Millisecond)

	// 3. HAS (should be true)
	fmt.Printf("\n[3] Client sending HAS for key '%s' to coordinator %s\n", testKey, coordinatorNode.GetID())
	has, err = coordinatorNode.Has(testKey)
	if err != nil {
		fmt.Println("HAS check failed:", err)
	} else if has {
		fmt.Printf("HAS check successful! Key exists.\n")
	} else {
		fmt.Printf("HAS check failed! Key should exist now.\n")
	}
	time.Sleep(500 * time.Millisecond)

	// 4. GET
	fmt.Printf("\n[4] Client sending GET for key '%s' to coordinator %s\n", testKey, coordinatorNode.GetID())
	retrievedValue, err := coordinatorNode.Get(testKey)
	if err != nil {
		fmt.Println("GET failed:", err)
	} else {
		fmt.Printf("GET successful! Value: %s\n", string(retrievedValue))
	}
	time.Sleep(500 * time.Millisecond)

	// 5. DELETE
	fmt.Printf("\n[5] Client sending DELETE for key '%s' to coordinator %s\n", testKey, coordinatorNode.GetID())
	if err := coordinatorNode.Delete(testKey); err != nil {
		fmt.Println("DELETE failed:", err)
	} else {
		fmt.Println("DELETE successful!")
	}
	time.Sleep(500 * time.Millisecond)

	// 6. HAS (should be false again)
	fmt.Printf("\n[6] Client sending HAS for key '%s' to coordinator %s\n", testKey, coordinatorNode.GetID())
	has, err = coordinatorNode.Has(testKey)
	if err != nil {
		// Again, expecting an error for a non-existent key.
		fmt.Printf("HAS check completed (key does not exist).\n")
	} else if !has {
		fmt.Printf("HAS check successful! Key does not exist anymore.\n")
	} else {
		fmt.Printf("HAS check failed! Key should be deleted.\n")
	}

	fmt.Println("\n--- Test finished ---")

	// Wait for shutdown signal
	select {
	case <-sigCh:
		println("\nReceived interrupt signal, shutting down...")
	case err := <-errCh:
		println("Error occurred:", err.Error())
	}

	// Gracefully close all nodes
	for _, nd := range nodes {
		if err := nd.StopReceiving(); err != nil {
			println("Error closing node:", err.Error())
		}
	}
	println("All nodes closed successfully")
}
