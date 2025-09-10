package llm

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"orkestra-api/config"
	"strings"
	"time"

	"google.golang.org/genai"
)

type Service struct {
	db     *sql.DB
	apiKey string
	apiURL string
	cfg    config.Config
}

func NewService(db *sql.DB, cfg config.Config) *Service {
	return &Service{
		db:     db,
		apiKey: cfg.LLMApiKey,
		apiURL: cfg.LLMUrl,
		cfg:	cfg,
	}
}

func (s *Service) ProcessQuery(ctx context.Context, req QueryRequest) (*QueryResponse, error) {
	// Primer, detectem si la pregunta és estratègica
	/*isStrategic, err := s.isStrategicQuestion(ctx, req.Question)
	if err != nil {
		return nil, fmt.Errorf("error detecting question type: %w", err)
	}

	// Obtenim dades rellevants de la base
	data, err := s.getRelevantDataForQuestion(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving relevant data: %w", err)
	}

	if isStrategic {
		// Genera resposta raonada basada en dades
		answer, err := s.generateStrategicAnswer(ctx, req.Question, data)
		if err != nil {
			return nil, fmt.Errorf("error generating strategic answer: %w", err)
		}
		return &QueryResponse{
			Answer:    answer,
			Data:      data,
			Timestamp: time.Now(),
		}, nil
	}*/

	// Si no és estratègica, fem el flux clàssic SQL
	schema, err := s.getDatabaseSchema()
	if err != nil {
		return nil, fmt.Errorf("error getting database schema: %w", err)
	}

	sqlQuery, err := s.generateSQLQuery(ctx, req.Question, schema)
	if err != nil {
		return nil, fmt.Errorf("error generating SQL query: %w", err)
	}

	sqlData, err := s.executeQuery(ctx, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	answer, err := s.generateAnswer(ctx, req.Question, sqlQuery, sqlData)
	if err != nil {
		return nil, fmt.Errorf("error generating answer: %w", err)
	}

	return &QueryResponse{
		Answer:    answer,
		Query:     sqlQuery,
		Data:      sqlData,
		Timestamp: time.Now(),
	}, nil
}

// --------------------------------------------------
// Detecció tipus de pregunta
// --------------------------------------------------
func (s *Service) isStrategicQuestion(ctx context.Context, question string) (bool, error) {
	prompt := fmt.Sprintf(`Classifica la següent pregunta com:
1) SQL
2) Estratègica/optimització de projectes/recursos

Pregunta: "%s"

Respón només amb "SQL" o "ESTRATÈGICA".`, question)

	resp, err := s.callLLM(ctx, prompt)
	if err != nil {
		return false, err
	}
	resp = strings.ToUpper(strings.TrimSpace(resp))
	return resp == "ESTRATÈGICA", nil
}

// --------------------------------------------------
// Recuperació de dades rellevants per preguntes estratègiques
// --------------------------------------------------
func (s *Service) getRelevantDataForQuestion(ctx context.Context) ([]map[string]interface{}, error) {
	// Exemples de dades que podríem voler recuperar:
	// - operaris i costos
	// - projectes i dates
	// - assignacions actuals
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, surname, cost, color FROM operators;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	cols, _ := rows.Columns()
	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		row := make(map[string]interface{})
		for i, col := range cols {
			if b, ok := values[i].([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = values[i]
			}
		}
		results = append(results, row)
	}

	return results, nil
}

// --------------------------------------------------
// Generació de resposta estratègica
// --------------------------------------------------
func (s *Service) generateStrategicAnswer(ctx context.Context, question string, data []map[string]interface{}) (string, error) {
	dataJSON, _ := json.MarshalIndent(data, "", "  ")
	prompt := fmt.Sprintf(`Ets un assistent que ajuda a interpretar dades i prendre decisions d'assignació de projectes i operaris.

Pregunta: %s

Dades disponibles:
%s

Proposa una assignació òptima, ordre d'execució i explicació raonada. Si hi ha limitacions, explica-les de manera clara.`, question, string(dataJSON))

	return s.callLLM(ctx, prompt)
}

func (s *Service) getDatabaseSchema() (*QueryContext, error) {
	return &QueryContext{
		Tables: []TableInfo{
			{
				Name:        "operators",
				Description: "Operaris de l'empresa amb els seus costos",
				Columns: []ColumnInfo{
					{Name: "id", Type: "uuid", Description: "Identificador únic de l'operari"},
					{Name: "name", Type: "string", Description: "Nom de l'operari"},
					{Name: "surname", Type: "string", Description: "Cognoms de l'operari"},
					{Name: "cost", Type: "decimal", Description: "Cost per hora de l'operari"},
					{Name: "color", Type: "string", Description: "Color assignat a l'operari"},
				},
			},
			{
				Name:        "projects",
				Description: "Projectes de l'empresa amb informació financera",
				Columns: []ColumnInfo{
					{Name: "id", Type: "uuid", Description: "Identificador únic del projecte"},
					{Name: "description", Type: "string", Description: "Descripció del projecte"},
					{Name: "start_date", Type: "timestamp", Description: "Data d'inici del projecte"},
					{Name: "end_date", Type: "timestamp", Description: "Data de finalització del projecte"},
					{Name: "customer_id", Type: "uuid", Description: "Identificador del client"},
					{Name: "amount", Type: "decimal", Description: "Import total del projecte"},
					{Name: "estimated_cost", Type: "decimal", Description: "Cost estimat del projecte"},
				},
			},
			{
				Name:        "operators_to_projects",
				Description: "Assignacions d'operaris a projectes amb dates i dedicació",
				Columns: []ColumnInfo{
					{Name: "id", Type: "uuid", Description: "Identificador únic de l'assignació"},
					{Name: "operator_id", Type: "uuid", Description: "Identificador de l'operari"},
					{Name: "project_id", Type: "uuid", Description: "Identificador del projecte"},
					{Name: "cost", Type: "decimal", Description: "Cost de l'operari en aquest projecte"},
					{Name: "dedication_percent", Type: "decimal", Description: "Percentatge de dedicació"},
					{Name: "start_date", Type: "timestamp", Description: "Data d'inici de l'assignació"},
					{Name: "end_date", Type: "timestamp", Description: "Data de finalització de l'assignació"},
				},
			},
			{
				Name:        "tasks",
				Description: "Tasques assignades als projectes",
				Columns: []ColumnInfo{
					{Name: "id", Type: "uuid", Description: "Identificador únic de la tasca"},
					{Name: "description", Type: "string", Description: "Descripció de la tasca"},
					{Name: "user_id", Type: "uuid", Description: "Identificador de l'usuari assignat"},
					{Name: "project_id", Type: "uuid", Description: "Identificador del projecte"},
					{Name: "status", Type: "string", Description: "Estat de la tasca"},
					{Name: "start_date", Type: "timestamp", Description: "Data d'inici de la tasca"},
					{Name: "end_date", Type: "timestamp", Description: "Data de finalització de la tasca"},
				},
			},
			{
				Name:        "customers",
				Description: "Clients de l'empresa",
				Columns: []ColumnInfo{
					{Name: "id", Type: "uuid", Description: "Identificador únic del client"},
					{Name: "comercial_name", Type: "string", Description: "Nom comercial del client"},
					{Name: "vat_number", Type: "string", Description: "NIF/CIF del client"},
					{Name: "phone_number", Type: "string", Description: "Telèfon del client"},
				},
			},
		},
		Relationships: []Relationship{
			{FromTable: "projects", FromColumn: "customer_id", ToTable: "customers", ToColumn: "id", Type: "many_to_one"},
			{FromTable: "operator_to_projects", FromColumn: "operator_id", ToTable: "operators", ToColumn: "id", Type: "many_to_one"},
			{FromTable: "operator_to_projects", FromColumn: "project_id", ToTable: "projects", ToColumn: "id", Type: "many_to_one"},
			{FromTable: "tasks", FromColumn: "project_id", ToTable: "projects", ToColumn: "id", Type: "many_to_one"},
		},
	}, nil
}

func (s *Service) generateSQLQuery(ctx context.Context, question string, schema *QueryContext) (string, error) {
	prompt := s.buildSQLPrompt(question, schema)
	
	response, err := s.callLLM(ctx, prompt)
	if err != nil {
		return "", err
	}

	// Extract SQL from response
	sqlQuery := s.extractSQL(response)
	return sqlQuery, nil
}

func (s *Service) executeQuery(ctx context.Context, query string) ([]map[string]interface{}, error) {
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return results, nil
}

func (s *Service) generateAnswer(ctx context.Context, question string, query string, data []map[string]interface{}) (string, error) {
	prompt := s.buildAnswerPrompt(question, query, data)
	return s.callLLM(ctx, prompt)
}

func (s *Service) buildSQLPrompt(question string, schema *QueryContext) string {
	var sb strings.Builder
	
	sb.WriteString("Ets un expert en SQL i bases de dades. Genera una consulta SQL basada en la pregunta de l'usuari.\n\n")
	sb.WriteString("ESQUEMA DE LA BASE DE DADES:\n")
	
	for _, table := range schema.Tables {
		sb.WriteString(fmt.Sprintf("\nTaula: %s - %s\n", table.Name, table.Description))
		for _, col := range table.Columns {
			sb.WriteString(fmt.Sprintf("  - %s (%s): %s\n", col.Name, col.Type, col.Description))
		}
	}
	
	sb.WriteString("\nRELACIONS:\n")
	for _, rel := range schema.Relationships {
		sb.WriteString(fmt.Sprintf("- %s.%s -> %s.%s (%s)\n", rel.FromTable, rel.FromColumn, rel.ToTable, rel.ToColumn, rel.Type))
	}
	
	sb.WriteString(fmt.Sprintf("\nPREGUNTA: %s\n\n", question))
	sb.WriteString("Genera NOMÉS la consulta SQL, sense explicacions addicionals. La consulta ha de ser compatible amb PostgreSQL.\n")
	sb.WriteString("Si la pregunta fa referència a 'avui', utilitza CURRENT_DATE.\n")
	sb.WriteString("Si la pregunta demana informació sobre marges de benefici, calcula: amount - estimated_cost.\n")
	
	return sb.String()
}

func (s *Service) buildAnswerPrompt(question string, query string, data []map[string]interface{}) string {
	dataJSON, _ := json.MarshalIndent(data, "", "  ")
	
	return fmt.Sprintf(`Ets un assistent que ajuda a interpretar resultats de consultes de base de dades.

PREGUNTA ORIGINAL: %s

CONSULTA SQL EXECUTADA: %s

RESULTATS:
%s

Proporciona una resposta clara i concisa en català que respongui a la pregunta original basant-te en els resultats obtinguts. Si no hi ha resultats, explica-ho de manera amigable.`, 
		question, query, string(dataJSON))
}

func (s *Service) callLLM(ctx context.Context, prompt string) (string, error) {
    switch s.cfg.LLMProvider {
	case "apifreellm":
		return s.callfreeApi(ctx, prompt)
	case "gemini":
		return s.callGemini(ctx, prompt)
	default:
		return s.callOpenAi(ctx, prompt)
	}
    
}

func (s *Service) callOpenAi(ctx context.Context, prompt string)(string, error){
	client := &http.Client{Timeout: 30 * time.Second}
	body, _ := json.Marshal(map[string]interface{}{
        "model": "gpt-3.5-turbo",
        "messages": []map[string]string{
            {"role": "user", "content": prompt},
        },
        "max_tokens":  1000,
        "temperature": 0.1,
    })
    req, _ := http.NewRequestWithContext(ctx, "POST", s.apiURL, bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+s.apiKey)

    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    raw, _ := io.ReadAll(resp.Body)
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("OpenAI API error: %s", string(raw))
    }

    var response struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        } `json:"choices"`
    }
    if err := json.Unmarshal(raw, &response); err != nil {
        return "", err
    }
    if len(response.Choices) == 0 {
        return "", fmt.Errorf("no response from LLM")
    }

    return response.Choices[0].Message.Content, nil
}
func (s *Service) callfreeApi(ctx context.Context, prompt string)(string, error){
	client := &http.Client{Timeout: 30 * time.Second}
	body, _ := json.Marshal(map[string]string{
            "message": prompt,
        })
        req, _ := http.NewRequestWithContext(ctx, "POST", s.apiURL, bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")

        resp, err := client.Do(req)
        if err != nil {
            return "", err
        }
        defer resp.Body.Close()

        raw, _ := io.ReadAll(resp.Body)
        if resp.StatusCode != http.StatusOK {
            return "", fmt.Errorf("apifreellm error: %s", string(raw))
        }

        var result map[string]interface{}
        if err := json.Unmarshal(raw, &result); err != nil {
            return "", err
        }
        if status, ok := result["status"].(string); ok && status == "success" {
            return fmt.Sprintf("%v", result["response"]), nil
        }
        return "", fmt.Errorf("error: %v", result["error"])
}

func (s *Service) callGemini(ctx context.Context, prompt string)(string, error){
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: s.apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %w", err)
	}

	response, err := client.Models.GenerateContent(
		ctx,
        "gemini-2.5-flash-lite",
        genai.Text(prompt),
        nil,
	)
	if err  != nil {
		return "", fmt.Errorf("Gemini API error: %w", err)
	}
	return response.Text(), nil	
}


func (s *Service) extractSQL(response string) string {
	// Look for SQL between ```sql and ``` or just return the response if it looks like SQL
	if strings.Contains(response, "```sql") {
		start := strings.Index(response, "```sql") + 6
		end := strings.Index(response[start:], "```")
		if end != -1 {
			return strings.TrimSpace(response[start : start+end])
		}
	}
	
	// If no code blocks, assume the entire response is SQL
	return strings.TrimSpace(response)
}
