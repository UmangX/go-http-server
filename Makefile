run:
	go run ./app --directory ./files/
gzip:
	curl -v -H "Accept-Encoding: invalid-encoding-1, gzip, invalid-encoding-2" http://localhost:4221/echo/abc
post:
	curl -v --data "12345" -H "Content-Type: application/octet-stream" http://localhost:4221/files/file_123
file:
	echo 'Hlo, World!' > /tmp/sd & curl -i http://localhost:4221/files/sd
empty:
	curl -v http://localhost:4221/echo/hello
test:
	curl -v --header "user-agent: foobar/1.2.3" http://localhost:4221/user-agent
