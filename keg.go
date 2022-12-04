package keg

import (
	"bufio"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/fs"
	_fs "github.com/rwxrob/fs"
	"github.com/rwxrob/fs/dir"
	"github.com/rwxrob/fs/file"
	"github.com/rwxrob/keg/kegml"
	"github.com/rwxrob/term"
	"github.com/rwxrob/to"
)

// NodePaths returns a list of node directory paths contained in the
// keg root directory path passed. Paths returns are fully qualified and
// cleaned. Only directories with valid integer node IDs will be
// recognized. Empty slice is returned if kegroot doesn't point to
// directory containing node directories with integer names.
//
// The lowest and highest integer names are returned as well to help
// determine what to name a new directory.
//
// File and directories that do not have an integer name will be
// ignored.
var NodePaths = _fs.IntDirs

var LatestDexEntryExp = regexp.MustCompile(
	`^\* (\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\dZ) \[(.*)\]\(\.\./(\d+)\)$`,
)

// ParseDex parses any input valid for to.String into a Dex pointer.
// FIXME: replace regular expression with pegn.Scanner instead
func ParseDex(in any) (*Dex, error) {
	dex := Dex{}
	s := bufio.NewScanner(strings.NewReader(to.String(in)))
	for line := 1; s.Scan(); line++ {
		f := LatestDexEntryExp.FindStringSubmatch(s.Text())
		if len(f) != 4 {
			return nil, fmt.Errorf("bad line in changes.md: %v", line)
		}
		if t, err := time.Parse(IsoDateFmt, string(f[1])); err != nil {
			return nil, err
		} else {
			if i, err := strconv.Atoi(f[3]); err != nil {
				return nil, err
			} else {
				dex = append(dex, &DexEntry{U: t, T: f[2], N: i})
			}
		}
	}
	return &dex, nil
}

// ReadDex reads an existing dex/changes.md dex and returns it.
func ReadDex(kegdir string) (*Dex, error) {
	f := filepath.Join(kegdir, `dex`, `changes.md`)
	buf, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return ParseDex(buf)
}

// ScanDex takes the target path to a keg root directory returns a
// Dex object.
func ScanDex(kegdir string) (*Dex, error) {
	var dex Dex
	dirs, _, _ := NodePaths(kegdir)
	sort.Slice(dirs, func(i, j int) bool {
		_, iinfo := _fs.LatestChange(dirs[i].Path)
		_, jinfo := _fs.LatestChange(dirs[j].Path)
		return iinfo.ModTime().After(jinfo.ModTime())
	})
	for _, d := range dirs {
		_, i := _fs.LatestChange(d.Path)
		title, _ := kegml.ReadTitle(d.Path)
		id, err := strconv.Atoi(d.Info.Name())
		if err != nil {
			continue
		}
		entry := &DexEntry{U: i.ModTime().UTC(), T: title, N: id}
		dex = append(dex, entry)
	}
	return &dex, nil
}

// MakeDex calls ScanDex and writes (or overwrites) the output to the
// reserved dex node file within the kegdir passed. File-level
// locking is attempted using the go-internal/lockedfile (used by Go
// itself). Both a friendly markdown file reverse sorted by time of last
// update (changes.md) and a tab-delimited file sorted numerically by
// node ID (nodes.tsv) are created.
func MakeDex(kegdir string) error {
	_dex, err := ScanDex(kegdir)
	if err != nil {
		return err
	}

	// remove any empties that might have crept in
	dex := Dex{}
	for _, entry := range *_dex {
		d := filepath.Join(kegdir, entry.ID())
		if dir.IsEmpty(d) {
			log.Println("deleting", d)
			if err := os.RemoveAll(d); err != nil {
				return err
			}
			continue
		}
		dex = append(dex, entry)
	}

	// markdown is first since reverse chrono of updates is default
	mdpath := filepath.Join(kegdir, `dex`, `changes.md`)
	if err := file.Overwrite(mdpath, dex.MD()); err != nil {
		return err
	}

	tsvpath := filepath.Join(kegdir, `dex`, `nodes.tsv`)
	if err := file.Overwrite(tsvpath, dex.ByID().TSV()); err != nil {
		return err
	}

	return UpdateUpdated(kegdir)
}

