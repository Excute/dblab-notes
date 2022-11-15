$fullFileName=$args[0]
# $csvFileName=$args[1]

$fileName=$fullFileName.Substring(0, $fullFileName.LastIndexOf("."))
$extName=$fullFileName.Substring($fullFileName.LastIndexOf("."))

# echo $fileName $extName

$i=1

Import-Csv "$($fileName).csv" | ForEach-Object {
	$HH=[Math]::Floor([int]$_.time/10000)
	$MM=[Math]::Floor([int]$_.time%10000/100)
	$SS=[Math]::Floor([int]$_.time%100)

	$endTime=(Get-Date -Hour $HH -Minute $MM -Second $SS)
	$startTime=$endTime.AddSeconds(-10)

	# echo "ffmpeg -i $($fullFileName) -ss $($startTime.ToString("HH:mm:ss")) -t 00:00:10 -c:v copy -c:a copy $($fileName)-$($startTime.ToString("HHmmss"))-$($endTime.ToString("HHmmss"))$($extName)"
	# echo "mpv $($fullFileName) --no-audio --start=$($startTime.ToString("HH:mm:ss")) --frames=1 --sid=no -o $($fileName.Substring($fileName.LastIndexOf("-"))+1)-01.png"
	$aCarDir="./$($fileName)/$(([string]$i).PadLeft(4,"0"))"
	New-Item $aCarDir -ItemType Directory

	ffmpeg -i "$($fullFileName)" -ss $($startTime.ToString("HH:mm:ss")) -t 00:00:10 -c:v copy -c:a copy "./$($aCarDir)/$($fileName)-$($startTime.ToString("HHmmss"))-$($endTime.ToString("HHmmss"))$($extName)"
	mpv "$($fullFileName)" --no-audio --start=$($endTime.ToString("HH:mm:ss")) --frames=1 --sid=no -o "./$($aCarDir)/01.png"
	
	$i=$i+1
}