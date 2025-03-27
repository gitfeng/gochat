package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gochat/internal/middleware"
	"gochat/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	rdb *redis.Client
	ctx = context.Background()
)

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("配置文件读取失败: %s", err))
	}
}

func initMySQL() {
	dsn := viper.GetString("database.mysql.dsn")
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("MySQL连接失败: " + err.Error())
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("database.redis.addr"),
		Password: viper.GetString("database.redis.password"),
		DB:       viper.GetInt("database.redis.db"),
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("Redis连接失败: " + err.Error())
	}
}

// 修改main函数尾部
func main() {
	initConfig()
	initMySQL()
	// initRedis()

	r := gin.Default()

	// 全局中间件
	r.Use(
		gin.Logger(),
		gin.Recovery(),
		DatabaseMiddleware(db, rdb),
		middleware.TraceMiddleware(),
	)

	/*
		// 内存模式限流（每秒10请求，突发20）
		r.Use(middleware.RateLimitMiddleware(middleware.RateLimiterConfig{
			Rate:        10,
			Burst:       20,
			LimiterType: "ip",
		}))

		// Redis分布式限流（需要传入Redis实例）
		r.Use(middleware.RateLimitMiddleware(middleware.RateLimiterConfig{
			Rate:        100,
			Burst:       200,
			Redis:       rdb,
			KeyPrefix:   "rate_limit:",
			LimiterType: "global",
		}))
	*/
	router.InitRouter(r)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":" + viper.GetString("server.port"),
		Handler: r,
	}

	// 启动服务协程
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务强制关闭: %v", err)
	}
	log.Println("服务已正常退出")
}

// 数据库中间件
func DatabaseMiddleware(db *gorm.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Set("Redis", rdb)
		c.Next()
	}
}
