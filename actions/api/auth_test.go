package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gobuffalo/suite/v4"
	"catter-cms/actions"
	"catter-cms/models"
)

// AuthSuite 认证测试套件
type AuthSuite struct {
	*suite.Action
}

// SetupTest 设置测试
func (as *AuthSuite) SetupTest() {
	suite.Action.SetupTest(as)
}

// TestLoginHandler 测试登录接口
func (as *AuthSuite) TestLoginHandler() {
	// 创建测试用户
	user := &models.User{
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // password: password
		Role:         "admin",
	}
	if err := models.DB.Create(user); err != nil {
		as.T().Fatal(err)
	}

	// 测试登录请求
	req := as.JSON(http.MethodPost, "/api/auth/login").SetJSON(map[string]string{
		"email":    "test@example.com",
		"password": "password",
	})

	res, err := req.Bytes()
	as.NoError(err)
	as.Equal(http.StatusOK, req.Code)

	// 验证响应
	var response LoginResponse
	if err := json.Unmarshal(res, &response); err != nil {
		as.T().Fatal(err)
	}

	as.NotEmpty(response.Token)
	as.Equal(user.Email, response.User.Email)
	as.Equal(user.Role, response.User.Role)
}

// TestLoginHandlerInvalidCredentials 测试登录接口 - 无效凭证
func (as *AuthSuite) TestLoginHandlerInvalidCredentials() {
	// 测试登录请求 - 无效密码
	req := as.JSON(http.MethodPost, "/api/auth/login").SetJSON(map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	})

	res, err := req.Bytes()
	as.NoError(err)
	as.Equal(http.StatusUnauthorized, req.Code)
}

// TestMeHandler 测试获取当前用户信息接口
func (as *AuthSuite) TestMeHandler() {
	// 创建测试用户
	user := &models.User{
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // password: password
		Role:         "admin",
	}
	if err := models.DB.Create(user); err != nil {
		as.T().Fatal(err)
	}

	// 先登录获取 token
	loginReq := as.JSON(http.MethodPost, "/api/auth/login").SetJSON(map[string]string{
		"email":    "test@example.com",
		"password": "password",
	})

	loginRes, err := loginReq.Bytes()
	as.NoError(err)
	as.Equal(http.StatusOK, loginReq.Code)

	var loginResponse LoginResponse
	if err := json.Unmarshal(loginRes, &loginResponse); err != nil {
		as.T().Fatal(err)
	}

	// 使用 token 测试获取用户信息
	meReq := as.GET("/api/auth/me").Set("Authorization", "Bearer "+loginResponse.Token)

	meRes, err := meReq.Bytes()
	as.NoError(err)
	as.Equal(http.StatusOK, meReq.Code)

	// 验证响应
	var meResponse models.User
	if err := json.Unmarshal(meRes, &meResponse); err != nil {
		as.T().Fatal(err)
	}

	as.Equal(user.Email, meResponse.Email)
	as.Equal(user.Role, meResponse.Role)
}

// TestMeHandlerUnauthorized 测试获取当前用户信息接口 - 未授权
func (as *AuthSuite) TestMeHandlerUnauthorized() {
	// 测试未授权请求
	req := as.GET("/api/auth/me")

	res, err := req.Bytes()
	as.NoError(err)
	as.Equal(http.StatusUnauthorized, req.Code)
}

// Test_AuthSuite 运行认证测试套件
func Test_AuthSuite(t *testing.T) {
	action, err := suite.NewAction(actions.App())
	if err != nil {
		t.Fatal(err)
	}

	as := &AuthSuite{
		Action: action,
	}
	suite.Run(t, as)
}
