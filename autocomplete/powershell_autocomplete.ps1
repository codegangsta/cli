$fn = $($MyInvocation.MyCommand.Name)
$name = $fn -replace "(.*)\.ps1$", '$1'
Register-ArgumentCompleter -Native -CommandName $name -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    # Like the bash and fish scripts, ask for completions of the words before
    # the cursor and pass the word being completed only when it is a flag.
    # Otherwise a fully typed subcommand is completed as if a space followed it.
    $words = @($commandAst.CommandElements |
        Where-Object { $_.Extent.EndOffset -lt $cursorPosition } |
        ForEach-Object { $_.Extent.Text })
    if ($wordToComplete.StartsWith('-')) {
        $words += $wordToComplete
    }
    $other = "& $($words -join ' ') --generate-shell-completion"
    Invoke-Expression $other | ForEach-Object {
        $parts = $_.Split(':', 2)
        if ($parts.Count -eq 2) {
            $completion = $parts[0].Trim()
            $description = $parts[1].Trim()
        } else {
            $completion = $_
            $description = $_
        }
        # PowerShell shows the results of a native completer as they are, so
        # keep only the ones that start with the word being completed.
        if ($completion.StartsWith($wordToComplete, [System.StringComparison]::OrdinalIgnoreCase)) {
            [System.Management.Automation.CompletionResult]::new($completion, $completion, 'ParameterValue', $description)
        }
    }
}