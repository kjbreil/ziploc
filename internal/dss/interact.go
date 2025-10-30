package dss

import "strings"

type Entries []*Table

func (d *DSS) Get(script string) Entries {
	var entries Entries
	for _, e := range d.Data {
		if strings.EqualFold(e.Script, script) {
			entries = append(entries, e)
		}
	}
	return entries
}

func (d *DSS) Matches(o *Table) bool {
	entries := d.Get(o.Script)
	for _, e := range entries {
		if e.Signature == o.Signature {
			return true
		}

	}
	return false

}

func (d *DSS) Merge(n *DSS) {
	d.Data = append(d.Data, n.Data...)
	d.SIL.View.Data = append(d.SIL.View.Data, n.SIL.View.Data...)
}

func (d *DSS) Overwrite(entries Entries) {
entryLoop:
	for _, e := range entries {
		// if the entry is already  in data replace the signature
		for _, de := range d.Data {
			if strings.EqualFold(de.Script, e.Script) {
				de.Signature = e.Signature
				continue entryLoop
			}
		}
		panic("entry not found cannot continue")
		d.Data = append(d.Data, e)
	}
	d.SIL.View.Data = make([]interface{}, 0, len(d.Data))
	for _, e := range d.Data {
		d.SIL.View.Data = append(d.SIL.View.Data, *e)
	}
}

func (d *DSS) Delete(script string) {
	for i := 0; i < len(d.Data); i++ {
		if strings.EqualFold(d.Data[i].Script, script) {
			d.Data = append(d.Data[:i], d.Data[i+1:]...)
			d.SIL.View.Data = append(d.SIL.View.Data[:i], d.SIL.View.Data[i+1:]...)
			i--
		}
	}
}
