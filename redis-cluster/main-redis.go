package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
    "strconv"
    // "testing"
    "time"
)

func main() {
    client := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{
            "localhost:7009",
            "localhost:7001",
            "localhost:7002",
            "localhost:7003",
            "localhost:7004",
            "localhost:7005",
        },
		// Username:	   "default",
        // Password:      "redispw", 
        RouteRandomly: true,
    })
	
    start := time.Now()
    for i := 0; i < 100000; i++ {
        ctx := context.Background()
        err := client.Set(ctx, "name"+strconv.Itoa(i), "Hsiao"+strconv.Itoa(i), 0).Err()
        if err != nil {
            panic(err)
        }

        val, err := client.Get(ctx, "name"+strconv.Itoa(i)).Result()
        if err != nil {
            panic(err)
        }
        fmt.Println("key", val)
    }
	
    end := time.Now()
    fmt.Printf("TestRedis Runtime: %v\n", end.Sub(start))
}