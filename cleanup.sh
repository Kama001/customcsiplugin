docker_image_name="stark985/bsos:0.0.9"
kubectl delete pvc bsos-claim
kubectl delete deployment bsos
kubectl delete pod task-pv-pod
kubectl delete daemonsets.apps node-plugin
# for node in $(kind get nodes --name three-node); do
#   docker exec "$node" ctr -n k8s.io images rm docker.io/$docker_image_name
# done
# rm bsos
