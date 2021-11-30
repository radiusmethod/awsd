#!/usr/bin/env bash
mkdir -p "${PREFIX}}/bin"
PREFIX="$(cd -P -- "${PREFIX}" && pwd)"
echo "Installing into ${PREFIX}" | sed "s#$HOME#~#g"

mkdir -p ${PREFIX}
mkdir -p "${PREFIX}/bin"

RC="${SHELL##*/}rc"

GOOS= GOARCH= GOARM= GOFLAGS= go build -o "${PREFIX}/bin/_awsd_prompt"

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

chmod +x "${PREFIX}/bin/_awsd"

cat >> "$HOME/.$RC" <<EOF
alias awsd="source _awsd"
EOF
