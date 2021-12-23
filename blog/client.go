package blog

import (
	"context"
	"io"
	"log"

	"github.com/jwmwalrus/grpc-go/blog/blogpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DoCreateBlog -
func DoCreateBlog(c blogpb.BlogServiceClient) {
	log.Println("Starting to do a unary RPC for CreateBlog...")

	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			Title:    "A Post",
			AuthorId: "JohnM",
			Content:  "A little bit of this, a little bit of that.",
		},
	}
	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		if s.Code() == codes.InvalidArgument {
			log.Printf("\tError creating blog: %v", s.Message())
		} else if s.Code() == codes.Internal {
			log.Printf("\tInternal error creating blog: %v", s.Message())
		} else {
			log.Fatalf("\tError obtaining response: %v", s.Message())
		}
	}
	log.Printf("\tID: %s", res.Blog.GetId())
}

// DoReadBlog -
func DoReadBlog(c blogpb.BlogServiceClient) {
	log.Println("Starting to do a unary RPC for ReadBlog...")

	req := &blogpb.ReadBlogRequest{
		// BlogId: "61c20d0558075306f71b49d5g", //invalid
		// BlogId: "61c20d0558075306f71b49d9", //not found
		BlogId: "61c20d0558075306f71b49d5", //valid
	}
	res, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		if s.Code() == codes.InvalidArgument {
			log.Fatalf("\tWrong ID provided: %v ", s.Message())
		} else if s.Code() == codes.NotFound {
			log.Fatalf("\tDocument not found: %v", s.Message())
		} else {
			log.Fatalf("\tError obtaining response: %v", s.Message())
		}
	}
	log.Printf("\tID: %v", res.Blog)
}

// DoUpdateBlog -
func DoUpdateBlog(c blogpb.BlogServiceClient) {
	log.Println("Starting to do a unary RPC for UpdateBlog...")

	req := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			// Id: "61c20d0558075306f71b49d5g", //invalid
			// Id: "61c20d0558075306f71b49d9", //not found
			Id:       "61c20d0558075306f71b49d5", //valid
			AuthorId: "SomeoneElse",
			Title:    "I've changed my mind",
			Content:  "Not really.",
		},
	}
	res, err := c.UpdateBlog(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		if s.Code() == codes.InvalidArgument {
			log.Fatalf("\tWrong ID provided: %v ", s.Message())
		} else if s.Code() == codes.NotFound {
			log.Fatalf("\tDocument not found: %v", s.Message())
		} else {
			log.Fatalf("\tError obtaining response: %v", s.Message())
		}
	}
	log.Printf("\tID: %v", res.Blog)
}

// DoDeleteBlog -
func DoDeleteBlog(c blogpb.BlogServiceClient) {
	log.Println("Starting to do a unary RPC for DoDeleteBlog...")

	req := &blogpb.DeleteBlogRequest{
		// BlogId: "61c20d0558075306f71b49d5g", //invalid
		// BlogId: "61c20d0558075306f71b49d9", //not found
		BlogId: "61c20d0558075306f71b49d5", //valid
	}
	res, err := c.DeleteBlog(context.Background(), req)
	if err != nil {
		s := status.Convert(err)
		if s.Code() == codes.InvalidArgument {
			log.Fatalf("\tWrong ID provided: %v ", s.Message())
		} else if s.Code() == codes.NotFound {
			log.Fatalf("\tDocument not found: %v", s.Message())
		} else {
			log.Fatalf("\tError obtaining response: %v", s.Message())
		}
	}
	log.Printf("\tDeleted ID: %v", res.BlogId)
}

// DoListBlog -
func DoListBlog(c blogpb.BlogServiceClient) {
	log.Println("Starting to do a streaming server RPC for DoListBlog...")

	req := &blogpb.ListBlogRequest{}

	stream, err := c.ListBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling ListBlog RPC:%v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while obtaining blog item from RPC call:%v", err)
		}
		log.Printf("\tOne item ID is: %v", res.GetBlog().GetId())
	}
}
