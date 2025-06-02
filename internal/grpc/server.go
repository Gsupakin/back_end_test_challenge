package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/Gsupakin/back_end_test_challeng/internal/domain"
	"github.com/Gsupakin/back_end_test_challeng/pkg/validator"
	pb "github.com/Gsupakin/back_end_test_challeng/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserServer implements the gRPC UserService
type UserServer struct {
	pb.UnimplementedUserServiceServer
	userRepo domain.UserRepository
}

// NewUserServer creates a new UserServer instance
func NewUserServer(userRepo domain.UserRepository) *UserServer {
	return &UserServer{
		userRepo: userRepo,
	}
}

// CreateUser implements the CreateUser RPC method
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Validate input
	if err := validator.ValidateUserInput(req.Name, req.Email, req.Password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check if user exists
	_, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
	}

	_, err = s.userRepo.FindByName(ctx, req.Name)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "user with this name already exists")
	}

	// Create user
	now := time.Now()
	user := domain.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		CreatedAt: now,
		UpdatedAt: &now,
	}

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &pb.CreateUserResponse{
		Id: id.Hex(),
	}, nil
}

// GetUser implements the GetUser RPC method
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Parse ID
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:        user.ID.Hex(),
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(*user.UpdatedAt),
		},
	}, nil
}

// AuthInterceptor implements gRPC interceptor for authentication
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Skip auth for CreateUser
	if info.FullMethod == "/user.UserService/CreateUser" {
		return handler(ctx, req)
	}

	// Get metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	// Get token
	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	// TODO: Validate token
	// For now, just check if token exists
	if tokens[0] == "" {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
}
