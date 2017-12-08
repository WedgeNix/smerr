package sms

import (
	"fmt"
	"net/smtp"
	"strconv"
	"sync"
)

func toString(number uint64) string {
	return strconv.FormatUint(number, 10)
}

func Verizon(number uint64) Gateway {
	return Gateway(toString(number) + "@vtext.com")
}

func ATT(number uint64) Gateway {
	return Gateway(toString(number) + "@txt.att.net")
}

func TMobile(number uint64) Gateway {
	return Gateway(toString(number) + "@tmomail.net")
}

func Sprint(number uint64) Gateway {
	return Gateway(toString(number) + "@messaging.sprintpcs.com")
}

func MetroPCS(number uint64) Gateway {
	return Gateway(toString(number) + "@mymetropcs.com")
}

func BoostMobile(number uint64) Gateway {
	return Gateway(toString(number) + "@myboostmobile.com")
}

func VirginMobile(number uint64) Gateway {
	return Gateway(toString(number) + "@vmobl.com")
}

func Cricket(number uint64) Gateway {
	return Gateway(toString(number) + "@mms.mycricket.com")
}

func USCellular(number uint64) Gateway {
	return Gateway(toString(number) + "@email.uscc.net")
}

var carriers = map[string][]string{
	"7-11 Speakout (USA GSM)":              {"number@cingularme.com"},
	"Alaska Communications Systems":        {"number@msg.acsalaska.com"},
	"Alltel Wireless":                      {"number@message.alltel.com"},
	"AT&T Mobility (formerly Cingular)":    {"number@mms.att.net", "number@txt.att.net", "number@mmode.com", "number@cingularme.com"},
	"Bell Mobility & Solo Mobile (Canada)": {"number@txt.bell.ca"},
	"Boost Mobile":                         {"number@myboostmobile.com"},
	"Cellular One (Dobson)":                {"number@mobile.celloneusa.com"},
	"Cingular (Postpaid)":                  {"number@cingularme.com"},
	"Centennial Wireless":                  {"number@cwemail.com"},
	"Cingular (GoPhone prepaid)":           {"number@cingularme.com"},
	"Claro (Nicaragua)":                    {"number@ideasclaro-ca.com"},
	"Comcel":                               {"number@comcel.com.co"},
	"Cricket":                              {"number@mms.mycricket.com"},
	"CTI":                                  {"number@sms.ctimovil.com.ar"},
	"Emtel (Mauritius)":                    {"number@emtelworld.net"},
	"Fido (Canada)":                        {"number@fido.ca"},
	"Globalstar":                           {"number@msg.globalstarusa.com"},
	"Helio":                                {"number@messaging.sprintpcs.com"},
	"Illinois Valley Cellular":             {"number@ivctext.com"},
	"IT Company Australia":                 {"number@itcompany.com.au"},
	"Iridium (satellite)":                  {"number@msg.iridium.com"},
	"Meteor (Ireland)":                     {"number@sms.mymeteor.ie"},
	"MetroPCS":                             {"number@mymetropcs.com"},
	"Movicom":                              {"number@movimensaje.com.ar"},
	"Movistar (Colombia)":                  {"number@movistar.com.co"},
	"MTN (South Africa)":                   {"number@sms.co.za"},
	"MTS (Canada)":                         {"number@text.mtsmobility.com"},
	"Nextel (Argentina)":                   {"TwoWay.11number@nextel.net.ar"},
	"Personal (Argentina)":                 {"11number@personal-net.com.ar"},
	"Plus GSM (Poland)":                    {"+48number@text.plusgsm.pl"},
	"President's Choice (Canada)":          {"number@txt.bell.ca"},
	"Qwest":                      {"number@qwestmp.com"},
	"Rogers (Canada)":            {"number@pcs.rogers.com"},
	"Sasktel (Canada)":           {"number@sms.sasktel.com"},
	"Setar Mobile email (Aruba)": {"297+number@mas.aw"},
	"SMSGlobal":                  {"number@sms.smsglobal.com.au"},
	"Sprint (PCS)":               {"number@messaging.sprintpcs.com", "number@pm.sprint.com"},
	"Sprint (Nextel)":            {"number@page.nextel.com", "number@messaging.nextel.com"},
	"Suncom":                     {"number@tms.suncom.com"},
	"T-Mobile":                   {"number@tmomail.net"},
	"T-Mobile (Austria)":         {"number@sms.t-mobile.at"},
	"Telus Mobility (Canada)":    {"number@msg.telus.com"},
	"Tigo (Formerly Ola)":        {"number@sms.tigo.com.co"},
	"Tracfone (prepaid)":         {"number@cingularme.com", "number@tmomail.net", "number@vtext.com", "number@email.uscc.net", "number@message.alltel.com"},
	"Unicel":                     {"number@utext.com"},
	"US Cellular":                {"number@email.uscc.net", "number@mms.uscc.net"},
	"Verizon":                    {"number@vtext.com", "number@vzwpix.com"},
	"Virgin Mobile (Canada)":     {"number@vmobile.ca"},
	"Virgin Mobile (USA)":        {"number@vmobl.com"},
	"Vodacom (South Africa)":     {"number@voda.co.za"},
	"YCC": {"number@sms.ycc.ru"},
	"B2sms (International)":         {"number@b2sms.com"},
	"CardBoardFish (International)": {"number@username.etexting.com"},
	"Club4sms (Pakistan)":           {"number@club4sms.com"},
	"Esendex (AU, ES, FR, IE, UK)":  {"number@esendex.net"},
	"Ipipi.com":                     {"number@opensms.ipipi.com"},
	"Kapow! SMS Gateway":            {"number@kapow.co.uk"},
	"Letxt (International)":         {"number@sms.letxt.com.au"},
	"Me2mobile (Australia)":         {"number@me2mobile.com"},
	"Mobe.Net":                      {"number@mobe.net"},
	"pktpix.com (International)":    {"number@pktpix.com"},
	"Red Oxygen (International)":    {"number@redoxygen.net"},
	"Soprano (Australia)":           {"number@soprano.com.au"},
	"TellusTalk":                    {"number@esms.nu"},
	"ToText.net":                    {"number@totext.net"},
	"Txtlocal.com":                  {"number@txtlocal.co.uk"},
	"ViaNett":                       {"number@sms.vianett.no"},
	"Webtext":                       {"number@webtext.com"},
	"ABTXT.COM":                     {"number@abtxt.com"},
	"MOBILEMAIL.RU":                 {"number@mobilemail.ru"},
}

