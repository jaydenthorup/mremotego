# Convert mRemoteNG XML to MremoteGO YAML format
param(
    [string]$SourceXml = "Devops-Mremote Config.xml",
    [string]$OutputYaml = "connections.yaml"
)

function Convert-ProtocolToMremoteGO {
    param([string]$protocol)
    
    switch ($protocol) {
        "SSH2" { return "ssh" }
        "RDP" { return "rdp" }
        "VNC" { return "vnc" }
        "HTTP" { return "http" }
        "HTTPS" { return "https" }
        "Telnet" { return "telnet" }
        default { return "unknown" }
    }
}

function Convert-ConnectionNode {
    param($node, $indent = 0)
    
    $spaces = "  " * $indent
    $yaml = ""
    
    # Clean up name - remove trailing spaces
    $cleanName = $node.Name.Trim()
    $yaml += "$spaces- name: `"$cleanName`"`n"
    
    if ($node.Type -eq "Container") {
        $yaml += "$spaces  type: folder`n"
        if ($node.Descr) {
            $yaml += "$spaces  description: `"$($node.Descr)`"`n"
        }
        
        # Process children
        $children = $node.SelectNodes("Node")
        if ($children.Count -gt 0) {
            $yaml += "$spaces  children:`n"
            foreach ($child in $children) {
                $yaml += Convert-ConnectionNode -node $child -indent ($indent + 2)
            }
        }
    }
    elseif ($node.Type -eq "Connection") {
        $yaml += "$spaces  type: connection`n"
        
        # Protocol
        $protocol = Convert-ProtocolToMremoteGO -protocol $node.Protocol
        $yaml += "$spaces  protocol: $protocol`n"
        
        # Host
        if ($node.Hostname) {
            $yaml += "$spaces  host: `"$($node.Hostname)`"`n"
        }
        
        # Port
        if ($node.Port) {
            $yaml += "$spaces  port: $($node.Port)`n"
        }
        
        # Username
        if ($node.Username) {
            $yaml += "$spaces  username: `"$($node.Username)`"`n"
        }
        
        # Password - skipped, use 1Password or set manually
        # if ($node.Password) {
        #     $yaml += "$spaces  password: `"`"`n"
        # }
        
        # Domain
        if ($node.Domain) {
            $yaml += "$spaces  domain: `"$($node.Domain)`"`n"
        }
        
        # Description
        if ($node.Descr) {
            # Clean up description - remove newlines
            $cleanDescr = $node.Descr -replace '[\r\n]+', ' '
            $cleanDescr = $cleanDescr.Trim()
            $yaml += "$spaces  description: `"$cleanDescr`"`n"
        }
        
        # RDP specific settings
        if ($protocol -eq "rdp") {
            if ($node.UseCredSsp -eq "true") {
                $yaml += "$spaces  use_credssp: true`n"
            }
            if ($node.Resolution) {
                $yaml += "$spaces  resolution: `"$($node.Resolution)`"`n"
            }
        }
    }
    
    return $yaml
}

# Load the XML
[xml]$xml = Get-Content $SourceXml

# Get all connection nodes
$connections = $xml.SelectNodes("//Node[@Type='Connection']")

Write-Host "Found $($connections.Count) connections"

# Start building YAML
$yaml = "version: `"1.0`"`n"
$yaml += "connections:`n"

# Convert each connection
foreach ($conn in $connections) {
    $yaml += Convert-ConnectionNode -node $conn -indent 1
}

# Save to file
$yaml | Out-File -FilePath $OutputYaml -Encoding UTF8

Write-Host "Converted to YAML: $OutputYaml"
Write-Host ""
Write-Host "NOTE: Passwords are encrypted in mRemoteNG format and need to be decrypted."
Write-Host "You'll need to manually set passwords or use 1Password integration."
