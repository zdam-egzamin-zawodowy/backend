package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	graphqlhttpdelivery "github.com/zdam-egzamin-zawodowy/backend/internal/graphql/delivery/http"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/directive"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/resolvers"

	"github.com/pkg/errors"
	"github.com/zdam-egzamin-zawodowy/backend/internal/auth/jwt"
	authusecase "github.com/zdam-egzamin-zawodowy/backend/internal/auth/usecase"
	"github.com/zdam-egzamin-zawodowy/backend/internal/db"
	"github.com/zdam-egzamin-zawodowy/backend/internal/gin/middleware"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/dataloader"
	professionrepository "github.com/zdam-egzamin-zawodowy/backend/internal/profession/repository"
	professionusecase "github.com/zdam-egzamin-zawodowy/backend/internal/profession/usecase"
	qualificationrepository "github.com/zdam-egzamin-zawodowy/backend/internal/qualification/repository"
	qualificationusecase "github.com/zdam-egzamin-zawodowy/backend/internal/qualification/usecase"
	questionrepository "github.com/zdam-egzamin-zawodowy/backend/internal/question/repository"
	questionusecase "github.com/zdam-egzamin-zawodowy/backend/internal/question/usecase"
	userrepository "github.com/zdam-egzamin-zawodowy/backend/internal/user/repository"
	userusecase "github.com/zdam-egzamin-zawodowy/backend/internal/user/usecase"

	ginlogrus "github.com/Kichiyaki/gin-logrus"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/filestorage"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/mode"
	envutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/env"
)

func init() {
	os.Setenv("TZ", "UTC")

	if mode.Get() == mode.DevelopmentMode {
		godotenv.Load(".env.local")
	}

	setupLogger()
}

func main() {
	fileStorage := filestorage.New(&filestorage.Config{
		BasePath: os.Getenv("FILE_STORAGE_PATH"),
	})

	dbConn, err := db.New(&db.Config{
		LogQueries: envutils.GetenvBool("LOG_DB_QUERIES"),
	})
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "Couldn't connect to the db"))
	}
	logrus.Info("Connection with the database has been established")

	//repositories
	userRepository, err := userrepository.NewPGRepository(&userrepository.PGRepositoryConfig{
		DB: dbConn,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	professionRepository, err := professionrepository.NewPGRepository(&professionrepository.PGRepositoryConfig{
		DB: dbConn,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	qualificationRepository, err := qualificationrepository.NewPGRepository(&qualificationrepository.PGRepositoryConfig{
		DB: dbConn,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	questionRepository, err := questionrepository.NewPGRepository(&questionrepository.PGRepositoryConfig{
		DB:          dbConn,
		FileStorage: fileStorage,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	//usecases
	authUsecase, err := authusecase.New(&authusecase.Config{
		UserRepository: userRepository,
		TokenGenerator: jwt.NewTokenGenerator(os.Getenv("ACCESS_SECRET")),
	})
	if err != nil {
		logrus.Fatal(err)
	}
	userUsecase, err := userusecase.New(&userusecase.Config{
		UserRepository: userRepository,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	professionUsecase, err := professionusecase.New(&professionusecase.Config{
		ProfessionRepository: professionRepository,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	qualificationUsecase, err := qualificationusecase.New(&qualificationusecase.Config{
		QualificationRepository: qualificationRepository,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	questionUsecase, err := questionusecase.New(&questionusecase.Config{
		QuestionRepository: questionRepository,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	router := setupRouter()
	graphql := router.Group("")
	graphql.Use(
		middleware.GinContextToContext(),
		middleware.DataLoaderToContext(dataloader.Config{
			ProfessionRepo:    professionRepository,
			QualificationRepo: qualificationRepository,
		}),
		middleware.Authenticate(authUsecase),
	)
	graphqlhttpdelivery.Attach(graphql, graphqlhttpdelivery.Config{
		Resolver: &resolvers.Resolver{
			AuthUsecase:          authUsecase,
			UserUsecase:          userUsecase,
			ProfessionUsecase:    professionUsecase,
			QualificationUsecase: qualificationUsecase,
			QuestionUsecase:      questionUsecase,
		},
		Directive: &directive.Directive{},
	})
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
	logrus.Info("Server exiting")
}

func setupLogger() {
	if mode.Get() == mode.DevelopmentMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	timestampFormat := "2006-01-02 15:04:05"
	if mode.Get() == mode.ProductionMode {
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

func setupRouter() *gin.Engine {
	router := gin.New()

	router.Use(ginlogrus.Logger(logrus.WithField("hostname", "api")), gin.Recovery())
	if mode.Get() == mode.DevelopmentMode {
		router.Use(cors.New(cors.Config{
			AllowOriginFunc: func(string) bool {
				return true
			},
			AllowCredentials: true,
			ExposeHeaders:    []string{"Authorization"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			AllowWebSockets:  false,
		}))
	}

	return router
}
