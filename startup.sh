docker_image_name="stark985/bsos:0.0.10"
go build -o bsos .
docker build -t $docker_image_name .
kind load docker-image $docker_image_name --name three-node
kubectl apply -f manifests/controller.yaml