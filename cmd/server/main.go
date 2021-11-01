package main

import (
	"context"
	"github.com/Kichiyaki/appmode"
	"github.com/Kichiyaki/chilogrus"
	"github.com/Kichiyaki/goutil/envutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-pg/pg/v10"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/zdam-egzamin-zawodowy/backend/internal/auth"
	"github.com/zdam-egzamin-zawodowy/backend/internal/chi/middleware"
	graphqlhttpdelivery "github.com/zdam-egzamin-zawodowy/backend/internal/graphql/delivery/httpdelivery"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/directive"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/resolvers"
	"github.com/zdam-egzamin-zawodowy/backend/internal/profession"
	"github.com/zdam-egzamin-zawodowy/backend/internal/qualification"
	"github.com/zdam-egzamin-zawodowy/backend/internal/question"
	"github.com/zdam-egzamin-zawodowy/backend/internal/user"

	"github.com/pkg/errors"

	"github.com/zdam-egzamin-zawodowy/backend/internal/auth/jwt"
	authusecase "github.com/zdam-egzamin-zawodowy/backend/internal/auth/usecase"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/dataloader"
	"github.com/zdam-egzamin-zawodowy/backend/internal/postgres"
	professionrepository "github.com/zdam-egzamin-zawodowy/backend/internal/profession/repository"
	professionusecase "github.com/zdam-egzamin-zawodowy/backend/internal/profession/usecase"
	qualificationrepository "github.com/zdam-egzamin-zawodowy/backend/internal/qualification/repository"
	qualificationusecase "github.com/zdam-egzamin-zawodowy/backend/internal/qualification/usecase"
	questionrepository "github.com/zdam-egzamin-zawodowy/backend/internal/question/repository"
	questionusecase "github.com/zdam-egzamin-zawodowy/backend/internal/question/usecase"
	userrepository "github.com/zdam-egzamin-zawodowy/backend/internal/user/repository"
	userusecase "github.com/zdam-egzamin-zawodowy/backend/internal/user/usecase"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/zdam-egzamin-zawodowy/backend/fstorage"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func init() {
	os.Setenv("TZ", "UTC")

	if appmode.Equals(appmode.DevelopmentMode) {
		godotenv.Load(".env.local")
	}

	prepareLogger()
}

func main() {
	fileStorage := fstorage.New(&fstorage.Config{
		BasePath: envutil.GetenvString("FILE_STORAGE_PATH"),
	})

	dbConn, err := postgres.Connect(&postgres.Config{
		LogQueries: envutil.GetenvBool("LOG_DB_QUERIES"),
	})
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "Couldn't connect to the db"))
	}

	repos, err := prepareRepositories(dbConn, fileStorage)
	if err != nil {
		logrus.Fatal(err)
	}

	ucases, err := prepareUsecases(repos)
	if err != nil {
		logrus.Fatal(err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: prepareRouter(repos, ucases),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalln("listen:", err)
		}
	}()
	logrus.Info("Server is listening on the port 8080")

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalln("Server Shutdown:", err)
	}
	logrus.Info("Server exiting")
}

func prepareLogger() {
	if appmode.Equals(appmode.DevelopmentMode) {
		logrus.SetLevel(logrus.DebugLevel)
	}

	timestampFormat := "2006-01-02 15:04:05"
	if appmode.Equals(appmode.ProductionMode) {
		customFormatter := new(logrus.JSONFormatter)
		customFormatter.TimestampFormat = timestampFormat
		logrus.SetFormatter(customFormatter)
	} else {
		customFormatter := new(logrus.TextFormatter)
		customFormatter.TimestampFormat = timestampFormat
		customFormatter.FullTimestamp = true
		logrus.SetFormatter(customFormatter)
	}
}

type repositories struct {
	userRepository          user.Repository
	professionRepository    profession.Repository
	qualificationRepository qualification.Repository
	questionRepository      question.Repository
}

