curl -X POST http://localhost:8080/signup -d '{ "username": "Etz", "password": "Pass1234" }'
curl -X POST http://localhost:8080/login -d '{ "username": "Etz", "password": "Pass1234" }'
etz api --exec=scripts/exec.yaml -w=5 -r=5 -d=3s
