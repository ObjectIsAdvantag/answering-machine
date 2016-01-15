## GOLAM

a #GOLang Answering Machine backed by Cisco Tropo Communication Services 

The AnsweringMachine is a customizable software that:
   - plays a welcome message at incoming calls, 
   - records voice messages,
   - let an authorized user check for new messages,
   - sends transcripts of recordings to your mailbox.

The project was started as a proof of concept to learn #golang and #tropo APIs.

You may find it worth giving it a try for fun and no profit, and who knows... fork this github repo to build your own.


## Walkthrough

1. Signup at http://tropo.com
   - Note: In the header bar, close to your firstname, your tropo account number is displayed, 7 digits, (5048353) in my case  
   - Note: Your login/password credentials will be used to authenticate against the Tropo REST API (see provisonning / postman collection below) 

2. Register the AnsweringMachine application
   - Go to "MY APPS", click "Create new app"
   - Fill the "Create new app" form with
      - name: Answering Machine (or whatever)
      - select Web (HTTP) API
      - script: fill in your future AnsweringMachine internet facing endpoint
         - example: https://myansweringmachine.localtunnel.me/tropo
         - Note: you'll be able to change it later
      - Click "create app"
   - On the next page, we'll reserve and associate a phone number to our AnsweringMachine
      - scroll down to the Numbers section
      - choose a country in the list, and a region 
      - Note : tropo is free for development, you won't be charged before opting to go to tropo production environment
      - Note : if your country is not listed, simply opt for a country to which you can initiate calls for free, so that you can play with your AnsweringMachine (US and UK went out fine in my case)
      - /!\ Write down your brand new AnsweringMachine phone number
         
3. Install your AnsweringMachine
   - download a Linux or Windows package from [releases](https://github.com/ObjectIsAdvantag/answering-machine/releases)
   - or use docker (see Docker below)
   - or git clone and build your own answering machine (see Contribute below)
   
4. Configure your AnsweringMachine
   - copy env-tropofs.json to env.json
       - the tropofs configuration makes Tropo store your recordings. A standalone service (#golang recorder-server) is provided to your convenience if you prefer hosting your recordings (audio files) in your data center.
   - customize the entries by editing the env.json file
       - Note: do not prefix any phone numbers with "+" nor "00", simply start with your country code (336780078XX in my case)
       - Note: if worried about personal or sensitive info disclosure, you can mixt data configuration in the env.json file and set others via environment variables (such as passwords, phone numbers or email). Simply name your environement variables the same as in the env.json file.

5. Launch your AnsweringMachine
   - check everything went well by calling a few AnsweringMachine endpoints in a Web browser or via curl:
   - Note: it is possible to override configuration by setting environment variables

   ``` bash
   > GOLAM_RECORDING_PASSWORD=YYYYYYYYY
   > GOLAM_CHECKER_NUMBER=336780078XX
   > ./answering-machine --port=8080 -logtostderr=true -v=5 --env=env-tropofs.json --messages=messages-fr.json
   ...
   > CURL -X GET http://localhost:8080/ping
   > CURL -X GET http://localhost:8080/conf
   ```

6. If your AnsweringMachine host is not visible on the internet, install a tunneler ([localtunnel](http://localtunnel.me/) in my case)

   ``` bash
   > npm install -g localtunnel
   > lt -p 8080 -l 127.0.0.1 -s myansweringmachine
   your url is: https://myansweringmachine.localtunnel.me
   ```

   Again, check everything went well by calling a few endpoints in a Web browser or via curl
      - https://myansweringmachine.localtunnel.me/ping
      - https://myansweringmachine.localtunnel.me/conf


7. Call your AnsweringMachine and leave a message
   - dial in your AnsweringMachine phone number
   - listen to your welcome message
   - after the beep leave a message
   - check your email for a transcript 
   - call again with the number specified as checker in env.json to check your new message
   - visit https://myansweringmachine.localtunnel.me/messages to get a global view of your recorded messages (and their evolving states from NEW to CHECKED)

     ``` bash
    // https://myansweringmachine.localtunnel.me/messages
    [
      {
        "CallID": "e916188203beb99ae0c4677af17ea84e",
        "CreatedAt": "2016-01-13T23:10:4567261876Z",
        "CallerNumber": "+3395421XXXX",
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

[v0.5]: in progress
   - fixed transcripts not transmitted
   - handles disconnect properly
   - documentation updates

[v0.4](https://github.com/ObjectIsAdvantag/answering-machine/milestones/v0.4): Hosting & Packaging
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
 
Pull the [AnsweringMachine image from DockerHub](https://hub.docker.com/r/objectisadvantag/answeringmachine/) and run your answering machine.
Thanks to #golang capacity to generate autonomous binaries, the image is pretty small (about 4 MB compressed).
Note: the image size is half size (2 MB) if we do not onboard the recorder-server, check tag 0.4-small

   ``` bash
   > docker pull objectisadvantag/answeringmachine
   > docker run -e XXXXX -it -p 8080:8080 objectisadvantag/answeringmachine
   ```

To override default configuration (which can be checked at the /conf endpoint), replace XXXXX in the command line above with any configuration variable you wish to overload:
   ``` bash
   -e GOLAM_RECORDER_USERNAME={tropo_account}
   -e GOLAM_RECORDER_PASSWORD={tropo_password}
   -e GOLAM_AUDIO_ENDPOINT=http://hosting.tropo.com/{tropo_account_number}/www/audio
   -e GOLAM_TRANSCRIPTS_EMAIL={your email}
   -e GOLAM_CHECKER_NAME={your firstname}
   -e GOLAM_CHECKER_NUMBER={your phone number without + prefix}
   ``` 

Example : 
   ``` bash
   -e GOLAM_RECORDER_USERNAME=ObjectIsAdvantag 
   -e GOLAM_RECORDER_PASSWORD=MonMotDePasse
   -e GOLAM_AUDIO_ENDPOINT=http://hosting.tropo.com/5048353/www/audio
   -e GOLAM_TRANSCRIPTS_EMAIL=steve.sfartz@gmail.com
   -e GOLAM_CHECKER_NAME=Steve
   -e GOLAM_CHECKER_NUMBER=336780078XX
   ``` 


### Building a Docker image
You can build your own Docker image by cloning this repo.

``` bash
> git clone https://github.com/ObjectIsAdvantag/answering-machine
> make docker
> docker run -e XXXXX -it -p 8080:8080 <image>  
```


## Want to contribute 

Your contributions are welcome, simply create Issues to discusss a design point, ask for an evolution, or push pull requests 
   - linux and mac testers wanted
   - check the Triage and vNext milestones 

To start with, fork the repository and run the make command.
Note : the whole project was developed on a windows platform, with a goal of full interoperability with linux and macOS. Drop me an email if you encounter any issue.


## License

MIT, see license file.

Feel free to use, reuse, extend, and contribute



