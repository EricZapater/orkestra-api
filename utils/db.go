package utils

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Utils struct {
	db *sql.DB
}

func NewUtils(db *sql.DB) *Utils {
	return &Utils{
		db:db,
	}
}

func(u *Utils)CheckDatabase(dbName string)(bool, error){
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1);`
	err := u.db.QueryRow(query, dbName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *Utils) CreateDatabase(dbName string) error{ 	
	_, err := u.db.Exec(fmt.Sprintf(`CREATE DATABASE "%s";`, dbName))
	if err != nil {
		return fmt.Errorf("no s'ha pogut crear la base de dades: %w", err)
	}

	fmt.Printf("✅ Base de dades %s creada correctament\n", dbName)
	return nil
}

func(u *Utils) RunMigrations(dbName string, migrationsDir string) error{
	_, err := u.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		filename TEXT PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT now()
	);`)
	if err != nil {
		return fmt.Errorf("error creant taula de migracions: %w", err)
	}

	// Llegeix els fitxers
	absPath, err := filepath.Abs(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("no puc obtenir path absolut: %w", err)
	}
	fmt.Printf("🔍 Buscant fitxers a: %s\n", absPath)
	
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("no es poden llegir fitxers de migració: %w", err)
	}
	fmt.Printf("📄 Fitxers trobats: %v\n", files) 
	sort.Strings(files) // Per ordre

	for _, file := range files {
		filename := filepath.Base(file)
		fmt.Printf("%s\n", filename)

		var exists bool
		err := u.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename = $1)`, filename).Scan(&exists)
		if err != nil {
			return fmt.Errorf("error consultant schema_migrations: %w", err)
		}

		if exists {
			fmt.Printf("⏩ %s ja aplicat, saltant...\n", filename)
			continue
		}
		
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("no puc llegir el fitxer %s: %w", filename, err)
		}		
		sqlContent := string(sqlBytes)
		if strings.TrimSpace(sqlContent) == "" {
			fmt.Printf("⚠️  %s està buit, saltant...\n", filename)
			continue
		}

		fmt.Printf("▶️  Aplicant %s...\n", filename)
		queries := strings.Split(string(sqlContent), ";")
		for _, q := range queries {
			q = strings.TrimSpace(q)
			if q == "" {
				continue
			}
			_, err := u.db.Exec(q)
			if err != nil {
				fmt.Printf("❌ Error a '%s': %v\n", q, err)
				return err
			}
		}

		_, err = u.db.Exec(`INSERT INTO schema_migrations (filename, applied_at) VALUES ($1, $2)`, filename, time.Now())
		if err != nil {
			return fmt.Errorf("error enregistrant %s: %w", filename, err)
		}
	}

	fmt.Println("✅ Migracions aplicades correctament.")
	return nil
}