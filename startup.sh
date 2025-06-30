go build -o bsos .
docker build -t stark985/bsos:fix_idempotency_issue .
kind load docker-image stark985/bsos:fix_idempotency_issue --name three-node
kubectl apply -f manifests/deployment.yaml