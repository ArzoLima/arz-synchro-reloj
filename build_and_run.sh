#!/bin/bash

echo "Compilando programas..."

# Compilar el programa principal (console)
echo "Compilando programa principal..."
go build -o console cmd/console/main.go
if [ $? -ne 0 ]; then
    echo "Error compilando el programa principal"
    exit 1
fi

# Compilar el scheduler
echo "Compilando scheduler..."
go build -o scheduler cmd/scheduler/main.go
if [ $? -ne 0 ]; then
    echo "Error compilando el scheduler"
    exit 1
fi

echo "Compilación completada exitosamente"
echo ""
echo "Para ejecutar el scheduler que correrá todos los días a las 9am:"
echo "./scheduler"
echo ""
echo "Para ejecutar el programa principal una sola vez:"
echo "./console" 