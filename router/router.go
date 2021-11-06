package router

import (
	"MyGram/controller"

	"github.com/gin-gonic/gin"
)

func MainRouter() *gin.Engine {

	r := gin.Default()

	userRouter := r.Group("/users")
	{
		userRouter.POST("/login", controller.Login)
		userRouter.POST("/register", controller.Register)
		userRouter.PUT("/")
		userRouter.DELETE("/")
	}

	photosRouter := r.Group("/photos")
	{
		photosRouter.POST("/")
		photosRouter.GET("/")
		photosRouter.GET("/:photo_id")
		photosRouter.PUT("/:photo_id")
		photosRouter.DELETE("/:photo_id")
	}

	commentsRouter := r.Group("/comments")
	{
		commentsRouter.POST("/")
		commentsRouter.GET("/")
		commentsRouter.PUT("/:comment_id")
		commentsRouter.DELETE("/:comment_id")
	}

	socialMediasRouter := r.Group("/socialmedias")
	{
		socialMediasRouter.POST("/")
		socialMediasRouter.GET("/")
		socialMediasRouter.PUT("/:socialmedia_id")
		socialMediasRouter.DELETE("/:socialmedia_id")
	}

	return r

}
