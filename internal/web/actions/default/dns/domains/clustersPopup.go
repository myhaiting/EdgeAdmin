package domains

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type ClustersPopupAction struct {
	actionutils.ParentAction
}

func (this *ClustersPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *ClustersPopupAction) RunGet(params struct {
	DomainId int64
}) {
	// 域名信息
	domainResp, err := this.RPC().DNSDomainRPC().FindEnabledBasicDNSDomain(this.AdminContext(), &pb.FindEnabledBasicDNSDomainRequest{
		DnsDomainId: params.DomainId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	domain := domainResp.DnsDomain
	if domain == nil {
		this.NotFound("dnsDomain", params.DomainId)
		return
	}

	this.Data["domain"] = domain.Name

	// 集群
	clustersResp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClustersWithDNSDomainId(this.AdminContext(), &pb.FindAllEnabledNodeClustersWithDNSDomainIdRequest{DnsDomainId: params.DomainId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterMaps := []maps.Map{}
	for _, cluster := range clustersResp.NodeClusters {
		isOk := false
		if len(cluster.Name) > 0 {
			checkResp, err := this.RPC().DNSDomainRPC().ExistDNSDomainRecord(this.AdminContext(), &pb.ExistDNSDomainRecordRequest{
				DnsDomainId: params.DomainId,
				Name:        cluster.DnsName,
				Type:        "A",
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			isOk = checkResp.IsOk
		}

		clusterMaps = append(clusterMaps, maps.Map{
			"id":      cluster.Id,
			"name":    cluster.Name,
			"dnsName": cluster.DnsName,
			"isOk":    isOk,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}