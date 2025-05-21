#!/bin/bash

set -e  # Detener si algo falla

echo "🧨 Eliminando clúster anterior (si existe)..."
kind delete cluster --name sistema-personas

echo "🚀 Creando clúster nuevo desde kind-config.yaml..."
kind create cluster --config k8s/kind-config.yaml

echo "🌐 Instalando Ingress NGINX..."
kubectl apply -f k8s/ingress-nginx.yaml

echo "⏳ Esperando a que el controlador Ingress esté listo..."
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

# ------------------ MONGO ------------------
echo "🛠️ Desplegando base de datos MongoDB..."
kubectl apply -f k8s/mongo/namespace-mongo.yaml
kubectl apply -f k8s/mongo/pv-mongo.yaml
kubectl apply -f k8s/mongo/pvc-mongo.yaml
kubectl apply -f k8s/mongo/secrets-mongo.yaml
kubectl apply -f k8s/mongo/deployment-mongo.yaml
kubectl apply -f k8s/mongo/service-mongo.yaml

kubectl wait --namespace=mongo-ns \
  --for=condition=available deployment/mongo-deployment \
  --timeout=90s

# ------------------ CREATE ------------------
echo "🧩 Desplegando p-go-create..."
./inicializar-create.sh

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
