package repository

import (
	"path/filepath"

	"github.com/zdam-egzamin-zawodowy/backend/pkg/fstorage/fstorageutil"

	"github.com/99designs/gqlgen/graphql"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/fstorage"
)

type repository struct {
	fileStorage fstorage.FileStorage
}

func (repo *repository) saveImages(destination *model.Question, input *model.QuestionInput) {
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
				filenames[index] = fstorageutil.GenerateFilename(filepath.Ext(file.Filename))
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

func (repo *repository) deleteImagesBasedOnInput(question *model.Question, input *model.QuestionInput) {
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

func (repo *repository) getAllImagesAndDelete(questions []*model.Question) {
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
