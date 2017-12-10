# quickpic-cloud-linux-downloader

Quickpic has ability easily backup photos to the cloud and restore
them back to phone/tablet. But sometimes may be useful to restore the
photos directly to PC. The cloud used by Quickpic offers such feature
but requires use their own downloader client for Windows. You can
avoid it. Just open the file with the list of URLs that you can get as
described
at
[cloud.cmcm.com/p/offline/download_description.html](https://cloud.cmcm.com/p/offline/download_description.html),
open it and load the URLs with `curl` or `wget` in a `bash`
loop. There is a simple script in Go that does this thing for you.

