package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/micro/go-micro"
	pb "github.com/soypita/shippy-service-vessel/proto/vessel"
)

const (
	defaultHost = "datastore:27017"
)

func createDummyData(repo Repository) {
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Kane's Salty Secret", MaxWeight: 200000, Capacity: 500},
	}
	for _, v := range vessels {
		repo.Create(context.Background(), MarshallVessel(v))
	}
}

func main() {
	srv := micro.NewService(
		micro.Name("shippy.service.vessel"),
	)
	srv.Init()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}
	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.Background())

	vesselCollection := client.Database("shippy").Collection("vessels")

	repo := &MongoRepository{vesselCollection}
	createDummyData(repo)

	h := &handler{repo}
	pb.RegisterVesselServiceHandler(srv.Server(), h)

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
