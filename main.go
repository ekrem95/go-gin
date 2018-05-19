package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
)

func main() {
	if err := testSQLConnection(); err != nil {
		log.Fatal(err)
	}

	r := router()

	r.Run(":8080")
}

func router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	r.Use(sessions.Sessions("session", store))
	r.LoadHTMLGlob("./app/templates/*")
	r.StaticFS("/src", http.Dir("./app/src"))
	r.StaticFile("/favicon.ico", "./app/templates/favicon.ico")

	// socketio
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so.Join("chat")

		so.On("msg", func(msg *Message) {
			so.BroadcastTo("chat", "dist", msg)
			RedisSaveMsg(msg)
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	r.GET("/", common)
	r.GET("/signup", common)
	r.GET("/login", common)
	r.GET("/add", common)
	r.GET("/upload", common)
	r.GET("/user", getUser)
	r.GET("/messages", RedisGetMsgs)
	r.GET("/api/posts", getPosts)
	r.GET("/api/postbyid/:id", getPostByID)
	r.GET("/api/commentsbyid/:id", getCommentsByID)
	r.GET("/p/*all", common)
	r.GET("/myposts", common)
	r.GET("/api/getpostbyusername/:name", getPostsByUsername)
	r.GET("/edit/:id", common)
	r.GET("/changepassword", common)
	r.GET("/get_likes/:id", getLikes)

	r.POST("/signup", signup)
	r.POST("/login", login)
	r.POST("/logout", logout)
	r.POST("/add", addPost)
	r.POST("/comment", postComment)
	r.POST("/upload", uploadFile)
	r.POST("/edit/:id", editPost)
	r.POST("/delete/:id", deletePostByID)
	r.POST("/changepassword", changePassword)
	r.POST("/post_likes", postLikes)

	// socketio
	r.GET("/socket.io/", gin.WrapH(server))
	r.POST("/socket.io/", gin.WrapH(server))

	return r
}
