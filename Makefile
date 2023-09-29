restart:
    @echo "Restarting the server..."
    @docker-compose down
    @docker rmi -f exo5-app
    @docker-compose up

logs:
    @echo "Showing logs..."
    @docker logs -f exo5-app-1`