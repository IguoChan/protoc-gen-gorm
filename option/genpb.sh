set -ex

curPath="$(cd "$(dirname "$0")";pwd -P)"
pppPath="$curPath""/../third_party/proto/"

protoc -I./ \
  -I"${pppPath}" \
  --go_out="./" --go_opt paths=source_relative \
  gorm.proto
