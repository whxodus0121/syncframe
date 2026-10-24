$ErrorActionPreference = 'Stop'
Add-Type -AssemblyName System.IO.Compression
Add-Type -AssemblyName System.IO.Compression.FileSystem
$qaRoot = Join-Path (Split-Path -Parent $PSScriptRoot) 'artifacts\browser-qa'
$utf8 = [System.Text.UTF8Encoding]::new($false, $true)
$changedEntries = 0
foreach ($traceFile in Get-ChildItem -LiteralPath $qaRoot -Recurse -Filter trace.zip -File) {
    $zip = [System.IO.Compression.ZipFile]::Open($traceFile.FullName, [System.IO.Compression.ZipArchiveMode]::Update)
    try {
        foreach ($entry in @($zip.Entries)) {
            $stream = $entry.Open()
            $buffer = [System.IO.MemoryStream]::new()
            try { $stream.CopyTo($buffer) } finally { $stream.Dispose() }
            $bytes = $buffer.ToArray(); $buffer.Dispose()
            try { $text = $utf8.GetString($bytes) } catch [System.Text.DecoderFallbackException] { continue }
            $safe = [regex]::Replace($text, '(?i)(?<![a-f0-9])[a-f0-9]{64}(?![a-f0-9])', '[REDACTED_TEST_CAPABILITY]')
            if ($safe -ceq $text) { continue }
            $entryName = $entry.FullName
            $entry.Delete()
            $replacement = $zip.CreateEntry($entryName, [System.IO.Compression.CompressionLevel]::Optimal)
            $writer = [System.IO.StreamWriter]::new($replacement.Open(), [System.Text.UTF8Encoding]::new($false))
            try { $writer.Write($safe) } finally { $writer.Dispose() }
            $changedEntries++
        }
    } finally { $zip.Dispose() }
}
Write-Host "Redacted local test capabilities in $changedEntries trace entries. DOM, actions and timings retained."
