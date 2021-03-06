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
		photosRouter.GET("/:photo_id", middleware.Authentication(), controller.GetPhotoByID)
		photosRouter.PUT("/:photo_id", middleware.Authentication(), controller.EditPhotoByID)
		photosRouter.DELETE("/:photo_id", middleware.Authentication(), controller.DeletePhotoByID)
	}

	commentsRouter := r.Group("/comments")
	{
		commentsRouter.POST("/", middleware.Authentication(), controller.CreateComment)
		commentsRouter.GET("/", middleware.Authentication(), controller.GetComment)
		commentsRouter.PUT("/:comment_id", middleware.Authentication(), controller.EditCommentByID)
		commentsRouter.DELETE("/:comment_id", middleware.Authentication(), controller.DeleteCommentByID)
	}

	socialMediasRouter := r.Group("/socialmedias")
	{
		socialMediasRouter.POST("/", middleware.Authentication(), controller.PostSocialMedia)
		socialMediasRouter.GET("/", middleware.Authentication(), controller.GetSocialMedias)
		socialMediasRouter.PUT("/:socialmedia_id", middleware.Authentication(), controller.EditSocialMediaByID)
		socialMediasRouter.DELETE("/:socialmedia_id", middleware.Authentication(), controller.DeleteSocialMediaByID)
	}

	return r

}
