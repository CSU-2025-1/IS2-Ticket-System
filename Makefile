rebuild_compose:
	sudo docker-compose down && sudo docker-compose up -d --build

restart_compose:
	sudo docker-compose down && sudo docker-compose up -d