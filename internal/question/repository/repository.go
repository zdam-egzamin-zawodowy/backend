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

func (repo *repository) deleteImages(images []string) {
	for _, image := range images {
		repo.fileStorage.Remove(image)
	}
}

func (repo *repository) deleteImagesBasedOnInput(question *models.Question, input *models.QuestionInput) {
	images := []string{}

	if input.DeleteImage != nil &&
		*input.DeleteImage &&
		input.Image == nil &&
		question.Image != "" {
		images = append(images, question.Image)
		question.Image = ""
	}

	if input.DeleteAnswerAImage != nil &&
		*input.DeleteAnswerAImage &&
		input.AnswerAImage == nil &&
		question.AnswerAImage != "" {
		images = append(images, question.AnswerAImage)
		question.AnswerAImage = ""
	}

	if input.DeleteAnswerBImage != nil &&
		*input.DeleteAnswerBImage &&
		input.AnswerBImage == nil &&
		question.AnswerBImage != "" {
		images = append(images, question.AnswerBImage)
		question.AnswerBImage = ""
	}

	if input.DeleteAnswerCImage != nil &&
		*input.DeleteAnswerCImage &&
		input.AnswerCImage == nil &&
		question.AnswerCImage != "" {
		images = append(images, question.AnswerCImage)
		question.AnswerCImage = ""
	}

	if input.DeleteAnswerDImage != nil &&
		*input.DeleteAnswerDImage &&
		input.AnswerDImage == nil &&
		question.AnswerDImage != "" {
		images = append(images, question.AnswerDImage)
		question.AnswerDImage = ""
	}

	repo.deleteImages(images)
}

func (repo *repository) getAllImagesAndDelete(questions []*models.Question) {
	images := []string{}

	for _, question := range questions {
		if question.Image != "" {
			images = append(images, question.Image)
		}
		if question.AnswerAImage != "" {
			images = append(images, question.AnswerAImage)
		}
		if question.AnswerBImage != "" {
			images = append(images, question.AnswerBImage)
		}
		if question.AnswerCImage != "" {
			images = append(images, question.AnswerCImage)
		}
		if question.AnswerDImage != "" {
			images = append(images, question.AnswerDImage)
		}
	}

	repo.deleteImages(images)
}
