// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

package tropo


import (
	"encoding/json"
	"testing"
)


func TestDecoding (t *testing.T) {

	bytes := []byte(`{"session":{"id":"16b48066f456157ce88e5378bf9482c3","accountId":"5048353","timestamp":"2016-01-05T09:52:56.110Z","userType":"HUMAN","initialText":null,"callId":"3a90fa5e9bba9ff5ab46a3c96f4c7fbd","to":{"id":"33182882566","name":null,"channel":"VOICE","network":"SIP"},"from":{"id":"33954218763","name":null,"channel":"VOICE","network":"SIP"},"headers":{"Max-Forwards":"66","x-sid":"f8b23dacb4dab0778a6373b8b11f220c","Record-Route":"<sip:198.11.254.99:5060;transport=udp;lr>","Content-Length":"321","Contact":"<sip:+33954218763@81.201.82.106:5060;transport=udp>","To":"<sip:33182882566@sip-trunk-voxbone.tropo.com>","CSeq":"102 INVITE","User-Agent":"Vox Callcontrol","Via":"SIP/2.0/UDP 198.11.254.99:5060;branch=z9hG4bKbjwfdev3yjdu;rport=5060;received=10.108.198.68","Call-ID":"XJI67ZEDKJFVVFCSLZIUJ6JCLA@81.201.82.106","Content-Type":"application/sdp","From":"<sip:33954218763@voxbone.com>;tag=as5e853893"}}}`)
	var sw sessionWrapper
	json.Unmarshal(bytes, &sw)
	if sw.Session.ID != "16b48066f456157ce88e5378bf9482c3" {
		t.Error("Wrong session ID, expected 16b48066f456157ce88e5378bf9482c3, got ", sw.Session.ID)
	}
}


func TestDecodingWithNestedStructs (t *testing.T) {

	type SessionWrapper struct {
		Session struct {
					ID        string `json:"id"`
					Accountid string `json:"accountId"`
				} `json:"session"`
	}

	bytes := []byte(`{"session":{"id":"16b48066f456157ce88e5378bf9482c3","accountId":"5048353" }}`)
	var sw SessionWrapper
	json.Unmarshal(bytes, &sw)
	if sw.Session.ID != "16b48066f456157ce88e5378bf9482c3" {
		t.Error("Wrong session ID, expected 16b48066f456157ce88e5378bf9482c3, got ", sw.Session.ID)
	}
}

func TestDecodingWithNestedStructs2 (t *testing.T) {


	type Session struct {
		ID        string `json:"id"`
		AccountID string `json:"accountId"`
	}

	type SessionWrapper struct {
		Session `json:"session"`
	}

	bytes := []byte(`{"session":{"id":"16b48066f456157ce88e5378bf9482c3","accountId":"5048353" }}`)
	var sw SessionWrapper
	json.Unmarshal(bytes, &sw)
	if sw.Session.ID != "16b48066f456157ce88e5378bf9482c3" {
		t.Error("Wrong session ID, expected 16b48066f456157ce88e5378bf9482c3, got ", sw.Session.ID)
	}
}


