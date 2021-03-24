package querycomplexity

import (
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/internal/profession"
	"github.com/zdam-egzamin-zawodowy/backend/internal/qualification"
	"github.com/zdam-egzamin-zawodowy/backend/internal/question"
	"github.com/zdam-egzamin-zawodowy/backend/internal/user"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/utils"
)

const (
	complexityLimit                    = 10000
	countComplexity                    = 1000
	professionsTotalFieldComplexity    = 100
	qualificationsTotalFieldComplexity = 100
	questionsTotalFieldComplexity      = 300
	usersTotalFieldComplexity          = 50
)

func GetComplexityLimitExtension() *extension.ComplexityLimit {
	return extension.FixedComplexityLimit(complexityLimit)
}

func GetComplexityRoot() generated.ComplexityRoot {
	complexityRoot := generated.ComplexityRoot{}

	complexityRoot.Profession.Qualifications = func(childComplexity int) int {
		return 10 + childComplexity
	}
	complexityRoot.ProfessionList.Total = getCountComplexity
	complexityRoot.Query.Professions = func(
		childComplexity int,
		filter *models.ProfessionFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, profession.FetchDefaultLimit),
			professionsTotalFieldComplexity,
			1,
		)
	}

	complexityRoot.QualificationList.Total = getCountComplexity
	complexityRoot.Query.Qualifications = func(
		childComplexity int,
		filter *models.QualificationFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, qualification.FetchDefaultLimit),
			qualificationsTotalFieldComplexity,
			1,
		)
	}

	complexityRoot.QuestionList.Total = getCountComplexity
	complexityRoot.Query.Questions = func(
		childComplexity int,
		filter *models.QuestionFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, question.FetchDefaultLimit),
			questionsTotalFieldComplexity,
			1,
		)
	}
	complexityRoot.Query.GenerateTest = func(childComplexity int, qualificationIDs []int, limit *int) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, question.TestMaxLimit),
			0,
			3,
		)
	}

	complexityRoot.UserList.Total = getCountComplexity
	complexityRoot.Query.Users = func(
		childComplexity int,
		filter *models.UserFilter,
		limit *int,
		offset *int,
		sort []string,
	) int {
		return computeComplexity(
			childComplexity,
			utils.SafeIntPointer(limit, user.FetchMaxLimit),
			usersTotalFieldComplexity,
			1,
		)
	}

	complexityRoot.Mutation.CreateProfession = func(childComplexity int, input models.ProfessionInput) int {
		return (complexityLimit / 5) + childComplexity
	}

	complexityRoot.Mutation.CreateQualification = func(
		childComplexity int,
		input models.QualificationInput,
	) int {
		return (complexityLimit / 5) + childComplexity
	}

	complexityRoot.Mutation.CreateQuestion = func(childComplexity int, input models.QuestionInput) int {
		return (complexityLimit / 4) + childComplexity
	}

	complexityRoot.Mutation.CreateUser = func(childComplexity int, input models.UserInput) int {
		return (complexityLimit / 5) + childComplexity
	}

	complexityRoot.Mutation.SignIn = func(
		childComplexity int,
		email string,
		password string,
		staySignedIn *bool,
	) int {
		return (complexityLimit / 2) + childComplexity
	}

	complexityRoot.Mutation.UpdateManyUsers = func(
		childComplexity int,
		ids []int,
		input models.UserInput,
	) int {
		return (complexityLimit / 5) + childComplexity
	}

	complexityRoot.Mutation.UpdateProfession = func(
		childComplexity int,
		id int,
		input models.ProfessionInput,
	) int {
		return (complexityLimit / 5) + childComplexity
	}

	complexityRoot.Mutation.UpdateQualification = func(
		childComplexity int,
		id int,
		input models.QualificationInput,
	) int {
		return (complexityLimit / 5) + childComplexity
	}

	complexityRoot.Mutation.UpdateQuestion = func(
		childComplexity int,
		id int,
		input models.QuestionInput,
	) int {
		return (complexityLimit / 4) + childComplexity
	}

	complexityRoot.Mutation.UpdateUser = func(childComplexity int, id int, input models.UserInput) int {
		return (complexityLimit / 5) + childComplexity
	}

	complexityRoot.Mutation.DeleteProfessions = func(childComplexity int, ids []int) int {
		return (complexityLimit / 5) + childComplexity
	}

	complexityRoot.Mutation.DeleteQualifications = func(childComplexity int, ids []int) int {
		return (complexityLimit / 5) + childComplexity
	}

	complexityRoot.Mutation.DeleteQuestions = func(childComplexity int, ids []int) int {
		return (complexityLimit / 4) + childComplexity
	}

	complexityRoot.Mutation.DeleteUsers = func(childComplexity int, ids []int) int {
		return (complexityLimit / 5) + childComplexity
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
