#!/bin/bash

set -e  # Detener si algo falla

echo "🧨 Eliminando clúster anterior (si existe)..."
kind delete cluster --name sistema-personas

echo "🚀 Creando clúster nuevo desde kind-config.yaml..."
kind create cluster --config k8s/kind-config.yaml

# ------------------ PERMISOS ------------------
echo "🔧 Asegurando permisos de ejecución para scripts..."
chmod +x k8s/scripts/inicializar-create.sh
chmod +x ../p-go-search/k8s/inicializar-search.sh
chmod +x ../p-go-list/k8s/inicializar-list.sh
chmod +x ../p-go-update/k8s/inicializar-update.sh
chmod +x ../p-go-delete/k8s/inicializar-delete.sh

# ------------------ INGRESS ------------------
echo "🌐 Instalando Ingress NGINX..."
kubectl apply -f k8s/ingress-nginx.yaml

echo "⏳ Esperando a que el controlador Ingress esté listo..."
if ! kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=180s; then
  echo "⚠️ Advertencia: Ingress NGINX no alcanzó el estado 'Ready' en el tiempo esperado, pero podría estar funcionando."
else
  echo "✅ Ingress NGINX está listo y en ejecución."
fi

# ------------------ MONGO ------------------
echo "🛠️ Desplegando base de datos MongoDB..."
kubectl apply -f k8s/mongo/namespace-mongo.yaml
kubectl apply -f k8s/mongo/pv-mongo.yaml
kubectl apply -f k8s/mongo/pvc-mongo.yaml
kubectl apply -f k8s/mongo/secrets-mongo.yaml
kubectl apply -f k8s/mongo/deployment-mongo.yaml
kubectl apply -f k8s/mongo/service-mongo.yaml

echo "⏳ Esperando a que Mongo esté listo..."
if ! kubectl wait --namespace=mongo-ns \
  --for=condition=available deployment/mongo-deployment \
  --timeout=180s; then
  echo "⚠️ Advertencia: Mongo no alcanzó el estado 'Available' en el tiempo esperado, pero el pod podría estar corriendo."
else
  echo "✅ Mongo está disponible y funcionando correctamente."
fi

# ------------------ CREATE ------------------
echo "🧩 Desplegando p-go-create..."
./k8s/scripts/inicializar-create.sh

# ------------------ SEARCH ------------------
echo "🔍 Desplegando p-go-search..."
(cd ../p-go-search && ./k8s/inicializar-search.sh)

# ------------------ LIST ------------------
echo "📋 Desplegando p-go-list..."
(cd ../p-go-list && ./k8s/inicializar-list.sh)

# ------------------ UPDATE ------------------
echo "✏️ Desplegando p-go-update..."
(cd ../p-go-update && ./k8s/inicializar-update.sh)

# ------------------ DELETE ------------------
echo "❌ Desplegando p-go-delete..."
(cd ../p-go-delete && ./k8s/inicializar-delete.sh)

# ------------------ FINAL ------------------
echo -e "\n✅ Todo el clúster ha sido desplegado correctamente 🎉"
kubectl get all -A
