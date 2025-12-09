# PowerShell script to create new migration files
# Usage: .\scripts\create-migration.ps1 migration_name

param(
    [Parameter(Mandatory=$true)]
    [string]$MigrationName
)

$MigrationsDir = ".\migrations"

# Get the highest migration number
$HighestNumber = Get-ChildItem $MigrationsDir -Filter "*.sql" | 
    ForEach-Object { 
        if ($_.Name -match '^(\d+)_') { 
            [int]$matches[1] 
        } 
    } | 
    Measure-Object -Maximum | 
    Select-Object -ExpandProperty Maximum

if ($null -eq $HighestNumber) {
    $NewNumber = "000001"
} else {
    $NewNumber = ($HighestNumber + 1).ToString("D6")
}

$UpFile = Join-Path $MigrationsDir "${NewNumber}_${MigrationName}.up.sql"
$DownFile = Join-Path $MigrationsDir "${NewNumber}_${MigrationName}.down.sql"

# Create migration files
"-- Add your UP migration here" | Out-File -FilePath $UpFile -Encoding UTF8
"-- Add your DOWN migration here" | Out-File -FilePath $DownFile -Encoding UTF8

Write-Host "Created migration files:" -ForegroundColor Green
Write-Host "  $UpFile"
Write-Host "  $DownFile"
