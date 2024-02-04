package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	// sessionの作成(Mustはエラーが発生した場合にパニックを起こす)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		// 設定や認証情報の集まり
		Profile:           "di",
		//  共有設定ファイル (~/.aws/config) および共有資格情報ファイル (~/.aws/credentials) から設定を読み込む
		SharedConfigState: session.SharedConfigEnable,
	}))

	// S3オブジェクトを書き込むファイルの作成
	f, err := os.Create("sample.txt")
	if err != nil {
		log.Fatal(err)
	}

	bucketName := "xxx-bucket"
	objectKey := "xxx-key"

	// Downloaderを作成し、S3オブジェクトをダウンロード
	downloader := s3manager.NewDownloader(sess)
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("DownloadedSize: %d byte", n)
}