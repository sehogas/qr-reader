package util

import (
	"log"
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CheckConfig(m map[string]string) {

	if m == nil {
		log.Fatal("Invalid environment file")
	}

	v, e := m["CLIENT_ID"]
	if !e || v == "" {
		log.Fatal("CLIENT_ID is a required parameter")
	}

	v, e = m["ZONE_ID"]
	if !e || v == "" {
		log.Fatal("ZONE_ID is a required parameter")
	}

	v, e = m["EVENT_ID"]
	if !e || v == "" {
		log.Fatal("EVENT_ID is a required parameter")
	}
	if v == "I" || v == "O" {
		if v == "I" {
			m["EVENT_ID"] = "E"
		}
		if v == "O" {
			m["EVENT_ID"] = "S"
		}
	} else {
		log.Fatal("EVENT_ID is a required parameter")
	}

	v, e = m["MODE"]
	if !e || v == "" {
		log.Fatal("MODE is a required parameter")
	}
	if !(v == "1" || v == "2") {
		log.Fatal("The value MODE parameter is invalid")
	}

	v, e = m["URL_GET_CARDS"]
	if !e || v == "" {
		log.Fatal("URL_GET_CARDS is a required parameter")
	}

	v, e = m["URL_POST_ACCESS"]
	if !e || v == "" {
		log.Fatal("URL_POST_ACCESS is a required parameter")
	}

	v, e = m["API_KEY"]
	if !e || v == "" {
		log.Fatal("API_KEY is a required parameter")
	}

	v, e = m["DEVICE_ID"]
	if !e || v == "" {
		log.Fatal("DEVICE_ID is a required parameter")
	}

	_, e = m["RTSP"]
	if !e {
		log.Fatal("RTSP is a required parameter")
	}

	v, e = m["DB"]
	if !e || v == "" {
		log.Fatal("DB is a required parameter")
	}

	v, e = m["PATH_WAV_GRANTED"]
	if !e || v == "" {
		log.Fatal("PATH_WAV_GRANTED is a required parameter")
	}
	if !FileExists(v) {
		log.Fatal("File configured in parameter PATH_WAV_GRANTED does not exist")
	}

	v, e = m["PATH_WAV_DENIED"]
	if !e || v == "" {
		log.Fatal("PATH_WAV_DENIED is a required parameter")
	}
	if !FileExists(v) {
		log.Fatal("File configured in parameter PATH_WAV_DENIED does not exist")
	}

	v, e = m["PATH_BACKGROUND_GRANTED"]
	if !e || v == "" {
		log.Fatal("PATH_BACKGROUND_GRANTED is a required parameter")
	}
	if !FileExists(v) {
		log.Fatal("File configured in parameter PATH_BACKGROUND_GRANTED does not exist")
	}

	v, e = m["PATH_BACKGROUND_DENIED"]
	if !e || v == "" {
		log.Fatal("PATH_BACKGROUND_DENIED is a required parameter")
	}
	if !FileExists(v) {
		log.Fatal("File configured in parameter PATH_BACKGROUND_DENIED does not exist")
	}
}
