build:
	go build
	mkdir -p output
	mv products-api output/api

image:
	docker build . -t osvaldomargato/products-api:1.0.0
