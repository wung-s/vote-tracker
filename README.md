## Database Setup

### Create Databases

    $ buffalo db create -a

### Run Migration

    $ buffalo db migrate

### Seed Database

    $ buffalo task db:seed

## Starting the Application

    $ PORT=4000 buffalo dev

If you point your browser to [http://127.0.0.1:4000](http://127.0.0.1:4000) you should see a "Welcome to Buffalo!" page.

[Powered by Buffalo](http://gobuffalo.io)
