package ipadmin

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type AllowListAction struct {
	actionutils.ParentAction
}

func (this *AllowListAction) Init() {
	this.Nav("", "setting", "allowList")
	this.SecondMenu("waf")
}

func (this *AllowListAction) RunGet(params struct {
	ServerId         int64
	FirewallPolicyId int64
}) {
	this.Data["featureIsOn"] = true
	this.Data["firewallPolicyId"] = params.FirewallPolicyId

	listId, err := dao.SharedIPListDAO.FindAllowIPListIdWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建
	if listId == 0 {
		listId, err = dao.SharedIPListDAO.CreateIPListForServerId(this.AdminContext(), params.ServerId, "white")
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["listId"] = listId

	// 数量
	countResp, err := this.RPC().IPItemRPC().CountIPItemsWithListId(this.AdminContext(), &pb.CountIPItemsWithListIdRequest{IpListId: listId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	// 列表
	itemsResp, err := this.RPC().IPItemRPC().ListIPItemsWithListId(this.AdminContext(), &pb.ListIPItemsWithListIdRequest{
		IpListId: listId,
		Offset:   page.Offset,
		Size:     page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	itemMaps := []maps.Map{}
	for _, item := range itemsResp.IpItems {
		expiredTime := ""
		if item.ExpiredAt > 0 {
			expiredTime = timeutil.FormatTime("Y-m-d H:i:s", item.ExpiredAt)
		}

		itemMaps = append(itemMaps, maps.Map{
			"id":          item.Id,
			"ipFrom":      item.IpFrom,
			"ipTo":        item.IpTo,
			"expiredTime": expiredTime,
			"reason":      item.Reason,
		})
	}
	this.Data["items"] = itemMaps

	this.Show()
}
