build: 
	GOOS=linux go build -o app 
	docker build -t restanrm/bell . 
	rm app
