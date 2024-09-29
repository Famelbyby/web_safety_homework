run:
	docker-compose up -d

break:
	docker-compose down

gen_certs:
	./gen_ca.sh
	touch certbundle.pem
	./gen_cert.sh $(NAME) $(RANDOM_INT) > certbundle.pem
	mv ca.key certs/ca.key
	cat ca.crt >> certbundle.pem
	mv ca.crt certs/ca.crt
	mv cert.key certs/cert.key
	mv certbundle.pem certs/certbundle.pem
	mv certs API_server/certs
	cp -r API_server/certs proxy_server/certs/