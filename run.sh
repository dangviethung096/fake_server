docker pull dangviethung096/fake_server
docker run -p 10015:10015 --name fake -d --mount type=bind,source=/home/hungdv39gec/data,target=/app/data dangviethung096/fake_server