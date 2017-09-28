package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youyo/hidy/lib/hidy"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get from Parameter-store",
	Long:  `Get from Parameter-store`,
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
		params := &ssm.GetParameterInput{
			Name:           aws.String(name),
			WithDecryption: aws.Bool(true),
		}
		r, err := ssmClient.GetParameterWithContext(ctx, params)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(*r.Parameter.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&name, "name", "n", "", "Parameter name")
}
