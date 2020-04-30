package main

import (
	"context"
	"log"

	pb "github.com/soypita/shippy-service-vessel/proto/vessel"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Vessel struct {
	ID        string `json:"id"`
	Capacity  int32  `json:"capacity"`
	MaxWeight int32  `json:"max_weight"`
	Name      string `json:"name"`
	Available bool   `json:"available"`
	OwnerID   string `json:"owner_id"`
}

type Specification struct {
	Capacity  int32 `json:"capacity"`
	MaxWeight int32 `json:"max_weight"`
}

func MarshalSpecification(spec *pb.Specification) *Specification {
	return &Specification{
		Capacity:  spec.Capacity,
		MaxWeight: spec.MaxWeight,
	}
}

func MarshallVessel(vessel *pb.Vessel) *Vessel {
	return &Vessel{
		ID:        vessel.Id,
		Capacity:  vessel.Capacity,
		MaxWeight: vessel.MaxWeight,
		Name:      vessel.Name,
		Available: vessel.Available,
		OwnerID:   vessel.OwnerId,
	}
}

func UnmarshallVessel(vessel *Vessel) *pb.Vessel {
	return &pb.Vessel{
		Id:        vessel.ID,
		Capacity:  vessel.Capacity,
		MaxWeight: vessel.MaxWeight,
		Name:      vessel.Name,
		Available: vessel.Available,
		OwnerId:   vessel.OwnerID,
	}
}

type Repository interface {
	FindAvailable(context.Context, *Specification) (*Vessel, error)
	Create(context.Context, *Vessel) error
}

type MongoRepository struct {
	collection *mongo.Collection
}

func (m *MongoRepository) FindAvailable(ctx context.Context, spec *Specification) (*Vessel, error) {
	filter := bson.M{
		"capacity":  bson.M{"$gte": spec.Capacity},
		"maxweight": bson.M{"$gte": spec.MaxWeight},
	}

	var vessel Vessel
	log.Println("Collection is", m.collection)
	err := m.collection.FindOne(ctx, filter).Decode(&vessel)
	if err != nil {
		return nil, err
	}
	return &vessel, nil
}

func (m *MongoRepository) Create(ctx context.Context, vessel *Vessel) error {
	_, err := m.collection.InsertOne(ctx, vessel)
	return err
}
