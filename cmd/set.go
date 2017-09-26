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
	"context"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/k0kubun/pp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youyo/hidy/lib/hidy"
)

var (
	name  string
	file  string
	str   string
	value string
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set to Parameter-store",
	Long:  `Set to Parameter-store`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := hidy.NewConfig()
		awsProfile := viper.Get("aws_profile").(string)
		cfg.SetProfileName(awsProfile)
		_ = cfg.FetchArn()
		s, _ := hidy.NewSession(cfg.SourceProfile)
		svc := hidy.NewService(s)
		resp, _ := svc.AssumingRole(cfg)
		creds := credentials.NewStaticCredentials(
			*resp.Credentials.AccessKeyId,
			*resp.Credentials.SecretAccessKey,
			*resp.Credentials.SessionToken,
		)
		session, _ := session.NewSession(&aws.Config{Credentials: creds})

		ssmClient := ssm.New(session)
		ctx := context.Background()
		if file != "" {
			b, _ := ioutil.ReadFile(file)
			value = string(b)
		}
		if str != "" {
			value = str
		}
		params := &ssm.PutParameterInput{
			Name:  aws.String(name),
			Type:  aws.String("SecureString"),
			Value: aws.String(value),
		}
		r, err := ssmClient.PutParameterWithContext(ctx, params)
		pp.Print(err)
		pp.Print(r)
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
	setCmd.Flags().StringVarP(&name, "name", "n", "", "Parameter name")
	setCmd.Flags().StringVarP(&file, "file", "f", "", "A file to read the string to be set in the parameter store")
	setCmd.Flags().StringVarP(&str, "string", "s", "", "String to set for the parameter store")
}
