# chatapp
Steps to run program
1.) Create and start docker container [ONLY NEED TO RUN STEP 1 ONCE, FOR FUTURE RUN OF THE PROGRAM YOU JUST NEED TO START DOCKER AND RUN THE chatapp container]
    a.) Open command prompt
    b.) cd into chatapp folder
    c.) Create and Start the docker using the command
        docker-compose up -d
        
2.) Open consumer.go in another visual studio code window.
    go run consumer.go

3.) Open producer.go in a visual studio code window.
    open producer folder in visual studio code
    go run producer.go

4.) NSQ admin
   http://localhost:4171


5.) Create Docker Container for sql database
    Working: https://www.youtube.com/watch?v=X8W5Xq9e2Os
    a.) Download Docker Image
        docker pull mysql/mysql-server:latest

    b.) Run docker container
        docker run -p 3307:3306 --name mysql -e MYSQL_ROOT_PASSWORD=password -d mysql

    c.) Execute container
        docker exec -it mysql /bin/bash

    d.) Login to mysql
        mysql -uroot -p -A
        enter password:password

    e.) To create a table in sql database in container
 	    1.) USE mysql;
	    2.) Press enter
 	    3.) CREATE TABLE Users (UserName VARCHAR(30) NOT NULL PRIMARY KEY, Password VARCHAR(256), FirstName VARCHAR(30), LastName VARCHAR(30), Language VARCHAR(30));
	    4.) Press enter

    f.) To check content of table:
        USE mysql; SELECT * FROM Users;
