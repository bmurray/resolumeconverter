# Resolume Converter

This is a bit of a hack to quickly bulk-add video files to Resolume Arena, for use with Denon DJ. Instead of clicking a dozen things just to import files, you can run this command to bulk add entire directories of files quickly. 

For Resolume's official support for Denon DJ players, see [their official site](https://resolume.com/support/en/sync-to-denon-players)

## Usage

Install [ffmpeg](https://ffmpeg.org) and ffprobe. It doesn't really matter where, as long as it's in your path. `echo $PATH` to find candidate locations. [or click here](https://google.com/search?q=how+to+install+ffmpeg)  

**ffprobe** is critical. It is used to read the metadata. 

See the Tips and Tricks below to make performing easier on Resolume. 

Follow the official Resolume guide above to link your deck to specific layers. 
   
1. Enable the Resolume HTTP API. This is only needed during import and not during performance. Arena settings -> Webserver -> Enable Webserver. Use IP 0.0.0.0, port 8089. This is hard coded in the command, and I can't be bothered to fix that. YMMV. 
2.  Pick a layer that you want to splat your video files on. Remember this layer. They're 1 indexed; the bottom layer is 1, second is 2, etc.
3.  Dump your mp4 files into a directory. I personally get mine from  Xtendamix, but any should work. Note; they MUST contain the title metadata. Use ffprobe to verify. 
4.  Run the command: 
    `./converter input-audio <*dir with your mp4 files*> <*dir where you want to store your m4a audio files for Engine*>`
5. Wait. A while. It's stripping the audio out of your music videos so Engine can read them. Don't worry, it's only copying the audio, so you won't lose quality. 
6. Open Resolume Alley. Drag all of your mp4 music videos into Alley. Click Convert at the bottom. **Uncheck Audio** (This is **critical**). Change the output folder to somewhere that you can access from Resolume Arena. Use the DXV3 codec for the best performance. This will work with other, smaller codecs, but gets "jumpy", so buy another drive and use DXV3. 
7. Click Queue. And wait. Even longer this time. It's converting your files to an optimal file format. But, it can do hundreds at a time, so if you have a lot of files, go get a coffee.
8. You're now ready to import everything. Run the command `./converter import <*dir where you exported your Resolume dxv3 files*> <*layer*>`
9.  Import all of your m4a files into Engine. Add a beatgrid, and transfer them to your Engine DJ gear. **Caution:** Changing the Title can BREAK the association. Try at your peril. This works over a network connection to your desktop version of Engine, or USB, or internal disk. And probably others. 
10. ...
11. Profit. 

### Alternate method
Once you're at stage 8, you CAN just drag all of the video files into Resolume. But, once imported, you need to select them all, right click, select Transport -> Denon DJ. Right click again and select Target -> Denon Player Determined. You can skip step 8. Note: this method does NOT prevent you from adding the same video file more than once and causing all kinds of havok. The import command checks for existing instances of the file in the composition and skips them if they exist. 

## That sounds complicated

I don't make the rules. Suggest an enhancement to Resolume. They seem like cool dudes. 

## Why don't you just use Rekordbox / etc, that have video built in?

I enjoy the pain. 

## Tips and Tricks

On the toolbar (with BPM/Resync/Etc), right click your player and select a layer for each deck. These should be dedicated to each deck. If you have 4 decks like on a Prime 4, dedicate 4 layers. 

I personally group track layers inside a group and roll them up. I don't generally care which tracks are where, and this is the most compact way to hide them. It also gives you one-click access to Bypass/Solo all music videos.

I dedicated one layer per deck. This will "steal" videos from other layers if that layer is playing it. Eg, if you have Song A on layer 4, but you play the song on a deck tied to Layer 5, it will play the video on Layer 5. This makes it super easy to Video DJ with music videos. 

I set a Auto-Size Fill onto each of the video layers. You may opt for Fit. You may also select Stretch if you want drunk people badgering you all night about how the video looks awful. That's up to you. 

## Why does this work? 

Resolume does matching based on the "ID3" tags of the audio. Kinda. But also not really at all. Audio Video files contain a bunch of metadata like the Title, Artist, etc. Often these titles do not match with the filename, which is why it sometimes works to just drag video files into Resolume, but usually doesn't. When you import a file into Resolume, it sets a name tag to filename. When you import the file into Engine, it disregards the filename for the most part, and uses the embeded ID3 tags. When you play the file on Engine, it transmits the ID3 data over StageLinq. Resolume then tries to match it with a filename, and if it matches, then it uses that as a Denon sync'ed track. BUT, if your filename started with an artist name, and then the track, it won't match. Too bad. 

This solves that issue in two steps. First, the `input` command (`./converter input <*path to mp4 files*>`), renames all of the MP4 files to match the Title field of the metadata. 

The next step, `audio` (`./convert audio <*path to MP4 files*> <*audio storage folder*>`) converts the audio by simply copying it. There is NO transcoding, only audio copying. There will be NO audio quality loss. If it has already been converted, it won't convert it again, so it's safe to run this as often as you want on the same directory. 

*note: the `import-audio` command does both of these steps, so it's marginally easier.*

Finally, the `import` step uses the Resolume API to drop the files into Resolume. This can be done manually in bulk too. However, the import command checks to see if the file is already in the composition, and skips it if it already is. This makes it safe to run it multiple times in a row as you add more files. 

## Support / Waranty / Contributing

None. Zero. This is a free project I whipped up in an evening on stream while trying to figure it out. If you have a suggestion for it, submit a pull request and I'll probably blindly merge it. 

"But I get XYZ error". No idea. Might be ffmpeg not installed. Might be resolume webserver disabled. It could be the impending heat death of the universe. 