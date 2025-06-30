kubectl delete pvc bsos-claim
kubectl delete deployment bsos
for node in $(kind get nodes --name three-node); do
  docker exec "$node" ctr -n k8s.io images rm docker.io/stark985/bsos:fix_idempotency_issue
done
rm bsos
