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

func (dm diskMap) checksum() int {
	sum, pos := 0, 0
	for e := dm.first; ; e = e.next {
		for i := 0; i < e.fileBlocks; i++ {
			sum += (i + pos) * e.fileID
		}
		pos += e.fileBlocks
		if e.next == nil {
			break
		}
	}
	return sum
}

func (dm diskMap) fullyDefragged() bool {
	for e := dm.first; e.next != nil; e = e.next {
		if e.emptyBlocks > 0 && e.next != nil {
			return false
		}
	}
	return true
}

func (dm diskMap) defragOnce() bool {
	if dm.fullyDefragged() {
		return true
	}
	var first, last, prevLast *diskMapEntry
	for last = dm.first; last.next != nil; last = last.next {
		prevLast = last
	}
	for first = dm.first; first.emptyBlocks == 0; first = first.next {
	}
	// log.Printf("first:%v last:%v", first, last)

	if first.next.fileID != last.fileID || first.emptyBlocks > 0 {
		// insert new entry
		existingNext := first.next
		first.next = &diskMapEntry{
			fileID:      last.fileID,
			fileBlocks:  0,
			emptyBlocks: first.emptyBlocks,
			next:        existingNext,
		}
		first.emptyBlocks = 0
	}
	first.next.fileBlocks++
	first.next.emptyBlocks--
	last.fileBlocks--
	last.emptyBlocks++
	if last.fileBlocks == 0 {
		prevLast.next = nil
		prevLast.emptyBlocks += last.emptyBlocks
	}
	/*
		if last.fileID == prevLast.fileID {
			// Coalesce the last two blocks.  todo; Probably need to do this throughout
			prevLast.fileBlocks += last.fileBlocks
			prevLast.emptyBlocks += last.emptyBlocks
			prevLast.next = nil
		}
	*/
	return false
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

	for i := 0; ; i++ {
		log.Printf("Iter %d: read %v", i, entries)
		if entries.defragOnce() {
			break
		}
	}

	log.Printf("Checksum: %d", entries.checksum())
}
