package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

func main() {
	rpcServer := rpc.NewServer()

	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcServer.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	tweety := NewTweetyService()

	err := rpcServer.RegisterService(tweety, "")
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.Handle("/", rpcServer)

	http.ListenAndServe(":1337", router)
}

type TweetyService struct{}

type HelloResponse struct {
	Result string
}

type HelloRequest struct {
	Subject, Content string
}

func (t *TweetyService) Hello(r *http.Request, req *HelloRequest, result *HelloResponse) error {
	*result = HelloResponse{Result: fmt.Sprintf("Hello subject and content: %s ; %s", req.Subject, req.Content)}
	return nil
}

func NewTweetyService() *TweetyService {
	return &TweetyService{}
}
