package s3

import (
	"fmt"
	"io"
	"log"
	"main/src/config"
	"main/src/helpers"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

func UploadToS3(f *multipart.FileHeader, path string) (string, string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           "di",
		SharedConfigState: session.SharedConfigEnable,
	}))

	file, err := f.Open()
	helpers.CheckAndPrintErr(err, "Failed to open file")
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	helpers.CheckAndPrintErr(err, "Failed to read file")

	fileType := http.DetectContentType(buff) // 与えられたデータのContent-Typeを決定する。最大でデータの最初の512バイトを引数にする。
	file.Seek(0, io.SeekStart) // 自動でfileの先頭には戻らないのでここで戻す

	directory := path + uuid.NewRandom().String() + "/"
	var fileName string
	re := regexp.MustCompile("(?i)^[a-z0-9._-]+$")
	// ファイル名のバリデーション
	if re.MatchString(f.Filename) && len(strings.Split(f.Filename, ".")) == 2 && len(directory+f.Filename) < 255 {
		fileName = f.Filename
	} else {
		ext := filepath.Ext(f.Filename)
		fileName = fmt.Sprintf("file%s", ext)
	}
	finalPath := directory + fileName

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:        file,
		Bucket:      aws.String(config.S3Bucket),
		Key:         aws.String(finalPath),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(fileType),
	})
	helpers.CheckAndPrintErr(err, "Failed to upload")

	log.Println("Successfully uploaded to", result.Location)

	return result.Location, fileType
}