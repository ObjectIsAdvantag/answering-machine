# Goal

A #GOLang Answering Machine backed by Cisco Tropo Communication Services 


# How to use it

0. Download binaries from releases, see https://github.com/ObjectIsAdvantag/answering-machine/releases
   - or git clone and build (make all)
   
1. Signup at http://tropo.com
   - Note : Your login/password credentials will be used to authenticate against the Tropo REST API (provisonning) 

2. Run provision.sh
   - TROPO_USER
   - TROPO_PASSWORD
   - TROPO_COUNTRY_PREFIX       # Example : 1, 44, ... note that France, Germany need escaladation to Cisco teams
   - GOLAM_ENDPOINT             # Example : http://mygolam.localtunnel.me 
   - GOLAM_EMAIL_TRANSCRIPT     # Email where to send your messages
 
/!\ Write down your brand new Answering Machine phone number
   
3. Launch your answering machine
   - WINDOWS > golam.exe -port 8080 -logtostderr=true -v=5
   - Linux   > golam -port 8080 -logtostderr=true -v=5 
   - DOCKER  # docker run -d -port 8080:8080 ObjectIsAdvantag/golam

5. If your host is not visible on the internet, install localtunnel
   > npm install -g localtunnel
   > lt -p 8080 -d mygolam

6. Call your answering machine and leave a message
   - check your email
   - browse through your messages

7. Optional : launch the messages recorder
   - GOLAM_RECORDER
   

# Roadmap

Check Releases and Milestones for more details

FUTURE : see vNext Milestone
   
IN PROGRESS: see v0.3 Milestone
   [ ] installation guidelines 
   [ ] recordings persistance (BoltDB)
   [ ] API to browse messages (date, duration, recording, transcript if available)
   [X] Enhanced Tropo encapsulation (TropoVoice)

v0.2: MVP with Cisco Tropo's communication backend (Call/Voice/Recording)
   - golang encapsulation of the Tropo Web API (Session, Record, Say)  
   - provisioning env (registration, phone number) via Tropo REST API (see postman collection)
   - local tests via localtunnel 
   
v0.1: back-end API structure
   - args, version, glog
   - healthcheck
   - error structure
      

# Want to contribute 

Take your dev hat and ...

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


# License

MIT, see license file.

Feel free to use, reuse, extend, and contribute



