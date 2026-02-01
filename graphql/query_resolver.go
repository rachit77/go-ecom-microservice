package main

import "context"

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pg *PaginationInput, id *string) ([]*Account, error) {
	return nil, nil
}

func (r *queryResolver) Products(ctx context.Context, pg *PaginationInput, query *string, id *string) ([]*Product, error) {
	return nil, nil
}