var (
	Email string

	Pass string

	Port = 25

	Host string

	sms SMS
)

// Gateway is an SMS email wrapper for a phone number.
type Gateway string

// Errors stores a mapping of gateways to a specific error.
type Errors map[Gateway]error

// Gateways pulls a slice of numbers from the map.
func (e Errors) Gateways() []Gateway {
	var gates []Gateway
	for gate := range e {
		gates = append(gates, gate)
	}
	return gates
}

// SMS controls repeated SMS error messaging.
type SMS struct {
	once sync.Once

	Email string
	Pass  string
	Host  string
	Port  int
}

func (s *SMS) prep() {
	if len(s.Email) == 0 {
		s.Email = Email
	}
	if len(s.Pass) == 0 {
		s.Pass = Pass
	}
	if len(s.Host) == 0 {
		s.Host = Host
	}
	if s.Port == 0 {
		s.Port = Port
	}
}

// Text sends a text message to all the numbers.
func Text(v interface{}, g ...Gateway) Errors {
	return sms.Text(v, g...)
}

// Text sends a text message to all the numbers.
func (s *SMS) Text(v interface{}, g ...Gateway) Errors {
	s.once.Do(s.prep)

	auth := smtp.PlainAuth(
		"",
		s.Email,
		s.Pass,
		s.Host,
	)

	errs := Errors{}

	for _, gate := range g {
		email := string(gate)

		if err := smtp.SendMail(
			s.Host+":"+strconv.Itoa(s.Port),
			auth,
			s.Email,
			[]string{email},
			[]byte(fmt.Sprintf(`From: %s\nTo: %s\n\n%s\n`, s.Email, email, fmt.Sprint(v))),
		); err != nil {
			errs[gate] = err
		}
	}

	return errs
}

