# Golang-And-MongoDB-Image-Gallery

A REST API that stores and retrieves posts that contain images as "thumbnails".

## Run.
Use 'go run' with these environment vars.
PORT, DB_URL, DB_TIMEOUT

Or use docker.
```
docker-compose up --build
```

## Test.
```
go test ./server
```

## Use.
- **Create a new post using multipart form-data:**
  - POST /posts
- **Get a post:**
  - GET /posts/{postID}
- **Get many posts:**
  - GET /posts/{offset}/{limit}
- **Delete a post:**
  - DELETE /posts/{postID}
- **Get the thumbnail for a post:**
  - GET /posts/{postID}/thumbnail
