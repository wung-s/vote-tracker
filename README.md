## Database Setup

### Create Databases

    $ buffalo db create -a

### Run Migration

    $ buffalo db migrate

### Seed Database

**Important**: Please create the user manually from Auth0 and make the necessary
changes to `grifts/db.go` before running the seed command

    $ buffalo task db:seed

## Set Application Service Keys

The keys for services like Auth0 & Twilio can be set in the `.env` file found in the root of the
project

```
//example content of .env

TWILIO_AC_SID=AC23a19xxxxxxxxxxb60783eccfa4c2
TWILIO_AUTH_TOKEN=f57fbe7303e5dxxxxxxx4eafbf6796
TWILIO_NO=+0123456789
AUTH0_API_AUDIENCE_VAL1=https://dev.rallychq.com
AUTH0_API_ISSUER=https://gotv.auth0.com/
JWKS_URI=https://gotv.auth0.com/.well-known/jwks.json
```

## Starting the Application

    $ PORT=4000 buffalo dev

If you point your browser to [http://127.0.0.1:4000](http://127.0.0.1:4000) you should see a "Welcome to Buffalo!" page.

## Heroku

### Set Environment Variables

    $ heroku config:set TWILIO_AC_SID=AC23a19xxxxxxxxxxb60783eccfa4c2
    $ heroku config:set TWILIO_AUTH_TOKEN=f57fbe7303e5dxxxxxxx4eafbf6796
    $ heroku config:set TWILIO_NO=+0123456789
    $ heroku config:set AUTH0_API_AUDIENCE_VAL1=https://dev.rallychq.com
    $ heroku config:set AUTH0_API_ISSUER=https://gotv.auth0.com/
    $ heroku config:set JWKS_URI=https://gotv.auth0.com/.well-known/jwks.json
    $ heroku config:set AUTH0_CLIENT_ID=nlKFW5maxxxXouVfuVXDs8duFFpH1
    $ heroku config:set AUTH0_CLIENT_SECRET=BcDvk7wxxxxxchm1J2PYwTpa-TMICbvYf1Vc8nCeiNxGao1D8tNL62o

### Deployment

    $ heroku container:login
    $ heroku container:push web
    $ heroku run /bin/app migrate
    $ heroku run /bin/app  task db:seed

[Powered by Buffalo](http://gobuffalo.io)
