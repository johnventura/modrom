package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// FileDiff is the structure of an individual patch
type FileDiff struct {
	// Start is an offset inside a file where data is changed
	Start int `json:"start"`
	// Offset is added to Start to define where to modify data
	// Sometimes, Offset is necessary when you are near the
	// end of a block of code and you are using shasum
	Offset int `json:"offset"`
	// NewBytes is the new data being inserted
	NewBytes string `json:"newbytes"`
	// Shasum is the SHA256 hash of a block of binary data
	Shasum string `json:"shasum"`
	// Comment should include a description of the change
	Comment string `json:"comment"`
}

// PatchFile allows for multiple differences to be describe in one file
type PatchFile struct {
	// Patch includes an array of FileDiff file changes
	Patch []FileDiff `json:"patch"`
}

// ReadFileToMemory reads a file into memory and returns it as a []byte
func ReadFileToMemory(fileName string) ([]byte, error) {
	fileBuf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return fileBuf, nil

}

// WriteBytesToDisk writes a chunk of memory to a file
func WriteBytesToDisk(fileName string, fileBuf []byte) error {
	err := ioutil.WriteFile(fileName, fileBuf, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetOffsetFromHash returns an offset of a block of data that has a specified sha256 hash
func GetOffsetFromHash(memBuf []byte, inHash []byte) (int, error) {
	parseLimit := len(memBuf) - sha256.BlockSize
	for i := 0; i < parseLimit; i++ {
		shaSum := sha256.Sum256(memBuf[i : i+sha256.BlockSize])
		if bytes.Equal(inHash[:], shaSum[:]) {
			return i, nil
		}
	}
	// we couldn't find the block with the given hash
	return 0, errors.New("Data not found")
}

// GetHashFromOffset returns a sha256 hash for a block of data in a []byte array
func GetHashFromOffset(memBuf []byte, offset int) ([32]byte, error) {
	if (offset + sha256.BlockSize) > len(memBuf) {
		return [32]byte{}, errors.New("Out of bounds for ROM")
	}
	shaSum := sha256.Sum256(memBuf[offset : offset+sha256.BlockSize])
	return shaSum, nil
}

// FindDifference compares two byte arrays and returns difference as []FileDiff array
func FindDifference(bufSrc []byte, bufAlt []byte, outputType int8) []FileDiff {
	diffStart := 0
	diffStartOld := -1
	rDiffs := []FileDiff{}
	// compare each byte.
	for i := 0; i < len(bufSrc); i++ {
		// if the bytes don't match at a given offset, figure out where
		// the difference occured
		if (bufSrc[i] != bufAlt[i]) && diffStart != diffStartOld {
			diffStart = i
			diffStartOld = i
		}
		// calculate the SHA of the block where the change started
		if (bufSrc[i] == bufAlt[i]) && diffStart == diffStartOld {
			diffEnd := i
			newDiff := FileDiff{}
			// determine offset, if we are in the last block
			if i > (len(bufSrc) - sha256.BlockSize) {
				newDiff.Start = len(bufSrc) - sha256.BlockSize
				newDiff.Offset = diffStart - newDiff.Start
			} else {
				// when we are not in the last block use simple offset
				newDiff.Start = diffStart
				newDiff.Offset = 0
			}
			newDiff.NewBytes = hex.EncodeToString(bufAlt[diffStart:diffEnd])
			// get hash of original block
			if (outputType & DiffsWithSHA) != 0 {
				newHash, err := GetHashFromOffset(bufSrc, newDiff.Start)
				if err != nil {
					newDiff.Shasum = ""
				}
				newDiff.Shasum = hex.EncodeToString(newHash[:])
			}
			// if user does not want offsets for some reason
			if (outputType & DiffsWithOffset) == 0 {
				newDiff.Start = 0
				newDiff.Offset = 0
			}
			rDiffs = append(rDiffs, newDiff)
			diffStartOld = -1
		}
	}
	return rDiffs
}

// PatchByteArray reads a parsed Patchfile, and returns byte array copy with patches applied
func PatchByteArray(bufSrc []byte, patchFile PatchFile) ([]byte, error) {
	newBuf := make([]byte, len(bufSrc))
	copy(newBuf, bufSrc)
	for _, patch := range patchFile.Patch {
		newBytes, err := hex.DecodeString(patch.NewBytes)
		if err != nil {
			return newBuf, err
		}
		if patch.Shasum == "" {
			copy(newBuf[patch.Start:(patch.Start+len(newBytes))], newBytes)

		} else {
			for i := 0; i <= len(bufSrc)-sha256.BlockSize; i++ {
				newHash, err := GetHashFromOffset(bufSrc, i)
				if err != nil {
					return newBuf, err
				}
				displayHash := hex.EncodeToString(newHash[:])
				if displayHash == patch.Shasum {
					patchLocation := i + patch.Offset
					copy(newBuf[patchLocation:(patchLocation+len(newBytes))], newBytes)
				}
			}
		}
		if patch.Comment != "" {
			fmt.Printf("%v\n", patch.Comment)
		}

	}
	return newBuf, nil
}

// ReadPatchJson parses a Json formatted patch file and returns a struct
func ReadPatchJson(inJson string) (PatchFile, error) {

	var patchFile PatchFile
	err := json.Unmarshal([]byte(inJson), &patchFile)
	if err != nil {
		return patchFile, err
	}

	return patchFile, nil
}

// cmdDiff reads two filenames and returns a Json formatted diff
func cmdDiff(fileNameSrc string, fileNameAlt string, OutputFormat int8) (string, error) {
	fileBuf, err := ReadFileToMemory(fileNameSrc)
	if err != nil {
		return "", err
	}

	fileBufAlt, err := ReadFileToMemory(fileNameAlt)
	if err != nil {
		return "", err
	}

	diffStr, err := DiffsToJson(fileBuf, fileBufAlt, OutputFormat)
	if err != nil {
		return "", err
	}
	return diffStr, nil
}

// cmdPatch apply binary patches and return patched buffer
func cmdPatch(fileNameSrc string, patchFileName string) ([]byte, error) {
	rBytes := []byte{}

	patchFileBuf, err := ReadFileToMemory(patchFileName)
	if err != nil {
		return rBytes, err
	}
	originalFileBuf, err := ReadFileToMemory(fileNameSrc)
	if err != nil {
		return rBytes, err
	}

	patchFile, err := ReadPatchJson(string(patchFileBuf))
	if err != nil {
		return rBytes, err
	}

	rBytes, err = PatchByteArray(originalFileBuf, patchFile)
	if err != nil {
		return rBytes, err
	}

	return rBytes, nil
}

// these constants describe what patch data should be included
var (
	DiffsWithSHA    int8 = 1
	DiffsWithOffset int8 = 2
	DiffsPretty     int8 = 4
)

// DiffsToJson compares two byte arrays and describes the difference as Json
func DiffsToJson(fileBuf []byte, fileBufAlt []byte, OutputType int8) (string, error) {
	var out bytes.Buffer
	fileDiffs := FindDifference(fileBuf, fileBufAlt, OutputType)
	var patchFile PatchFile
	patchFile.Patch = fileDiffs

	patchJson, err := json.Marshal(patchFile)
	if err != nil {
		return "", err
	}

	rBytes := []byte{}
	if (OutputType & DiffsPretty) == 0 {
		rBytes = patchJson
	} else {
		json.Indent(&out, patchJson, "", "  ")
		rBytes = out.Bytes()
	}

	return string(rBytes), nil
}

// usage helps the user by displaying usage information
func usage(progName string) {
	fmt.Printf("%v is a utility for patching ROMs.\nUsage:\n", progName)
	fmt.Printf("%v [command] filenames...\n", progName)
	fmt.Printf("%v diff [oringial file] [modified file]\n", progName)
	fmt.Printf("\t diff creates a Json formatted description of differences between the files\n")
	fmt.Printf("%v patch [original file] [Json patch file] [output file name]\n", progName)
	fmt.Printf("\t patch reads a Json formatted patch file and creates a modfied copy\n")
	os.Exit(1)
}

func main() {
	ArgC := len(os.Args)
	// user type command with no paramaters
	if ArgC <= 2 {
		usage(os.Args[0])
	}
	command := os.Args[1]

	switch command {
	// user wants to compare two files
	case "diff":
		if ArgC <= 3 {
			usage(os.Args[0])
		}
		// optional commands for patch file output
		diffMode := DiffsWithSHA | DiffsWithOffset
		if ArgC > 4 {
			for _, directive := range os.Args[4:] {
				switch directive {
				case "nosha":
					diffMode = diffMode ^ DiffsWithSHA
				case "nooffset":
					diffMode = diffMode ^ DiffsWithOffset
				case "pretty":
					diffMode = diffMode | DiffsPretty
				default:
					fmt.Printf("unknown directive %v\n", directive)
					usage(os.Args[0])
				}
			}
		}
		if (diffMode & (DiffsWithSHA | DiffsWithOffset)) == 0 {
			fmt.Printf("The \"nosha\" and \"nooffset\" directives are incompatible")
			usage(os.Args[0])
		}
		OriginalFileName := os.Args[2]
		ModFileName := os.Args[3]
		jsonDiffOut, err := cmdDiff(OriginalFileName, ModFileName, diffMode)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
		fmt.Printf("%s", jsonDiffOut)
	// user wants to apply a Json formatted patch
	case "patch":
		OriginalFileName := os.Args[2]
		diffFileName := os.Args[3]
		outBuf, err := cmdPatch(OriginalFileName, diffFileName)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
		outFileName := OriginalFileName + ".patched.bin"
		if ArgC <= 3 {
			outFileName = os.Args[4]
		}

		err = WriteBytesToDisk(outFileName, outBuf)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	default:
		usage(os.Args[0])
	}
}
