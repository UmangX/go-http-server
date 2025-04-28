run:
	go run ./app 
empty:
	curl -v http://localhost:4221/
test:
	curl -v --header "User-Agent: foobar/1.2.3" http://localhost:4221/user-agent
	
