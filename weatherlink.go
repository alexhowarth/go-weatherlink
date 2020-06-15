// Package weatherlink provides a client to the Davis Weatherlink weather station API
package weatherlink

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

const apiURL string = "api.weatherlink.com"
const apiVersion string = "v2"

const (
	keyParam    string = "api-key"
	secretParam string = "api-secret"
	sigParam    string = "api-signature"
	tParam      string = "t"
)

const (
	stationsPathFmt   string = "/stations/%v"
	sensorsPathFmt    string = "/sensors/%v"
	sensorCatalogPath string = "/sensor-catalog"
	currentPathFmt    string = "/current/%v"
	historicPathFmt   string = "/historic/%v?start-timestamp=%v&end-timestamp=%v"
)

type Config struct {
	Client *http.Client
	Key    string
	Secret string
}

type Client struct {
	client *http.Client
	config *Config
}

type SignatureParams map[string]string

// NewClient returns a WeatherLink client for interacting with the API
func (conf *Config) NewClient() *Client {
	if conf.Key == "" || conf.Secret == "" {
		panic("Key and Secret required.")
	}
	wl := &Client{
		client: &http.Client{},
		config: conf,
	}
	if conf.Client != nil {
		wl.client = conf.Client
	}
	return wl
}

// MakeSignatureParams contains the common signature parameters
func (w *Client) MakeSignatureParams() SignatureParams {
	p := make(SignatureParams)
	p[keyParam] = w.config.Key
	p[tParam] = strconv.FormatInt(time.Now().Unix(), 10)
	return p
}

func (w *Client) BuildURL(s string, p SignatureParams) string {

	u, err := url.Parse(s)
	if err != nil {
		println(err)
	}

	if u.Path == "" {
		panic("Path required.")
	}

	u.Scheme = "https"
	u.Host = apiURL
	u.Path = path.Join(apiVersion + u.Path)

	q := u.Query()
	for k := range q {
		p[k] = q[k][0]
	}

	q = u.Query()
	q.Add(sigParam, p.Signature(w.config.Secret))
	q.Add(keyParam, w.config.Key)
	q.Add(tParam, p[tParam])
	u.RawQuery = q.Encode()

	return u.String()
}

// Add a kv parameter to the signature
func (s SignatureParams) Add(key string, value string) {
	s[key] = value
	return
}

// Signature encodes and returns the HMAC hexidecimal string
func (s SignatureParams) Signature(secret string) string {
	return encode(secret, s.String())
}

// String is the unencoded signature
func (s SignatureParams) String() string {
	var buf strings.Builder
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		buf.WriteString(k)
		buf.WriteString(s[k])
	}
	return buf.String()
}

func (w *Client) get(url string, params SignatureParams) (*http.Response, error) {
	if params == nil {
		params = w.MakeSignatureParams()
	}
	return w.client.Get(w.BuildURL(url, params))
}

// encode returns an hexidecimal HMAC string (used for the signature)
func encode(secret string, msg string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}

// StationsResponse represents data from the /stations endpoint
type StationsResponse struct {
	Stations []struct {
		StationID           int     `json:"station_id"`
		StationName         string  `json:"station_name"`
		GatewayID           int     `json:"gateway_id"`
		GatewayIDHex        string  `json:"gateway_id_hex"`
		ProductNumber       string  `json:"product_number"`
		Username            string  `json:"username"`
		UserEmail           string  `json:"user_email"`
		CompanyName         string  `json:"company_name"`
		Active              bool    `json:"active"`
		Private             bool    `json:"private"`
		RecordingInterval   int     `json:"recording_interval"`
		FirmwareVersion     string  `json:"firmware_version"`
		Meid                string  `json:"meid"`
		RegisteredDate      int     `json:"registered_date"`
		SubscriptionEndDate int     `json:"subscription_end_date"`
		TimeZone            string  `json:"time_zone"`
		City                string  `json:"city"`
		Region              string  `json:"region"`
		Country             string  `json:"country"`
		Latitude            float64 `json:"latitude"`
		Longitude           float64 `json:"longitude"`
		Elevation           float64 `json:"elevation"`
	} `json:"stations"`
	GeneratedAt int `json:"generated_at"`
}

// AllStations gets all weather stations associated with your API Key
func (w *Client) AllStations() (sr StationsResponse, err error) {
	return w.Stations(nil)
}

// Stations gets weather stations for one or more station IDs provided
func (w *Client) Stations(stations []int) (sr StationsResponse, err error) {

	csv := intArrToCSV(stations)

	sp := w.MakeSignatureParams()
	if stations != nil {
		sp.Add("station-ids", csv)
	}

	resp, err := w.get(fmt.Sprintf(stationsPathFmt, csv), sp)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error making Stations request. Got status: %v", resp.Status)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&sr)
	if err != nil {
		return
	}

	return sr, nil
}

