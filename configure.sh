#!/usr/bin/env bash

PREFIX=${1:="/usr/local/bin"}
RC="${SHELL##*/}rc"

cat > "${PREFIX}/_awsd" <<EOF
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

chmod +x "${PREFIX}/_awsd"

cat >> "$HOME/.$RC" <<EOF
alias awsd="source _awsd"
EOF
