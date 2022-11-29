package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"

	"todolist.go/db"
	"todolist.go/service"
)

const port = 8000

func main() {
	// initialize DB connection
	dsn := db.DefaultDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}

	// initialize Gin engine
	engine := gin.Default()
	engine.LoadHTMLGlob("views/*.html")

	// prepare session
    store := cookie.NewStore([]byte("my-secret"))
    engine.Use(sessions.Sessions("user-session", store))

	// routing
	engine.Static("/assets", "./assets")
	engine.GET("/", service.Home)
	engine.GET("/list/:page", service.LoginCheck, service.TaskList)

	taskGroup := engine.Group("/task")
    taskGroup.Use(service.LoginCheck)
    {
        taskGroup.GET("/new", service.NewTaskForm)
        taskGroup.POST("/new", service.RegisterTask)

		taskGroup.Use(service.SecureTaskGroup)
		{
			taskGroup.GET("/:id", service.ShowTask)
		}

		taskGroup.Use(service.SecureTask)
		{
        	taskGroup.GET("/edit/:id", service.EditTaskForm)
        	taskGroup.POST("/edit/:id", service.UpdateTask)
        	taskGroup.GET("/delete/:id", service.DeleteTask)
		}
    }

	// ユーザ登録
    engine.GET("/user/new", service.NewUserForm)
    engine.POST("/user/new", service.RegisterUser)
	// ログイン
	engine.GET("/login", service.Login)
	engine.POST("/login", service.Login)
	// ログアウト
	engine.GET("/logout", service.Logout)

	// アカウント編集
	engine.GET("/user", service.LoginCheck, service.ShowUser)
	engine.GET("/user/edit", service.LoginCheck, service.EditUser)
	engine.POST("/user/edit", service.LoginCheck, service.UpdateUser)
	// アカウント削除
	engine.GET("/user/delete", service.LoginCheck, service.DeleteUser)

	GroupGroup := engine.Group("/group")
	GroupGroup.Use(service.LoginCheck)
	{
		//グループ確認
		GroupGroup.GET("", service.LoginCheck, service.ShowGroup)
		//グループ登録
		GroupGroup.GET("/new", service.NewGroupForm)
		GroupGroup.POST("/new", service.RegisterGroup)
		//グループ登録完了
		GroupGroup.GET("/new/complete/:id", service.Complete)
		//グループ参加
		GroupGroup.GET("/belong", service.Belong)
		GroupGroup.POST("/belong", service.Belong)
		//グループ脱退
		GroupGroup.GET("/leave", service.Leave)
		//グループ編集
		GroupGroup.GET("/edit/:id", service.EditGroupForm)
		GroupGroup.POST("/edit/:id", service.UpdateGroup)
	}
	

	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}
