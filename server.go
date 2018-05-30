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

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	req, err := http.NewRequest("GET", url, nil)
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
	p.Thumbnail = fmt.Sprintf("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/%d.png", pokemon1.Id)
	p.Image = fmt.Sprintf("https://img.pokemondb.net/artwork/%s.jpg", name)

	//TODO: Get types
	// var types = [];
	// pokemon.types.forEach(element => {
	//     types.push(element.type.name);
	// });

	//TODO: Get species
	// var species = pokemon.species.name;

	//TODO: Call species API  https://pokeapi.co/api/v2/pokemon-species/...

	//TODO: Get hatitats and flavorText
	// var habitatats = ""
	// var flavorText = "";

	// if(pokemonSpecies != null) {
	//   if(pokemonSpecies.habitat != null) {
	//     habitatats = pokemonSpecies.habitat.name;
	//   }

	//   var flavors = pokemonSpecies.flavor_text_entries;

	//   flavors.forEach(element => {
	//       if(element.language.name === "en") {
	//           flavorText = element.flavor_text.replace(/(?:\r\n|\r|\n|\f)/g, ' ');
	//       }
	//   });

	//   pokemon.flavorText = flavorText;

	return p, nil
	//return p, errors.New("notFound")
}
