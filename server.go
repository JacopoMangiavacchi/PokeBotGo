package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/JacopoMangiavacchi/PokeBotGo/pokemon"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

//SearchPokemon(*PokeInput, PokeBot_SearchPokemonServer) error
// SearchPokemon implements PokeBotServer
func (s *server) SearchPokemon(input *pb.PokeInput, stream pb.PokeBot_SearchPokemonServer) error {
	fmt.Println(input.Name)
	//p := pb.Pokemon{}
	var p pb.Pokemon
	p.Name = "Jacopo"
	stream.Send(&p)

	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPokeBotServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
