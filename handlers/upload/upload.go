package upload

import (
	/*
		"bytes"
		"context"
		"fmt"
		"io"

		"mime"
		"net/http"
		"time"
	*/
	"mime/multipart"

	"github.com/gin-gonic/gin"
	/*
		Models "github.com/restdoc/restdoc-models"
		pb "restdoc/internal/proto/base"
	*/)

type MailUpload struct {
	Attachments []*multipart.FileHeader `json:"file,omitempty" form:"file" binding:"required"` // change form to "attachments[]" take no effect also.
}

func Upload(c *gin.Context) {
}
