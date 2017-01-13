cloud · vet
===========

Running large-scale distributed systems like [Cloud Foundry][cf]
is hard enough, and then some goes and adds in a bunch of data
services to the mix.  Patches and upgrades threaten the stability
of those data services, and at the end of the day, what
application owners care most about is **their data**.

You need tests.  _Smoke tests_.

[Cloud · Vet][code] brings _operational clarity_ to running smoke
tests against a Cloud Foundry instance.  It provides direct,
unambiguous diagnostic output about what parts of your data
service layer it is testing, and what it found.

No more tests written by developers, in unit-test framework of the
week.  No more "it expected to exit, but did not - no furture
change possible".  What does that _even mean?!?_

Instead, we want simple, easy to understand output of all the
interesting ways in which to use a data service.

Getting Started
---------------

Since I assume you're already using BOSH to deploy and manage your
Cloud Foundry (and if you aren't, you should be!), so Cloud · Vet
is packaged as an easy-to-use [BOSH release][release].  Add this
to your CF deployment manifest:

```
instance_groups:
  - name: vet
    lifecycle: errand
    instances: 1

    azs: [z1]
    networks:
      - name: default

    vm_type: default
    stemcell: default

    jobs:
      - name:    tests
        release: cloudvet
        properties:
          cf:
            api_url:  https://api.system.10.244.0.34.netip.cc
            skip_ssl_verify: true

            username: admin
            password: secret

            domain:   run.10.244.0.34.netip.cc
            org:      system
            space:    cloudvet

          test:
            redis:
              service: redis
              plan:    small
            mongo:
              service: mongo
              plan:    small
            mysql:
              service: mysql
              plan:    shared
            # etc.
```

And then run the errand with `bosh run errand vet`

Supported Data Services
-----------------------

Cloud · Vet supports the following data services, with more coming
as soon as someone implements the support and submits a pull
request!

### Redis

Set the `test.redis.service` and `test.redis.plan` BOSH manifest
manifest properties, and you'll get smoke-tests for Redis!

```
starting redis smoke tests...
parsing VCAP_SERVICES env var to find our Redis endpoint... OK
searching VCAP_SERVICES for our 'redis' service... OK
connecting to 10.244.62.2:6379... OK
storing a value... OK
retrieving the stored value... OK
updating the stored a value... OK
retrieving the updated value... OK


REDIS TESTS PASSED!
```

### MongoDB

Set the `test.mongo.service` and `test.mongo.plan` BOSH manifest
manifest properties, and you'll get smoke-tests for MongoDB!


```
starting mongo smoke tests...
parsing VCAP_SERVICES env var to find our MongoDB endpoint... OK
searching VCAP_SERVICES for our 'mongo' service... OK
connecting to
mongodb://EbNsNWyUWjfQh0un:imwEFfWhHJjPvDQHtCCgijW5@10.244.61.2...
OK
using vet-db-ihRWHmrUQMqLxcJE/vet-col-vukvWTrTyBevgsuh
inserting a tracer document... OK
retrieving tracer document... OK
removing tracer document... OK
```

### MySQL

Set the `test.mysql.service` and `test.mysql.plan` BOSH manifest
manifest properties, and you'll get smoke-tests for MySQL /
MariaDB!

```
(session transcript forthcoming)
```

Contributing
------------

If you want smoke tests and can't find them, make some!  All you
have to do is register a new endpoint in the API, and implement
the guts of your testing requirements.  The cloudvet codebase
provides some neat helper utilies for managing output and
interacting with the `VCAP_SERVICES` environment variable, so you
can focus on your data system.

If you do write a new one, I'd love to merge it.  Submit a PR via
Github, and we'll go from there.  If you can also supply any
necessary bits for testing your pull request, that'd be awesome
(i.e. BOSH release tarballs, manifest files, etc.)

Happy Hacking!


[cf]: https://cloudfoundry.org
[code]: https://github.com/jhunt/cloudvet
[release]: https://github.com/jhunt/cloudvet-boshrelease
