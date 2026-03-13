package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/dukerupert/skalkaho/internal/auth"
	"github.com/dukerupert/skalkaho/internal/database"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := database.Open(databaseURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() { _ = db.Close() }()

	queries := repository.New(db)
	ctx := context.Background()

	switch os.Args[1] {
	case "create":
		createUser(ctx, queries, os.Args[2:])
	case "list":
		listUsers(ctx, queries)
	case "delete":
		deleteUser(ctx, queries, os.Args[2:])
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Skalkaho User Management

Usage:
  go run ./cmd/seed create --email <email> --password <password> --name <name>
  go run ./cmd/seed list
  go run ./cmd/seed delete --email <email>

Note: Run the server first (go run ./cmd/server) to ensure migrations are applied.
`)
}

func createUser(ctx context.Context, queries *repository.Queries, args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	email := fs.String("email", "", "User email (required)")
	password := fs.String("password", "", "User password (required)")
	name := fs.String("name", "", "User name (required)")
	if err := fs.Parse(args); err != nil {
		log.Fatal(err)
	}

	if *email == "" || *password == "" || *name == "" {
		fmt.Fprintln(os.Stderr, "Error: --email, --password, and --name are all required")
		fs.Usage()
		os.Exit(1)
	}

	// Check if user already exists
	_, err := queries.GetUserByEmail(ctx, *email)
	if err == nil {
		log.Fatalf("User with email %q already exists", *email)
	}
	if err != sql.ErrNoRows {
		log.Fatalf("Checking for existing user: %v", err)
	}

	passwordHash, err := auth.HashPassword(*password)
	if err != nil {
		log.Fatalf("Hashing password: %v", err)
	}

	id := generateID()
	user, err := queries.CreateUser(ctx, repository.CreateUserParams{
		ID:           id,
		Email:        *email,
		PasswordHash: passwordHash,
		Name:         *name,
		Status:       "active",
	})
	if err != nil {
		log.Fatalf("Creating user: %v", err)
	}

	fmt.Printf("Created user: %s (%s) [id: %s]\n", user.Name, user.Email, user.ID)
}

func listUsers(ctx context.Context, queries *repository.Queries) {
	users, err := queries.ListUsers(ctx)
	if err != nil {
		log.Fatalf("Listing users: %v", err)
	}

	if len(users) == 0 {
		fmt.Println("No users found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tEMAIL\tNAME\tSTATUS\tCREATED")
	for _, u := range users {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", u.ID, u.Email, u.Name, u.Status, u.CreatedAt.Format("2006-01-02"))
	}
	w.Flush()
}

func deleteUser(ctx context.Context, queries *repository.Queries, args []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	email := fs.String("email", "", "User email (required)")
	if err := fs.Parse(args); err != nil {
		log.Fatal(err)
	}

	if *email == "" {
		fmt.Fprintln(os.Stderr, "Error: --email is required")
		fs.Usage()
		os.Exit(1)
	}

	user, err := queries.GetUserByEmail(ctx, *email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("User with email %q not found", *email)
		}
		log.Fatalf("Looking up user: %v", err)
	}

	if err := queries.DeleteUser(ctx, user.ID); err != nil {
		log.Fatalf("Deleting user: %v", err)
	}

	fmt.Printf("Deleted user: %s (%s)\n", user.Name, user.Email)
}

func generateID() string {
	return uuid.New().String()[:20]
}
