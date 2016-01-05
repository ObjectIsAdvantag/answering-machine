# Goal

The cross device answering machine can be used to unify voice messages among all your devices


# Scenarios

Reserve a phone number from a Voice Service, ie your new answering machine

[Optional] Transfer incoming phone calls to your answering machine phone number (wire, mobile phones)

Define notification channels : SMS, mail, instant messengers to be informed when you receive new phone calls


# Implementation

[X] back-end API structure
   - args, version, glog
   - healthcheck
   - error structure
   - v0.1

[ ] MVP with a Cisco backend (Tropo, Spark)
   - Provisioning scripts (postman collection)
   - local hosting tunnelled 
   - Tropo voice services

[ ] Notification to a Spark room

[ ] Hosting on google app engine
   
[ ] Switch hosting to Cisco shipped


# Bootstrapping 

## Take your dev hat and ...

1. From Tropo.com, create a developper account.

   Note : Your login/password credentials are used to authenticate against the Tropo API via BasicAuth
   
2. [optional] Try the Tropo provisioning Postman collection

   - start or install postman
   - load the "Tropo provisonning test suite" collection : https://www.getpostman.com/collections/78022a95468b8ef01de9
   - create an environnement, and add the tropo_key variable with value Basic XXXXXXXXXX= (your base64 encrypted HTTP credentials)
   - run the Test suite
   
3. Register the answering machine

   - start or install postman
   - load the "Answering Machine" collection : https://www.getpostman.com/collections/682f3c9f46d74e7ed85f
   - create an environnement, and add the tropo_key variable with value Basic XXXXXXXXXX= (your base64 encrypted HTTP credentials)
   - run the Provision commands
   - go to the environment to retreive your phone number
   - note :you may also connect to the Tropo portal and check your application has been created

4. Test the answering machine service against your local env 

In a terminal, run localtunnel

```
> npm install -g localtunnel
note : a go version also exists, see https://github.com/NoahShen/gotunnelme

> lt --port 8080 --subdomain answeringmachine
your url is: https://answeringmachine.localtunnel.me
```

In another terminal, launch the service 
```
> make dev 
or
> go build .
> ./answering-machine.exe -logtostderr=true -v=5
```

In a third terminal, or postman, call the answering machine healthcheck endpoint
```
> curl -v -X GET https://answeringmachine.localtunnel.me/ping
```

## Traffic inspection (Tropo <=> Answering machine)

### 1. Incoming call (to the answering machine)

check API reference : https://www.tropo.com/docs/webapi/session 

#### Captured from Phono (Web emulator)

{"session":{"id":"42d71155dc14c8ad94af1fa24c652324","accountId":"5048353","timestamp":"2016-01-05T01:47:37.395Z","userType":"HUMAN","initialText":null,"callId":"74b2d7a1530e301f2dd021e46914f5a9","to":{"id":"9999556971","name":null,"channel":"V
OICE","network":"SIP"},"from":{"id":"tropo.com phono","name":null,"channel":"VOICE","network":"SIP"},"headers":{"Record-Route":"<sip:198.11.254.102:5060;transport=udp;lr>","Content-Length":"314","To":"sip:9999556971@sip.tropo.com","Contact":"<
sip:54.208.174.25:5060;transport=udp>","User-Agent":"Phono","Max-Forwards":"68","x-sid":"12d84bdefaf15f942e22ea2cdb8ecf8d","CSeq":"1 INVITE","Via":"SIP/2.0/UDP 198.11.254.102:5060;branch=z9hG4bK17hwja9as4gov;rport=5060;received=10.108.198.74",
"x-phono-sessionid":"8eed2c2c-d00f-4cf9-ac42-64ff5cf19a0d@pgw-v11g.phono.com","Call-ID":"1639jlhk6ch5x","Content-Type":"application/sdp","From":"<sip:tropo.com%20phono@pgw-v11g.phono.com>;tag=3c7dny90z5si"}}}

#### Captured from a Device (calling device)

{"session":{"id":"5176f28e6ba0453c30ba71e8e53ffacc","accountId":"5048353","timestamp":"2016-01-05T01:52:11.039Z","userType":"HUMAN","initialText":null,"callId":"ebc7bdc19ce4a4c77d4fdf5c5bd5311f","to":{"id":"3474913912","name":null,"channel":"V
OICE","network":"SIP"},"from":{"id":"33954218763","name":null,"channel":"VOICE","network":"SIP"},"headers":{"Record-Route":"<sip:198.11.254.102:5060;transport=udp;lr>","Content-Length":"329","To":"<sip:3474913912@192.168.3.83>","Contact":"<sip
:+33954218763@67.231.1.115:5060>","Max-Forwards":"66","x-sid":"98ab30c2c814b70096b479a4cfdeab13","Allow":"INVITE","CSeq":"213969 INVITE","Via":"SIP/2.0/UDP 198.11.254.102:5060;branch=z9hG4bKk50xh1h5jfr7;rport=5060;received=10.108.198.74","Call
-ID":"292328416_134037636@67.231.1.115","Content-Type":"application/sdp","Accept":"application/sdp","remote-party-id":"<sip:+33954218763@67.231.1.115:5060>;privacy=off;screen=no","From":"<sip:33954218763@67.231.1.115>;tag=gK0c018915"}}}


# License

MIT, see license file.

Feel free to use, reuse, extend, and contribute



