package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youyo/hidy/lib/hidy"
)

/*
var (
	name    string
	file    string
	str     string
	value   string
	profile string
)
*/

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set to Parameter-store",
	Long:  `Set to Parameter-store`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := hidy.NewConfig()
		awsProfile, err := func() (p string, err error) {
			if profile != "" {
				p = profile
			} else if viper.IsSet("aws_profile") {
				p = viper.GetString("aws_profile")
			} else {
				err = errors.New("aws profile is not set.\nuse '-p' option or use environment 'HIDY_AWS_PROFILE'")
			}
			return
		}()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
		value, err := readValue(str, file)
		if err != nil {
			fmt.Println(err)
		}
		params := &ssm.PutParameterInput{
			Name:  aws.String(name),
			Type:  aws.String("SecureString"),
			Value: aws.String(value),
		}
		if _, err = ssmClient.PutParameterWithContext(ctx, params); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("success")
		}
	},
}

func readValue(str, file string) (value string, err error) {
	if str != "" {
		value = str
	} else if file != "" {
		b, _ := ioutil.ReadFile(file)
		value = string(b)
	} else {
		err = errors.New("value is not set.")
	}
	return
}

func init() {
	RootCmd.AddCommand(setCmd)
	setCmd.Flags().StringVarP(&name, "name", "n", "", "Parameter name")
	setCmd.Flags().StringVarP(&file, "file", "f", "", "A file to read the string to be set in the parameter store")
	setCmd.Flags().StringVarP(&str, "string", "s", "", "String to set for the parameter store")
}
