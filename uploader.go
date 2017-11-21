package uploader

import (
	"path"
	"github.com/gin-gonic/gin"
	"github.com/axetroy/go-fs"
)

type FileConfig struct {
	Path      string   // 普通文件的存放目录
	MaxSize   int      // 普通文件上传的限制大小，单位byte
	AllowType []string // 允许上传的文件后缀名
}

type ImageConfig struct {
	Path      string // 图片存储路径
	MaxSize   int    // 最大图片上传限制，单位byte
	Thumbnail ThumbnailConfig
}

type ThumbnailConfig struct {
	Path      string // 缩略图存放路径
	MaxWidth  int    // 缩略图最大宽度
	MaxHeight int    // 缩略图最大高度
}

type TConfig struct {
	Path      string //文件上传的根目录
	UrlPrefix string // api的url前缀
	File      FileConfig
	Image     ImageConfig
}

var Config TConfig

func Init() (err error) {
	if err = fs.EnsureDir(Config.Path); err != nil {
		return
	}
	if err = fs.EnsureDir(path.Join(Config.Path, Config.File.Path)); err != nil {
		return
	}
	if err = fs.EnsureDir(path.Join(Config.Path, Config.Image.Path)); err != nil {
		return
	}
	if err = fs.EnsureDir(path.Join(Config.Path, Config.Image.Thumbnail.Path)); err != nil {
		return
	}
	return
}

func Resolve(e *gin.Engine, config TConfig) (err error, uploader *gin.RouterGroup, downloader *gin.RouterGroup) {
	Config = config
	if err = Init(); err != nil {
		return
	}
	// upload all
	uploader = e.Group(Config.UrlPrefix + "/upload")

	uploader.POST("/image", UploaderImage)
	uploader.POST("/file", UploadFile)
	uploader.GET("/example", UploaderTemplate("image"))

	// download all
	downloader = e.Group(Config.UrlPrefix + "/download")

	// download file
	uploadFile := downloader.Group("/file")
	uploadFile.GET("/raw/:filename", GetFileRaw)
	uploadFile.GET("/download/:filename", DownloadFile)

	downloader.Use(func(context *gin.Context) {
		header := context.Writer.Header()
		// alone dns prefect
		header.Set("X-DNS-Prefetch-Control", "on")
		// IE No Open
		header.Set("X-Download-Options", "noopen")
		// not cache
		header.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		header.Set("Expires", "max-age=0")
		// Content Security Policy
		header.Set("Content-Security-Policy", "default-src 'self'")
		// xss protect
		// it will caught some problems is old IE
		header.Set("X-XSS-Protection", "1; mode=block")
		// Referrer Policy
		header.Set("Referrer-Header", "no-referrer")
		// cros frame, allow same origin
		header.Set("X-Frame-Options", "SAMEORIGIN")
		// HSTS
		header.Set("Strict-Transport-Security", "max-age=5184000;includeSubDomains")
		// no sniff
		header.Set("X-Content-Type-Options", "nosniff")
	})

	{
		// download image
		downloadImage := downloader.Group("/image")
		downloadImage.GET("/thumbnail/:filename", GetThumbnailImage)
		downloadImage.GET("/origin/:filename", GetOriginImage)
	}

	return
}
