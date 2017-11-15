### go文件上传模块

- [x] 文件上传
- [x] 限制上传文件的后缀名
- [x] 限制上传文件的大小
- [x] Hash去重，防止重复上传
- [x] 图片自动生成缩略图
- [x] 全由配置

### 如何运行

```bash
go get -v github.com/axetroy/gin-uploader
```

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/axetroy/gin-uploader"
)

func main() {
	Router := gin.Default()

	uploader.Resolve(Router, uploader.TConfig{
		Path: "upload",
		File: uploader.FileConfig{
			Path:      "files",
			MaxSize:   10485760,
			AllowType: []string{},
		},
		Image: uploader.ImageConfig{
			Path:    "images",
			MaxSize: 10485760,
			Thumbnail: uploader.ThumbnailConfig{
				Path:      "thumbnail",
				MaxWidth:  300,
				MaxHeight: 300,
			},
		},
	})

	if err := Router.Run("localhost:9090"); err != nil {
		panic(err)
		return
	}

})
```

### API

```bash

# 上传相关
[POST]  /upload/image                           # 图片上传
[POST]  /upload/file                            # 其他文件上传
[GET]   /upload/example                         # 上传demo，仅在开发模式下
# 下载相关
[GET]   /download/image/origin/:filename        # 获取上传的原始图片
[GET]   /download/image/thumbnail/:filename     # 获取上传的缩略图片
[GET]   /download/file/download/:filename       # 下载文件
[GET]   /download/file/raw/:filename            # 获取文件
```

## License

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Faxetroy%2Fgo-upload.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Faxetroy%2Fgo-upload?ref=badge_large)
