package sms

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"sync"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func toString(number uint64) string {
	return strconv.FormatUint(number, 10)
}

func Verizon(number uint64) Number {
	return Number(toString(number) + "@vtext.com")
}

func ATT(number uint64) Number {
	return Number(toString(number) + "@mmc.att.net")
}

func Sprint(number uint64) Number {
	return Number(toString(number) + "@messaging.sprintpcc.com")
}

func TMobile(number uint64) Number {
	return Number(toString(number) + "@tmomail.net")
}

var carriers = map[string][]string{
	"7-11 Speakout (USA GSM)":              {"number@cingularme.com"},
	"Alaska Communications Systems":        {"number@msg.acsalaska.com"},
	"Alltel Wireless":                      {"number@message.alltel.com"},
	"AT&T Mobility (formerly Cingular)":    {"number@mmc.att.net", "number@txt.att.net", "number@mmode.com", "number@cingularme.com"},
	"Bell Mobility & Solo Mobile (Canada)": {"number@txt.bell.ca"},
	"Boost Mobile":                         {"number@myboostmobile.com"},
	"Cellular One (Dobson)":                {"number@mobile.celloneusa.com"},
	"Cingular (Postpaid)":                  {"number@cingularme.com"},
	"Centennial Wireless":                  {"number@cwemail.com"},
	"Cingular (GoPhone prepaid)":           {"number@cingularme.com"},
	"Claro (Nicaragua)":                    {"number@ideasclaro-ca.com"},
	"Comcel":                               {"number@comcel.com.co"},
	"Cricket":                              {"number@mmc.mycricket.com"},
	"CTI":                                  {"number@smc.ctimovil.com.ar"},
	"Emtel (Mauritius)":                    {"number@emtelworld.net"},
	"Fido (Canada)":                        {"number@fido.ca"},
	"Globalstar":                           {"number@msg.globalstarusa.com"},
	"Helio":                                {"number@messaging.sprintpcc.com"},
	"Illinois Valley Cellular":             {"number@ivctext.com"},
	"IT Company Australia":                 {"number@itcompany.com.au"},
	"Iridium (satellite)":                  {"number@msg.iridium.com"},
	"Meteor (Ireland)":                     {"number@smc.mymeteor.ie"},
	"MetroPCS":                             {"number@mymetropcc.com"},
	"Movicom":                              {"number@movimensaje.com.ar"},
	"Movistar (Colombia)":                  {"number@movistar.com.co"},
	"MTN (South Africa)":                   {"number@smc.co.za"},
	"MTS (Canada)":                         {"number@text.mtsmobility.com"},
	"Nextel (Argentina)":                   {"TwoWay.11number@nextel.net.ar"},
	"Personal (Argentina)":                 {"11number@personal-net.com.ar"},
	"Plus GSM (Poland)":                    {"+48number@text.plusgsm.pl"},
	"President's Choice (Canada)":          {"number@txt.bell.ca"},
	"Qwest":                      {"number@qwestmp.com"},
	"Rogers (Canada)":            {"number@pcc.rogerc.com"},
	"Sasktel (Canada)":           {"number@smc.sasktel.com"},
	"Setar Mobile email (Aruba)": {"297+number@mac.aw"},
	"SMSGlobal":                  {"number@smc.smsglobal.com.au"},
	"Sprint (PCS)":               {"number@messaging.sprintpcc.com", "number@pm.sprint.com"},
	"Sprint (Nextel)":            {"number@page.nextel.com", "number@messaging.nextel.com"},
	"Suncom":                     {"number@tmc.suncom.com"},
	"T-Mobile":                   {"number@tmomail.net"},
	"T-Mobile (Austria)":         {"number@smc.t-mobile.at"},
	"Telus Mobility (Canada)":    {"number@msg.teluc.com"},
	"Tigo (Formerly Ola)":        {"number@smc.tigo.com.co"},
	"Tracfone (prepaid)":         {"number@cingularme.com", "number@tmomail.net", "number@vtext.com", "number@email.uscc.net", "number@message.alltel.com"},
	"Unicel":                     {"number@utext.com"},
	"US Cellular":                {"number@email.uscc.net", "number@mmc.uscc.net"},
	"Verizon":                    {"number@vtext.com", "number@vzwpix.com"},
	"Virgin Mobile (Canada)":     {"number@vmobile.ca"},
	"Virgin Mobile (USA)":        {"number@vmobl.com"},
	"Vodacom (South Africa)":     {"number@voda.co.za"},
	"YCC": {"number@smc.ycc.ru"},
	"B2sms (International)":         {"number@b2smc.com"},
	"CardBoardFish (International)": {"number@username.etexting.com"},
	"Club4sms (Pakistan)":           {"number@club4smc.com"},
	"Esendex (AU, ES, FR, IE, UK)":  {"number@esendex.net"},
	"Ipipi.com":                     {"number@opensmc.ipipi.com"},
	"Kapow! SMS Gateway":            {"number@kapow.co.uk"},
	"Letxt (International)":         {"number@smc.letxt.com.au"},
	"Me2mobile (Australia)":         {"number@me2mobile.com"},
	"Mobe.Net":                      {"number@mobe.net"},
	"pktpix.com (International)":    {"number@pktpix.com"},
	"Red Oxygen (International)":    {"number@redoxygen.net"},
	"Soprano (Australia)":           {"number@soprano.com.au"},
	"TellusTalk":                    {"number@esmc.nu"},
	"ToText.net":                    {"number@totext.net"},
	"Txtlocal.com":                  {"number@txtlocal.co.uk"},
	"ViaNett":                       {"number@smc.vianett.no"},
	"Webtext":                       {"number@webtext.com"},
	"ABTXT.COM":                     {"number@abtxt.com"},
	"MOBILEMAIL.RU":                 {"number@mobilemail.ru"},
}

var (
	Email string

	Folder = "sms_credentials"

	s sms
)

// Errors maps numbers to errors.
type Errors map[Number]error

// Numbers gathers numbers only from the errors.
func (e Errors) Numbers() []Number {
	var nums []Number
	for num := range e {
		nums = append(nums, num)
	}
	return nums
}

// Text sends a text message to the phone numbers.
func Text(v interface{}, numbers ...Number) Errors {
	errs := Errors{}
	for _, num := range numbers {
		if err := num.Text(v); err != nil {
			errs[num] = err
		}
	}
	return errs
}

// Number is an SMS email wrapper for a phone number.
type Number string

// Text sends a text message to the phone number.
func (n Number) Text(v interface{}) error {
	return s.email(string(n), fmt.Sprint(v))
}

type sms struct {
	Email string

	once sync.Once
	srv  *gmail.Service
}

func (s *sms) prep() {
	if len(s.Email) == 0 {
		s.Email = Email
	}

	key, err := ioutil.ReadFile(Folder + "/client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	cfg, err := google.ConfigFromJSON(key, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	s.srv, err = gmail.New(getClient(context.Background(), cfg))
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}
}

func (s *sms) email(to, body string) error {
	s.once.Do(s.prep)

	boundary := "__Server_Mailing__"

	rawMsg := []byte(
		`Content-Type: multipart/mixed; boundary=` + boundary + `
MIME-Version: 1.0
to: ` + to + `
from: ` + s.Email + `
subject: 

--` + boundary + `
Content-Type: text/html; charset="UTF-8"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit

` + body + `

--` + boundary + `--`,
	)

	var msg gmail.Message
	msg.Raw = base64.URLEncoding.EncodeToString(rawMsg)

	if _, err := s.srv.Users.Messages.Send("me", &msg).Do(); err != nil {
		return err
	}

	return nil
}
