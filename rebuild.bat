copy /y service\service.go %GOPATH%\src\github.com\NYTimes\mock-ec2-metadata\service\service.go
del bin\m2.exe
%GOPATH%/bin/gox -output "bin/m2" -os="windows" ./
cd bin 
m2