package main

import (
	"context"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	pb "grpc-go/grpc-go/myservice"
)

type server struct {
	pb.UnimplementedMyServiceServer
	collection *mongo.Collection
}

func NewServer() (*server, error) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://kundan16239:kundan16239@cluster0.89z140s.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	collection := client.Database("grpc-go").Collection("dummy")

	return &server{collection: collection}, nil
}

func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	result, err := s.collection.InsertOne(ctx, bson.M{"name": req.Name})
	if err != nil {
		return nil, err
	}

	insertedID := result.InsertedID.(primitive.ObjectID).Hex()
	return &pb.CreateResponse{Id: insertedID}, nil
}

func (s *server) FindOne(ctx context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var result bson.M

	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	name := result["name"].(string)
	return &pb.FindOneResponse{Name: name}, nil
}

func main() {
	// Initialize the gRPC server
	grpcServer, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to initialize gRPC server: %v", err)
	}

	// Set up gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMyServiceServer(s, grpcServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
