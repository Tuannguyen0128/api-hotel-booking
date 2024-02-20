package app

import (
	"api-hotel-booking/internal/app/server"
	database "api-hotel-booking/internal/app/thirdparty/database"
	"api-hotel-booking/internal/app/thirdparty/mongo"
	"time"
)

type (
	Config struct {
		Environment   string        `mapstructure:"Environment"`
		AppID         string        `mapstructure:"AppID"`
		AppAbbr       string        `mapstructure:"AppAbbr"`
		ReqAppID      string        `mapstructure:"ReqAppID"`
		Server        server.Config `mapstructure:"Server"`
		Service       Service       `mapstructure:"Service"`
		Mongo         mongo.Config  `mapstructure:"Mongo"`
		MySQL         database.MysqlConf  `mapstructure:"Mysql"`
		Email         Email         `mapstructure:"Email"`
		HTTPTransport HTTPTransport `mapstructure:"HttpTransport"`
		IsDebug       bool          `mapstructure:"DebugMode"`
		CarLog        string        `mapstructure:"CarLogId"`
		AWSS3Config   S3            `mapstructure:"S3"`
		MasterFiles   MasterFiles   `mapstructure:"MasterFiles"`
	}

	Service struct {
		Jwt         JwtConfig      `mapstructure:"Jwt"`
		Session     SessionConfig  `mapstructure:"Session"`
		ForgetPW    ForgetPWConfig `mapstructure:"ForgetPassword"`
		Car         CarConfig      `mapstructure:"Car"`
		Prefix      string         `mapstructure:"Prefix"`
		RequestSize int64          `mapstructure:"RequestSize"`
		API         EndPoints      `mapstructure:"API"`
		Cache       Cache          `mapstructure:"Cache"`
	}
	JwtConfig struct {
		SigningKey string `mapstructure:"SigningKey"`
	}
	SessionConfig struct {
		Timeout                   int    `mapstructure:"TimeoutSec"`
		IsRequireEncryption       bool   `mapstructure:"IsRequireEncryption"`
		WrongPasswordLockLimit    int    `mapstructure:"WrongPasswordLockLimit"`
		WrongPasswordLockToStatus string `mapstructure:"WrongPasswordLockToStatus"`
		AllowLoginStatus          string `mapstructure:"AllowLoginStatus"`
	}
	ForgetPWConfig struct {
		SigningKey string `mapstructure:"SigningKey"`
		Timeout    int    `mapstructure:"TimeoutSec"`
	}
	CarConfig struct {
		PublicUrlTimeHour int `mapstructure:"PublicUrlTimeHour"`
	}
	Cache struct {
		Company struct {
			CacheTimeInSec    int `mapstructure:"CacheTimeInSec"`
			CacheSize         int `mapstructure:"CacheSize"`
			MemorizerTimeInMs int `mapstructure:"MemorizerTimeInMs"`
		} `mapstructure:"Company"`
	}

	EndPoints struct {
		Web struct {
			V1 struct {
				Login             API `mapstructure:"Login"`
				ExtendSession     API `mapstructure:"ExtendSession"`
				Logout            API `mapstructure:"Logout"`
				CreateUser        API `mapstructure:"CreateUser"`
				EditUser          API `mapstructure:"EditUser"`
				GetUser           API `mapstruture:"GetUser"`
				ChangePassword    API `mapstructure:"ChangePassword"`
				ListUser          API `mapstructure:"ListUser"`
				DeleteUser        API `mapstructure:"DeleteUser"`
				ForgetPassword    API `mapstructure:"ForgetPassword"`
				ResetPassword     API `mapstructure:"ResetPassword"`
				CreateCompany     API `mapstructure:"CreateCompany"`
				GetCompany        API `mapstructure:"GetCompany"`
				EditCompany       API `mapstructure:"EditCompany"`
				ListCompany       API `mapstructure:"ListCompany"`
				CreateCar         API `mapstructure:"CreateCar"`
				ListCar           API `mapstructure:"ListCar"`
				CountCar          API `mapstructure:"CountCar"`
				GetCar            API `mapstructure:"GetCar"`
				SubmitCar         API `mapstructure:"SubmitCar"`
				CancelCar         API `mapstructure:"CancelCar"`
				DisapproveCar     API `mapstructure:"DisapproveCar"`
				ApproveCar        API `mapstructure:"ApproveCar"`
				RevokeCar         API `mapstructure:"RevokeCar"`
				EditCar           API `mapstructure:"EditCar"`
				ListCarModel      API `mapstructure:"ListCarModel"`
				ListCarDetails    API `mapstructure:"ListCarDetails"`
				UploadImages      API `mapstructure:"UploadImages"`
				DeleteImages      API `mapstructure:"DeleteImages"`
				CreateCarFromFile API `mapstructure:"CreateCarFromFile"`
			} `mapstructure:"V1"`
		} `mapstructure:"Web"`
	}

	API struct {
		URL     string `mapstructure:"URL"`
		Method  string `mapstructure:"Method"`
		Enabled bool   `mapstructure:"Enabled"`
	}

	Email struct {
		Address string `mapstructure:"Address"`
		Auth    struct {
			Username string `mapstructure:"Username"`
			Password string `mapstructure:"Password"`
			Host     string `mapstructure:"Host"`
		} `mapstructure:"Auth"`
		Sender             string `mapstructure:"Sender"`
		ConnectionPoolSize int    `mapstructure:"ConnectionPoolSize"`
		SendTimeout        int    `mapstructure:"SendTimeout"`
		RetryAttempts      int    `mapstructure:"RetryAttempts"`
		RetryTimeout       int    `mapstructure:"RetryTimeout"`
		Url                string `mapstructure:"Url"`
	}

	HTTPTransport struct {
		MaxIdleConns          int           `mapstructure:"MaxIdleConns"`
		MaxIdleConnsPerHost   int           `mapstructure:"MaxIdleConnsPerHost"`
		IdleConnTimeout       time.Duration `mapstructure:"IdleConnTimeout"`
		TLSHandshakeTimeout   time.Duration `mapstructure:"TlsHandshakeTimeout"`
		ResponseHeaderTimeout time.Duration `mapstructure:"ResponseHeaderTimeout"`
		ExpectContinueTimeout time.Duration `mapstructure:"ExpectContinueTimeout"`
	}
	S3 struct {
		Region       string      `mapstructure:"Region"`
		Bucket       string      `mapstructure:"Bucket"`
		Acl          string      `mapstructure:"Acl"`
		Encryption   string      `mapstructure:"Encryption"`
		StorageClass string      `mapstructure:"StorageClass"`
		Credential   Credentials `mapstructure:"Credential"`
	}
	Credentials struct {
		SecretAccessKey string `mapstructure:"SecretAccessKey"`
		AccessKeyId     string `mapstructure:"AccessKeyId"`
	}

	MasterFiles struct {
		Car struct {
			Path     string `mapstructure:"Path"`
			FileName struct {
				Brand             string `mapstructure:"Brand"`
				BrandModelSub     string `mapstructure:"BrandModelSub"`
				BrandModelYears   string `mapstructure:"BrandModelYears"`
				BrandModelYearSub string `mapstructure:"BrandModelYearSub"`
				Province          string `mapstructure:"Province"`
				VehicleType       string `mapstructure:"VehicleType"`
				BodyType          string `mapstructure:"BodyType"`
				Color             string `mapstructure:"Color"`
				Fuel              string `mapstructure:"Fuel"`
				EngineBrand       string `mapstructure:"EngineBrand"`
				Gear              string `mapstructure:"Gear"`
			}
			FileType string `mapstructure:"FileType"`
		}
	}
)

var CFG = &Config{}
