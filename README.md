# modrom

Re-mix classic video games (and other software) while respecting intellectual property rights!

version 0.1

This program allows developers to share patches for ROM files without sharing any part of the original ROM files themselves.
Many people, including several people in the legal and software development communities, believe that some ROM files are owned by individuals, corporations, or other groups. Those people believe that the distribution of those ROM files is limited to their owners.
This utility allows for the sharing of patches to ROM files in a way that does not include any of the original data.

This program was originally developed as part of an effort to expand game play for Atari 2600 ROMs to include their modification or "remixing."
It may have other uses as well.

Developers have a choice to define a patch with an offset from the beginning of a file, or they can have the patching utility search for a block of data based on a SHA256 hash. This latter method may allow developers to create a single patch that can be applied to multiple variants of a single piece of software.

## Examples

### Recording Changes
The modrom utility can document changes applied to a ROM file:

* modrom diff polybius.bin polybius.fixed.bin
This command will create a Json formatted description of the changes between a file named polybius.bin and a slightly different file named polybius.fixed.bin.

* modrom diff polybius.bin polybius.fixed.bin nosha
This command will also create a Json formatted description of the changes between these files, but it will not include the SHA256 hashes of the original blocks of data that were changed. 

* modrom diff polybius.bin polybius.fixed.bin nooffset
This command will also create a Json formatted description of the changes between these files, but it will not include the offsets of the original blocks of data that were changed.

### Applying Changes
The modrom utility can read a Json formatted description of changes made to a ROM file and create a new ROM with the changes applied.

* modrom patch dinnerwithandre.rom patches.json dinnerwithandre.mod.rom
This command will apply the patches described in "patches.json" to a ROM file named "dinnerwithandre.rom" and apply create a new file that will be identical to the original file except where specified in the "patches.json" patch file.

* modrom patch dinnerwithandre.rom patches.json
This command will also apply the changes to the "dinnerwithandre.rom" file, but it will create a new file named "mydinnerwithandre.rom.patched.bin". The utility should not modify the original ROM, unless instructed to do so. Please do not modify the original ROM.


## File Format
This utility records differences between original ROMs and modified ROMs as Json formatted text. The standard for these files includes an array labeled "patch" at the top level. This "patch" array includes several fields that describe each change. Those fields include:

* start - This is a numerical offset from the beginning of he original file where a block of data has been modified.
* offset - Although the "start" field defines the beginning of a block of changed code, the changes may not begin with the first byte of that block. The offset value defines the number of bytes after the start of the block where new data is substituted as part of a patch.
* newbytes - This field is an ASCII representation of a hexadecimal block defining new data to be written over the original data as part of a patch.
* shasum - This is a SHA256 hash of the original block of data. 
* comment - This is an optional field, which developers should use to include descriptive comments about a patch.
### Example Json
```
{
  "patch": [
    {
      "start": 412,
      "offset": 0,
      "newbytes": "EAEAEAEA",
      "shasum": "843c295955090eb4b25b449f8b03b9ab88027d83ca888209eab597fc5149a67c",
      "comment": "This change stops the explosion flicker effect."
    },
    {
      "start": 1640,
      "offset": 46,
      "newbytes": "636c6f73",
      "shasum": "88cbe53cbf609c75fbc3ac3aa1310fd25f73d7db82db80773df9203dbb2a9b99",
      "comment": "This change modifies a sprite to adhere to community standards of decency."
    }
  ]
}
```
