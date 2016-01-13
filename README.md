## Goal

A #GOLang Answering Machine backed by Cisco Tropo Communication Services 


## Give it a try

1. Signup at http://tropo.com
   - Note : Your login/password credentials can also be used to authenticate against the Tropo REST API (see provisonning) 


2. Register the AnsweringMachine application
   - Go to "MY APPS", click "Create new app"
   - Fill the "Create new app" form with
      - name : Answering Machine (or whatever)
      - select Web (HTTP) API
      - script : fill in your answering machine public endpoint
         - example : https://myansweringmachine.localtunnel.me/tropo
      - Click "create app"
   - On the next screen, we'll get a phone number
      - Choose a country in the list, and a region 
      - Note : don't worry, you won't be billed, dev is free on tropo
      - Note : if your country is not listed, simply opt for a country to which you can initiate calls for free to test your answering machine (US and UK went out fine in my case)
      - /!\ Write down your brand new Answering Machine phone number
         
         
3. Install your AnsweringMachine
   - Download a Linux or Windows package from [releases](https://github.com/ObjectIsAdvantag/answering-machine/releases)
   
   - Or use docker (see Docker below)
   - Or git clone and build your own answering machine (see Contribute below)
   
   
4. Configure your AnsweringMachine
   - Copy env-tropofs.json to env.json
   - Customize the entries
   - Note that you should not prefix any phone number with + nor 00, simply start with your country code


5. Launch your answering machine

   ``` bash
   > ./answering-machine --port=8080 -logtostderr=true -v=5 --env=env-tropofs.json --messages=messages-fr.json
   ```

   You can check everything went well by calling a few endpoints in a Web browser or via curl:

   - http://localhost:8080/ping
   - http://localhost:8080/conf

   Note that it is possible to override configuration by setting environment variables
   Ex : GOLAM_CHECKER_NAME=Stève, GOLAM_CHECKER_NUMBER=33678007835


6. If your host is not visible on the internet, install localtunnel

   ``` bash
   > npm install -g localtunnel
   > lt -p 8080 -l 127.0.0.1 -s myansweringmachine
   your url is: https://myansweringmachine.localtunnel.me
   ```

    You can check everything went well by calling a few endpoints in a Web browser or via curl

   - http://myansweringmachine.localtunnel.me/ping
   - http://myansweringmachine.localtunnel.me/conf


7. Call your AnsweringMachine and leave a message
   - dial in your AnsweringMachine phone number
   - listen to your message
   - after the beep leave a message
   - check your email for a transcript 
   - call again with the number specified as checker in env.json to check your new message
   - visit http://myansweringmachine.localtunnel.me/messages to have a global view of your recorded messages (and their evolving states from NEW to CHECKED)

     ``` bash
    // https://answeringmachine.localtunnel.me/messages
    [
      {
        "CallID": "e916188203beb99ae0c4677af17ea84e",
        "CreatedAt": "2016-01-13T23:10:4567261876Z",
        "CallerNumber": "+33954218763",
        "Progress": "RECORDED",
        "Recording": "ftp://ftp.tropo.com/www/audio/e916188203beb99ae0c4677af17ea84e.wav",
        "Duration": 4400,
        "Transcript": "",
        "Status": "NEW",
        "CheckedAt": "0001-01-01T00:00:00Z"
      }
    ]
    ```
 


## Roadmap

Check [vNext](https://github.com/ObjectIsAdvantag/answering-machine/milestones/vNext) and [Triage](https://github.com/ObjectIsAdvantag/answering-machine/milestones/Triage) for non priorized issues.

[v0.4](https://github.com/ObjectIsAdvantag/answering-machine/milestones/v0.4) : Hosting & Packaging
   - i18n messages
   - distinct messages & env conf
   - docker support
   - configuration endpoint /conf
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


## Docker
 
TODO: Publish the answering machine image to docker

You can build your own Docker image by cloning this repo, and running make docker

``` bash
> git clone https://github.com/ObjectIsAdvantag/answering-machine
> make docker
> docker run -e XXXXX -it -p 8080:8080 <image>  
```

Please find hereafter an example of configuration variable,
``` bash
-e GOLAM_RECORDER_USERNAME=ObjectIsAdvantag 
-e GOLAM_RECORDER_PASSWORD=XXXXX
-e GOLAM_AUDIO_ENDPOINT=http://hosting.tropo.com/5048353/www/audio
-e GOLAM_TRANSCRIPTS_EMAIL=steve.sfartz@gmail.com
-e GOLAM_CHECKER_NAME=Stève
-e GOLAM_CHECKER_NUMBER=33678007833
``` 


## Want to contribute 

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


## License

MIT, see license file.

Feel free to use, reuse, extend, and contribute



