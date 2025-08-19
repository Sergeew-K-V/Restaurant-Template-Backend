package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Migration struct {
	Version  string
	Filename string
	SQL      string
}

func RunMigrations(db *sql.DB) error {
	log.Println("Starting database migrations...")

	// Создаем таблицу для отслеживания выполненных миграций
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Получаем список выполненных миграций
	executedMigrations, err := getExecutedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get executed migrations: %v", err)
	}

	// Загружаем все миграции из папки migrations
	migrations, err := loadMigrationsFromFiles()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %v", err)
	}

	// Сортируем миграции по версии
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Выполняем миграции по порядку
	for _, migration := range migrations {
		if !isMigrationExecuted(executedMigrations, migration.Version) {
			log.Printf("Executing migration %s (%s)...", migration.Version, migration.Filename)

			if err := executeMigration(db, migration); err != nil {
				return fmt.Errorf("failed to execute migration %s: %v", migration.Version, err)
			}

			log.Printf("Migration %s completed successfully", migration.Version)
		} else {
			log.Printf("Migration %s already executed, skipping", migration.Version)
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}

func loadMigrationsFromFiles() ([]Migration, error) {
	migrationsDir := "src/database/migrations"

	// Получаем список всех .sql файлов в папке migrations
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %v", err)
	}

	var migrations []Migration

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			// Извлекаем версию из имени файла (например, "1.0" из "1.0-create_users_table.sql")
			version := extractVersionFromFilename(file.Name())
			if version == "" {
				log.Printf("Warning: skipping file %s - could not extract version", file.Name())
				continue
			}

			// Читаем содержимое SQL файла
			filePath := filepath.Join(migrationsDir, file.Name())
			sqlContent, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
			}

			migration := Migration{
				Version:  version,
				Filename: file.Name(),
				SQL:      string(sqlContent),
			}

			migrations = append(migrations, migration)
		}
	}

	return migrations, nil
}

func extractVersionFromFilename(filename string) string {
	re := regexp.MustCompile(`^v?(\d+\.\d+(?:\.\d+)?)-.*\.sql$`)
	matches := re.FindStringSubmatch(filename)

	if len(matches) >= 2 {
		return matches[1]
	}

	return ""
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			version VARCHAR(20) PRIMARY KEY,
			filename VARCHAR(255) NOT NULL,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func getExecutedMigrations(db *sql.DB) ([]string, error) {
	query := "SELECT version FROM migrations ORDER BY version"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, nil
}

func isMigrationExecuted(executedMigrations []string, version string) bool {
	for _, executed := range executedMigrations {
		if executed == version {
			return true
		}
	}
	return false
}

func executeMigration(db *sql.DB, migration Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(migration.SQL); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute SQL: %v", err)
	}

	insertQuery := "INSERT INTO migrations (version, filename) VALUES ($1, $2)"
	if _, err := tx.Exec(insertQuery, migration.Version, migration.Filename); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to mark migration as executed: %v", err)
	}

	return tx.Commit()
}

func GetMigrationStatus(db *sql.DB) error {
	query := `
		SELECT version, filename, executed_at 
		FROM migrations 
		ORDER BY version
	`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("Version\t\tFilename\t\t\tExecuted At")
	fmt.Println("-------\t\t--------\t\t\t-----------")

	for rows.Next() {
		var version, filename, executedAt string
		if err := rows.Scan(&version, &filename, &executedAt); err != nil {
			return err
		}
		fmt.Printf("%s\t\t%s\t\t%s\n", version, filename, executedAt)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}