// SensorsResponse represents data from the /sendors endpoint
type SensorsResponse struct {
	Sensors []struct {
		Lsid              int         `json:"lsid"`
		SensorType        int         `json:"sensor_type"`
		Category          string      `json:"category"`
		Manufacturer      string      `json:"manufacturer"`
		ProductName       string      `json:"product_name"`
		ProductNumber     string      `json:"product_number"`
		RainCollectorType int         `json:"rain_collector_type"`
		Active            bool        `json:"active"`
		CreatedDate       int         `json:"created_date"`
		ModifiedDate      int         `json:"modified_date"`
		StationID         int         `json:"station_id"`
		StationName       string      `json:"station_name"`
		ParentDeviceType  string      `json:"parent_device_type"`
		ParentDeviceName  string      `json:"parent_device_name"`
		ParentDeviceID    int         `json:"parent_device_id"`
		ParentDeviceIDHex string      `json:"parent_device_id_hex"`
		PortNumber        int         `json:"port_number"`
		Latitude          float64     `json:"latitude"`
		Longitude         float64     `json:"longitude"`
		Elevation         float64     `json:"elevation"`
		TxID              interface{} `json:"tx_id"`
	} `json:"sensors"`
	GeneratedAt int `json:"generated_at"`
}

// AllSensors gets all sensors attached to all weather stations associated with your API Key
func (w *Client) AllSensors() (sr SensorsResponse, err error) {
	return w.Sensors(nil)
}

// Sensors gets sensors for one or more sensor IDs provided
func (w *Client) Sensors(sensors []int) (sr SensorsResponse, err error) {

	csv := intArrToCSV(sensors)

	sp := w.MakeSignatureParams()
	if sensors != nil {
		sp.Add("sensor-ids", csv)
	}

	resp, err := w.get(fmt.Sprintf(sensorsPathFmt, csv), sp)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error making Sensors request. Got status: %v", resp.Status)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&sr)
	if err != nil {
		return
	}

	return sr, nil
}

// CurrentResponse represents data from the /current endpoint
type CurrentResponse struct {
	StationID int `json:"station_id"`
	Sensors   []struct {
		Lsid              int `json:"lsid"`
		SensorType        int `json:"sensor_type"`
		DataStructureType int `json:"data_structure_type"`
		Data              []struct {
			Ts                int         `json:"ts"`
			BarTrend          float64     `json:"bar_trend"`
			Bar               float64     `json:"bar"`
			TempIn            float64     `json:"temp_in"`
			HumIn             float64     `json:"hum_in"`
			TempOut           float64     `json:"temp_out"`
			WindSpeed         float64     `json:"wind_speed"`
			WindSpeed10MinAvg float64     `json:"wind_speed_10_min_avg"`
			WindDir           float64     `json:"wind_dir"`
			TempExtra1        interface{} `json:"temp_extra_1"`
			TempExtra2        interface{} `json:"temp_extra_2"`
			TempExtra3        interface{} `json:"temp_extra_3"`
			TempExtra4        interface{} `json:"temp_extra_4"`
			TempExtra5        interface{} `json:"temp_extra_5"`
			TempExtra6        interface{} `json:"temp_extra_6"`
			TempExtra7        interface{} `json:"temp_extra_7"`
			TempSoil1         interface{} `json:"temp_soil_1"`
			TempSoil2         interface{} `json:"temp_soil_2"`
			TempSoil3         interface{} `json:"temp_soil_3"`
			TempSoil4         interface{} `json:"temp_soil_4"`
			TempLeaf1         interface{} `json:"temp_leaf_1"`
			TempLeaf2         interface{} `json:"temp_leaf_2"`
			TempLeaf3         interface{} `json:"temp_leaf_3"`
			TempLeaf4         interface{} `json:"temp_leaf_4"`
			HumOut            float64     `json:"hum_out"`
			HumExtra1         interface{} `json:"hum_extra_1"`
			HumExtra2         interface{} `json:"hum_extra_2"`
			HumExtra3         interface{} `json:"hum_extra_3"`
			HumExtra4         interface{} `json:"hum_extra_4"`
			HumExtra5         interface{} `json:"hum_extra_5"`
			HumExtra6         interface{} `json:"hum_extra_6"`
			HumExtra7         interface{} `json:"hum_extra_7"`
			RainRateClicks    float64     `json:"rain_rate_clicks"`
			RainRateIn        float64     `json:"rain_rate_in"`
			RainRateMm        float64     `json:"rain_rate_mm"`
			Uv                interface{} `json:"uv"`
			SolarRad          interface{} `json:"solar_rad"`
			RainStormClicks   float64     `json:"rain_storm_clicks"`
			RainStormIn       float64     `json:"rain_storm_in"`
			RainStormMm       float64     `json:"rain_storm_mm"`
			RainDayClicks     float64     `json:"rain_day_clicks"`
			RainDayIn         float64     `json:"rain_day_in"`
			RainDayMm         float64     `json:"rain_day_mm"`
			RainMonthClicks   float64     `json:"rain_month_clicks"`
			RainMonthIn       float64     `json:"rain_month_in"`
			RainMonthMm       float64     `json:"rain_month_mm"`
			RainYearClicks    float64     `json:"rain_year_clicks"`
			RainYearIn        float64     `json:"rain_year_in"`
			RainYearMm        float64     `json:"rain_year_mm"`
			EtDay             float64     `json:"et_day"`
			EtMonth           float64     `json:"et_month"`
			EtYear            float64     `json:"et_year"`
			MoistSoil1        interface{} `json:"moist_soil_1"`
			MoistSoil2        interface{} `json:"moist_soil_2"`
			MoistSoil3        interface{} `json:"moist_soil_3"`
			MoistSoil4        interface{} `json:"moist_soil_4"`
			WetLeaf1          interface{} `json:"wet_leaf_1"`
			WetLeaf2          interface{} `json:"wet_leaf_2"`
			WetLeaf3          interface{} `json:"wet_leaf_3"`
			WetLeaf4          interface{} `json:"wet_leaf_4"`
		} `json:"data"`
	} `json:"sensors"`
	GeneratedAt int `json:"generated_at"`
}

