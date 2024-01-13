echo 'docker build -t fake_server:latest -f account_service/dockerfile .'
docker build -t dangviethung096/fake_server:latest -f account_service/dockerfile .

echo 'docker build -t fake_web:latest -f fake_web/dockerfile .'
docker build -t dangviethung096/fake_web:latest -f fake_web/dockerfile .

echo 'docker push dangviethung096/fake_web:latest'
docker push dangviethung096/fake_web:latest

echo 'docker pull dangviethung096/fake_web:latest'
docker pull dangviethung096/fake_web:latest