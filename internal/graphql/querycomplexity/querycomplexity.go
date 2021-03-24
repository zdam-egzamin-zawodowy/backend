package querycomplexity

import (
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

const (
	countComplexity = 1000
)

func GetComplexityRoot() generated.ComplexityRoot {
	complexityRoot := generated.ComplexityRoot{}
	complexityRoot.Query.GenerateTest = func(childComplexity int, qualificationIDs []int, limit *int) int {
		return 300 + childComplexity
	}
	complexityRoot.Query.Professions = func(
		childComplexity int,
		filter *models.ProfessionFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return 200 + childComplexity
	}
	complexityRoot.Query.Qualifications = func(
		childComplexity int,
		filter *models.QualificationFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return 200 + childComplexity
	}
	complexityRoot.Query.Questions = func(
		childComplexity int,
		filter *models.QuestionFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return 200 + childComplexity
	}
	complexityRoot.Query.Users = func(
		childComplexity int,
		filter *models.UserFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return 200 + childComplexity
	}
	complexityRoot.Mutation.CreateProfession = func(childComplexity int, input models.ProfessionInput) int {
		return 200 + childComplexity
	}
	complexityRoot.Mutation.CreateQualification = func(
		childComplexity int,
		input models.QualificationInput,
	) int {
		return 200 + childComplexity
	}
	complexityRoot.Mutation.CreateQuestion = func(childComplexity int, input models.QuestionInput) int {
		return 400 + childComplexity
	}
	complexityRoot.Mutation.CreateUser = func(childComplexity int, input models.UserInput) int {
		return 200 + childComplexity
	}
	complexityRoot.Mutation.SignIn = func(
		childComplexity int,
		email string,
		password string,
		staySignedIn *bool,
	) int {
		return 400 + childComplexity
	}
	complexityRoot.Mutation.DeleteProfessions = func(childComplexity int, ids []int) int {
		return 200 + childComplexity
	}
	complexityRoot.Mutation.DeleteQualifications = func(childComplexity int, ids []int) int {
		return 200 + childComplexity
	}
	complexityRoot.Mutation.DeleteQuestions = func(childComplexity int, ids []int) int {
		return 400 + childComplexity
	}
	complexityRoot.Mutation.DeleteUsers = func(childComplexity int, ids []int) int {
		return 200 + childComplexity
	}
	complexityRoot.Mutation.UpdateManyUsers = func(
		childComplexity int,
		ids []int,
		input models.UserInput,
	) int {
		return 200 + childComplexity
	}
	complexityRoot.Mutation.UpdateProfession = func(
		childComplexity int,
		id int,
		input models.ProfessionInput,
	) int {
		return 200 + childComplexity
	}
	complexityRoot.Mutation.UpdateQualification = func(
		childComplexity int,
		id int,
		input models.QualificationInput,
	) int {
		return 200 + childComplexity
	}
	complexityRoot.Mutation.UpdateQuestion = func(
		childComplexity int,
		id int,
		input models.QuestionInput,
	) int {
		return 400 + childComplexity
	}
	complexityRoot.Mutation.UpdateUser = func(childComplexity int, id int, input models.UserInput) int {
		return 200 + childComplexity
	}
	complexityRoot.Profession.Qualifications = func(childComplexity int) int {
		return 50 + childComplexity
	}
	return complexityRoot
}

func computeComplexity(childComplexity, limit, totalFieldComplexity, multiplyBy int) int {
	complexity := 0
	if childComplexity >= countComplexity {
		childComplexity -= countComplexity
		complexity += totalFieldComplexity
	}
	return limit*childComplexity*multiplyBy + complexity
}

func getCountComplexity(childComplexity int) int {
	return countComplexity + childComplexity
}
