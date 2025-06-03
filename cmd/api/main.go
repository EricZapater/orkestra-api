// main.go
package main

import (
	"log"
	"orkestra-api/config"
	"orkestra-api/server"
)

func main() {
    // Carregar configuració
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }
    
    // Connectar a la base de dades
    db, err := config.ConnectDB(cfg)
    if err != nil {
        log.Fatalf("failed to connect to db: %v", err)
    }
    defer db.Close()
    
    // Inicialitzar el servidor
    srv := server.NewServer(cfg, db)
    
    // Configurar middlewares i rutes
    if err := srv.Setup(); err != nil {
        log.Fatalf("failed to set up middlewares: %v", err)
    }
    
    // Iniciar el servidor
    if err := srv.Run(); err != nil {
        log.Fatalf("failed to start server: %v", err)
    }
}