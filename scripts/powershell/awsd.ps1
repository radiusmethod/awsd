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

$awsd_profile = $null
$awsd_region = $null
$awsd_has_region = $false

if (Test-Path "$env:USERPROFILE\.awsd")
{
    $lines = Get-Content "$env:USERPROFILE\.awsd"
    $hasKV = $false
    foreach ($line in $lines) {
        if ($line -match '=') { $hasKV = $true; break }
    }

    if ($hasKV) {
        foreach ($line in $lines) {
            if ($line -match '^\s*([^=]+?)\s*=\s*(.*)$') {
                $key = $Matches[1]
                $val = $Matches[2].Trim()
                switch ($key) {
                    'profile' { $awsd_profile = $val }
                    'region'  { $awsd_region = $val; $awsd_has_region = $true }
                }
            }
        }
    } else {
        # Legacy single-line format: whole file is a profile name.
        $awsd_profile = ($lines | Out-String).Trim()
    }
}

if (-not $awsd_profile)
{
    $env:AWS_PROFILE = $null
}
else
{
    $env:AWS_PROFILE = $awsd_profile
}

if ($awsd_has_region) {
    if (-not $awsd_region) {
        $env:AWS_REGION = $null
        $env:AWS_DEFAULT_REGION = $null
    } else {
        $env:AWS_REGION = $awsd_region
        $env:AWS_DEFAULT_REGION = $awsd_region
    }
}
