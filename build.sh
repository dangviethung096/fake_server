echo -n "Enter version: "
read VERSION

if [ -z "$VERSION" ]
then
    echo 'docker build -t dangviethung096/fake_server:latest -f dockerfile .'
    docker build -t dangviethung096/fake_server -f dockerfile .
    docker push dangviethung096/fake_server
else
    echo 'docker build -t dangviethung096/fake_server:$VERSION -f dockerfile .'
    docker build -t dangviethung096/fake_server:$VERSION -f dockerfile .
    docker push dangviethung096/fake_server:$VERSION
fi