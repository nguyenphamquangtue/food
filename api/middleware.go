package api

import (
	"fmt"
	"food/model"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// authentication middleware
func authen(c *gin.Context) {
	ip := c.Request.Header.Get("ip")
	user, err := model.GetUserByIp(ip)
	if (err != nil) || ((user.Role != model.RoleAdmin) && (user.Role != model.RoleScan)) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	c.Next()
}

// get request ip middleware
func ip(c *gin.Context) {
	ip, _ := getRequestIP(c)
	c.Request.Header.Set("ip", ip)
	c.Next()
}

func getRequestIP(c *gin.Context) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := c.Request.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := c.Request.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	return "", fmt.Errorf("No valid ip found")
}
