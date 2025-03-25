docker run \
  --rm \
  --user "$(id -u):$(id -g)" \
  -v $PWD:/local openapitools/openapi-generator-cli:v7.12.0 generate \
  -i /local/purchase-cart-openapi.yaml \
  -g go-server \
  -o /local/internal/api \
  --additional-properties packageName=api,outputAsLibrary=true,onlyInterfaces=true,disallowAdditionalPropertiesIfNotPresent=false
