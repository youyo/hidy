// Copyright © 2017 youyo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youyo/hidy/lib/hidy"
)

var encFile string

// encCmd represents the enc command
var encCmd = &cobra.Command{
	Use:   "enc",
	Short: "Encrypt",
	Long:  `Encrypt`,
	Run: func(cmd *cobra.Command, args []string) {
		p := viper.Get("aws_profile").(string)
		b := viper.Get("s3_bucket").(string)
		keyId := viper.Get("kms_key_id").(string)
		cfg, _ := hidy.NewConfig()
		cfg.SetProfileName(p)
		_ = cfg.FetchArn()
		s, _ := hidy.NewSession(cfg.SourceProfile)
		svc := hidy.NewService(s)
		resp, _ := svc.AssumingRole(cfg)
		creds := credentials.NewStaticCredentials(
			*resp.Credentials.AccessKeyId,
			*resp.Credentials.SecretAccessKey,
			*resp.Credentials.SessionToken,
		)
		ss, _ := session.NewSession(&aws.Config{Credentials: creds})

		kmsClient := kms.New(ss)
		secretBytes, _ := ioutil.ReadFile(encFile)
		params := &kms.EncryptInput{
			KeyId:     &keyId,
			Plaintext: secretBytes,
		}
		r, _ := kmsClient.Encrypt(params)

		s3Client := s3.New(ss)
		input := &s3.PutObjectInput{
			Body:                 aws.ReadSeekCloser(bytes.NewReader(r.CiphertextBlob)),
			Bucket:               aws.String(b),
			Key:                  aws.String("hidy_" + encFile),
			ServerSideEncryption: aws.String("AES256"),
			Tagging:              aws.String("CreatedBy=hidy"),
		}
		_, err := s3Client.PutObject(input)
		if err != nil {
			fmt.Println("Object name is " + encFile)
		}
	},
}

func init() {
	RootCmd.AddCommand(encCmd)
	encCmd.Flags().StringVarP(&encFile, "file", "f", "", "Target file")
}
