# Goal

A #GOLang Answering Machine backed by Cisco Tropo Communication Services 


# How to use it

1. Signup at http://tropo.com
   - Note : Your login/password credentials can also be used to authenticate against the Tropo REST API (see provisonning) 


2. Register the AnsweringMachine application
   - Go to "MY APPS", click "Create new app"
   - Fill the "Create new app" form with
      - name : Answering Machine (or whatever)
      - select Web (HTTP) API
      - script : fill in your answering machine public endpoint, example : https://myansweringmachine.localtunnel.me/tropo
      - Click create app
   - On the next screen, you'll create a phone number
      - Choose a country in the list, and a region 
      - Don't worry, you won't be billed, dev is free on tropo
      - /!\ Write down your brand new Answering Machine phone number
         
         
3. Download an answering machine packages 
   - Go to [releases](https://github.com/ObjectIsAdvantag/answering-machine/releases)
   
   - Or use docker (see Docker hereafter)
   - Or git clone and build your own answering machine (see Contribute hereafter)
   
   
4. Create an configuration with your personal data
   - Copy env-tropofs.json to env.json
   - Customize the entries


5. Launch your answering machine

``` bash
> ./answering-machine --port=8080 -logtostderr=true -v=5 --env=env-tropofs.json --messages=messages-fr.json
```

You can check everything went well by calling a few endpoints in a Web browser or via curl

   - http://localhost:8080/ping
   - http://localhost:8080/conf


6. If your host is not visible on the internet, install localtunnel

``` bash
> npm install -g localtunnel
> lt -p 8080 -s myansweringmachine
your url is: https://myansweringmachine.localtunnel.me
```

You can check everything went well by calling a few endpoints in a Web browser or via curl

   - http://myansweringmachine.localtunnel.me/ping
   - http://myansweringmachine.localtunnel.me/conf


7. Call your answering machine and leave a message
   - dial in your answering machine phone number
   - listen to your message
   - after the beep leave a message
   - check your email for a transcript 
   - call again with the number specified as checker in env.json to check your new message
   - visit http://myansweringmachine.localtunnel.me/messages to have a global view of your recorded messages (and their evolving states from NEW to CHECKED)
  

# Roadmap

Check Releases and Milestones for more details

FUTURE : see milestones [vNext](https://github.com/ObjectIsAdvantag/answering-machine/milestones/vNext) and [Triage](https://github.com/ObjectIsAdvantag/answering-machine/milestones/Triage) for non priorized issues

[v0.4](https://github.com/ObjectIsAdvantag/answering-machine/milestones/v0.4) : Hosting & Packaging
   - i18n messages
   - distinct messages & env conf
   - Docker support
   - Configuration endpoint /conf
   - installation guidelines
   
[v0.3](https://github.com/ObjectIsAdvantag/answering-machine/milestones/v0.3): Full MVP
   - welcome message, record, check messages
   - on-disk recordings storage (BoltDB)
   - admin API to browse voice messages
   - upload and download of recordings via Tropo File Storage or Standalone Recorder/AudioServer (sse Recorder)
   - enhanced Tropo WebAPI golang client
    
v0.2: MVP with Cisco Tropo's communication backend (Say/Record)
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



