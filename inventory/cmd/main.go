package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	inventory_v1 "github.com/ChopX4/raketka/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const grpcPort = 50051

type InventoryService struct {
	inventory_v1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventory_v1.Part
}

func NewInventoryService() *InventoryService {
	service := &InventoryService{
		parts: make(map[string]*inventory_v1.Part),
	}

	id1 := "550e8400-e29b-41d4-a716-446655440000"
	service.parts[id1] = &inventory_v1.Part{
		Uuid:          id1,
		Name:          "Turbocharger GT28",
		Description:   "High-performance ball bearing turbocharger",
		Price:         1250.00, // float64
		StockQuantity: 15,      // int64
	}

	id2 := "67218051-7c5d-4e1a-9f5c-2e3b2e3b2e3b"
	service.parts[id2] = &inventory_v1.Part{
		Uuid:          id2,
		Name:          "Brake Discs Carbon",
		Description:   "Carbon-ceramic brake discs for racing",
		Price:         850.45,
		StockQuantity: 4,
	}

	return service
}

func (s *InventoryService) GetPart(_ context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventory_v1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *InventoryService) ListParts(_ context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	filter := req.GetFilter()

	s.mu.RLock()
	defer s.mu.RUnlock()

	storage := make([]*inventory_v1.Part, 0)

	for _, part := range s.parts {
		if filter == nil {
			storage = append(storage, part)
			continue
		}

		if len(filter.GetUuids()) > 0 {
			found := false
			for _, v := range filter.GetUuids() {
				if part.GetUuid() == v {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if len(filter.GetNames()) > 0 {
			found := false
			for _, v := range filter.GetNames() {
				if part.GetName() == v {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if len(filter.GetCategories()) > 0 {
			found := false
			for _, v := range filter.GetCategories() {
				if part.GetCategory() == v {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if len(filter.GetManufacturerCountries()) > 0 {
			found := false
			for _, v := range filter.GetManufacturerCountries() {
				if part.Manufacturer.GetCountry() == v {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if len(filter.GetTags()) > 0 {
			found := false
			for _, v := range filter.GetTags() {
				for _, tag := range part.GetTags() {
					if tag == v {
						found = true
						break
					}
				}
			}
			if !found {
				continue
			}
		}

		storage = append(storage, part)
	}

	return &inventory_v1.ListPartsResponse{
		Parts: storage,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()
	service := NewInventoryService()

	inventory_v1.RegisterInventoryServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
