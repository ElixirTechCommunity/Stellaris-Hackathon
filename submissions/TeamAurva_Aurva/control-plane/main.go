package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aurva/compliance-engine/control-plane/api"
	"github.com/aurva/compliance-engine/control-plane/server"
	pb "github.com/aurva/compliance-engine/proto"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

const (
	grpcPort = ":50055"
	httpPort = ":9090"
	dbConnString = "host=localhost port=5432 user=aurva password=aurva_dev_password dbname=aurva sslmode=disable"
)

func main() {
	// Connect to PostgreSQL
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = dbConnString
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	log.Println("✓ Connected to PostgreSQL")

	// Create gRPC server
	grpcServer := grpc.NewServer()
	complianceServer := server.NewComplianceServer(db)
	pb.RegisterComplianceServiceServer(grpcServer, complianceServer)
	
	// Create HTTP server
	httpServer := api.NewHTTPServer(db, complianceServer)
	mux := http.NewServeMux()
	httpServer.RegisterRoutes(mux)
	
	// Add CORS middleware
	handler := corsMiddleware(mux)
	
	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	var wg sync.WaitGroup
	
	// Start gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("Failed to listen on %s: %v", grpcPort, err)
		}
		log.Printf("✓ gRPC server listening on %s", grpcPort)
		
		go func() {
			<-ctx.Done()
			grpcServer.GracefulStop()
		}()
		
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()
	
	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv := &http.Server{
			Addr:    httpPort,
			Handler: handler,
		}
		log.Printf("✓ HTTP server listening on %s", httpPort)
		
		go func() {
			<-ctx.Done()
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()
			srv.Shutdown(shutdownCtx)
		}()
		
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()
	
	log.Println("✓ Aurva Control Plane is running")
	log.Println("  - gRPC API: localhost:50055")
	log.Println("  - HTTP API: http://localhost:9090")
	log.Println("  - Health check: http://localhost:9090/health")
	
	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	
	log.Println("Shutting down gracefully...")
	cancel()
	wg.Wait()
	log.Println("✓ Shutdown complete")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
