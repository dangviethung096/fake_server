echo 'docker build -t dangviethung096/fake_web:latest -f fake_web/dockerfile .'
docker build -t dangviethung096/fake_web:latest -f fake_web/dockerfile .

echo 'docker push dangviethung096/fake_web:latest'
docker push dangviethung096/fake_web:latest
