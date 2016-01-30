# paste

paste is a simple (very simple) paste service in go

it expects a raw post (use curl) and will serve the file back as plain text

# install

    go get github.com/sigmonsays/paste/pasted


    mkdir -pv /srv/pastes
    pasted -data /srv/pastes -bindaddr :5555

You now have a HTTP paste service running on port 5555

Visit the home page for the bash client

    http://localhost:5555/

