package pb

import (
	"context"
	"fmt"
	"github.com/keets2012/Micro-Go-Practrise/ch9-rpc/pb"
	"google.golang.org/grpc"
)

func main(){
	serviceAddress := "127.0.0.1:1234"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())\
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewStringServiceClient(conn)
	stringReq := &pb.StringRequest{A: "Hello", B: "World"}
	reply, _ := client.Concat(context.Background(), stringReq)
	fmt.Printf("%s\n", reply.Ret)
}