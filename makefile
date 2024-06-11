apis:
	cd .\api\ && go run  .\cmd\main.go 

migrate_up: 
	cd .\migrations_up\ && go run .\main.go

run:
	echo "RUN"