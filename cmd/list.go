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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List parameters",
	Long:  `List parameters`,
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
		params := &ssm.DescribeParametersInput{
			MaxResults: aws.Int64(50),
		}
		r, err := ssmClient.DescribeParametersWithContext(ctx, params)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, v := range r.Parameters {
				fmt.Println(*v.Name)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&profile, "profile", "p", "", "aws profile name")
}
