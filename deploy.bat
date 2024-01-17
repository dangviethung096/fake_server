docker stop fake_web
docker rm fake_web
docker run -d -p 80:80 -p 443:443 --name fake_web dangviethung096/fake_web:latest