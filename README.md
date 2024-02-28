# scalable-websockets

## Scalable WebSocket with Redis Pub/Sub and Traefik Gateway
This documentation outlines a scalable WebSocket implementation leveraging Go, Redis Pub/Sub, and Traefik as a reverse proxy gateway. The architecture is designed to support scaling within a Docker Swarm environment, ensuring efficient load handling and dynamic service discovery.

**WebSocket in Go:** Utilizes the Go programming language to establish and manage WebSocket connections, providing a robust and efficient foundation for real-time communication.

**Redis Pub/Sub:** Implements Redis Publish/Subscribe mechanisms for message distribution, ensuring scalable and efficient message delivery across WebSocket connections.

**Traefik as Gateway:** Employs Traefik to serve as a reverse proxy gateway, facilitating routing, load balancing, and SSL termination for WebSocket connections.

### get started

#### compose mode
```sh
docker-compose up -d
```

#### compose down
```sh
docker-compose down
```

#### swarm mode
```sh
docker stack deploy -c stack-deploy.yml ws
```

#### stack remove ws
```sh
docker stack rm ws
```

### scale up
```sh
docker service scale ws_ws=4
``` 

### scale down
```sh
docker service scale ws_ws=1
``` 

### watch log
```sh
docker service logs ws_ws
```