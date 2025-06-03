package searches

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type SearchRepository interface {
	GetByText(ctx context.Context, request *SearchRequest)([]SearchResult, error)
}

type searchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) SearchRepository {
	return &searchRepository{db: db}
}

func (r *searchRepository) GetByText(ctx context.Context, request *SearchRequest)([]SearchResult, error){
	query := `
		SELECT tipus, id, text, link
		FROM search_index
		WHERE group_id = ANY($1)
		  AND search_vector @@ to_tsquery($2 || ':*')
	`
	rows, err := r.db.QueryContext(ctx, query, pq.Array(request.GroupIDS), request.Text)
	if err != nil {
		return []SearchResult{}, err
	}
	defer rows.Close()

	results := []SearchResult{}
	for rows.Next() {
		var r SearchResult
		err := rows.Scan(&r.Tipus, &r.ID, &r.Text, &r.Link)
		if err != nil {
			return []SearchResult{}, err
		}
		results = append(results, r)
	}

	return results, nil
}