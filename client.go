package awqatsalah

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	// The Base URL to send requests to. In regular usage, this is
	// https://awqatsalah.diyanet.gov.tr/
	BaseUrl string

	// The credentials we are sending to get access token
	Credentials Credentials

	accessToken string

	// Net/http Client for contacting the Awqat Salah API.
	c *http.Client
}

type Credentials struct {
	Email    string
	Password string
}

type AwqatResponse[TData any] struct {
	Data    *TData `json:"data"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (err *AwqatResponse[nil]) Error() string {
	if err.Message == "" {
		return fmt.Sprintf("%d API error: %s", err.Message)
	}
	return fmt.Sprintf("%d (%s) API error: %s", err.Success, err, err.Message)
}

type AuthResponse struct {
	AccessToken  string
	RefreshToken string
}

type Location struct {
	ID   int
	Code string
	Name string
}

type DailyContent struct {
	ID           int
	DayOfYear    int
	Verse        string
	VerseSource  string
	Hadith       string
	HadithSource string
	Pray         string
	PraySource   string
}

type PrayerTime struct {
	ShapeMoonUrl              string
	Fajr                      string
	Sunrise                   string
	Dhuhr                     string
	Asr                       string
	Maghrib                   string
	Isha                      string
	AstronomicalSunset        string
	AstronomicalSunrise       string
	HijriDateShort            string
	HijriDateShortIso8601     string
	HijriDateLongIso8601      string
	HijriDateLong             string
	QiblaTime                 string
	GregorianDateShort        string
	GregorianDateShortIso8601 string
	GregorianDateLong         string
	GregorianDateLongIso8601  string
}

type PrayerTimeEid struct {
	EidAlAdhaHijri string
	EidAlAdhaTime  string
	EidAlAdhaDate  string
	EidAlFitrHijri string
	EidAlFitrTime  string
	EidAlFitrDate  string
}

func newDefaultHTTPClient() *http.Client { return &http.Client{Timeout: time.Minute} }

func New(credentials Credentials) (*Client, error) {
	c := &Client{
		BaseUrl:     "https://awqatsalah.diyanet.gov.tr/",
		Credentials: credentials,
		c:           newDefaultHTTPClient(),
	}

	err := c.auth()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) auth() error {
	res := &AwqatResponse[AuthResponse]{}
	if err := c.execute(http.MethodPost, "auth/login", "", c.Credentials, res); err != nil {
		return err
	}

	c.accessToken = res.Data.AccessToken

	return nil
}

func (c *Client) Countries() ([]Location, error) {
	res := &AwqatResponse[[]Location]{}
	if err := c.execute(http.MethodGet, "api/place/countries", "", nil, &res); err != nil {
		return nil, err
	}
	return *res.Data, nil
}

func (c *Client) States() ([]Location, error) {
	res := &AwqatResponse[[]Location]{}
	if err := c.execute(http.MethodGet, "api/place/states", "", nil, &res); err != nil {
		return nil, err
	}
	return *res.Data, nil
}

func (c *Client) Cities() ([]Location, error) {
	res := &AwqatResponse[[]Location]{}
	if err := c.execute(http.MethodGet, "api/place/cities", "", nil, &res); err != nil {
		return nil, err
	}
	return *res.Data, nil
}

func (c *Client) StatesByCountryID(countryID string) ([]Location, error) {
	res := &AwqatResponse[[]Location]{}
	if err := c.execute(http.MethodGet, "api/place/states", countryID, nil, &res); err != nil {
		return nil, err
	}
	return *res.Data, nil
}

func (c *Client) CitiesByStateID(stateID string) ([]Location, error) {
	res := &AwqatResponse[[]Location]{}
	if err := c.execute(http.MethodGet, "api/place/cities/", stateID, nil, &res); err != nil {
		return nil, err
	}
	return *res.Data, nil
}

func (c *Client) DailyContent() (*DailyContent, error) {
	res := &AwqatResponse[DailyContent]{}
	if err := c.execute(http.MethodGet, "api/DailyContent", "", nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) PrayerTimeDailyByCityID(cityID string) ([]PrayerTime, error) {
	res := &AwqatResponse[[]PrayerTime]{}
	if err := c.execute(http.MethodGet, "api/PrayerTime/Daily", cityID, nil, &res); err != nil {
		return nil, err
	}
	return *res.Data, nil
}

func (c *Client) PrayerTimeWeeklyByCityID(cityID string) ([]PrayerTime, error) {
	res := &AwqatResponse[[]PrayerTime]{}
	if err := c.execute(http.MethodGet, "api/PrayerTime/Weekly", cityID, nil, &res); err != nil {
		return nil, err
	}
	return *res.Data, nil
}

func (c *Client) PrayerTimeMonthlyByCityID(cityID string) ([]PrayerTime, error) {
	res := &AwqatResponse[[]PrayerTime]{}
	if err := c.execute(http.MethodGet, "api/PrayerTime/Monthly", cityID, nil, &res); err != nil {
		return nil, err
	}
	return *res.Data, nil
}

func (c *Client) PrayerTimeEidByCityID(cityID string) (*PrayerTimeEid, error) {
	res := &AwqatResponse[PrayerTimeEid]{}
	if err := c.execute(http.MethodGet, "api/PrayerTime/Eid", cityID, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) PrayerTimeRamadanByCityID(cityID string) ([]PrayerTime, error) {
	res := &AwqatResponse[[]PrayerTime]{}
	if err := c.execute(http.MethodGet, "api/PrayerTime/Ramadan", cityID, nil, &res); err != nil {
		return nil, err
	}
	return *res.Data, nil
}

func (c *Client) execute(method string, endp string, path string, body interface{}, expectedResponse interface{}) error {
	u, err := url.Parse(c.BaseUrl)
	if err != nil {
		return err
	}
	u = u.ResolveReference(&url.URL{Path: endp})

	var req *http.Request

	if method == http.MethodPost {
		bodyJson, err := json.Marshal(body)

		req, err = http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(bodyJson))
		if err != nil {
			return err
		}
	} else if method == http.MethodGet {
		req, err = http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return err
		}

		req.URL = req.URL.JoinPath(path)
	}

	if len(c.accessToken) > 0 {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		if err := json.NewDecoder(res.Body).Decode(expectedResponse); err != nil {
			return err
		}

		return nil
	case 400, 401, 403, 404, 415, 500:
		var errRes AwqatResponse[interface{}]
		if err := json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			return err
		}

		return &errRes
	default:
		return fmt.Errorf("unexpected HTTP response status code: %d", res.StatusCode)
	}
}
