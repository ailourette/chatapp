# chatapp
Steps to run program
1.) Start the docker using the command
    docker-compose up -d

2.) Open consumer.go in another visual studio code window. Run the nsqServer.go 
    go run consumer.go

3.) Open producer.go in a visual studio code window. Run the nsqClient.go
    cd producer
    go run producer.go

4.)NSQ admin
   http://localhost:4171