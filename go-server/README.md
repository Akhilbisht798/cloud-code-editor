### Containers.
- minikube tunnel
- docker run -p 9000:9000 -p 9001:9001 -e "MINIO_ROOT_USER=ROOT" -e "MINIO_ROOT_PASSWORD=password" minio/minio server /data --console-address ":9001"
- docker run --name mysql-container -e MYSQL_ROOT_PASSWORD=rootpassword -e MYSQL_DATABASE=db -p 3306:3306 -d mysql:latest
