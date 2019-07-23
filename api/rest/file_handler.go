package rest

import (
	"github.com/labstack/echo"

	log "github.com/sirupsen/logrus"
)

func (api *Server) NewFileHandler(g *echo.Group) {
	g.GET("/:bucket/:key", api.GetFile)
}

// swagger:operation GET /{bucket}/{key} get file by name
// ---
// summary: Get a File by Filename
// description: If the file will not be found, a 404 will be returned
// parameters:
// - name: fileName
//   in: key
//   description: name of file
//   type: string
//   required: true
// - name: bucketName
//   in: bucket
//   description: name of bucket
//   type: string
//   required: true
func (api *Server) GetFile(c echo.Context) error {
	bucket := c.Param("bucket")
	key := c.Param("key")
	log.Debug("bucket=", bucket, " key=", key)
	file := api.services.userService.GetFile(bucket, key)
	return c.File(file.Name())
}
