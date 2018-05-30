package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	pb "github.com/JacopoMangiavacchi/PokeBotGo/pokemon"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement PokeBotServer.
type server struct{}

// SearchPokemon(*PokeInput, PokeBot_SearchPokemonServer) error
// SearchPokemon implements PokeBotServer
func (s *server) SearchPokemon(input *pb.PokeInput, stream pb.PokeBot_SearchPokemonServer) error {
	fmt.Println(input.Name)
	p, error := getPokemon(input.Name)

	if error == nil {
		stream.Send(&p)
	} else {

	}

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

type pokemon struct {
	Id     int32  `json:"id"`
	Name   string `json:"name"`
	Height int32  `json:"height"`
	Weight int32  `json:"weight"`
}

func getPokemon(name string) (pb.Pokemon, error) {
	var p pb.Pokemon

	req, err := http.NewRequest("GET", "https://pokeapi.co/api/v2/pokemon/bulbasaur", nil)
	if err != nil {
		return p, errors.New("wrongUrl")
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return p, errors.New("notFound")
	}

	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
		return p, errors.New("wrongBody")
	}

	pokemon1 := pokemon{}
	jsonErr := json.Unmarshal(body, &pokemon1)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return p, errors.New("jsonUnmarshal")
	}

	p.Name = pokemon1.Name
	p.Id = pokemon1.Id
	p.Height = pokemon1.Height
	p.Weight = pokemon1.Weight

	return p, nil
	//return p, errors.New("notFound")
}
