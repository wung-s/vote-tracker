## Database Setup

### Create Databases

    $ buffalo db create -a

### Run Migration

    $ buffalo db migrate

### Seed Database

**Important**: Please create the user manually from Auth0 and make the necessary
changes to `grifts/db.go` before running the seed command

    $ buffalo task db:seed

## Starting the Application

    $ PORT=4000 buffalo dev

If you point your browser to [http://127.0.0.1:4000](http://127.0.0.1:4000) you should see a "Welcome to Buffalo!" page.

## Deploy to Heroku

    $ heroku container:login
    $ heroku container:push web
    $ heroku run /bin/app migrate
    $ heroku run /bin/app  task db:seed

[Powered by Buffalo](http://gobuffalo.io)
