# check if $1 is empty
if (-not $args)
{
    # no argument passed
    Set-Variable -Name "AWS_PROFILE" -Value "$env:AWS_PROFILE"
    awsd_prompt
}
else
{
    # argument passed, assume it's a profile name
    Set-Variable -Name "AWS_PROFILE" -Value "$env:AWS_PROFILE"
    awsd_prompt $args
}

$selected_profile = Get-Content "$env:USERPROFILE\.awsd"

if (-not $selected_profile)
{
    $env:AWS_PROFILE = $null
}
else
{
    $env:AWS_PROFILE = $selected_profile
}
