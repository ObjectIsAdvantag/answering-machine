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
   - Tropo voice services
   - local hosting tunneled 
   
[ ] Notification to a Spark room

[ ] Hosting on google app engine
   
[ ] Switch hosting to Cisco shipped


# Bootstrapping 

## Take your dev hat and ...

1. From Tropo.com, create a developper account.
Note : Your login/password credentials are used to authenticate against the Tropo API via BasicAuth

2. Load the Postman collection : https://www.getpostman.com/collections/147b4e86dba33b6af8f5
Note : create a new environnement, where you'll store your BasicAuth credentials as tropo_key

3. Run a test from the Postman collection


# License

MIT, see license file.

Feel free to use, reuse, extend, and contribute



