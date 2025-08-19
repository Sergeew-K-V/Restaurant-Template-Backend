package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"restaurant-backend/src/config"
	"restaurant-backend/src/database"
)

func main() {
	// CLI flags
	var runMigrations = flag.Bool("migrate", false, "Run database migrations")
	var showMigrationStatus = flag.Bool("migrate-status", false, "Show migration status")
	var showHelp = flag.Bool("help", false, "Show this help message")
	var showVersion = flag.Bool("version", false, "Show version information")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Database Migration Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  go run migrate.go [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Available Commands:\n")
		fmt.Fprintf(os.Stderr, "  -migrate         Run database migrations\n")
		fmt.Fprintf(os.Stderr, "  -migrate-status  Show migration status\n")
		fmt.Fprintf(os.Stderr, "  -version         Show version information\n")
		fmt.Fprintf(os.Stderr, "  -help            Show this help message\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  go run migrate.go        # Run migrations (default)\n")
		fmt.Fprintf(os.Stderr, "  go run migrate.go -migrate # Run migrations explicitly\n")
		fmt.Fprintf(os.Stderr, "  go run migrate.go -migrate-status # Check migration status\n")
		fmt.Fprintf(os.Stderr, "  go run migrate.go -version # Show version\n")
	}

	flag.Parse()

	// Handle help and version flags
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Println("Migration Tool v1.0.0")
		os.Exit(0)
	}

	// Load configuration
	envConfig := config.LoadGlobalConfig()

	// Get database connection
	db, err := database.GetDBConnection(envConfig.DB)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	defer database.CloseDB(db)

	// Handle CLI commands
	if *runMigrations {
		fmt.Println("Running database migrations...")
		if err := database.RunMigrations(db); err != nil {
			log.Fatal("Error running migrations:", err)
		}
		fmt.Println("Migrations completed successfully")
		os.Exit(0)
	}

	if *showMigrationStatus {
		fmt.Println("Migration status:")
		if err := database.GetMigrationStatus(db); err != nil {
			log.Fatal("Error getting migration status:", err)
		}
		os.Exit(0)
	}

	// Default behavior - run migrations
	fmt.Println("Running database migrations...")
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Error running migrations:", err)
	}
	fmt.Println("Migrations completed successfully")
}
