package router

import (
	"MyGram/controller"
	"MyGram/middleware"

	"github.com/gin-gonic/gin"
)

func MainRouter() *gin.Engine {

	r := gin.Default()

	userRouter := r.Group("/users")
	{
		userRouter.POST("/login", controller.Login)
		userRouter.POST("/register", controller.Register)
		userRouter.PUT("/:user_id", middleware.Authentication(), controller.EditUser)
		userRouter.DELETE("/", middleware.Authentication(), controller.DeleteUser)
	}

	photosRouter := r.Group("/photos")
	{
		photosRouter.POST("/",middleware.Authentication(), controller.AddPhoto)
		photosRouter.GET("/",middleware.Authentication() , controller.GetPhotos)
		photosRouter.GET("/:photo_id")
		photosRouter.PUT("/:photo_id")
		photosRouter.DELETE("/:photo_id")
	}

	commentsRouter := r.Group("/comments")
	{
		commentsRouter.POST("/", middleware.Authentication(), controller.CreateComment)
		commentsRouter.GET("/")
		commentsRouter.PUT("/:comment_id")
		commentsRouter.DELETE("/:comment_id")
	}

	socialMediasRouter := r.Group("/socialmedias")
	{
		socialMediasRouter.POST("/", middleware.Authentication(), controller.PostSocialMedia)
		socialMediasRouter.GET("/")
		socialMediasRouter.PUT("/:socialmedia_id")
		socialMediasRouter.DELETE("/:socialmedia_id")
	}

	return r

}
