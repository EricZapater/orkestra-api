package searches

import (
	"context"
)

type SearchService interface {
	GetByText(ctx context.Context, request *SearchRequest)([]SearchResult, error)
}

type searchService struct{
	searchRepository SearchRepository	
}

func NewSearchService(searchRepository SearchRepository) SearchService {
	return &searchService{
		searchRepository: searchRepository,		
	}
}

func (s *searchService) GetByText(ctx context.Context, request *SearchRequest)([]SearchResult, error){	
	return s.searchRepository.GetByText(ctx, request)
}