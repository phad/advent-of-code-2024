package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

/* input format
2333133121414131402
represents
00...111...2...333.44.5555.6666.777.888899
*/

type diskMapEntry struct {
	fileID      int
	fileBlocks  int
	emptyBlocks int
	next        *diskMapEntry
}

func (dme *diskMapEntry) String() string {
	if dme.next != nil {
		return fmt.Sprintf("<id:%d #blk:%d #emp:%d id_next:%d>", dme.fileID, dme.fileBlocks, dme.emptyBlocks, dme.next.fileID)
	}
	return fmt.Sprintf("<id:%d #blk:%d #emp:%d (last)>", dme.fileID, dme.fileBlocks, dme.emptyBlocks)
}

type diskMap struct {
	first *diskMapEntry
}

func (dm diskMap) String() string {
	var s strings.Builder
	for e := dm.first; ; e = e.next {
		s.WriteString(strings.Repeat(fmt.Sprintf("%d", e.fileID), e.fileBlocks))
		s.WriteString(strings.Repeat(".", e.emptyBlocks))
		if e.next == nil {
			break
		}
	}
	return s.String()
}

func tri(n int64) int64 {
	return n * (n + 1) / 2
}

func (dm diskMap) checksum() int64 {
	sum, pos := int64(0), int64(-1)
	prevID := -1
	for e := dm.first; ; e = e.next {
		if e.fileID == prevID {
			log.Printf("found two adjacent entries with same file ID %d!", prevID)
		}
		prevID = e.fileID
		sum += int64(e.fileID) * (tri(pos+int64(e.fileBlocks)) - tri(pos))
		pos += int64(e.fileBlocks)
		if e.next == nil {
			break
		}
	}
	return sum
}

func (dm diskMap) fullyDefragged() bool {
	for e := dm.first; ; e = e.next {
		if e.emptyBlocks > 0 && e.next != nil {
			return false
		}
		if e.next == nil {
			break
		}
	}
	return true
}

func (dm diskMap) defragOnce() bool {
	if dm.fullyDefragged() {
		return true
	}
	var first, last, prevLast *diskMapEntry
	// Find the last entry, and the one previous to it
	for last = dm.first; last.next != nil; last = last.next {
		prevLast = last
	}
	// Find the first entry with empty blocks to fill
	for first = dm.first; first.emptyBlocks == 0; first = first.next {
	}
	//log.Printf("first:%v first.next:%v prevLast:%v last:%v", first, first.next, prevLast, last)

	// Otherwise prepare to move a file block from the last entry to
	// either 'first' (with empty, if the file ID matches), or to insert
	// a block (if the file ID doesn't match).
	var dest *diskMapEntry
	if first.fileID == last.fileID && first.emptyBlocks > 0 {
		dest = first
	} else {
		//log.Printf("Inserting new file entry before first.next=%v for file ID %d", first.next, last.fileID)
		// insert new entry
		existingNext := first.next
		first.next = &diskMapEntry{
			fileID:      last.fileID,
			fileBlocks:  0,
			emptyBlocks: first.emptyBlocks,
			next:        existingNext,
		}
		first.emptyBlocks = 0
		dest = first.next
	}
	// When we reach the end, we need to shuffle the prevLast's empty blocks
	// over to last's empties.
	if first == prevLast {
		last.emptyBlocks += prevLast.emptyBlocks
		prevLast.emptyBlocks = 0
	}

	// Adjust the block counts in the block we're moving the entry to, and from.
	//log.Printf("Moving a block for file ID=%d\nfirst=%v\nfirst.next=%v\nlast=%v\n", first.next.fileID, first, first.next, last)
	if dest != nil {
		dest.fileBlocks++
		dest.emptyBlocks--
		last.fileBlocks--
		last.emptyBlocks++
	}
	if last.fileBlocks == 0 {
		prevLast.next = nil
		prevLast.emptyBlocks += last.emptyBlocks
	}

	return false
}

type counts struct {
	file, empty int
}

func (dm diskMap) fileSummary() map[int]counts {
	summ := make(map[int]counts)
	for e := dm.first; ; e = e.next {
		c, ok := summ[e.fileID]
		if !ok {
			c = counts{}
			summ[e.fileID] = c
		}
		c.file += e.fileBlocks
		c.empty += e.emptyBlocks
		summ[e.fileID] = c
		if e.next == nil {
			break
		}
	}
	return summ
}

func (dm diskMap) diskSummary() counts {
	fileSumm := dm.fileSummary()
	var total counts
	for _, c := range fileSumm {
		total.file += c.file
		total.empty += c.empty
	}
	return total

}

func main() {
	log.Println("AoC-2024-day09-part1")
	if len(os.Args) < 2 {
		log.Fatal("Usage: main <in file>")
	}
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	if len(lines) != 1 {
		log.Fatalf("Too many input lines, got %d want 1", len(lines))
	}

	serializedDiskMap := lines[0]
	var entries diskMap

	var prev *diskMapEntry
	for i := 0; i <= len(serializedDiskMap); i += 2 {
		ce := &diskMapEntry{
			fileID:     i / 2,
			fileBlocks: int(mustParseInt(serializedDiskMap[i : i+1])),
		}
		if i < len(serializedDiskMap)-1 {
			ce.emptyBlocks = int(mustParseInt(serializedDiskMap[i+1 : i+2]))
		}
		if prev == nil {
			entries.first = ce
		} else {
			prev.next = ce
		}
		prev = ce
	}

	//log.Printf("Before:\nfsummary: %v\ndsummary: %v", entries.fileSummary(), entries.diskSummary())
	for i := 0; ; i++ {
		//log.Printf("\n\nIter %d: read\n%v\nfsummary: %v\ndsummary: %v", i, entries, entries.fileSummary(), entries.diskSummary())
		if entries.defragOnce() {
			break
		}
	}
	//log.Printf("Final:\n%v", entries)
	//log.Printf("After:\nfsummary: %v\ndsummary: %v", entries.fileSummary(), entries.diskSummary())
	log.Printf("Checksum: %d", entries.checksum())
}