// package sms

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"sync"
// )

// const (
// 	// EnvSID is the env name for storing a Twilio SID.
// 	EnvSID = "TWILIO_SID"

// 	// EnvAuth is the env name for storing a Twilio auth token.
// 	EnvAuth = "TWILIO_AUTH"
// )

// var sms SMS

// type SMS struct {
// 	once      sync.Once
// 	sid, auth string
// 	url       *url.URL
// }

// type response struct {
// 	SID                 string  `json:"sid"`
// 	DateCreated         string  `json:"date_created"`
// 	DateUpdated         string  `json:"date_updated"`
// 	DateSent            string  `json:"date_sent"`
// 	AccountSID          string  `json:"account_sid"`
// 	To                  string  `json:"to"`
// 	From                string  `json:"from"`
// 	MessagingServiceSID string  `json:"messaging_service_sid"`
// 	Body                string  `json:"body"`
// 	Code                int     `json:"code"`
// 	Message             string  `json:"message"`
// 	MoreInfo            string  `json:"more_info"`
// 	Status              string  `json:"status"`
// 	NumSegments         string  `json:"num_segments"`
// 	NumMedia            string  `json:"num_media"`
// 	Direction           string  `json:"direction"`
// 	APIVersion          string  `json:"api_version"`
// 	Price               float64 `json:"price"`
// 	PriceUnit           string  `json:"price_unit"`
// 	ErrorCode           int     `json:"error_code"`
// 	ErrorMessage        string  `json:"error_message"`
// 	URL                 string  `json:"uri"`
// 	SubresourceURIs     []struct {
// 		Media string `json:"media"`
// 	} `json:"subresource_uris"`
// }

// func env() (sid, auth string, err error) {
// 	sid, found := os.LookupEnv(EnvSID)
// 	if !found {
// 		err = errors.New("EnvSID(" + EnvSID + ") not found in env")
// 		return
// 	}

// 	auth, found = os.LookupEnv(EnvAuth)
// 	if !found {
// 		err = errors.New("EnvAuth(" + EnvAuth + ") not found in env")
// 		return
// 	}

// 	return
// }

// func (s *SMS) prep() {
// 	var err error

// 	s.sid, s.auth, err = env()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	s.url, err = url.Parse("https://api.twilio.com/2010-04-01/Accounts/" + s.sid + "/Messages.json")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// }

// // Text sends a text message to all the numbers.
// func Text(v interface{}, numbers ...string) Errors {
// 	return sms.Text(v, numbers...)
// }

// // Text sends a text message to all the numbers.
// func (s *SMS) Text(v interface{}, numbers ...string) Errors {
// 	s.once.Do(s.prep)

// 	errs := Errors{}

// 	for _, num := range numbers {
// 		vals := url.Values{}
// 		vals.Set("To", num)
// 		vals.Set("From", "+18583844354")
// 		vals.Set("Body", fmt.Sprint(v))

// 		req, err := http.NewRequest(http.MethodPost, s.url.String(), strings.NewReader(vals.Encode()))
// 		if err != nil {
// 			errs[num] = err
// 			continue
// 		}
// 		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 		req.SetBasicAuth(s.sid, s.auth)

// 		cl := http.Client{}

// 		resp, err := cl.Do(req)
// 		if err != nil {
// 			errs[num] = err
// 			continue
// 		}

// 		var r response
// 		err = json.NewDecoder(resp.Body).Decode(&r)
// 		resp.Body.Close()
// 		if err != nil {
// 			if resp.StatusCode >= 300 {
// 				errs[num] = errors.New(resp.Status)
// 				continue
// 			}
// 			errs[num] = err
// 			continue
// 		}

// 		if status, err := strconv.Atoi(r.Status); err == nil && status >= 300 {
// 			errs[num] = errors.New(r.Message)
// 			continue
// 		}
// 		if r.ErrorCode != 0 {
// 			errs[num] = errors.New(r.ErrorMessage)
// 			continue
// 		}
// 	}

// 	return errs
// }