// Current gets current conditions data for one station
func (w *Client) Current(station int) (cr CurrentResponse, err error) {

	sp := w.MakeSignatureParams()
	sp.Add("station-id", strconv.Itoa(station))

	resp, err := w.get(fmt.Sprintf(currentPathFmt, station), sp)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error making Current request. Got status: %v", resp.Status)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&cr)
	if err != nil {
		return
	}

	return cr, nil
}

// HistoricResponse represents historic data for one station ID within a given timerange
type HistoricResponse struct {
	Sensors []struct {
		Lsid int `json:"lsid"`
		Data []struct {
			Ts               int     `json:"ts"`
			ArchInt          int     `json:"arch_int"`
			RevType          int     `json:"rev_type"`
			TempOut          float64 `json:"temp_out"`
			TempOutHi        float64 `json:"temp_out_hi"`
			TempOutLo        float64 `json:"temp_out_lo"`
			TempIn           float64 `json:"temp_in"`
			HumIn            float64 `json:"hum_in"`
			HumOut           float64 `json:"hum_out"`
			RainfallIn       float64 `json:"rainfall_in"`
			RainfallClicks   float64 `json:"rainfall_clicks"`
			RainfallMm       float64 `json:"rainfall_mm"`
			RainRateHiIn     float64 `json:"rain_rate_hi_in"`
			RainRateHiClicks float64 `json:"rain_rate_hi_clicks"`
			RainRateHiMm     float64 `json:"rain_rate_hi_mm"`
			Et               float64 `json:"et"`
			Bar              float64 `json:"bar"`
			WindNumSamples   float64 `json:"wind_num_samples"`
			WindSpeedAvg     float64 `json:"wind_speed_avg"`
			WindSpeedHi      float64 `json:"wind_speed_hi"`
			WindDirOfHi      float64 `json:"wind_dir_of_hi"`
			WindDirOfPrevail float64 `json:"wind_dir_of_prevail"`
			ForecastRule     float64 `json:"forecast_rule"`
			AbsPress         float64 `json:"abs_press"`
			BarNoaa          float64 `json:"bar_noaa"`
			DewPointOut      float64 `json:"dew_point_out"`
			DewPointIn       float64 `json:"dew_point_in"`
			Emc              float64 `json:"emc"`
			HeatIndexOut     float64 `json:"heat_index_out"`
			HeatIndexIn      float64 `json:"heat_index_in"`
			WindChill        float64 `json:"wind_chill"`
			WindRun          float64 `json:"wind_run"`
			DegDaysHeat      float64 `json:"deg_days_heat"`
			DegDaysCool      float64 `json:"deg_days_cool"`
			ThwIndex         float64 `json:"thw_index"`
			WetBulb          float64 `json:"wet_bulb"`
		} `json:"data"`
		SensorType        int `json:"sensor_type"`
		DataStructureType int `json:"data_structure_type"`
	} `json:"sensors"`
	GeneratedAt int `json:"generated_at"`
	StationID   int `json:"station_id"`
}

// Historic gets historic data for one station ID within a given timerange
func (w *Client) Historic(station int, start time.Time, end time.Time) (hr HistoricResponse, err error) {

	sp := w.MakeSignatureParams()
	sp.Add("station-id", strconv.Itoa(station))

	resp, err := w.get(fmt.Sprintf(historicPathFmt, station, start.Unix(), end.Unix()), sp)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error making Historic request. Got status: %v", resp.Status)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&hr)
	if err != nil {
		return
	}

	return hr, nil
}

// SensorCatalog saves a catalogue of all types of sensors to file
func (w *Client) SensorCatalog(path string) (err error) {

	resp, err := w.get(sensorCatalogPath, nil)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error making SensorCatalog request. Got status: %v", resp.Status)
		return
	}

	out, err := os.Create(path)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}

	return
}

func intArrToCSV(i []int) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(i)), ","), "[]")
}

func dumpResponse(resp *http.Response) {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n\nResponse: %q", dump)
}
