package hidy

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type (
	service struct {
		*sts.STS
	}
)

func NewSession(sourceProfile string) (s *session.Session, err error) {
	cred := credentials.NewSharedCredentials("", sourceProfile)
	s, err = session.NewSession(&aws.Config{Credentials: cred})
	return
}

func NewService(s *session.Session) (svc *service) {
	svc = &service{sts.New(s)}
	return
}

func (svc *service) AssumingRole(cfg *Config) (resp *sts.AssumeRoleOutput, err error) {
	roleSessionName := extractRoleSessionName(cfg.ARN)
	params := func() *sts.AssumeRoleInput {
		if cfg.MfaSerial == "" {
			return &sts.AssumeRoleInput{
				RoleArn:         aws.String(cfg.ARN),
				RoleSessionName: aws.String(roleSessionName),
			}
		}
		return &sts.AssumeRoleInput{
			RoleArn:         aws.String(cfg.ARN),
			RoleSessionName: aws.String(roleSessionName),
			SerialNumber:    aws.String(cfg.MfaSerial),
			TokenCode:       aws.String(cfg.MfaCode),
		}
	}()
	return svc.AssumeRole(params)
}

func extractRoleSessionName(arn string) (roleSessionName string) {
	roleSessionName = strings.Split(arn, "/")[1] + "@hidy"
	return
}
