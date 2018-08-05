IMPORTANT: Internet connection is necessary in order to get the project up and running

## Database Setup

### Create Databases

    $ buffalo db create -a

### Run Migration

    $ buffalo db migrate

### Seed Database

    $ buffalo task db:seed

Note: A `manager` user will be created through this task

### Import Electoral District and Polling Division

- manually import `electoral_districts.csv` directly into electoral_districts table
- manually import `polling_divisions.csv` directly into polling_divisions table

Note: both the files are stored in AmazonS3 > rallychq

## Set Application Service Keys

The keys for services like Twilio can be set in the `.env` file found in the root of the
project

```
// example content of .env

TWILIO_AC_SID=AC23a19xxxxxxxxxxb60783eccfa4c2
TWILIO_AUTH_TOKEN=f57fbe7303e5dxxxxxxx4eafbf6796
TWILIO_NO=+0123456789
```

## Starting the Application

    $ PORT=4000 buffalo dev

If you point your browser to [http://127.0.0.1:4000](http://127.0.0.1:4000) you should see a "Welcome to Buffalo!" page.

## Heroku

### Set Environment Variables

    $ heroku config:set GO_ENV=production
    $ heroku config:set TWILIO_AC_SID=AC23a19xxxxxxxxxxb60783eccfa4c2
    $ heroku config:set TWILIO_AUTH_TOKEN=f57fbe7303e5dxxxxxxx4eafbf6796
    $ heroku config:set TWILIO_NO=+0123456789
    $ heroku config:set GOOGLE_MAPS_KEY=xxxxxxx-xxxxxxxxxxxxx-xxxxxxxxxxxx
    $ heroku config:set MASTER_USER_EMAIL=test1@test.com
    $ heroku config:set MASTER_USER_PW=ffffff
    $ heroku config:set JWT_SIGN_KEY=jwt_sign_key

### Enable [Postgis](https://postgis.net/install/)

    $ heroku pg:psql
    $ CREATE EXTENSION postgis;
    $ CREATE EXTENSION postgis_topology;

### Deployment

    $ heroku container:login
    $ heroku container:push web
    $ heroku container:release web
    $ heroku run /bin/app migrate
    $ heroku run /bin/app  task db:seed

[Powered by Buffalo](http://gobuffalo.io)
