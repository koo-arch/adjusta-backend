package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/koo-arch/adjusta-backend/configs"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/api/handlers"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/scheduler"
	"github.com/koo-arch/adjusta-backend/api/middlewares"
	_ "github.com/koo-arch/adjusta-backend/ent/runtime"
	_ "github.com/lib/pq"

)

func main() {
	// 環境変数の読み込み
	configs.LoadEnv()

	// DB接続
	DBUser := configs.GetEnv("DB_USER")
	DBName := configs.GetEnv("DB_NAME")
	DBPass := configs.GetEnv("DB_PASS")
	DBHost := configs.GetEnv("DB_HOST")
	DBPort := configs.GetEnv("DB_PORT")
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DBUser, DBPass, DBHost, DBPort, DBName)
	
	client, err := ent.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalf("failed closing connection to postgres: %v", err)
		}
	}()
	
	// データベースのスキーマを更新
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// JWTキーの起動時生成
	keyManager := auth.NewKeyManager(client)
	if err := keyManager.InitializeJWTKeys(ctx); err != nil {
		log.Fatalf("failed to initialize JWT")
	}

	// スケジューラーのセットアップ
	s := scheduler.NewScheduler(client)
	s.SetupJobs(ctx)
	s.Start()
	defer s.Stop()

	
	//Ginフレームワークのデフォルトの設定を使用してルータを作成
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session", store))
	
	// ルートハンドラの定義
	router.GET("/auth/google/login", handlers.GoogleLoginHandler)
	router.GET("/auth/google/callback", handlers.GoogleCallbackHandler(client))
	router.GET("/auth/logout", handlers.LogoutHandler)

	// 認証が必要なAPIグループ
	auth := router.Group("/api").Use(middlewares.SessionRenewalMiddleware(), middlewares.AuthMiddleware(client))
	{
		auth.GET("/users/me", handlers.GetCurrentUserHandler(client))
	}
	
	// サーバー起動
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
