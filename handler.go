package main

import (
	"context"

	pb "github.com/soypita/shippy-service-vessel/proto/vessel"
)

type handler struct {
	repo Repository
}

func (s *handler) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	vessel, err := s.repo.FindAvailable(ctx, MarshalSpecification(req))
	if err != nil {
		return err
	}

	res.Vessel = UnmarshallVessel(vessel)
	return nil
}

func (s *handler) Create(ctx context.Context, req *pb.Vessel, res *pb.Response) error {
	if err := s.repo.Create(ctx, MarshallVessel(req)); err != nil {
		return err
	}

	res.Vessel = req
	return nil
}
