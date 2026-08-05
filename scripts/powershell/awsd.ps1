# Deprecated. Prefer the generated integration, which also wires up completion
# and picks up the active profile in new shells. Add to your $PROFILE:
#
#   awsd init powershell | Out-String | Invoke-Expression
#
# This file is kept so existing setups keep working. The ~/.awsd parsing that
# used to live here now lives in the binary, behind `awsd_prompt shellenv`.

& awsd_prompt @args
if ($LASTEXITCODE -eq 0) {
    & awsd_prompt shellenv powershell | Out-String | Invoke-Expression
}
