package alipay

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"gozore-mall/common/utils"
	"time"
)

type AliPayService struct {
	Ctx          context.Context
	AppId        string
	PrivateKey   string
	CertPathApp  string
	CertPathRoot string
	CertPathAli  string
}

type AlipayOauthTokenResponse struct {
	ErrorResponse                  ErrorResponse                  `json:"error_response"`
	AlipaySystemOauthTokenResponse AlipaySystemOauthTokenResponse `json:"alipay_system_oauth_token_response"`
	Sign                           string                         `json:"sign"`
}

type AlipayUserInfoResponse struct {
	ErrorResponse               ErrorResponse               `json:"error_response"`
	AlipayUserInfoShareResponse AlipayUserInfoShareResponse `json:"alipay_user_info_share_response"`
	Sign                        string                      `json:"sign"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

type AlipayUserInfoShareResponse struct {
	Code                      string `json:"code"`
	Msg                       string `json:"msg"`
	UserNation                string `json:"user_nation"`
	IdentityCardAddress       string `json:"identity_card_address"`
	IdentityCardProvince      string `json:"identity_card_province"`
	IdentityCardCity          string `json:"identity_card_city"`
	IdentityCardArea          string `json:"identity_card_area"`
	UserId                    string `json:"user_id"`
	TaobaoId                  string `json:"taobao_id"`
	Avatar                    string `json:"avatar"`
	PersonCertIssueDate       string `json:"person_cert_issue_date"`
	Phone                     string `json:"phone"`
	MemberGrade               string `json:"member_grade"`
	CountryCode               string `json:"country_code"`
	City                      string `json:"city"`
	Area                      string `json:"area"`
	Address                   string `json:"address"`
	Zip                       string `json:"zip"`
	NickName                  string `json:"nick_name"`
	IsBalanceFrozen           string `json:"is_balance_frozen"`
	IsStudentCertified        string `json:"is_student_certified"`
	UserType                  string `json:"user_type"`
	UserStatus                string `json:"user_status"`
	Email                     string `json:"email"`
	UserName                  string `json:"user_name"`
	Mobile                    string `json:"mobile"`
	IsCertified               string `json:"is_certified"`
	CertType                  string `json:"cert_type"`
	Province                  string `json:"province"`
	PersonBirthdayWithoutYear string `json:"person_birthday_without_year"`
	PersonPictures            []struct {
		PictureUrl  string `json:"picture_url"`
		PictureType string `json:"picture_type"`
	} `json:"person_pictures"`
	CertNo               string `json:"cert_no"`
	Gender               string `json:"gender"`
	PersonBirthday       string `json:"person_birthday"`
	Profession           string `json:"profession"`
	PersonCertExpiryDate string `json:"person_cert_expiry_date"`
	LicenseNo            string `json:"license_no"`
	BusinessScope        string `json:"business_scope"`
	LicenseExpiryDate    string `json:"license_expiry_date"`
	OrganizationCode     string `json:"organization_code"`
	FirmPictures         []struct {
		PictureUrl  string `json:"picture_url"`
		PictureType string `json:"picture_type"`
	} `json:"firm_pictures"`
	FirmType                      string `json:"firm_type"`
	FirmLegalPersonName           string `json:"firm_legal_person_name"`
	FirmLegalPersonCertNo         string `json:"firm_legal_person_cert_no"`
	FirmLegalPersonCertType       string `json:"firm_legal_person_cert_type"`
	FirmLegalPersonCertExpiryDate string `json:"firm_legal_person_cert_expiry_date"`
	FirmLegalPersonPictures       []struct {
		PictureUrl  string `json:"picture_url"`
		PictureType string `json:"picture_type"`
	} `json:"firm_legal_person_pictures"`
	FirmAgentPersonName           string `json:"firm_agent_person_name"`
	FirmAgentPersonCertNo         string `json:"firm_agent_person_cert_no"`
	FirmAgentPersonCertType       string `json:"firm_agent_person_cert_type"`
	FirmAgentPersonCertExpiryDate string `json:"firm_agent_person_cert_expiry_date"`
	DeliverAddresses              []struct {
		DeliverMobile         string `json:"deliver_mobile"`
		DeliverPhone          string `json:"deliver_phone"`
		Address               string `json:"address"`
		Zip                   string `json:"zip"`
		DeliverProvince       string `json:"deliver_province"`
		DeliverCity           string `json:"deliver_city"`
		DeliverArea           string `json:"deliver_area"`
		AddressCode           string `json:"address_code"`
		DeliverFullname       string `json:"deliver_fullname"`
		DefaultDeliverAddress string `json:"default_deliver_address"`
	} `json:"deliver_addresses"`
	CollegeName        string `json:"college_name"`
	Degree             string `json:"degree"`
	EnrollmentTime     string `json:"enrollment_time"`
	GraduationTime     string `json:"graduation_time"`
	EntLicenseProvince string `json:"ent_license_province"`
	EntLicenseCity     string `json:"ent_license_city"`
	EntLicenseArea     string `json:"ent_license_area"`
	EntLicenseAddress  string `json:"ent_license_address"`
	DisplayName        string `json:"display_name"`
	InstOrCorp         string `json:"inst_or_corp"`
	Age                string `json:"age"`
	IsAdult            string `json:"is_adult"`
	IsBlocked          string `json:"is_blocked"`
}

type AlipaySystemOauthTokenResponse struct {
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ReExpiresIn  string `json:"re_expires_in"`
	AuthStart    string `json:"auth_start"`
}

func NewAliPayService(ctx context.Context, appId, privateKey, certPathApp, certPathRoot, certPathAli string) *AliPayService {
	return &AliPayService{
		Ctx:          ctx,
		AppId:        appId,
		PrivateKey:   privateKey,
		CertPathApp:  certPathApp,
		CertPathRoot: certPathRoot,
		CertPathAli:  certPathAli,
	}
}

func (s *AliPayService) NewAliPayClient() *alipay.Client {
	client, err := alipay.NewClient(s.AppId, s.PrivateKey, true)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	certPathApp := s.CertPathApp
	certPathRoot := s.CertPathRoot
	certPathAli := s.CertPathAli
	//配置公共参数
	err = client.
		SetCertSnByPath(certPathApp, certPathRoot, certPathAli)
	println(err)

	return client
}

func (s *AliPayService) Transfer(amount float32, identity, realName, orderTitle string) (*alipay.FundTransUniTransferResponse, error) {
	client := s.NewAliPayClient()
	param := gopay.BodyMap{
		"out_biz_no":   fmt.Sprintf("%s%d%s", "51"+time.Now().Format("060102150405"), time.Now().UnixNano()%10000, utils.GetRandom(5, 1)),
		"trans_amount": utils.ToString(amount),
		"product_code": "TRANS_ACCOUNT_NO_PWD",
		"biz_scene":    "DIRECT_TRANSFER",
		"order_title":  orderTitle,
		"payee_info": map[string]string{
			"identity":      identity,
			"identity_type": "ALIPAY_LOGON_ID",
			"name":          realName,
		},
	}
	//配置公共参数
	return client.FundTransUniTransfer(s.Ctx, param)
}

func (s *AliPayService) TransferByUsrId(amount float32, identity, orderTitle string) (*alipay.FundTransUniTransferResponse, error) {
	client := s.NewAliPayClient()
	param := gopay.BodyMap{
		"out_biz_no":   fmt.Sprintf("%s%d%s", "52"+time.Now().Format("060102150405"), time.Now().UnixNano()%10000, utils.GetRandom(5, 1)),
		"trans_amount": utils.ToString(amount),
		"product_code": "TRANS_ACCOUNT_NO_PWD",
		"biz_scene":    "DIRECT_TRANSFER",
		"order_title":  orderTitle,
		"payee_info": map[string]string{
			"identity":      identity,
			"identity_type": "ALIPAY_USER_ID",
		},
	}
	//配置公共参数
	return client.FundTransUniTransfer(s.Ctx, param)
}

//alipay web
func (s *AliPayService) AliTradePagePay(amount float64, OrderSn, subject, notifyUrl, returnUrl, timeoutExpress string) (string, error) {
	client := s.NewAliPayClient()
	client.SetNotifyUrl(notifyUrl)
	client.SetReturnUrl(returnUrl)
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", subject)
	body.Set("out_trade_no", OrderSn)
	body.Set("total_amount", amount)
	body.Set("product_code", "FAST_INSTANT_TRADE_PAY")
	if timeoutExpress != "" {
		body.Set("timeout_express", timeoutExpress)
	}
	//电脑网站支付请求
	payUrl, err := client.TradePagePay(s.Ctx, body)
	return payUrl, err
}

//alipay APP
func (s *AliPayService) AliTradeAppPay(amount float64, OrderSn, subject, notifyUrl, returnUrl, timeoutExpress string) (string, error) {
	client := s.NewAliPayClient()
	client.SetNotifyUrl(notifyUrl)
	client.SetReturnUrl(returnUrl)
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", subject)
	body.Set("out_trade_no", OrderSn)
	body.Set("total_amount", amount)
	if timeoutExpress != "" {
		body.Set("timeout_express", timeoutExpress)
	}
	//电脑网站支付请求
	payParam, err := client.TradeAppPay(s.Ctx, body)
	return payParam, err
}

//alipay wap
func (s *AliPayService) AliTradeWapPay(amount float64, OrderSn, subject, notifyUrl, returnUrl, quitUrl, timeoutExpress string) (string, error) {
	client := s.NewAliPayClient()
	client.SetNotifyUrl(notifyUrl)
	client.SetReturnUrl(returnUrl)
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", subject)
	body.Set("out_trade_no", OrderSn)
	body.Set("total_amount", amount)
	body.Set("quit_url", quitUrl)
	body.Set("product_code", "QUICK_WAP_WAY")
	if timeoutExpress != "" {
		body.Set("timeout_express", timeoutExpress)
	}
	//电脑网站支付请求
	payUrl, err := client.TradeWapPay(s.Ctx, body)
	return payUrl, err
}

//退款操作
func (s *AliPayService) AliPayRefund(refundAmount float64, orderSn, refundSn, notifyUrl, returnUrl string) (*alipay.TradeRefundResponse, error) {
	client := s.NewAliPayClient()
	client.SetNotifyUrl(notifyUrl)
	client.SetReturnUrl(returnUrl)
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", orderSn)
	body.Set("out_request_no", refundSn)
	body.Set("refund_amount", refundAmount)
	return client.TradeRefund(s.Ctx, body)
}

//支付宝  查询订单详情
func (s *AliPayService) AliTradeQuery(orderSn, notifyUrl, returnUrl string) (*alipay.TradeQueryResponse, error) {
	client := s.NewAliPayClient()
	client.SetNotifyUrl(notifyUrl)
	client.SetReturnUrl(returnUrl)
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", orderSn)

	//电脑网站支付请求
	resp, err := client.TradeQuery(s.Ctx, body)
	return resp, err
}

//换取授权访问令牌
func (s *AliPayService) OauthToken(ctx context.Context, code string) (*AlipaySystemOauthTokenResponse, error) {
	client := s.NewAliPayClient()
	bm := make(gopay.BodyMap)

	// 自定义公共参数（根据自己需求，需要独立设置的自行设置，不需要单独设置的，共享client的配置）
	bm.Set("grant_type", "authorization_code")
	bm.Set("code", code)

	aliPsp := new(AlipayOauthTokenResponse)
	err := client.PostAliPayAPISelfV2(ctx, bm, "alipay.system.oauth.token", aliPsp)
	if err != nil {
		return nil, err
	}
	if len(aliPsp.ErrorResponse.Code) > 0 {
		return nil, errors.New(aliPsp.ErrorResponse.SubMsg)
	}
	return &aliPsp.AlipaySystemOauthTokenResponse, nil
}

//支付宝会员授权信息查询接口
func (s *AliPayService) Token2UserInfo(ctx context.Context, token string) (*AlipayUserInfoShareResponse, error) {
	client := s.NewAliPayClient()
	bm := make(gopay.BodyMap)

	// 自定义公共参数（根据自己需求，需要独立设置的自行设置，不需要单独设置的，共享client的配置）
	bm.Set("auth_token", token)

	aliPsp := new(AlipayUserInfoResponse)
	err := client.PostAliPayAPISelfV2(ctx, bm, "alipay.user.info.share", aliPsp)
	if err != nil {
		return nil, err
	}
	if len(aliPsp.ErrorResponse.Code) > 0 {
		return nil, errors.New(aliPsp.ErrorResponse.SubMsg)
	}
	return &aliPsp.AlipayUserInfoShareResponse, nil
}

//用户登录授权
func (s *AliPayService) OauthCode(ctx context.Context, state string) (string, error) {
	client := s.NewAliPayClient()
	bm := make(gopay.BodyMap)

	// 自定义公共参数（根据自己需求，需要独立设置的自行设置，不需要单独设置的，共享client的配置）
	bm.Set("scopes", "auth_base")
	bm.Set("state", state)

	htmlStr := ""
	err := client.PostAliPayAPISelfV2(ctx, bm, "alipay.user.info.auth", htmlStr)
	if err != nil {
		return "", err
	}
	return htmlStr, nil
}
