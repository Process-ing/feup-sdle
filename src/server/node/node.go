package node

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"sdle-server/communication/websocket"
	pb "sdle-server/proto"
	"sdle-server/ringview"
	"sdle-server/storage"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/pebbe/zmq4"
	"google.golang.org/protobuf/proto"
)

type Node struct {
	id         string
	addr       string
	wsAddr     string
	ringView   *ringview.RingView
	store      storage.Store
	repSock    *zmq4.Socket
	httpServer *http.Server
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

func New(id string, baseDir string) (*Node, error) {
	addr := idToAddr(id)

	// Configure websockets
	host, portStr, err := net.SplitHostPort(id)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID format: %w", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port in node ID: %w", err)
	}
	wsPort := port + 3000 // 5000 -> 8000, etc
	wsAddr := net.JoinHostPort(host, strconv.Itoa(wsPort))

	wsHandler := websocket.NewWebSocketHandler()
	mux := http.NewServeMux()
	mux.Handle("/ws", wsHandler)

	httpServer := &http.Server{
		Addr:    wsAddr,
		Handler: mux,
	}

	ringView := ringview.New()

	dir := filepath.Join(baseDir, id)
	store, err := storage.Open(dir)
	if err != nil {
		return nil, err
	}

	rep, err := zmq4.NewSocket(zmq4.REP)
	if err != nil {
		_ = store.Close()
		return nil, err
	}

	if err := rep.Bind(addr); err != nil {
		_ = rep.Close()
		_ = store.Close()
		return nil, err
	}

	return &Node{
			id:         id,
			addr:       addr,
			wsAddr:     wsAddr,
			ringView:   ringView,
			store:      *store,
			repSock:    rep,
			httpServer: httpServer,
			stopCh:     make(chan struct{}),
		},
		nil
}

func (n *Node) Start(errCh chan<- error) {
	n.wg.Add(2)
	go n.startReceiving(errCh)
	go func() {
		defer n.wg.Done()
		n.log("starting WebSocket server at " + n.wsAddr)
		if err := n.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("WebSocket server error: %w", err)
		}
	}()
}

func (n *Node) Stop() error {
	close(n.stopCh) // Signal all goroutines to stop

	// Shutdown websockets server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := n.httpServer.Shutdown(ctx); err != nil {
		log.Printf("WebSocket server shutdown error: %v", err)
	}

	n.wg.Wait() // Wait for all goroutines to finish

	var firstErr error
	if n.repSock != nil {
		if err := n.repSock.Close(); err != nil {
			firstErr = err
		}
	}

	if err := n.store.Close(); err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}

func idToAddr(id string) string {
	return "tcp://" + id
}

// Get the current ring view from a target node and update the local ring view
func (n *Node) UpdateRingView(targetAddr string) error {
	resp, err := n.sendFetchRing(targetAddr)

	if err != nil {
		return err
	}

	fetchRingResp := resp.GetFetchRing()
	if fetchRingResp == nil {
		return nil
	}

	// Create new RingView from received tokenToNode map
	newRingView := ringview.NewFromTokenMap(fetchRingResp.TokenToNode)
	n.ringView = newRingView

	n.log("Updated ring view: " + n.ringView.ToString())

	return nil
}

// Adds the new node to the ring - first get the current ring view from a target node, then imports the data for its tokens and finally adds itself to the ring (using gossip to inform other nodes)
func (n *Node) JoinToRing(targetAddr string) error {
	n.UpdateRingView(targetAddr)

	tokens := n.ringView.JoinToRing(n.GetID())

	n.log("joined the ring with tokens: " + fmt.Sprint(tokens))

	for _, token := range tokens {
		// TODO:
		// The code should find the node responsible for the given token,
		// request the data for that token from the responsible node,
		// and then store the received data in the local storage.

		n.log("would request data for token " + fmt.Sprint(token) + " from the responsible node.")

	}

	neighborsGossip := n.ringView.GetGossipNeighborsNodes(n.GetID())
	n.log("Starting gossip to inform other nodes about my joining. Neighbors: " + fmt.Sprint(neighborsGossip))

	for _, nodeId := range neighborsGossip {
		nodeAddr := idToAddr(nodeId)
		resp, err := n.sendJoinGossip(nodeAddr, n.GetID(), tokens)

		n.log("Gossip Response: Ok=" + fmt.Sprint(resp.Ok) + ", Error='" + fmt.Sprint(err) + "'")
	}

	return nil
}

func (n *Node) GetAddress() string {
	return n.addr
}

func (n *Node) GetID() string {
	return n.id
}

func (n *Node) log(msg string) {
	fmt.Printf("[Node %s]: %s\n", n.id, msg)
}

func (n *Node) startReceiving(errCh chan<- error) {
	defer n.wg.Done()
	n.log("ZMQ socket started at " + n.addr)

	poller := zmq4.NewPoller()
	poller.Add(n.repSock, zmq4.POLLIN)

	for {
		select {
		case <-n.stopCh:
			n.log("ZMQ receiver stopping.")
			return
		default:
		}

		sockets, err := poller.Poll(100 * time.Millisecond)
		if err != nil {
			// ETERM is expected on shutdown, so we don't send it to the error channel.
			if zmq4.AsErrno(err) != zmq4.ETERM {
				errCh <- fmt.Errorf("ZMQ poll error: %w", err)
			}
			return
		}

		if len(sockets) == 0 {
			continue // Poll timed out, loop again
		}

		msgBytes, err := n.repSock.RecvBytes(zmq4.DONTWAIT)
		if err != nil {
			// EAGAIN is expected when there's nothing to receive, just continue
			if zmq4.AsErrno(err) != zmq4.Errno(syscall.EAGAIN) {
				errCh <- fmt.Errorf("ZMQ receive error: %w", err)
			}
			continue
		}

		var req pb.Request
		if err := proto.Unmarshal(msgBytes, &req); err != nil {
			_ = n.sendResponseError("failed to unmarshal request: " + err.Error())
			continue
		}

		switch req.GetRequestType().(type) {
		case *pb.Request_Ping:
			n.handlePing(&req)
		case *pb.Request_FetchRing:
			n.handleFetchRing(&req)
		case *pb.Request_GetHashSpace:
			n.handleGetHashSpace(&req)
		case *pb.Request_GossipJoin:
			n.handleGossipJoin(&req)
		case *pb.Request_Get:
			n.handleGet(&req)
		case *pb.Request_Put:
			n.handlePut(&req)
		case *pb.Request_Delete:
			n.handleDelete(&req)
		case *pb.Request_Has:
			n.handleHas(&req)
		default:
			_ = n.sendResponseError("unknown request type")
		}
	}
}

func (n *Node) GetRingView() *ringview.RingView {
	return n.ringView
}
