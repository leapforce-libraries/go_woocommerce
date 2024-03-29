package woocommerce

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

const (
	apiName          string = "WooCommerce"
	apiPath          string = "wp-json/wc/v2"
	totalPagesHeader string = "X-WP-TotalPages"
	DateFormat       string = "2006-01-02T15:04:05"
)

// type
//
type Service struct {
	host        string
	token       string
	httpService *go_http.Service
}

type ServiceConfig struct {
	Host           string
	ConsumerKey    string
	ConsumerSecret string
}

func NewService(config *ServiceConfig) (*Service, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if config.Host == "" {
		return nil, errortools.ErrorMessage("Host not provided")
	}

	if config.ConsumerKey == "" {
		return nil, errortools.ErrorMessage("ConsumerKey not provided")
	}

	if config.ConsumerSecret == "" {
		return nil, errortools.ErrorMessage("ConsumerSecret not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		host:        config.Host,
		token:       base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", config.ConsumerKey, config.ConsumerSecret))),
		httpService: httpService,
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add authentication header
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Basic %s", service.token))
	(*requestConfig).NonDefaultHeaders = &header

	// add error model
	errorResponse := ErrorResponse{}
	(*requestConfig).ErrorModel = &errorResponse

	request, response, e := service.httpService.HttpRequest(requestConfig)
	if errorResponse.Message != "" {
		e.SetMessage(errorResponse.Message)
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s/%s", service.host, apiPath, path)
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.token
}

func (service *Service) ApiCallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) ApiReset() {
	service.httpService.ResetRequestCount()
}

func UIntArrayToString(unints []uint) string {
	ids := []string{}
	for _, include := range unints {
		ids = append(ids, fmt.Sprintf("%v", include))
	}

	return strings.Join(ids, ",")
}

func TotalPages(response *http.Response) (int, *errortools.Error) {
	if response == nil {
		return 0, nil
	}

	totalPages, err := strconv.Atoi(response.Header.Get(totalPagesHeader))
	if err != nil {
		return 0, errortools.ErrorMessage(fmt.Sprintf("Error while retrieving %s header (%s)", totalPagesHeader, err.Error()))
	}

	return totalPages, nil
}
