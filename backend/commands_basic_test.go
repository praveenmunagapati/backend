package backend

import (
	. "lime/backend/primitives"
	"reflect"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	data := `Hello world
Test
Goodbye world
`
	ed := GetEditor()
	w := ed.NewWindow()
	v := w.NewView()
	e := v.BeginEdit()
	v.Insert(e, 0, data)
	v.EndEdit(e)

	v.Sel().Clear()
	v.Sel().Add(Region{11, 11})
	v.Sel().Add(Region{16, 16})
	v.Sel().Add(Region{30, 30})
	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if v.Buffer().Data() != `Hello worl
Tes
Goodbye worl
` {
		t.Error(v.Buffer().Data())
	}
	ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": "a"})
	if d := v.Buffer().Data(); d != "Hello worla\nTesa\nGoodbye worla\n" {
		lines := strings.Split(v.Buffer().Data(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	v.Settings().Set("translate_tabs_to_spaces", true)
	ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": "\t"})
	if v.Buffer().Data() != "Hello worla \nTesa    \nGoodbye worla   \n" {
		lines := strings.Split(v.Buffer().Data(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}
	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if d := v.Buffer().Data(); d != "Hello worla\nTesa\nGoodbye worla\n" {
		lines := strings.Split(v.Buffer().Data(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if d := v.Buffer().Data(); d != "Hello worl\nTes\nGoodbye worl\n" {
		lines := strings.Split(v.Buffer().Data(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": "\t"})
	if d := v.Buffer().Data(); d != "Hello worl  \nTes \nGoodbye worl    \n" {
		lines := strings.Split(v.Buffer().Data(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if v.Buffer().Data() != "Hello worl\nTes\nGoodbye worl\n" {
		lines := strings.Split(v.Buffer().Data(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

}

func TestLeftDelete(t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()
	v := w.NewView()
	e := v.BeginEdit()
	v.Insert(e, 0, "12345678")
	v.EndEdit(e)
	v.Sel().Clear()
	v.Sel().Add(Region{1, 1})
	v.Sel().Add(Region{2, 2})
	v.Sel().Add(Region{3, 3})
	v.Sel().Add(Region{4, 4})
	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if d := v.buffer.Data(); d != "5678" {
		t.Error(d)
	}
}

func TestMove(t *testing.T) {
	ed := GetEditor()

	w := ed.NewWindow()
	v := w.NewView()
	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)

	type Test struct {
		in      []Region
		by      string
		extend  bool
		forward bool
		exp     []Region
		args    Args
	}

	tests := []Test{
		{
			[]Region{{1, 1}, {3, 3}, {6, 6}},
			"characters",
			false,
			true,
			[]Region{{2, 2}, {4, 4}, {7, 7}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {6, 6}},
			"characters",
			false,
			false,
			[]Region{{0, 0}, {2, 2}, {5, 5}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {10, 6}},
			"characters",
			false,
			true,
			[]Region{{2, 2}, {4, 4}, {7, 7}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {10, 6}},
			"characters",
			false,
			false,
			[]Region{{0, 0}, {2, 2}, {5, 5}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {10, 6}},
			"characters",
			true,
			true,
			[]Region{{1, 2}, {3, 4}, {10, 7}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {10, 6}},
			"characters",
			true,
			false,
			[]Region{{1, 0}, {3, 2}, {10, 5}},
			nil,
		},
		{
			[]Region{{1, 3}, {3, 5}, {10, 7}},
			"characters",
			true,
			true,
			[]Region{{1, 6}, {10, 8}},
			nil,
		},
		{
			[]Region{{1, 1}},
			"stops",
			true,
			true,
			[]Region{{1, 5}},
			Args{"word_end": true},
		},
		{
			[]Region{{1, 1}},
			"stops",
			false,
			true,
			[]Region{{6, 6}},
			Args{"word_begin": true},
		},
		{
			[]Region{{6, 6}},
			"stops",
			false,
			false,
			[]Region{{0, 0}},
			Args{"word_begin": true},
		},
	}
	for i, test := range tests {
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		args := Args{"by": test.by, "extend": test.extend, "forward": test.forward}
		if test.args != nil {
			for k, v := range test.args {
				args[k] = v
			}
		}
		ed.CommandHandler().RunTextCommand(v, "move", args)
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.exp) {
			t.Errorf("Move test %d failed: %v", i, sr)
		}
	}
}

func TestGlueCmds(t *testing.T) {
	ed := GetEditor()

	w := ed.NewWindow()
	v := w.NewView()
	v.SetScratch(true)
	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)
	v.SetScratch(false)
	v.RunCommand("mark_undo_groups_for_gluing", nil)
	v.RunCommand("insert", Args{"characters": "a"})
	v.RunCommand("insert", Args{"characters": "b"})
	v.RunCommand("insert", Args{"characters": "c"})
	v.RunCommand("glue_marked_undo_groups", nil)
	if v.undoStack.position != 1 {
		t.Error(v.undoStack.position)
	} else if d := v.Buffer().Data(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	v.RunCommand("undo", nil)
	if d := v.Buffer().Data(); d != "Hello World!\nTest123123\nAbrakadabra\n" {
		t.Error(d)
	}
	v.RunCommand("redo", nil)
	if d := v.Buffer().Data(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	if v.undoStack.position != 1 {
		t.Error(v.undoStack.position)
	} else if d := v.Buffer().Data(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	v.RunCommand("undo", nil)
	if d := v.Buffer().Data(); d != "Hello World!\nTest123123\nAbrakadabra\n" {
		t.Error(d)
	}

	v.RunCommand("maybe_mark_undo_groups_for_gluing", nil)
	v.RunCommand("insert", Args{"characters": "a"})
	v.RunCommand("maybe_mark_undo_groups_for_gluing", nil)
	v.RunCommand("insert", Args{"characters": "b"})
	v.RunCommand("maybe_mark_undo_groups_for_gluing", nil)
	v.RunCommand("insert", Args{"characters": "c"})
	v.RunCommand("maybe_mark_undo_groups_for_gluing", nil)
	v.RunCommand("glue_marked_undo_groups", nil)
	if v.undoStack.position != 1 {
		t.Error(v.undoStack.position)
	} else if d := v.Buffer().Data(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	v.RunCommand("undo", nil)
	if d := v.Buffer().Data(); d != "Hello World!\nTest123123\nAbrakadabra\n" {
		t.Error(d)
	}
	v.RunCommand("redo", nil)
	if d := v.Buffer().Data(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	if v.undoStack.position != 1 {
		t.Error(v.undoStack.position)
	} else if d := v.Buffer().Data(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
}

func TestInsert(t *testing.T) {
	ed := GetEditor()

	w := ed.NewWindow()
	v := w.NewView()
	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)

	type Test struct {
		in   []Region
		data string
		expd string
		expr []Region
	}

	tests := []Test{
		{
			[]Region{{1, 1}, {3, 3}, {6, 6}},
			"a",
			"Haelalo aWorld!\nTest123123\nAbrakadabra\n",
			[]Region{{2, 2}, {5, 5}, {9, 9}},
		},
		{
			[]Region{{1, 1}, {3, 3}, {6, 9}},
			"a",
			"Haelalo ald!\nTest123123\nAbrakadabra\n",
			[]Region{{2, 2}, {5, 5}, {9, 9}},
		},
	}
	for i, test := range tests {
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": test.data})
		if d := v.buffer.Data(); d != test.expd {
			t.Errorf("Insert test %d failed: %s", i, d)
		}
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.expr) {
			t.Errorf("Insert test %d failed: %v", i, sr)
		}
		v.RunCommand("undo", nil)
	}

}