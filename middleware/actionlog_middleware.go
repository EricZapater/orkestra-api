package middleware

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ActionLogMiddleware struct {
	db *sql.DB
}

func NewActionLogMiddleware(db *sql.DB) *ActionLogMiddleware {
	return &ActionLogMiddleware{db: db}
}

func (alm *ActionLogMiddleware) LogAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
        userIDRaw, exists := claims["id"]
		if !exists {
			userIDRaw = nil
		}
		actionType := c.Request.Method + " " + c.Request.URL.Path		
		timezone := "" // Aquí després pots millorar-ho si vols agafar el timezone del client
		performedAt := time.Now()

		// Inicialitzar metadata
		metadata := "{}"

		// Llegir i reconstruir body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err == nil && len(bodyBytes) > 0 {
			// Tornar a posar el body perquè altres handlers puguin llegir-lo
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			
			if c.ContentType() == "application/json" {
				
				var jsonBody map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {					
					delete(jsonBody, "password")					
					modifiedBodyBytes, err := json.Marshal(jsonBody)
					if err == nil {
						metadata = string(modifiedBodyBytes)
					} else {
						metadata = `{"note": "error marshalling JSON"}`
					}
				} else {
					metadata = `{"note": "invalid JSON body"}`
				}
			} else {
				metadata = `{"note": "non-JSON content type"}`
			}
		} else {
			// Si error llegint body, o està buit
			c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte{}))
		}

		// Transformar user_id
		userUUID := uuid.Nil
		if userIDRaw != nil {
			parsedUUID, err := uuid.Parse(userIDRaw.(string))
			if err == nil {
				userUUID = parsedUUID
			}
		}

		// Guardar log
		if err := alm.SaveActionLog(userUUID, actionType, metadata, timezone, performedAt); err != nil {
			// Potser vols fer un log aquí
			log.Printf("Error saving action log: %v", err)
		}

		c.Next()
	}
}


func (alm *ActionLogMiddleware) SaveActionLog(userID uuid.UUID, actionType, metadata, timezone string, performedAt time.Time) error {
		query := `INSERT INTO action_logs (user_id, action_type, metadata, timezone, performed_at) 
				VALUES ($1, $2, $3::jsonb, $4, $5)`
		_, err := alm.db.Exec(query, userID, actionType, metadata, timezone, performedAt)
		return err
	}
