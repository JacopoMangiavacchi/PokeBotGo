# PokeBotGo
PokeBot Go GRPC BFF


# Protoc compilation

    protoc -I pokemon/ pokemon/pokemon.proto --go_out=plugins=grpc:pokemon


