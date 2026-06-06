setup:
	sudo docker-compose up -d
	sleep 2
	sudo docker exec -it chattui-pgvector psql -U postgres -d vectordb -c "CREATE EXTENSION IF NOT EXISTS vector;"

down:
	sudo docker-compose down
