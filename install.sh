#!/usr/bin/env bash

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

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
KERN=$(uname -m)
RC="${SHELL##*/}rc"

if [ "$KERN" == "aarch64" ]; then
  ARCH="arm64"
elif [ "$KERN" == "i686" ]; then
  ARCH="386"
else
  ARCH="amd64"
fi

env GOOS="$OS" GOBIN=/usr/local/bin/  GOARCH="$ARCH" go install

touch /usr/local/bin/_awsd
chmod +x /usr/local/bin/_awsd

cat > "/usr/local/bin/_awsd" <<EOF
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
