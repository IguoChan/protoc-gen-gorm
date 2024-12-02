set -ex

curPath="$(cd "$(dirname "$0")";pwd -P)"
pppPath="$curPath""/../third_party/proto/"
optionPath="$curPath""/.."

protoc -I./ \
  -I"${pppPath}" \
  -I"${optionPath}" \
  --go_out="./model" --go_opt paths=source_relative \
  --gorm_out="./model" --gorm_opt paths=source_relative --gorm_opt with_gorm_option=true --gorm_opt with_gorm_dao=true\
  model.proto