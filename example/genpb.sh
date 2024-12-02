set -ex

curPath="$(cd "$(dirname "$0")";pwd -P)"
pppPath="$curPath""/../third_party/proto/"
optionPath="$curPath""/.."

protoc -I./ \
  -I"${pppPath}" \
  -I"${optionPath}" \
  --go_out="./" --go_opt paths=source_relative \
  --gorm_out="./" --gorm_opt paths=source_relative \
  model.proto