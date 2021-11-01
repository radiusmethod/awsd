#!/usr/bin/env bash

PREFIX=${PREFIX:="/usr/local"}
mkdir -p "${PREFIX}/bin"
PREFIX="$(cd -P -- "${PREFIX}" && pwd)"
echo "Installing into ${PREFIX}/bin" | sed "s#$HOME#~#g"

if ! command -v go version &> /dev/null
then
    echo " -=-=--=-=-=-=-=-=-=-=-=-=-=-=- "
    echo "                                "
    echo "              !!!               "
    echo "          go not found          "
    echo "      you can download here     "
    echo " https://golang.org/doc/install "
    echo "                                "
    echo " -=-=--=-=-=-=-=-=-=-=-=-=-=-=- "
    exit
fi

OS=$(go env GOOS)
ARCH=$(go env GOARCH)
RC="${SHELL##*/}rc"

env GOOS="$OS" GOBIN="${PREFIX}"/bin/  GOARCH="$ARCH" go install

touch "${PREFIX}"/bin/_awsd
chmod +x "${PREFIX}"/bin/_awsd

cat > "${PREFIX}/bin/_awsd" <<EOF
#!/usr/bin/env bash

if [[ "\$1" == "version" ]]; then
  _awsd_prompt version
  return
fi

AWS_PROFILE="\$AWS_PROFILE" _awsd_prompt

selected_profile="\$(cat ~/.awsd)"

if [ -z "\$selected_profile" ]
then
  unset AWS_PROFILE
else
  export AWS_PROFILE="\$selected_profile"
fi
EOF

echo 'alias awsd="source _awsd"' >> "$HOME/.$RC"
echo " -=-=--=-=-=-=-=-=-=-=-=-=-=-=- "
echo "                                "
echo "   To Finish Installation       "
echo "  open new terminal or run      "
echo "       source ~/.$RC            "
echo "                                "
echo " -=-=--=-=-=-=-=-=-=-=-=-=-=-=- "
