package hidy

import (
	"errors"
	"os"
	"path/filepath"

	ini "gopkg.in/ini.v1"
)

type (
	Config struct {
		Path          string
		Data          *ini.File
		ARN           string
		SourceProfile string
		MfaSerial     string
		MfaCode       string
		ProfileName   string
		S3Bucket      string
		KmsKeyId      string
	}
)

func NewConfig(path ...string) (cfg *Config, err error) {
	cp := func() string {
		if len(path) == 0 {
			return configPath()
		} else {
			return path[0]
		}
	}()
	cfg = &Config{
		Path: cp,
	}
	cfg.Data, err = loadConfig(cfg.Path)
	return
}

func configPath() (c string) {
	c = filepath.Join(os.Getenv("HOME"), ".aws/config")
	return
}

func loadConfig(configPath string) (cfg *ini.File, err error) {
	cfg, err = ini.Load(configPath)
	return
}

func (cfg *Config) SetProfileName(profileName string) {
	cfg.ProfileName = profileName
}

func (cfg *Config) SetS3Bucket(s3Bucket string) {
	cfg.S3Bucket = s3Bucket
}

func (cfg *Config) SetKmsKeyId(kmsKeyId string) {
	cfg.KmsKeyId = kmsKeyId
}

func (cfg *Config) SetMfaCode(mfaCode string) {
	cfg.MfaCode = mfaCode
}

func (cfg *Config) FetchArn() (err error) {
	s := "profile " + cfg.ProfileName
	cfg.ARN = cfg.Data.Section(s).Key("role_arn").String()
	cfg.SourceProfile = cfg.Data.Section(s).Key("source_profile").String()
	cfg.MfaSerial = cfg.Data.Section(s).Key("mfa_serial").String()
	if cfg.ARN == "" {
		err = errors.New("Could not fetch Arn")
		return
	}
	return
}
