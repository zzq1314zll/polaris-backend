cd "$(dirname "$0")"
cd ..

env=$1
if [ -z $env  ];then
  env="local"
fi
sleep 10
go run main.go --env=$env