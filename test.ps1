# Script to process Azure Policy folders with azopaform

# Define the root path
$rootPath = "d:\project\azure-policy\built-in-policies\policyDefinitions"

# Check if the root path exists
if (-not (Test-Path -Path $rootPath -PathType Container)) {
    Write-Error "Root path does not exist: $rootPath"
    exit 1
}

# Get all folders under the root path
$folders = Get-ChildItem -Path $rootPath -Directory
$totalFolders = $folders.Count
$processedCount = 0

Write-Host "Found $totalFolders folders to process."

# Process each folder
foreach ($folder in $folders) {
    $processedCount++
    $fullPath = $folder.FullName
    $folderName = $folder.Name

    Write-Host "[$processedCount/$totalFolders] Processing folder: $folderName"

    # Run azopaform command for the current folder
    try {
        azopaform -dir "$fullPath"
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  Success: $folderName" -ForegroundColor Green
        } else {
            Write-Host "  Failed: $folderName (Exit code: $LASTEXITCODE)" -ForegroundColor Red
        }
    } catch {
        Write-Host "  Error processing $folderName`: $_" -ForegroundColor Red
    }
    # Run git clean after processing all folders
    Write-Host "`nRunning git clean -fxd" -ForegroundColor Yellow
    git clean -fxd > $null
}


Write-Host "Process completed."