$i=1
Get-ChildItem -Path .\* -Include *.jpg | ForEach-Object {
	New-Item -ItemType "directory" -Name $(([string]$i).PadLeft(4,"0"))
	Move-Item -Path $_.Name -Destination $(([string]$i).PadLeft(4,"0"))
	$i=$i+1
}