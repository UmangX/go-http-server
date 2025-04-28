run:
	go run ./app 
empty:
	curl -v http://localhost:4221/echo/hello
test:
	curl -v --header "user-agent: foobar/1.2.3" http://localhost:4221/user-agent
fake:
	curl -v --header "user-agent: foobar/1.2.3" http://localhost:4221/gibberish
	
