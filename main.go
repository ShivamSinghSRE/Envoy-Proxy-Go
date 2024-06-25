package main

import (
    "context"
    "log"
    "net"
    "time"

    "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
    "github.com/envoyproxy/go-control-plane/pkg/server/v3"
    "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
    "github.com/envoyproxy/go-control-plane/pkg/test/v3"
    "google.golang.org/grpc"
)

func main() {
    // Create a snapshot cache
    snapshotCache := cache.NewSnapshotCache(true, cache.IDHash{}, nil)

    // Create a new gRPC server
    grpcServer := grpc.NewServer()

    // Create a new xDS server
    xdsServer := server.NewServer(context.Background(), snapshotCache, nil)

    // Register the xDS server with the gRPC server
    resource.RegisterServer(grpcServer, xdsServer)

    // Listen on a port
    lis, err := net.Listen("tcp", ":18000")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    // Start the gRPC server
    go func() {
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("Failed to serve: %v", err)
        }
    }()

    // Create a snapshot with a simple configuration
    nodeID := "test-node"
    snapshot := cache.NewSnapshot(
        "1", // version
        nil, // endpoints
        nil, // clusters
        nil, // routes
        nil, // listeners
        nil, // runtimes
        nil, // secrets
    )

    // Set the snapshot for the node
    if err := snapshotCache.SetSnapshot(context.Background(), nodeID, snapshot); err != nil {
        log.Fatalf("Failed to set snapshot: %v", err)
    }

    log.Println("xDS server is running...")

    // Keep the server running
    select {}
}
