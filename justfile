build:
	for file in `ls | grep cmd`; do\
		pushd $file;\
		docker build -t example/$file .;\
		popd;\
	done

load:
	for file in `ls | grep cmd`; do\
		kind load docker-image example/$file;\
	done

install:
	for file in `ls | grep yaml`; do\
		kubectl apply -f $file;\
	done

remove:
	kubectl delete -f registry.yaml

register:
	kubectl exec -n spire spire-server-0 -- \
		/opt/spire/bin/spire-server entry create \
		-spiffeID spiffe://example.org/ns/default/sa/default \
		-parentID spiffe://example.org/ns/spire/sa/spire-agent \
		-selector k8s:ns:default \
		-selector k8s:sa:default

list:
	kubectl get pods
