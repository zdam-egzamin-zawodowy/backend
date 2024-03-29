package resolvers

import (
	"github.com/zdam-egzamin-zawodowy/backend/internal/auth"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/profession"
	"github.com/zdam-egzamin-zawodowy/backend/internal/qualification"
	"github.com/zdam-egzamin-zawodowy/backend/internal/question"
	"github.com/zdam-egzamin-zawodowy/backend/internal/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AuthUsecase          auth.Usecase
	UserUsecase          user.Usecase
	ProfessionUsecase    profession.Usecase
	QualificationUsecase qualification.Usecase
	QuestionUsecase      question.Usecase
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type professionResolver struct{ *Resolver }
type questionResolver struct{ *Resolver }

func (r *Resolver) Mutation() generated.MutationResolver     { return &mutationResolver{r} }
func (r *Resolver) Query() generated.QueryResolver           { return &queryResolver{r} }
func (r *Resolver) Profession() generated.ProfessionResolver { return &professionResolver{r} }
func (r *Resolver) Question() generated.QuestionResolver     { return &questionResolver{r} }
