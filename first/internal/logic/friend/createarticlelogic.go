package friend

import (
	"context"

	"github.com/gnh123/redis-action-practice/first/internal/svc"
	"github.com/gnh123/redis-action-practice/first/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateArticleLogic {
	return &CreateArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateArticleLogic) CreateArticle(req *types.ArticleReq) (resp *types.ArticleResp, err error) {
	// todo: add your logic here and delete this line

	return
}
