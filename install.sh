#!/usr/bin/env bash

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
RC="${SHELL##*/}rc"

env GOOS="$OS" GOBIN=/usr/local/bin/  GOARCH=amd64 go install

touch /usr/local/bin/_awsd
chmod +x /usr/local/bin/_awsd

cat > "/usr/local/bin/_awsd" <<EOF
#!/usr/bin/env bash

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
echo " "
echo " To Finish Installation Run "
echo "       source ~/.$RC "
echo " "
