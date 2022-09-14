# Context:

   Application to run speedtest inside docker containers. Exposes matrix as json on port 8010

   By default it runs every 2 minutes

# Deployment:

   Deploys three containers:

   1 - MySQL as DB
   2 - Golang app which runs the cron and inserts results into the db.
   3 - WebUI which exposes the matrics from the DB and runs on port 8010.

# Steps:

Make sure docker and docker-compose is installed on your machine.

1- Clone the repo
      
      ```SHELL

      git clone https://github.com/rk280392/speedtest-golang.git

      ```
      
2 - Deploy this with docker-compose.

   ```SHELL

     docker-compose up -d --build

     ```

3 - Access the results in json through rest api by hitting endpoint 127.0.0.1:8010
