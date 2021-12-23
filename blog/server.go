package blog

import (
	"context"

	"github.com/jwmwalrus/grpc-go/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Server defines the blog server
type Server struct {
	Coll *mongo.Collection
}

// Item -
type Item struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

func (i *Item) toBlogpb() *blogpb.Blog {
	return &blogpb.Blog{
		Id:       i.ID.Hex(),
		Title:    i.Title,
		AuthorId: i.AuthorID,
		Content:  i.Content,
	}
}

// CreateBlog -
func (s *Server) CreateBlog(_ context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()

	data := Item{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	ior, err := s.Coll.InsertOne(context.Background(), data)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Cannot insert into database: %v", err)
	}

	id, ok := ior.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, grpc.Errorf(codes.Internal, "Cannot convert into OID")
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       id.Hex(),
			Title:    blog.Title,
			AuthorId: blog.AuthorId,
			Content:  blog.Content,
		},
	}, nil
}

// ReadBlog -
func (s *Server) ReadBlog(_ context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	blogID, err := primitive.ObjectIDFromHex(req.BlogId)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Cannot parse ID: %v", err)
	}

	data := &Item{}

	filter := bson.M{"_id": blogID}
	sr := s.Coll.FindOne(context.Background(), filter)
	if err := sr.Decode(data); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "Cannot find blog with specified ID: %v", err)
	}

	return &blogpb.ReadBlogResponse{Blog: data.toBlogpb()}, nil
}

// UpdateBlog -
func (s *Server) UpdateBlog(_ context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	blog := req.GetBlog()
	blogID, err := primitive.ObjectIDFromHex(blog.Id)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Cannot parse ID: %v", err)
	}

	data := &Item{}

	filter := bson.M{"_id": blogID}
	sr := s.Coll.FindOne(context.Background(), filter)
	if err := sr.Decode(data); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "Cannot find blog with specified ID: %v", err)
	}

	data.AuthorID = blog.GetAuthorId()
	data.Title = blog.GetTitle()
	data.Content = blog.GetContent()

	_, err = s.Coll.ReplaceOne(context.Background(), filter, data)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Cannot update blog: %v", err)
	}

	return &blogpb.UpdateBlogResponse{Blog: data.toBlogpb()}, nil
}

// DeleteBlog -
func (s *Server) DeleteBlog(_ context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	blogID, err := primitive.ObjectIDFromHex(req.BlogId)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Cannot parse ID: %v", err)
	}

	filter := bson.M{"_id": blogID}
	dr, err := s.Coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Cannot delete blog: %v", err)
	}
	if dr.DeletedCount < 1 {
		return nil, grpc.Errorf(codes.NotFound, "Cannot find blog entry to delete, with specified ID")
	}

	return &blogpb.DeleteBlogResponse{BlogId: req.BlogId}, nil
}

// ListBlog -
func (s *Server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	cur, err := s.Coll.Find(context.Background(), bson.D{{}})
	if err != nil {
		return grpc.Errorf(codes.Internal, "Error finding blog items: %v", err)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := &Item{}
		if err := cur.Decode(data); err != nil {
			return grpc.Errorf(codes.Internal, "Error decoding blog item: %v", err)
		}
		res := &blogpb.ListBlogResponse{Blog: data.toBlogpb()}
		if err := stream.Send(res); err != nil {
			return grpc.Errorf(codes.Internal, "Error decoding blog item: %v", err)
		}
	}
	if err := cur.Err(); err != nil {
		return grpc.Errorf(codes.Internal, "Cursor error: %v", err)
	}
	return nil
}
