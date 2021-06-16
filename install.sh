#!/bin/bash -e

processor=$(uname -m)

if [ "$processor" == "x86_64" ]; then
  arch="amd64"
else
  arch="386"
fi

case "$(uname -s)" in
  Darwin*)
    os="darwin_${arch}"
    ;;
  MINGW64*)
    os="windows_${arch}"
    ;;
  MSYS_NT*)
    os="windows_${arch}"
    ;;
  *)
    os="linux_${arch}"
    ;;
esac

echo "os=$os"

echo -e "\n\n===================================================="

get_latest_release() {
  curl --silent "https://api.github.com/repos/terraform-linters/tflint/releases/latest" | # Get latest release from GitHub api
    grep '"tag_name":' |                                                                  # Get tag line
    sed -E 's/.*"([^"]+)".*/\1/'                                                          # Pluck JSON value
}

if [ -z "${TFLINT_VERSION}" ] || [ "${TFLINT_VERSION}" == "latest" ]; then
  echo "Looking up the latest version ..."
  version=$(get_latest_release)
else
  version=${TFLINT_VERSION}
fi

echo "Downloading TFLint $version"
curl -L -o /tmp/tflint.zip "https://github.com/terraform-linters/tflint/releases/download/${version}/tflint_${os}.zip"
retVal=$?
if [ $retVal -ne 0 ]; then
  echo "Failed to download tflint_${os}.zip"
  exit $retVal
else
  echo "Download was successfully"
fi

echo -e "\n\n===================================================="
echo "Unpacking /tmp/tflint.zip ..."
unzip /tmp/tflint.zip -d /tmp/


echo "Cleaning /tmp/tflint.zip ..."
rm /tmp/tflint.zip 

echo -e "\n\n===================================================="
echo "Current tflint version"
/tmp/tflint -v
