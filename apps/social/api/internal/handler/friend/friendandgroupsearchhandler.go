package friend

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"imooc.com/easy-chat/apps/social/api/internal/logic/friend"
	"imooc.com/easy-chat/apps/social/api/internal/svc"
	"imooc.com/easy-chat/apps/social/api/internal/types"
)

// 搜索好友和群聊
func FriendAndGroupSearchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FriendAndGroupSearchReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := friend.NewFriendAndGroupSearchLogic(r.Context(), svcCtx)
		resp, err := l.FriendAndGroupSearch(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