func prepareRepositories(dbConn *pg.DB, fileStorage fstorage.FileStorage) (*repositories, error) {
	var err error
	repos := &repositories{}

	repos.userRepository, err = userrepository.NewPGRepository(&userrepository.PGRepositoryConfig{
		DB: dbConn,
	})
	if err != nil {
		return nil, errors.Wrap(err, "userRepository")
	}

	repos.professionRepository, err = professionrepository.NewPGRepository(&professionrepository.PGRepositoryConfig{
		DB: dbConn,
	})
	if err != nil {
		return nil, errors.Wrap(err, "professionRepository")
	}

	repos.qualificationRepository, err = qualificationrepository.NewPGRepository(&qualificationrepository.PGRepositoryConfig{
		DB: dbConn,
	})
	if err != nil {
		return nil, errors.Wrap(err, "qualificationRepository")
	}

	repos.questionRepository, err = questionrepository.NewPGRepository(&questionrepository.PGRepositoryConfig{
		DB:          dbConn,
		FileStorage: fileStorage,
	})
	if err != nil {
		return nil, errors.Wrap(err, "questionRepository")
	}

	return repos, nil
}

type usecases struct {
	authUsecase          auth.Usecase
	userUsecase          user.Usecase
	professionUsecase    profession.Usecase
	qualificationUsecase qualification.Usecase
	questionUsecase      question.Usecase
}

func prepareUsecases(repos *repositories) (*usecases, error) {
	var err error
	ucases := &usecases{}

	ucases.authUsecase, err = authusecase.New(&authusecase.Config{
		UserRepository: repos.userRepository,
		TokenGenerator: jwt.NewTokenGenerator(envutil.GetenvString("ACCESS_SECRET")),
	})
	if err != nil {
		return nil, errors.Wrap(err, "authUsecase")
	}

	ucases.userUsecase, err = userusecase.New(&userusecase.Config{
		UserRepository: repos.userRepository,
	})
	if err != nil {
		return nil, errors.Wrap(err, "userUsecase")
	}

	ucases.professionUsecase, err = professionusecase.New(&professionusecase.Config{
		ProfessionRepository: repos.professionRepository,
	})
	if err != nil {
		return nil, errors.Wrap(err, "professionUsecase")
	}

	ucases.qualificationUsecase, err = qualificationusecase.New(&qualificationusecase.Config{
		QualificationRepository: repos.qualificationRepository,
	})
	if err != nil {
		return nil, errors.Wrap(err, "qualificationUsecase")
	}

	ucases.questionUsecase, err = questionusecase.New(&questionusecase.Config{
		QuestionRepository: repos.questionRepository,
	})
	if err != nil {
		return nil, errors.Wrap(err, "questionUsecase")
	}

	return ucases, nil
}

func prepareRouter(repos *repositories, ucases *usecases) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimiddleware.RealIP)
	if envutil.GetenvBool("ENABLE_ACCESS_LOG") {
		r.Use(chilogrus.Logger(logrus.StandardLogger()))
	}
	r.Use(chimiddleware.Recoverer)

	if appmode.Equals(appmode.DevelopmentMode) {
		r.Use(cors.Handler(cors.Options{
			AllowOriginFunc: func(*http.Request, string) bool {
				return true
			},
			AllowCredentials: true,
			ExposedHeaders:   []string{"Authorization"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowedHeaders:   []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			MaxAge:           300,
		}))
	}

	r.Group(func(r chi.Router) {
		r.Use(
			middleware.DataLoaderToContext(dataloader.Config{
				ProfessionRepo:    repos.professionRepository,
				QualificationRepo: repos.qualificationRepository,
			}),
			middleware.Authenticate(ucases.authUsecase),
		)
		err := graphqlhttpdelivery.Attach(r, graphqlhttpdelivery.Config{
			Resolver: &resolvers.Resolver{
				AuthUsecase:          ucases.authUsecase,
				UserUsecase:          ucases.userUsecase,
				ProfessionUsecase:    ucases.professionUsecase,
				QualificationUsecase: ucases.qualificationUsecase,
				QuestionUsecase:      ucases.questionUsecase,
			},
			Directive: &directive.Directive{},
		})
		if err != nil {
			log.Fatalln(err)
		}
	})

	return r
}
