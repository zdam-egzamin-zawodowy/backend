package repository

import (
	"path/filepath"

	"github.com/99designs/gqlgen/graphql"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/filestorage"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/utils"
)

type repository struct {
	fileStorage filestorage.FileStorage
}

func (repo *repository) saveImages(destination *models.Question, input *models.QuestionInput) {
	images := [...]*graphql.Upload{
		input.Image,
		input.AnswerAImage,
		input.AnswerBImage,
		input.AnswerCImage,
		input.AnswerDImage,
	}
	filenames := [...]string{
		destination.Image,
		destination.AnswerAImage,
		destination.AnswerBImage,
		destination.AnswerCImage,
		destination.AnswerDImage,
	}

	for index, file := range images {
		if file != nil {
			generated := false
			if filenames[index] == "" {
				generated = true
				filenames[index] = utils.GenerateFilename(filepath.Ext(file.Filename))
			}
			err := repo.fileStorage.Put(file.File, filenames[index])
			if err != nil && generated {
				filenames[index] = ""
			}
		}
	}

	destination.Image = filenames[0]
	destination.AnswerAImage = filenames[1]
	destination.AnswerBImage = filenames[2]
	destination.AnswerCImage = filenames[3]
	destination.AnswerDImage = filenames[4]
}

func (repo *repository) deleteImages(filenames []string) {
	for _, filename := range filenames {
		repo.fileStorage.Remove(filename)
	}
}

func (repo *repository) deleteImagesBasedOnInput(question *models.Question, input *models.QuestionInput) {
	filenames := []string{}

	if input.DeleteImage != nil &&
		*input.DeleteImage &&
		input.Image == nil &&
		question.Image != "" {
		filenames = append(filenames, question.Image)
		question.Image = ""
	}

	if input.DeleteAnswerAImage != nil &&
		*input.DeleteAnswerAImage &&
		input.AnswerAImage == nil &&
		question.AnswerAImage != "" {
		filenames = append(filenames, question.AnswerAImage)
		question.AnswerAImage = ""
	}

	if input.DeleteAnswerBImage != nil &&
		*input.DeleteAnswerBImage &&
		input.AnswerBImage == nil &&
		question.AnswerBImage != "" {
		filenames = append(filenames, question.AnswerBImage)
		question.AnswerBImage = ""
	}

	if input.DeleteAnswerCImage != nil &&
		*input.DeleteAnswerCImage &&
		input.AnswerCImage == nil &&
		question.AnswerCImage != "" {
		filenames = append(filenames, question.AnswerCImage)
		question.AnswerCImage = ""
	}

	if input.DeleteAnswerDImage != nil &&
		*input.DeleteAnswerDImage &&
		input.AnswerDImage == nil &&
		question.AnswerDImage != "" {
		filenames = append(filenames, question.AnswerDImage)
		question.AnswerDImage = ""
	}

	repo.deleteImages(filenames)
}

func (repo *repository) getAllFilenamesAndDeleteImages(questions []*models.Question) {
	filenames := []string{}

	for _, question := range questions {
		if question.Image != "" {
			filenames = append(filenames, question.Image)
		}
		if question.AnswerAImage != "" {
			filenames = append(filenames, question.AnswerAImage)
		}
		if question.AnswerBImage != "" {
			filenames = append(filenames, question.AnswerBImage)
		}
		if question.AnswerCImage != "" {
			filenames = append(filenames, question.AnswerCImage)
		}
		if question.AnswerDImage != "" {
			filenames = append(filenames, question.AnswerDImage)
		}
	}

	repo.deleteImages(filenames)
}