// UpdateUpdated sets the updated YAML field in the keg info file.
func UpdateUpdated(kegpath string) error {
	kegfile := filepath.Join(kegpath, `keg`)
	updated := UpdatedString(kegpath)
	return file.ReplaceAllString(
		kegfile, `(^|\n)updated:.*(\n|$)`, `${1}updated: `+updated+`${2}`,
	)
}

// Updated parses the most recent change time in the dex/node.md file
// (the first line) and returns the time stamp it contains as
// a time.Time. If a time stamp could not be determined returns time.
func Updated(kegpath string) (*time.Time, error) {
	kegfile := filepath.Join(kegpath, `dex`, `changes.md`)
	str, err := file.FindString(kegfile, IsoDateExpStr)
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(IsoDateFmt, str)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// LastChanged parses and returns a DexEntry of the most recently
// updated node from first line of the dex/changes.md file. If cannot
// determine returns nil.
func LastChanged(kegpath string) *DexEntry {
	kegfile := filepath.Join(kegpath, `dex`, `changes.md`)
	lines, err := file.Head(kegfile, 1)
	if err != nil || len(lines) == 0 {
		return nil
	}
	dex, err := ParseDex(lines[0])
	if err != nil {
		return nil
	}
	return (*dex)[0]
}

// Last returns the last created content node. If cannot determine
// returns nil
func Last(kegpath string) *DexEntry {
	dex, err := ReadDex(kegpath)
	if err != nil {
		return nil
	}
	return dex.Last()
}

// Next returns a new DexEntry with its integer identify set to the next
// integer after Last and returns nil if cannot determine which is next.
// The updated time stamp is set to the current time even though the
// DexEntry may not have yet been written to disk and its time would be
// different from the actual time written. This is to save the overhead
// of grabbing it again once written.
func Next(kegpath string) *DexEntry {
	last := Last(kegpath)
	if last == nil {
		return nil
	}
	entry := new(DexEntry)
	entry.U = time.Now().UTC()
	entry.N = last.N + 1
	return entry
}

// UpdatedString returns Updated time as a string or an empty string if
// there is a error.
func UpdatedString(kegpath string) string {
	u, err := Updated(kegpath)
	if err != nil {
		log.Println(err)
		return ""
	}
	return (*u).Format(IsoDateFmt)
}

// Publish publishes the keg at kegpath location to its distribution
// targets listed in the keg file under "publish." Currently, this only
// involves looking for a .git directory and if found doing a git
// pull/add/commit/push. Git commit messages are always based on the
// latest node title without any verb.
func Publish(kegpath string) error {
	gitd, err := fs.HereOrAbove(`.git`)
	if err != nil {
		return nil
	}
	origd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(origd)
	os.Chdir(filepath.Dir(gitd))
	if err := Z.Exec(`git`, `-C`, kegpath, `pull`); err != nil {
		if _, is := err.(*exec.ExitError); is {
			return fmt.Errorf(
				"%vNo remote repo has been setup.%v First create it and git push to it.",
				term.Red, term.X,
			)
		}
	}
	if err := Z.Exec(`git`, `-C`, kegpath, `add`, `-A`, `.`); err != nil {
		return err
	}
	msg := "Publish changes"
	if n := Last(kegpath); n != nil {
		msg = n.T
	}
	if err := Z.Exec(`git`, `-C`, kegpath, `commit`, `-m`, msg); err != nil {
		return err
	}
	return Z.Exec(`git`, `-C`, kegpath, `push`)
}

// MakeNode examines the keg at kegpath for highest integer identifier
// and provides a new one returning a *DexEntry for it.
func MakeNode(kegpath string) (*DexEntry, error) {
	_, _, high := NodePaths(kegpath)
	if high < 0 {
		high = 0
	}
	high++
	path := filepath.Join(kegpath, strconv.Itoa(high))
	if err := dir.Create(path); err != nil {
		return nil, err
	}
	readme := filepath.Join(kegpath, `dex`, `README.md`)
	if err := file.Touch(readme); err != nil {
		return nil, err
	}
	return &DexEntry{N: high}, nil
}

// Edit calls file.Edit on the given node README.md file within the
// given kegpath.
func Edit(kegpath string, id int) error {
	node := strconv.Itoa(id)
	if node == "" {
		return fmt.Errorf(`node (%q) is not a valid node id`, id)
	}
	readme := filepath.Join(kegpath, node, `README.md`)
	return file.Edit(readme)
}

// DexUpdate first checks the keg at kegpath for an existing
// dex/changes.md file and if found loads it, if not, MakeDex is called
// to create it. Then DexUpdate examines the Dex for the DexEntry passed
// and if found updates it with the new information, otherwise, it will
// add the new entry without any further validation. The updated Dex is
// then written to the dex/changes.md file.
func DexUpdate(kegpath string, entry *DexEntry) error {

	if !HaveDex(kegpath) {
		if err := MakeDex(kegpath); err != nil {
			return err
		}
	}

	if err := entry.Update(kegpath); err != nil {
		return err
	}

	dex, err := ReadDex(kegpath)
	if err != nil {
		return err
	}

	found := dex.Lookup(entry.N)
	if found == nil {
		dex.Add(entry)
	} else {
		found.U = entry.U
		found.T = entry.T
	}

	return WriteDex(kegpath, dex)
}

// HaveDex returns true if keg at kegpath has a dex/changes.md file.
func HaveDex(kegpath string) bool {
	return file.Exists(filepath.Join(kegpath, `dex`, `changes.md`))
}

// WriteDex writes the dex/changes.md and dex/nodes.tsv files to the keg
// at kegpath and calls UpdateUpdated to keep keg info file in sync.
func WriteDex(kegpath string, dex *Dex) error {
	changes := filepath.Join(kegpath, `dex`, `changes.md`)
	nodes := filepath.Join(kegpath, `dex`, `nodes.tsv`)
	if err := file.Overwrite(changes, dex.ByChanges().MD()); err != nil {
		return err
	}
	if err := file.Overwrite(nodes, dex.ByID().TSV()); err != nil {
		return err
	}
	return UpdateUpdated(kegpath)
}

//go:embed testdata/samplekeg/1/README.md
var SampleNodeReadme string

// WriteSample writes the embedded SampleNodeReadme to the entry
// indicated in the keg specified by kegpath.
func WriteSample(kegpath string, entry *DexEntry) error {
	return file.Overwrite(
		filepath.Join(kegpath, entry.ID(), `README.md`),
		SampleNodeReadme,
	)
}

// Import imports the targets into the kegpath creating new, unique
// identifiers for each. If the target ends with an integer it is
// assumed to be a node directory. If not, it is assumed to contain node
// directories with integer identifiers. Currently, there is no
// resolution of any links contained within any node README.md file.
func Import(kegpath string, targets ...string) error {
	if !fs.IsDir(kegpath) {
		return fmt.Errorf("not a directory or does not exist: %v", kegpath)
	}
	for _, target := range targets {
		if fs.NameIsInt(target) {
			if err := ImportNode(kegpath, target); err != nil {
				return err
			}
			continue
		}
		dirs, _, _ := fs.IntDirs(target)
		for _, dir := range dirs {
			if err := ImportNode(kegpath, dir.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

// ImportNode imports a single specific directory into the kegpath by
// getting the next integer identifier and moving the target into the
// kegpath with an os.Rename (which has limitations based on the host
// operating system's handling of cross-file system boundaries).
func ImportNode(kegpath, target string) error {
	var err error

	next := Next(kegpath)
	if next == nil {
		return fmt.Errorf(`could not determine next node id for %v`, target)
	}

	next.T, err = kegml.ReadTitle(filepath.Join(target, `README.md`))
	if err != nil {
		return err
	}

	if err := os.Rename(target, filepath.Join(kegpath, next.ID())); err != nil {
		return err
	}

	return DexUpdate(kegpath, next)
}
