# chatapp
Steps to run program
1.) Create and start docker container [ONLY NEED TO RUN STEP 1 ONCE, FOR FUTURE RUN OF THE PROGRAM YOU JUST NEED TO START DOCKER AND RUN THE chatapp container]
    a.) Open command prompt
    b.) cd into chatapp folder
    c.) Create and Start the docker using the command
        docker-compose up -d

2.) Open nsqServer.go in another visual studio code window. Run the consumer.go 
    go run consumer.go

3.) Open nsqClient.go in a visual studio code window. Run the producer.go
    go run producer.go