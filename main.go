package main

import (
	"blog-service/global"
	"blog-service/internal/model"
	"blog-service/internal/routers"
	"blog-service/pkg/logger"
	"blog-service/pkg/setting"
	"blog-service/pkg/tracer"
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	port      string
	runMode   string
	config    string
	isVersion bool
)

func init() {
	if err := setupFlag(); err != nil {
		log.Fatalf("init.setupFlag err: %v", err)
	}
	if err := setupSetting(); err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	if err := setupLogger(); err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
	if err := setupTracer(); err != nil {
		log.Fatalf("init.setupTracer err: %v", err)
	}
	if err := setupDBEngine(); err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}
}

func main() {
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServer err: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")

}

func setupFlag() error {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "configs/", "指定要使用的配置文件路径")
	flag.BoolVar(&isVersion, "version", false, "编译信息")
	flag.Parse()

	return nil
}

func setupSetting() error {
	s, err := setting.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		return err
	}

	if err = s.ReadSection("Server", &global.ServerSetting); err != nil {
		return err
	}

	if err = s.ReadSection("App", &global.AppSetting); err != nil {
		return err
	}

	if err = s.ReadSection("Email", &global.EmailSetting); err != nil {
		return err
	}

	if err = s.ReadSection("Database", &global.DatabaseSetting); err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	if port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}

	return nil
}

func setupLogger() error {
	fileName := global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    500,
		MaxAge:     10,
		MaxBackups: 1024,
		LocalTime:  true,
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}

func setupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("blog-service", "192.168.119.128:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}

	// 初始化表
	if err = model.MysqlTables(global.DBEngine); err != nil {
		return err
	}

	return nil
}
