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
```

## Starting the Application

    $ PORT=4000 buffalo dev

If you point your browser to [http://127.0.0.1:4000](http://127.0.0.1:4000) you should see a "Welcome to Buffalo!" page.

## Heroku

### Set Environment Variables

    $ heroku config:set TWILIO_AC_SID=AC23a19xxxxxxxxxxb60783eccfa4c2
    $ heroku config:set TWILIO_AUTH_TOKEN=f57fbe7303e5dxxxxxxx4eafbf6796
    $ heroku config:set TWILIO_NO=+0123456789
    $ heroku config:set FB_SERVICE_AC_KEY=serviceAccountKey.json
    $ heroku config:set GOOGLE_MAPS_KEY=xxxxxxx-xxxxxxxxxxxxx-xxxxxxxxxxxx

### Deployment

    $ heroku container:login
    $ heroku container:push web
    $ heroku run /bin/app migrate
    $ heroku run /bin/app  task db:seed

[Powered by Buffalo](http://gobuffalo.io)
