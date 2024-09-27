package roothttp

import (
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/application/presenters/user"
	"nosebook/src/application/services/user"
	"nosebook/src/deps_root/http/exec"
	reqcontext "nosebook/src/deps_root/http/req_context"
	"nosebook/src/errors"
	"nosebook/src/lib/image"
	"nosebook/src/lib/secret"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type s3AvatarStorage struct {
	s3 *s3.S3
}

func (this *s3AvatarStorage) Upload(img *image.Image, userId uuid.UUID) (string, *errors.Error) {
	folder := "users/avatars/"
	filename := userId.String() + img.Extension()
	key := folder + filename

	_, err := errors.Using(this.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("nosebook"),
		Key:    aws.String(key),
		Body:   img.NewReader(),
	}))
	if err != nil {
		return "", err
	}

	downloadUrl := "https://storage.yandexcloud.net/nosebook/" + key
	return downloadUrl, nil
}

func (this *RootHTTP) addUserHandlers(userRepository user.UserRepository) {
	presenter := presenteruser.New(this.db)
	s3 := s3.New(session.Must(session.NewSession(&aws.Config{
		Endpoint:    aws.String("https://storage.yandexcloud.net"),
		Region:      aws.String("ru-central1"),
		Credentials: credentials.NewStaticCredentials(secret.YandexS3AccessKeyId, secret.YandexS3SecretAccessKey, ""),
	})))
	service := user.New(userRepository, &s3AvatarStorage{
		s3: s3,
	}, this.tracer)

	group := this.authRouter.Group("/users")

	group.POST("/change-avatar", execCommand(service.ChangeAvatar, this, exec.WithFileBinding))

	group.GET("/:id", func(ctx *gin.Context) {
		reqctx := reqcontext.From(ctx)
		id, ok := reqctx.ParamUUID("id")
		if !ok {
			return
		}

		users, ok := handle(presenter.FindByIds(ctx.Request.Context(), []uuid.UUID{id}))(reqctx)
		if !ok {
			return
		}

		var user *presenterdto.User
		if len(users) > 0 {
			user = users[id]
		}

		reqctx.SetResponseData(user)
		reqctx.SetResponseOk(true)
	})

	group.GET("", exec.Presenter(presenter.FindByText, map[string]exec.PresenterOption{
		"text": {
			Type: exec.STRING,
		},
		"next": {
			Type: exec.STRING,
		},
	}, &presenteruser.FindByTextInput{}, this.tracer))
}
