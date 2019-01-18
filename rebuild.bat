make deps
echo F | xcopy /y service\service.go %GOPATH%\src\github.com\NYTimes\mock-ec2-metadata\service\service.go
del bin\m2.exe
copy /y mock-ec2-metadata-config.json bin\mock-ec2-metadata-config.json
%GOPATH%/bin/gox -output "bin/m2" -os="windows" ./
bin\m2