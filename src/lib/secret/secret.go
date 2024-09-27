package secret

import "os"

var DBPassword = func() string {
	postgresPasswordBytes, err := os.ReadFile(os.Getenv("POSTGRES_PASSWORD_FILE"))
	if err != nil {
		panic(err)
	}

	return string(postgresPasswordBytes[:len(postgresPasswordBytes)-1])
}()

var YandexS3AccessKeyId = func() string {
	accessKeyIdBytes, err := os.ReadFile(os.Getenv("AWS_ACCESS_KEY_ID_FILE"))
	if err != nil {
		panic(err)
	}

	return string(accessKeyIdBytes[:len(accessKeyIdBytes)-1])
}()

var YandexS3SecretAccessKey = func() string {
	accessSecretKeyBytes, err := os.ReadFile(os.Getenv("AWS_ACCESS_SECRET_KEY_FILE"))
	if err != nil {
		panic(err)
	}

	return string(accessSecretKeyBytes[:len(accessSecretKeyBytes)-1])
}()
