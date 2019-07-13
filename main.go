package main

import (
	"caspar/gin-blog/pkg/setting"
	"caspar/gin-blog/routers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           setting.HTTPPort,
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func(){
		if err := s.ListenAndServe(); err != nil {
			fmt.Println("Listen:%s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")



	/*endless.DefaultReadTimeOut = setting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.HTTPPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string){
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err:%v", err)
	}*/



}
