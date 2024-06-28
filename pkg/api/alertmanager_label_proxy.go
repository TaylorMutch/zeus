package api

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/alertmanager/api/v2/models"
)

func NewAlertmanagerLabelProxy(serviceName, alertmanagerAddr, proxyLabel, tenantIDHeader string) (*gin.Engine, error) {
	router := New(serviceName)
	alertmanagerUrl, err := url.Parse(alertmanagerAddr)
	if err != nil {
		return nil, err
	}

	proxy := &AlertmanagerLabelProxyAPI{
		alertmanagerAddr: alertmanagerAddr,
		proxyLabel:       proxyLabel,
		tenantIDHeader:   tenantIDHeader,
		revproxy:         httputil.NewSingleHostReverseProxy(alertmanagerUrl),
	}

	// Set the director to modify the request before forwarding it
	proxy.revproxy.Director = func(req *http.Request) {
		req.URL.Scheme = alertmanagerUrl.Scheme
		req.URL.Host = alertmanagerUrl.Host
		req.Host = alertmanagerUrl.Host

		// Remove the tenant ID header before forwarding the request
		req.Header.Del(proxy.tenantIDHeader)
	}

	// TODO - check all paths that need to be posted from the alertmanager
	router.POST("/alertmanager", proxy.handleAlerts)
	return router, nil
}

type AlertmanagerLabelProxyAPI struct {
	alertmanagerAddr string
	proxyLabel       string
	tenantIDHeader   string
	revproxy         *httputil.ReverseProxy
}

func (p *AlertmanagerLabelProxyAPI) handleAlerts(c *gin.Context) {
	var alerts models.PostableAlerts
	err := c.BindJSON(&alerts)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse alerts"})
		return
	}

	tenantID := c.GetHeader(p.tenantIDHeader)
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tenant ID header"})
		return
	}

	alerts = addTenantLabelToAlerts(alerts, p.proxyLabel, tenantID)
	p.revproxy.ServeHTTP(c.Writer, c.Request)
}

func addTenantLabelToAlerts(alerts models.PostableAlerts, labelName, labelValue string) models.PostableAlerts {
	for _, alert := range alerts {
		alert.Labels[labelName] = labelValue
	}
	return alerts
}
