package client

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestClientCall(t *testing.T) {
	ctx := context.Background()
	client, err := NewClient(ctx, "whalethinker.test.consul")
	if err != nil {
		panic(err)
	}

	resp, err := client.Call("/consul_check_ping", http.MethodGet, map[string]string{}, make(map[string]string), "")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp))
}
