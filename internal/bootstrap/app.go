package bootstrap

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	authhandler "sys-admin-serve/internal/handler/auth"
	categoryhandler "sys-admin-serve/internal/handler/category"
	healthhandler "sys-admin-serve/internal/handler/health"
	"sys-admin-serve/internal/middleware"
	"sys-admin-serve/internal/pkg/cache"
	"sys-admin-serve/internal/pkg/config"
	"sys-admin-serve/internal/pkg/database"
	jwtutil "sys-admin-serve/internal/pkg/jwt"
	repositoryauth "sys-admin-serve/internal/repository/auth"
	repositorycategory "sys-admin-serve/internal/repository/category"
	"sys-admin-serve/internal/router"
	serviceauth "sys-admin-serve/internal/service/auth"
	servicecategory "sys-admin-serve/internal/service/category"
	appLogger "sys-admin-serve/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const shutdownTimeout = 10 * time.Second

// App 聚合应用启动所需的核心依赖，负责承载配置、日志、路由和基础设施客户端。
type App struct {
	Config *config.Config
	Logger *zap.Logger
	Engine *gin.Engine
	Server *http.Server
	DB     *gorm.DB
	Redis  *redis.Client
}

// New 按启动顺序完成配置加载、日志初始化、MySQL 和 Redis 连接建立，以及 HTTP Server 装配。
func New(configPath string) (*App, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

	log, err := appLogger.New(cfg.Logger, cfg.App)
	if err != nil {
		return nil, err
	}

	db, err := database.NewMySQL(cfg.MySQL)
	if err != nil {
		return nil, err
	}

	redisClient, err := cache.NewRedis(cfg.Redis)
	if err != nil {
		sqlDB, closeErr := db.DB()
		if closeErr == nil {
			_ = sqlDB.Close()
		}
		return nil, err
	}

	healthHandler := healthhandler.NewHandler(cfg)
	authRepository := repositoryauth.NewRepository(db)
	categoryRepository := repositorycategory.NewRepository(db)
	jwtManager := jwtutil.NewManager(cfg.JWT)
	authService := serviceauth.NewService(authRepository, jwtManager, log)
	categoryService := servicecategory.NewService(categoryRepository, log)
	authHandler := authhandler.NewHandler(authService)
	categoryHandler := categoryhandler.NewHandler(categoryService)
	authMiddleware := middleware.JWTAuth(jwtManager, log)

	engine := router.New(cfg, log, router.Dependencies{
		HealthHandler:   healthHandler,
		AuthHandler:     authHandler,
		CategoryHandler: categoryHandler,
		AuthMiddleware:  authMiddleware,
	})
	server := &http.Server{
		Addr:              cfg.Server.Address(),
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		Config: cfg,
		Logger: log,
		Engine: engine,
		Server: server,
		DB:     db,
		Redis:  redisClient,
	}, nil
}

// Run 启动 HTTP 服务，并在收到退出信号时执行优雅停机和基础设施资源释放。
func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		a.Logger.Info("shutting down server")
		if err := a.closeInfrastructure(); err != nil {
			a.Logger.Error("close infrastructure", zap.Error(err))
		}
		if err := a.Server.Shutdown(shutdownCtx); err != nil {
			errCh <- fmt.Errorf("shutdown server: %w", err)
		}
	}()

	a.Logger.Info(
		"server started",
		zap.String("address", a.Server.Addr),
		zap.String("mysql", a.Config.MySQL.DBName),
		zap.String("redis", a.Config.Redis.Address()),
	)
	if err := a.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen server: %w", err)
	}

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

// closeInfrastructure 统一关闭应用持有的外部资源，避免在停机时泄漏连接。
func (a *App) closeInfrastructure() error {
	var joinedErr error

	if a.Redis != nil {
		if err := a.Redis.Close(); err != nil {
			joinedErr = errors.Join(joinedErr, fmt.Errorf("close redis: %w", err))
		}
	}

	if a.DB != nil {
		sqlDB, err := a.DB.DB()
		if err != nil {
			joinedErr = errors.Join(joinedErr, fmt.Errorf("get mysql db: %w", err))
		} else if err = closeSQLDB(sqlDB); err != nil {
			joinedErr = errors.Join(joinedErr, err)
		}
	}

	return joinedErr
}

// closeSQLDB 负责关闭底层 MySQL 连接池，隔离具体关闭细节。
func closeSQLDB(sqlDB *sql.DB) error {
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("close mysql: %w", err)
	}

	return nil
}
