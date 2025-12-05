package node

import pb "sdle-server/proto"

func (n *Node) handlePing(req *pb.Request) error {
	println("Node " + n.addr + " received Ping from " + req.Origin)
	return n.SendReplyOK([]byte("Pong from " + n.addr))
}

func (n *Node) handleFetchRing(req *pb.Request) error {
	println("Node " + n.addr + " received FetchRing from " + req.Origin)
	return n.SendReplyError("FetchRing not implemented")
}

func (n *Node) handleGetHashSpace(req *pb.Request) error {
	println("Node " + n.addr + " received GetHashSpace from " + req.Origin)
	return n.SendReplyError("GetHashSpace not implemented")
}

func (n *Node) handleGossipJoin(req *pb.Request) error {
	println("Node " + n.addr + " received GossipJoin from " + req.Origin)
	return n.SendReplyError("GossipJoin not implemented")
}
