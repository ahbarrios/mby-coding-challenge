# Coding Challenge: AI Immigration Chatbot with Golang and Temporal


## Start the platform

```
docker-compose up -d
```

## Start **Temporal** worker
```
go run ./worker
```

## Start your http server
```
go run main.go
```

## Test it
```
curl -X POST http://localhost:3002/chat -d '{"user":"TOTO","content":"visa"}'
```

## TODO

- [ ] Manage conversational workflows to add context from previous chat history
- [ ] Add user authentication/identification to pick the right workflow/conversation with the bot
- [ ] Configure API, Keys and any other configuration parameter following [12factor](https://12factor.net/) principals
- [ ] Deploy somewhere 
