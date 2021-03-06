package resolvers

import "github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type questionResolver struct{ *Resolver }

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() generated.QueryResolver       { return &queryResolver{r} }
func (r *Resolver) Question() generated.QuestionResolver { return &questionResolver{r} }
